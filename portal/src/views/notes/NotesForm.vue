<template>
    <div class="notes-form">
        <div class="inner">
            <div class="header">
                <div class="name">{{ station.name }}</div>
                <div class="completed">{{ completed }}% Complete</div>
                <div class="buttons">
                    <button type="submit" class="button" v-on:click="onSave">Save</button>
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
                    <AuthenticatedPhoto :url="photo.url" />
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

import { Notes, AddedPhoto, NoteMedia } from "./model";

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
    methods: {
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

<style scoped lang="scss">
@import '../../scss/mixins';

.notes-form {
    text-align: left;
    display: flex;
    flex-direction: column;
    padding: 28px;

    @include bp-down($md) {
        padding: 25px 8px;
    }
}
.header {
    @include flex(center);
    padding-bottom: 11px;
    border-bottom: 1px solid #d8dce0;

    @include bp-down($md) {
        border: 0;
        padding: 0;
    }
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
    font-weight: 600;
}
.header .buttons {
    margin-left: auto;
    display: flex;
}
.site-notes {
    margin-top: 20px;
}
.button {
    padding: 0;
    width: 80px;
    height: 33px;
    border-radius: 3px;
    border: solid 1px #cccdcf;
    background-color: #ffffff;
    font-size: 14px;
    font-weight: 600;
    letter-spacing: 0.08px;
    color: #2c3e50;
    margin-left: 7px;
    margin-bottom: 0;
    @include flex(center, center);
}

.photos {
    display: flex;
    flex-wrap: wrap;
}

.photos .photo {
    flex-basis: 400px;
}

::v-deep .photos img {
    max-width: 400px;
    max-height: 400px;
    margin-right: 10px;
    margin-bottom: 10px;
}
</style>
