package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

func GetNetworkIP() error {
	log.Println("Network IPv4 addresses:")
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP

			}
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			log.Println(ip.String())
		}
	}
	return nil
}

// Middleware to handle CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, you can restrict this
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func defineRoutes(r *mux.Router, state *GameState) {
	/*
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "PLEASE don't fix IISc wifi security , it'll be so funny")
		})
	*/

	r.HandleFunc("/{game_type}/{server_type}/{component}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		log.Println(vars)
		switch vars["game_type"] {
		case "turn-based":
			switch vars["server_type"] {
			case "http":
				switch vars["component"] {
				case "server":
					TurnBasedHTTPServer(w, r, state)

				}

			}

		}

	})

}

func HandleVue(r *mux.Router) {
	//fs := http.FileServer(http.Dir("vue/dist"))
	//log.Println(fs)
	//http.Handle("/", fs)
	log.Println("reached handleVue")
	directoryPath := "./vue/dist"

	// Check if the directory exists
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		fmt.Printf("Directory '%s' not found.\n", directoryPath)
		return
	}

	fileServer := http.FileServer(http.Dir(directoryPath))
	// Create a new HTTP server and handle requests
	// r.Handle("/turn-based/http/players/", http.StripPrefix("/turn-based/http/players/", fileServer))

	// Standard StripPrefix does not have support for Pattern Matching, this custom function adds support for it.
	r.Handle("/turn-based/http/players/{id}", PatternStripPrefix("/[^/]+/[^/]+/players/[0-9]+", fileServer))
	//sr := r.PathPrefix("/players/{id}").Subrouter()
	//r.ServeFile("/", fileServer)

	//	do this for gorilla mux
	//	r.Handle("/test/", http.StripPrefix("/test/", fs))
	//	r.Handle("/test/{files}", http.StripPrefix("/test/", fs))

	fileServerStatic := http.FileServer(http.Dir("./vue/dist/static"))
	// in vue.config, if publicPath is '', then static files are served under /.../.../static/stuff
	// r.Handle("/players/static/{files}/{more}", http.StripPrefix("/players/static/", fileServerStatic))

	// if publicPath is '/', then static files are served under /static

	//r.Handle("/static/", http.StripPrefix("/static/", fileServerStatic))
	//r.Handle("/static/{f}/", http.StripPrefix("/static/", fileServerStatic))
	//r.Handle("/static/{f}/{g}", http.StripPrefix("/static/", fileServerStatic))
	//same as: (but files are not opening, file in the root is opening
	//)

	r.PathPrefix("/static").Handler(http.StripPrefix("/static", fileServerStatic))

	//this is working perfectly
	// meaning there was something in routes blocking the functioning, we moved this function before
	//the other routes function and its working
	//r.PathPrefix("/").Handler(http.StripPrefix("/", fileServerStatic))

}
func PatternStripPrefix(pattern string, h http.Handler) http.Handler {
	re := regexp.MustCompile(pattern)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Match the pattern with the request URL path
		matches := re.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			// Remove the matched prefix from the URL path
			newPath := strings.TrimPrefix(r.URL.Path, matches[0])
			if newPath == "" {
				newPath = "/"
			}
			// Update the request URL
			r.URL.Path = newPath
		}
		// Pass the request to the next handler
		h.ServeHTTP(w, r)
	})
}

func TurnBasedHTTPServer(w http.ResponseWriter, r *http.Request, state *GameState) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, you can restrict this
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	log.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		err := json.NewEncoder(w).Encode(state)
		if err != nil {
			panic(err)
		}
	case http.MethodPost:
		var sq ServerQuery

		err := json.NewDecoder(r.Body).Decode(&sq)
		if err != nil {
			panic(err)
		}
		//json.Unmarshal(r.Body,sq)
		log.Println(sq)
		switch sq.Action {
		case LEFT:
			PlayerAction(w, r, state, LEFT, sq.Player)
		case RIGHT:
			PlayerAction(w, r, state, RIGHT, sq.Player)
		case UP:
			PlayerAction(w, r, state, UP, sq.Player)
		case DOWN:
			PlayerAction(w, r, state, DOWN, sq.Player)

		}

	}

}

func startHttpServer(r *mux.Router) *http.Server {
	log.Println("Starting server on :8080")
	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}
func stopHttpServer(srv *http.Server) {
	log.Println("Stopping http server")
	srv.Shutdown(context.Background())
}
func HandleHTTPServer(r *mux.Router) {
	server := startHttpServer(r)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Server is waiting until an interrupt (SIGINT or SIGTERM) is done")
	sig := <-sigChan
	fmt.Printf("Received signal: %s\n", sig)
	stopHttpServer(server)

}
