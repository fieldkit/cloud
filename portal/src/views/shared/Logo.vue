<template>
    <i id="header-logo" role="img" :aria-label="altText" class="icon" :class="iconClass"></i>
</template>

<script lang="ts">
import Vue from "vue";
import { interpolatePartner } from "./partners";

export default Vue.extend({
    name: "Logo",
    data(): {
        iconClass: string;
        altText: string;
    } {
        return {
            iconClass: interpolatePartner("icon-logo-"),
            altText: interpolatePartner("layout.logo.") + ".alt",
        };
    },
});
</script>

<style lang="scss" scoped></style>
