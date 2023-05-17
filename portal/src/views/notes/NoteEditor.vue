<template>
    <div class="note-editor">
        <div class="title">
            <TextAreaField v-if="editingTitle" v-model="title" @input="v.$touch()" />
            <template v-else>{{ title }}</template>
            <a class="edit-btn" v-if="editableTitle" @click="editingTitle = !editingTitle">{{ $t("edit") }}</a>
        </div>
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
        editableTitle: {
            type: Boolean,
            default: false,
        },
    },
    data() {
        return {
            editingTitle: false,
        };
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
        title: {
            get(this: any) {
                console.log("Radoi note", this.note);
                if (this.note.title !== "") {
                    return this.note.title;
                } else {
                    return this.note.help.title;
                }
            },
            set(this: any, value) {
                this.$emit("change", this.note.withBody(this.note.body, value));
            },
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.attached-audio {
    display: flex;
    flex-wrap: wrap;
    align-items: baseline;
    margin-bottom: 10px;
    border: 1px solid var(--color-border);
    border-radius: 4px;
    background-color: #fcfcfc;
    padding: 8px;
}

.audio-title {
    font-size: 14px;
    font-weight: 500;
    margin-right: 10px;

    @include bp-down($xs) {
        flex-basis: 100%;
    }
}

.title {
    font-size: 16px;
    font-weight: 500;
    display: flex;
    align-items: center;

    @include bp-down($xs) {
        font-size: 14px;
    }
}

.no-data-yet {
    color: #6a6d71;
    font-size: 13px;
    padding-top: 0.5em;
}

.edit-btn {
    opacity: 0.4;
    font-size: 12px;
    cursor: pointer;
    margin-left: 8px;
}

::v-deep .title textarea {
    padding-top: 0;
}
</style>
