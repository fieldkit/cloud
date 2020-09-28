<template>
    <div class="double-header">
        <div class="heading">
            <div class="back" v-on:click="onBack" v-if="backTitle && backRoute">
                <span class="small-arrow">&lt;</span>
                {{ backTitle }}
            </div>
            <div class="one" v-if="title">{{ title }}</div>
            <div class="two" v-if="subtitle">{{ subtitle }}</div>
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
    },
    methods: {
        onBack() {
            this.$router.push({ name: this.backRoute });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

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
.double-header > div {
}
.back {
    font-size: 14px;
    letter-spacing: 0.06px;
    color: #2c3e50;
    margin-bottom: 20px;
    cursor: pointer;
}
.one {
    font-weight: 600;
    font-size: 24px;
    color: #2c3e50;
    margin-bottom: 1px;
}
.two {
    font-weight: 500;
    font-size: 18px;
    color: #6a6d71;
}
</style>
