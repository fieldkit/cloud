<template>
    <div class="note-editor">
        <div class="title">{{ note.help.title }}</div>
        <div class="field" v-if="!readOnly">
            <TextField v-model="body" @input="v.$touch()" />
        </div>
        <div class="field" v-if="readOnly">
            {{ note.body }}
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { required } from "vuelidate/lib/validators";

import FKApi from "@/api/api";

export default Vue.extend({
    model: {
        prop: "note",
        event: "change",
    },
    name: "NoteEditor",
    components: {
        ...CommonComponents,
    },
    props: {
        readOnly: {
            type: Boolean,
            default: false,
        },
        note: {
            type: Object,
            required: true,
        },
        v: {
            type: Object,
            required: true,
        },
    },
    computed: {
        body: {
            get(this: any) {
                return this.note.body;
            },
            set(this: any, value) {
                this.$emit("change", this.note.withBody(value));
            },
        },
    },
    methods: {},
});
</script>

<style scoped>
.note-editor {
}

.title {
    font-size: 16px;
    font-weight: 500;
}

.field {
    margin-top: 10px;
    margin-bottom: 10px;
}
</style>
