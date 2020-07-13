<template>
    <div id="app">
        <router-view />
    </div>
</template>

<script>
import * as ActionTypes from "./store/actions";

export default {
    mounted() {
        return this.$store.dispatch(ActionTypes.INITIALIZE).catch((err) => {
            console.log("initialize error", err, err.stack);
        });
    },
};
</script>
<style>
html {
    overflow-y: scroll;
}
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
body.blue-background {
    background-color: #1b80c9;
}
html.map-view {
    height: 100%;
}
body.map-view {
    height: 100%;
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
    color: #2c3e50;
}
.main-panel h1 {
    font-size: 36px;
    margin-top: 40px;
}

/* Deprecated in favor of TextField */

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
    font-size: 16px;
    width: inherit;
    border: none;
    text-align: left;
    border-bottom: 2px solid #d8dce0;
    padding-bottom: 4px;
}
.inputText:focus {
    border-bottom: 2px solid #52b5e4;
}
.floating-label {
    color: #6a6d71;
    position: relative;
    top: -24px;
    pointer-events: none;
    text-align: left;
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
