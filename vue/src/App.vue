<template>
  <ZekromServer v-if = "isServer" />
  <ZekromPlayer v-if="!isServer"/>
</template>

<script setup>
import ZekromPlayer  from "@/components/ZekromPlayer.vue";
import ZekromServer from "@/components/ZekromServer.vue";


const routePath = window.location.pathname
const segments = routePath.split("/")
var isServer = null
if (segments[segments.length - 1] === "http") {
  isServer = true
}else{
  isServer = false
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
