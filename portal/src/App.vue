<template>
    <div id="app">
        <router-view />
    </div>
</template>

<style>
html,
body,
#app {
    min-width: 1100px;
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
a {
    text-decoration: none;
    color: inherit;
}
.main-panel {
    width: auto;
    text-align: left;
    overflow: hidden;
    color: rgb(41, 61, 81);
}
.main-panel h1 {
    font-size: 38px;
    margin-top: 40px;
}
.main-panel .container {
    min-width: 800px;
    padding-bottom: 40px;
    border-bottom: 2px solid rgb(235, 235, 235);
}
</style>
