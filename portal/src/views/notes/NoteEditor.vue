<template>
    <div class="note-editor">
        <div class="title">{{ note.help.title }}</div>
        <div class="field" v-if="!readOnly">
            <TextAreaField v-model="body" @input="v.$touch()" />
        </div>
        <div class="field" v-if="readOnly">
            {{ note.body }}
        </div>
        <div class="attached-audio" v-for="audio in note.audio" v-bind:key="audio.key">
            <div class="audio-title">
                {{ audio.key }}
            </div>
            <VueAudio :file="makeAudioUrl(audio)" />
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import VueAudio from "vue-audio";

import { required } from "vuelidate/lib/validators";

import FKApi, { makeAuthenticatedApiUrl } from "@/api/api";

export default Vue.extend({
    model: {
        prop: "note",
        event: "change",
    },
    name: "NoteEditor",
    components: {
        ...CommonComponents,
        VueAudio,
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
    methods: {
        makeAudioUrl(audio) {
            return makeAuthenticatedApiUrl(audio.url);
        },
    },
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

.attached-audio {
    display: flex;
    align-items: baseline;
    margin-bottom: 10px;
    border: 2px solid #d8dce0;
    border-radius: 4px;
    background-color: #fcfcfc;
    padding: 8px;
}

.audio-title {
    font-size: 14px;
    font-weight: 500;
    margin-right: 10px;
}
</style>
