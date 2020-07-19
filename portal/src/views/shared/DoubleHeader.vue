<template>
    <div class="double-header">
        <div class="back" v-on:click="onBack" v-if="backTitle && backRoute">
            <span class="small-arrow">&lt;</span>
            {{ backTitle }}
        </div>
        <div class="one">{{ title }}</div>
        <div class="two" v-if="subtitle">{{ subtitle }}</div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "DoubleHeader",
    props: {
        title: {
            type: String,
            required: true,
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

<style>
.double-header {
    text-align: left;
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
    font-weight: 900;
    font-size: 24px;
    color: #2c3e50;
    margin-bottom: 10px;
}
.two {
    font-weight: 500;
    font-size: 18px;
    color: #6a6d71;
}
</style>
