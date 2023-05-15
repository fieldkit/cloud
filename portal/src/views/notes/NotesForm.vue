<template>
    <div class="notes-form">
        <div class="header">
            <div class="name">{{ $t("notes.title") }}</div>
            <div class="completed">{{ completed }}% {{ $t("notes.complete") }}</div>
            <div class="buttons" v-if="isAuthenticated">
                <button type="submit" class="button" @click="onSave">{{ $t("notes.btn.save") }}</button>
            </div>
        </div>
        <div class="site-notes">
            <form id="form">
                <NoteEditor v-model="form.studyObjective" :v="$v.form.studyObjective" :readonly="readonly" @change="onChange" />
                <NoteEditor v-model="form.sitePurpose" :v="$v.form.sitePurpose" :readonly="readonly" @change="onChange" />
                <NoteEditor v-model="form.siteCriteria" :v="$v.form.siteCriteria" :readonly="readonly" @change="onChange" />
                <NoteEditor v-model="form.siteDescription" :v="$v.form.siteDescription" :readonly="readonly" @change="onChange" />
            </form>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapGetters } from "vuex";
import CommonComponents from "@/views/shared";

import { mergeNotes, NoteMedia, Notes, PortalNoteMedia } from "./model";
import NoteEditor from "./NoteEditor.vue";
import { ActionTypes } from "@/store";
import { SnackbarStyle } from "@/store/modules/snackbar";

export default Vue.extend({
    name: "NotesForm",
    components: {
        ...CommonComponents,
        NoteEditor,
    },
    props: {
        station: {
            type: Object,
            required: true,
        },
        notes: {
            type: Object,
            required: true,
        },
        readonly: {
            type: Boolean,
            default: true,
        },
    },
    validations: {
        form: {
            studyObjective: {},
            sitePurpose: {},
            siteCriteria: {},
            siteDescription: {},
        },
    },
    data: () => {
        return {
            form: new Notes(),
            notesState: {
                dirty: false,
            },
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        media(): PortalNoteMedia[] {
            return this.$state.notes.media;
        },
        completed(this: any) {
            const notesProgress = this.form.progress;
            const anyPhotos = NoteMedia.onlyPhotos(this.form.addedPhotos).length + NoteMedia.onlyPhotos(this.notes.media).length > 0;
            const percentage = ((notesProgress.completed + anyPhotos) / (notesProgress.total + 1)) * 100;
            return percentage.toFixed(0);
        },
    },
    mounted(this: any) {
        this.form = Notes.createFrom(this.notes);
    },
    methods: {
        async onSave(): Promise<void> {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            const payload = mergeNotes({ notes: this.notes, media: this.media }, this.form);
            return this.$services.api
                .patchStationNotes(this.station.id, payload)
                .then(() => {
                    this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                        message: this.$tc("notes.updateSuccess"),
                        type: SnackbarStyle.success,
                    });
                    this.$emit('saved');
                })
                .catch(() => {
                  this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                        message: this.$tc("notes.updateFail"),
                        type: SnackbarStyle.fail,
                    });
                });
        },
        onChange(): void {
            this.notesState.dirty = true;
            this.$emit("change");
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/global";
@import "../../scss/notes";
</style>
