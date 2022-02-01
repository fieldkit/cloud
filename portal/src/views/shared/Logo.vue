<template>
    <i id="header-logo" role="img" :aria-label="altText" class="icon" :class="iconClass"></i>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "Logo",
    computed: {
        altText(): string {
            if (document.body.classList.contains("floodnet")) {
                return this.$tc("layout.logo.floodnet.alt");
            }
            return this.$tc("layout.logo.fieldkit.alt");
        },
        iconClass(): string {
            if (document.body.classList.contains("floodnet")) {
                return "icon-logo-floodnet";
            }
            return "icon-logo-fieldkit";
        },
    },
});
</script>

<style lang="scss" scoped></style>
