<template>
    <div class="double-header">
        <div class="back">
            <span class="small-arrow">&lt;</span>
            Back to dashboard
        </div>
        <div class="one">Field Notes</div>
        <div class="two">Project Name</div>
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
            default: "",
        },
        subtitle: {
            type: String,
            default: "",
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
