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
button {
    cursor: pointer;
}
.main-panel {
    width: auto;
    text-align: left;
    overflow: hidden;
    color: #2c3e50;
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
input:focus ~ .floating-label,
input:not(:focus):valid ~ .floating-label {
    top: -48px;
    font-size: 12px;
    opacity: 1;
}
input:invalid {
    box-shadow: none;
}
.inputText {
    outline: 0;
    color: #2c3e50;
    font-size: 14px;
    width: inherit;
    border: none;
    border-bottom: 2px solid rgb(235, 235, 235);
    font-size: 15px;
    padding-bottom: 4px;
}
.inputText:focus {
    border-bottom: 2px solid #52b5e4;
}
.floating-label {
    color: rgb(85, 85, 85);
    position: relative;
    top: -24px;
    pointer-events: none;
    transition: 0.2s ease all;
}
.disabled {
    background-color: rgb(210, 210, 210);
    cursor: not-allowed;
    border: none;
    color: white;
    font-size: 18px;
    border-radius: 5px;
}
</style>
