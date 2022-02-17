<template>
    <v-date-picker :value="value" :masks="{ input: format }" :min-date="minDate" @input="onInput">
        <template v-slot="{ inputValue, inputEvents }">
            <label class="has-float-label">
                <input type="text" :value="inputValue" v-on="inputEvents" :placeholder="label" />
                <span v-if="label">{{ label }}</span>
                <i class="icon icon-calendar"></i>
            </label>
        </template>
    </v-date-picker>
</template>

<script>
import Vue from "vue";
import moment from "moment";

export default Vue.extend({
    name: "DateField.vue",
    props: {
        value: {
            type: String,
            required: true,
        },
        format: {
            type: String,
            default: "M/D/YYYY",
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
        onInput(value) {
            this.$emit("input", moment(value).format(this.format));
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
