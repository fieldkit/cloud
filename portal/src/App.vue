<template>
    <div id="app">
        <router-view />
    </div>
</template>

<script>
import * as ActionTypes from "@/store/actions";
import { AuthenticationRequiredError } from "@/api/api";

export default {
    mounted() {
        return this.$store.dispatch(ActionTypes.INITIALIZE).catch((err) => {
            console.log("initialize error", err, err.stack);
        });
    },
    errorCaptured(err, vm, info) {
        if (err instanceof AuthenticationRequiredError) {
            this.$router.push({ name: "login", query: { after: this.$route.path } });
            return false;
        }
        return true;
    },
};
</script>
<style lang="scss">
@import 'scss/mixins';

html {
}
html,
body,
#app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
}
body {
    text-align: center;
    color: #2c3e50;
    margin: 0;
    padding: 0;

    flex-shrink: 0;

    * {
        font-family: "Avenir", Helvetica, Arial, sans-serif;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
    }
}
body:not(.disable-scrolling) {
    overflow-y: scroll;
}
body.disable-scrolling {
    margin-right: 14px; /* We need width of the scrollbars! */
}
body.blue-background {
    background-color: #1b80c9;

    @include bp-down($md) {
        background-color: #fff;
    }
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
</style>
