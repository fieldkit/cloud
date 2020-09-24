<template>
    <div class="notes-viewer">
        <div class="inner">
            <div class="header">
                <div class="name">{{ station.name }}</div>
                <div class="completed">{{ completed }}% Complete</div>
                <div class="buttons"></div>
            </div>
            <div class="site-notes">
                <NoteViewer :note="form.studyObjective" />
                <NoteViewer :note="form.sitePurpose" />
                <NoteViewer :note="form.siteCriteria" />
                <NoteViewer :note="form.siteDescription" />
            </div>
            <div class="photos">
                <div class="photo" v-for="photo in photos" v-bind:key="photo.key">
                    <AuthenticatedPhoto :url="photo.url" />
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { Notes, AddedPhoto, NoteMedia } from "./model";

const NoteViewer = Vue.extend({
    name: "NoteViewer",
    props: {
        note: {
            type: Object,
            required: true,
        },
    },
    computed: {},
    template: `
	<div class="note">
		<div class="title">{{ note.help.title }}</div>
		<div class="body">{{ note.body }}</div>
	</div>`,
});

export default Vue.extend({
    name: "NotesViewer",
    components: {
        ...CommonComponents,
        NoteViewer,
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
    },
    data: () => {
        return {
            form: new Notes(),
        };
    },
    computed: {
        photos(this: any) {
            return NoteMedia.onlyPhotos(this.notes.media);
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
    methods: {},
});
</script>

<style scoped>
.notes-viewer {
    text-align: left;
    display: flex;
    flex-direction: column;
    padding: 20px;
}
.header {
    display: flex;
    align-items: baseline;
    padding-bottom: 25px;
    border-bottom: 1px solid #d8dce0;
}
.header .name {
    color: #2c3e50;
    font-size: 20px;
    font-weight: 500;
}
.header .completed {
    margin-left: 10px;
    color: #0a67aa;
    font-size: 14px;
    font-weight: 900;
}
.header .buttons {
    margin-left: auto;
    display: flex;
}
.site-notes {
    margin-top: 20px;
}

/deep/ .photos img {
    max-width: 400px;
    max-height: 400px;
    margin-right: 10px;
    margin-bottom: 10px;
}

/deep/ .note .title {
    font-size: 18px;
    font-weight: 900;
    margin-bottom: 1em;
}

/deep/ .note .body {
    font-size: 14px;
    margin-bottom: 1em;
}
</style>
