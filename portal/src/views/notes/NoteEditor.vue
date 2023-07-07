<template>
    <div class="note-editor">
        <div class="title">{{ note.help.title }}</div>
        <div class="field" v-if="!readonly">
            <TextAreaField v-model="body" @input="v.$touch()" />
        </div>
        <div class="field" v-if="readonly">
            <template v-if="note.body">{{ note.body }}</template>
            <template v-if="!note.body">
                <div class="no-data-yet">{{ $t("notes.noFieldData") }}</div>
            </template>
        </div>
        <div class="attached-audio" v-for="audio in note.audio" v-bind:key="audio.key">
            <div class="audio-title">
                {{ audio.key }}
            </div>
            <AudioPlayer :url="audio.url" />
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import AudioPlayer from "./AudioPlayer.vue";

export default Vue.extend({
    model: {
        prop: "note",
        event: "change",
    },
    name: "NoteEditor",
    components: {
        ...CommonComponents,
        AudioPlayer,
    },
    props: {
        readonly: {
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
    mounted() {
        console.log("Radoi note", this.note);
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.attached-audio {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    margin-bottom: 10px;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    background-color: #fcfcfc;
    padding: 8px;
}

.audio-title {
    font-size: 14px;
    font-weight: 500;
    margin-right: auto;

    @include bp-down($xs) {
        flex-basis: 100%;
    }
}

.title {
    font-size: 16px;
    font-weight: 500;

    @include bp-down($xs) {
        font-size: 14px;
    }
}

.no-data-yet {
    color: #6a6d71;
    font-size: 13px;
    padding-top: 0.5em;
}
</style>
