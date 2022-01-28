<template>
    <i id="header-logo" role="img" :aria-label="altText"></i>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "Logo",
    computed: {
        altText(): string {
            if (document.body.classList.contains("floodnet")) {
                console.log("radoi true");

                return this.$tc("layout.logo.floodnet.alt");
            }
            console.log("radoi false");
            return this.$tc("layout.logo.fieldkit.alt");
        },
    },
});
</script>

<style lang="scss" scoped>
@import "../../scss/mixins";

#header-logo {
    width: 140px;
    display: block;
    background: url("../../assets/logo-fieldkit.svg") no-repeat center center;
    background-size: contain;

    @include bp-down($md) {
        @include position(fixed, null null null 50%);
        width: 130px;
        height: 50px;
        transform: translateX(-50%);
    }

    @include bp-down($xs) {
        width: 109px;
    }

    @at-root .floodnet & {
        width: 160px;
        background: url("../../assets/logo-floodnet.png") no-repeat center center;
        background-size: contain;

        @include bp-down($md) {
            width: 140px;
        }

        @include bp-down($xs) {
            width: 120px;
        }
    }
}
</style>
