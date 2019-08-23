<template>
    <div id="app">
        <router-view />
    </div>
</template>

<style>
html,
body,
#app {
    height: 100%;

    display: flex;
    flex-direction: column;
}
body {
    font-family: "Avenir", Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin: 0;
    padding: 0;

    flex-shrink: 0;
}
</style>
