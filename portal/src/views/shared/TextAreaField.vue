<template>
    <label class="has-float-label">
        <ResizeAuto>
            <template v-slot:default="{}">
                <textarea rows="2" :value="value" :type="type" :placeholder="placeholder || label" @input="(ev) => onInput(ev)" />
            </template>
        </ResizeAuto>
        <span v-if="label">{{ label }}</span>
    </label>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";

const ResizeAuto = Vue.extend({
    name: "ResizeAuto",
    methods: {
        resize(ev) {
            ev.target.style.height = "auto";
            ev.target.style.height = `${ev.target.scrollHeight}px`;
        },
    },
    mounted() {
        this.$nextTick(() => {
            const el: any = this.$el;
            el.setAttribute("style", "height:" + (this.$el.scrollHeight - 10) + "px");
        });

        this.$el.addEventListener("input", this.resize);
    },
    updated() {
        this.$nextTick(() => {
            const el: any = this.$el;
            el.setAttribute("style", "height:" + (this.$el.scrollHeight - 0) + "px");
        });
    },
    beforeDestroy() {
        this.$el.removeEventListener("input", this.resize);
    },
    render(this: any) {
        return this.$scopedSlots.default({
            resize: this.resize,
        });
    },
});

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
    box-sizing: border-box;
    font-size: inherit;
    padding-top: 0.5em;
    margin-bottom: 2px;
    border: 0;
    border-radius: 0;
    border-bottom: 2px solid rgba(0, 0, 0, 0.1);
}
.has-float-label textarea:focus {
    outline: none;
    border-color: rgba(0, 0, 0, 0.5);
}
</style>
