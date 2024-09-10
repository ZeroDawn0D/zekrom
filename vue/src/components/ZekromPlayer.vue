<script setup>
import {ref} from 'vue'


const player = ref(1)
const last_turn = ref(0)
//const boardPosition = ref( [ [1,0,0],[0,0,0],[0,0,0] ])
const turn_message = ref("It's your turn")
const routePath = window.location.pathname
const segments = routePath.split("/")
player.value = parseInt(segments[segments.length - 1])

var isPlayerTurn = true

const url = window.location.hostname
var port = window.location.port
port = 8080
const server = "http://" + url + ":" + port + "/turn-based/http/server"

async function playerAction(action){
  var serverQuery = {
    action: -1,
    player: player.value,
  }
  switch (action){
    case 'LEFT':
      serverQuery.action = 0
      break
    case 'RIGHT':
      serverQuery.action = 1
      break;
    case 'UP':
      serverQuery.action = 2
      break
    case 'DOWN':
      serverQuery.action = 3
      break
  }
  if(isPlayerTurn){
    const response = await fetch(server, {
      method: "POST",
      body:JSON.stringify(serverQuery)})
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`)
    }
    const json = await response.json()
    console.log(json)

  }
}

</script>

<template>
  Player: {{player}}
  <br />
  Last Turn: {{last_turn}}
  <br/>
  {{turn_message}}
  <br/>
  {{segments}}
  <br/>
  {{server}}
  <div class="grid-container">
    <div class="grid-item">1</div>
    <div class="grid-item direction" @click = "playerAction('UP')">UP</div>
    <div class="grid-item">3</div>
    <div class="grid-item direction" @click = "playerAction('LEFT')">LEFT</div>
    <div class="grid-item">5</div>
    <div class="grid-item direction" @click = "playerAction('RIGHT')">RIGHT</div>
    <div class="grid-item">7</div>
    <div class="grid-item direction" @click = "playerAction('DOWN')">DOWN</div>
    <div class="grid-item">9</div>
  </div>
</template>

<style scoped>
.grid-container {
  display: grid;
  position: relative;
  left: 50%;
  transform: translateX(-50%);
  grid-template-columns: repeat(3, 100px);
  grid-template-rows: repeat(3, 100px);
  gap: 5px; /* Optional: space between grid items */
  width: 315px; /* Adjust based on your grid size */
  height: 315px; /* Adjust based on your grid size */
}
.grid-item {
  border: 0px solid black;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;

}

.direction{
  border: 1px solid black;
}
.direction:hover {
  background-color: #f42069;
}
.direction:active {
  background-color: #b4da66;
}
</style>