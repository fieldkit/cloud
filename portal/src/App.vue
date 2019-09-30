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
a {
    text-decoration: none;
    color: inherit;
}
.header {
    width: 100%;
    height: 70px;
    border-bottom: 2px solid rgb(235, 235, 235);
}
.main-panel {
    position: absolute;
    left: 240px;
    color: rgb(41, 61, 81);
    text-align: left;
}
.main-panel h1 {
    font-size: 38px;
    margin-top: 40px;
}
.main-panel .container {
    width: 800px;
    padding-bottom: 40px;
    border-bottom: 2px solid rgb(235, 235, 235);
}
</style>
