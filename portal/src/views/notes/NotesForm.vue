<template>
    <div class="notes-form">
        <div class="inner">
            <div class="header">
                <div class="name">{{ station.name }}</div>
                <div class="completed">{{ form.completed }}% Complete</div>
                <div class="buttons">
                    <div class="button" v-on:click="onSave">Save</div>
                </div>
            </div>
            <div class="site-notes">
                <form id="form">
                    <NoteEditor v-model="form.studyObjective" :v="$v.form.studyObjective" @change="onChange" />
                    <NoteEditor v-model="form.sitePurpose" :v="$v.form.sitePurpose" @change="onChange" />
                    <NoteEditor v-model="form.siteCriteria" :v="$v.form.siteCriteria" @change="onChange" />
                    <NoteEditor v-model="form.siteDescription" :v="$v.form.siteDescription" @change="onChange" />
                </form>
            </div>
            <div class="photos">
                <div class="photo" v-for="photo in photos" v-bind:key="photo.key">
                    <img :src="makePhotoUrl(photo)" />
                </div>
                <div class="photo" v-for="photo in form.addedPhotos" v-bind:key="photo.key">
                    <img :src="photo.image" />
                </div>
                <ImageUploader @change="onImage" :placeholder="placeholder" :allowPreview="false" />
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import NoteEditor from "./NoteEditor.vue";

import NewPhoto from "../../assets/placeholder_station_thumbnail.png";

import { required } from "vuelidate/lib/validators";

import { Notes, AddedPhoto } from "./model";

import { makeAuthenticatedApiUrl } from "@/api/api";

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
            placeholder: NewPhoto,
        };
    },
    computed: {
        photos(this: any) {
            return this.notes.media;
        },
    },
    mounted(this: any) {
        this.form = Notes.createFrom(this.notes);
    },
    methods: {
        makePhotoUrl(photo) {
            return makeAuthenticatedApiUrl(photo.url);
        },
        onSave(this: any) {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }
            this.$emit("save", this.form);
        },
        onImage(this: any, image: any) {
            const reader = new FileReader();
            reader.readAsDataURL(image.file);
            reader.onload = (ev) => this.form.addedPhotos.push(new AddedPhoto(image.type, image.file, ev.target.result));
            this.$emit("change", this.form);
        },
        onChange() {
            this.$emit("change", this.form);
        },
    },
});
</script>

<style scoped>
.notes-form {
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
.button {
    margin-left: 20px;
    font-size: 12px;
    padding: 5px 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}

.photos {
    display: flex;
    flex-wrap: wrap;
}

.photos .photo {
    flex-basis: 400px;
}

/deep/ .photos img {
    max-width: 400px;
    max-height: 400px;
    margin-right: 10px;
    margin-bottom: 10px;
}
</style>
