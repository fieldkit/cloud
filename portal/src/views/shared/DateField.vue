<template>
    <v-date-picker :value="value" :masks="{ input: format }" :min-date="minDate">
        <template v-slot="{ inputValue, inputEvents }">
            <label class="has-float-label">
                <input type="text" :value="inputValue" @input="onInput" v-on="inputEvents" :placeholder="label" />
                <span v-if="label">{{ label }}</span>
                <i class="icon icon-calendar"></i>
            </label>
        </template>
    </v-date-picker>
</template>

<script>
import Vue from "vue";

export default Vue.extend({
    name: "DateField.vue",
    props: {
        value: {
            type: String,
            required: true,
        },
        format: {
            type: String,
            default: "MM/DD/YY",
        },
        label: {
            type: String,
            required: false,
        },
        minDate: {
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
@import "../../scss/float-label";

.icon-calendar {
    @include position(absolute, null 0 10px null);

    &:before {
        color: var(--color-dark);
    }
}
</style>
