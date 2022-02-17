<template>
    <label class="has-float-label">
        <input :value="value" :type="type" @input="onInput" :placeholder="placeholder || label" />
        <span v-if="label">{{ label }}</span>
    </label>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";

export default Vue.extend({
    name: "TextField",
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

<style scoped lang="scss">
@import '../../scss/float-label';

</style>
