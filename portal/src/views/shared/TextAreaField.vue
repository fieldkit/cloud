<template>
    <label class="has-float-label">
        <ResizeAuto>
            <template v-slot:default="{}">
                <textarea rows="1" :value="value" :type="type" :placeholder="placeholder || label" @input="(ev) => onInput(ev)" />
            </template>
        </ResizeAuto>
        <span v-if="label">{{ label }}</span>
    </label>
</template>

<script lang="ts">
import Vue from "vue";
import { ResizeAuto } from "./ResizeAuto";

export default Vue.extend({
    name: "TextAreaField",
    components: {
        ResizeAuto,
    },
    props: {
        type: {
            type: String,
            default: "text",
        },
        value: {
            type: String,
            required: true,
        },
        placeholder: {
            type: String,
            required: false,
            default: null,
        },
        label: {
            type: String,
            required: false,
        },
    },
    methods: {
        onInput(ev) {
            this.$emit("input", ev.target.value);
        },
    },
});
</script>

<style scoped>
.has-float-label textarea {
    width: 100%;
    font-family: "Avenir", Helvetica, Arial, sans-serif;
    font-size: inherit;
    padding-top: 0.5em;
    margin-bottom: 2px;
    box-sizing: border-box;
    border: 0;
    border-radius: 0;
    border-bottom: 2px solid rgba(0, 0, 0, 0.1);
    resize: none;
    overflow: hidden;
}
.has-float-label textarea:focus {
    outline: none;
    border-color: rgba(0, 0, 0, 0.5);
}
</style>
