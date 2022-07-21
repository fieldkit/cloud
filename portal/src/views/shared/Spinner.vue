<template>
    <div class="spinner loader">
        <div></div>
        <div></div>
        <div></div>
        <div></div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "Spinner",
});
</script>

<style scoped lang="scss">
@import "src/scss/_variables.scss";

.spinner {
    display: inline-block;
    position: relative;
    width: 20px;
    height: 20px;
}
.spinner div {
    box-sizing: border-box;
    display: block;
    position: absolute;
    width: 20px;
    height: 20px;
    border: 3px solid #fff;
    border-radius: 50%;
    animation: spinner 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
    border-color: #fff transparent transparent transparent;
}
.spinner div:nth-child(1) {
    animation-delay: -0.45s;
}
.spinner div:nth-child(2) {
    animation-delay: -0.3s;
}
.spinner div:nth-child(3) {
    animation-delay: -0.15s;
}
@keyframes spinner {
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
}
</style>
