<template>
    <span class="info" @click="$event.stopPropagation()">
        <span class="info-icon">?</span>
        <span class="info-content">{{ message }}</span>
    </span>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "InfoTooltip",
    props: {
        message: {
            type: String,
            required: true,
        },
    },
});
</script>

<style lang="scss" scoped>
@import "src/scss/mixins";

.info {
    position: relative;
    padding: 7px;
    margin-top: -3px;
    font-family: $font-family-bold;

    &:hover {
        .info-content {
            visibility: visible;
            opacity: 1;
        }
    }
}

.info-icon {
    display: block;
    background-color: #6a6d71;
    color: #fff;
    text-align: center;
    text-indent: -1px;
    line-height: 14px;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    opacity: 0.3;
    font-size: 10px;
    font-family: $font-family-medium;

    @include bp-up($sm) {
        display: none;
    }
}

.info-content {
    visibility: hidden;
    opacity: 0;
    padding: 14px;
    font-size: 10px;
    color: #6a6d71;
    position: absolute;
    top: 100%;
    border-radius: 2px;
    border: solid 1px #cccdcf;
    background-color: #fff;
    width: 200px;
    z-index: $z-index-top;
}
</style>
