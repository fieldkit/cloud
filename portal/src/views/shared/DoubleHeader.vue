<template>
    <div class="double-header">
        <div class="back">
            <span class="small-arrow">&lt;</span>
            {{ backTitle }}
        </div>
        <div class="one">{{ title }}</div>
        <div class="two">{{ subtitle }}</div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "DoubleHeader",
    components: {},
    props: {
        title: {
            type: String,
            required: true,
        },
        subtitle: {
            type: String,
            required: true,
        },
        backTitle: {
            type: String,
            required: true,
        },
    },
    methods: {},
});
</script>

<style>
.double-header {
    text-align: left;
    padding-left: 40px;
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
