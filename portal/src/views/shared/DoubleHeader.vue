<template>
    <div class="double-header">
        <div class="heading">
            <div class="back" v-on:click="onBack" v-if="backTitle">
                <span class="small-arrow">&lt;</span>
                {{ backTitle }}
            </div>
            <slot name="title">
                <div class="one" v-if="title">{{ title }}</div>
            </slot>
            <slot name="subtitle">
                <div class="two" v-if="subtitle">{{ subtitle }}</div>
            </slot>
        </div>
        <div class="actions">
            <slot></slot>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "DoubleHeader",
    props: {
        title: {
            type: String,
            default: null,
        },
        subtitle: {
            type: String,
            default: null,
        },
        backTitle: {
            type: String,
            default: null,
        },
        backRoute: {
            type: String,
            default: null,
        },
        backRouteParams: {
            type: Object,
            default: undefined,
        },
    },
    methods: {
        async onBack() {
            if (this.backRoute) {
                await this.$router.push({ name: this.backRoute, params: this.backRouteParams });
            } else {
                this.$emit("back");
            }
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.actions {
    display: flex;
}
.double-header {
    display: flex;
    text-align: left;
    flex-wrap: wrap;
}
.double-header .actions {
    margin-left: auto;
    margin-top: auto;

    @include bp-down($xs) {
        flex-basis: 100%;
        margin: 15px 0;
    }
}
.back {
    font-size: 14px;
    letter-spacing: 0.06px;
    margin-bottom: 20px;
    cursor: pointer;
}
.one {
    font-family: var(--font-family-bold);
    font-size: 24px;
    margin-bottom: 1px;

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
    }
}
.two {
    font-weight: 500;
    font-size: 18px;
    color: #6a6d71;
}
</style>
