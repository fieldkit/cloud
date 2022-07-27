<template>
    <aside :class="{ show: isVisible, [style]: true }" class="snackbar">{{ message }}</aside>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "SnackBar",
    data: () => {
        return {};
    },
    computed: {
        isVisible() {
            return this.$state.snackbar.visible;
        },
        message() {
            return this.$state.snackbar.message;
        },
        style() {
            return this.$state.snackbar.style;
        },
    },
});
</script>

<style scoped lang="scss">
@import "src/scss/variables";
.snackbar {
    visibility: hidden;
    min-width: 250px;
    margin-left: -125px;
    background-color: #333;
    color: $color-dark;
    text-align: center;
    border-radius: 4px;
    padding: 15px 20px;
    position: fixed;
    z-index: $z-index-top;
    left: 50%;
    bottom: 30px;

    &.show {
        visibility: visible;
        animation: fadein 0.5s, fadeout 0.5s 5s; // last value = persistence time
        animation-iteration-count: 1;
    }

    &.success {
        background-color: #7cdb938a;
    }
    &.fail {
        background-color: #f8d7da;
    }
}

@keyframes fadein {
    from {
        bottom: 0;
        opacity: 0;
    }
    to {
        bottom: 30px;
        opacity: 1;
    }
}

@keyframes fadeout {
    from {
        bottom: 30px;
        opacity: 1;
    }
    to {
        bottom: 0;
        opacity: 0;
    }
}
</style>
