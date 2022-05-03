<template>
    <div class="notes-form">
        <vue-confirm-dialog />
        <div class="notifications">
            <div v-if="notesState.failed" class="notification failed">{{ $tc("notes.failed") }}</div>
            <div v-if="notesState.success" class="notification success">{{ $tc("notes.success") }}</div>
        </div>
        <div class="header">
            <div class="name">{{ $t("notes.title") }}</div>
            <div class="completed">{{ completed }}% Complete</div>
            <div class="buttons" v-if="isAuthenticated">
                <button type="submit" class="button" @click="onSave">Save</button>
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
        <template v-if="readonly">
            <div class="photos">
                <div class="title">Photos</div>
                <div v-if="photos.length === 0" class="no-data-yet">No photos yet.</div>
                <div class="photo-container" v-for="photo in photos" v-bind:key="photo.key">
                    <AuthenticatedPhoto :url="photo.url" />
                </div>
            </div>
        </template>
        <template v-if="!readonly">
            <div class="photos">
                <div class="title">Photos</div>
                <div class="photo-container" v-for="photo in photos" v-bind:key="photo.key">
                    <AuthenticatedPhoto :url="photo.url" />
                    <ListItemOptions @listItemOptionClick="onPhotoOptionClick($event, photo)" :options="photoOptions"></ListItemOptions>
                </div>
                <div class="photo-container" v-for="photo in form.addedPhotos" v-bind:key="photo.key">
                    <img :src="photo.image" />
                    <ListItemOptions
                        @listItemOptionClick="onUnsavedPhotoOptionClick($event, photo)"
                        :options="photoOptions"
                    ></ListItemOptions>
                </div>
                <ImageUploader @change="onImage" :placeholder="placeholder" :allowPreview="false" />
            </div>
        </template>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import { mapGetters } from "vuex";
import CommonComponents from "@/views/shared";

import { Notes, AddedPhoto, NoteMedia, PortalNoteMedia, mergeNotes } from "./model";
import NoteEditor from "./NoteEditor.vue";
import ListItemOptions from "@/views/shared/ListItemOptions.vue";

import NewPhoto from "../../assets/image-placeholder.svg";
import { serializePromiseChain } from "@/utilities";
import { ActionTypes } from "@/store";

export default Vue.extend({
    name: "NotesForm",
    components: {
        ...CommonComponents,
        NoteEditor,
        ListItemOptions,
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
            default: false,
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
            photoOptions: [
                {
                    label: "Use as Station Image",
                    event: "use-as-station-image",
                },
                {
                    label: "Delete Image",
                    event: "delete-media",
                },
            ],
            notesState: {
                dirty: false,
                success: false,
                failed: false,
            },
        };
    },
    beforeRouteUpdate(to: never, from: never, next: any) {
        if (this.confirmLeave()) {
            next();
        }
    },
    beforeRouteLeave(to: never, from: never, next: any) {
        if (this.confirmLeave()) {
            next();
        }
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        photos(this: any) {
            return NoteMedia.onlyPhotos(this.notes.media);
        },
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
        onImage(this: any, image: any) {
            const reader = new FileReader();
            reader.readAsDataURL(image.file);
            reader.onload = (ev) => {
                if (ev?.target?.result) {
                    this.form.addedPhotos.push(new AddedPhoto(image.type, image.file, ev.target.result));
                }
            };
            this.$emit("change", this.form);
        },
        async onPhotoOptionClick(event: string, photo: PortalNoteMedia) {
            if (event === "delete-media") {
                await this.$confirm({
                    message: `Are you sure you want to delete this image?`,
                    button: {
                        no: "Cancel",
                        yes: "Delete",
                    },
                    callback: async (confirm) => {
                        if (confirm) {
                            await this.$services.api.deleteMedia(photo.id);
                            this.notes.media = this.notes.media.filter((media) => media.id !== photo.id);
                            return;
                        }
                    },
                });
            }

            if (event === "use-as-station-image") {
                this.$services.api
                    .setStationImage(this.station.id, photo.id)
                    .then((result) => {
                        this.notesState.success = true;
                    })
                    .catch((e) => {
                        this.notesState.failed = true;
                    });
            }
        },
        async onUnsavedPhotoOptionClick(event: string, photo: AddedPhoto) {
            if (event === "delete-media") {
                await this.$confirm({
                    message: `Are you sure you want to delete this image?`,
                    button: {
                        no: "Cancel",
                        yes: "Delete",
                    },
                    callback: async (confirm) => {
                        if (confirm) {
                            const photoIndex = this.form.addedPhotos.findIndex((addedPhoto) => addedPhoto.key === photo.key);
                            this.form.addedPhotos.splice(photoIndex, 1);
                            return;
                        }
                    },
                });
            }

            if (event === "use-as-station-image") {
                this.$services.api
                    .uploadStationMedia(this.station.id, photo.key, photo.file)
                    .then(async (uploadedPhoto: PortalNoteMedia) => {
                        await this.$services.api.setStationImage(this.station.id, uploadedPhoto.id);
                    });
            }
        },
        async deleteMedia(photo: PortalNoteMedia): Promise<void> {
            if (photo.id) {
                await this.$services.api.deleteMedia(photo.id);
                this.notes.media = this.notes.media.filter(
                    (media: { id: number; key: string; url: string; contentType: string }) => media.id !== photo.id
                );
                return;
            }
            const photoIndex = this.form.addedPhotos.findIndex((addedPhoto: NoteMedia) => addedPhoto.key === photo.key);
            this.form.addedPhotos.splice(photoIndex, 1);
            return;
        },

        async onSave(formNotes: Notes): Promise<void> {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            this.notesState.success = false;
            this.notesState.failed = false;

            await serializePromiseChain(formNotes.addedPhotos, (photo) => {
                return this.$services.api.uploadStationMedia(this.station.id, photo.key, photo.file).then(() => {
                    return [];
                });
            }).then(() => {
                const payload = mergeNotes({ notes: this.notes, media: this.media }, formNotes);
                return this.$services.api
                    .patchStationNotes(this.station.id, payload)
                    .then(
                        () => {
                            this.notesState.dirty = false;
                            this.notesState.success = true;
                        },
                        () => {
                            this.notesState.failed = true;
                        }
                    )
                    .catch((e) => {
                        console.log("radoi e", e);
                    })
                    .finally(() => {
                        this.$store.dispatch(ActionTypes.NEED_NOTES, { id: this.$route.params.stationId });
                    });
            });
        },
        onChange(): void {
            this.notesState.dirty = true;
        },
        confirmLeave(): boolean {
            if (this.notesState.dirty) {
                if (window.confirm("You may have unsaved changes, are you sure you'd like to leave?")) {
                    this.notesState.dirty = false;
                    return true;
                } else {
                    return false;
                }
            }
            return true;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/global";

.notes-form {
    text-align: left;
    display: flex;
    flex-direction: column;
    background: #fff;
    cursor: initial;

    @include bp-down($md) {
        padding: 25px 10px;
        border-bottom: 1px solid var(--color-border);
    }
}

.header {
    @include flex(center);
    padding-bottom: 15px;
    border-bottom: 1px solid var(--color-border);

    @include bp-down($md) {
        border: 0;
        padding: 0;
    }

    @include bp-down($xs) {
        flex-wrap: wrap;
    }
}

.header .name {
    color: #2c3e50;
    font-size: 20px;
    font-weight: 500;

    @include bp-down($xs) {
        font-size: 18px;
    }
}

.header .completed {
    margin-left: 10px;
    color: #0a67aa;
    font-size: 14px;
    font-weight: 600;

    @include bp-down($xs) {
        order: 3;
        margin-left: 0;
        flex-basis: 100%;
        font-size: 12px;
    }

    body.floodnet & {
        color: #6a6d71;
        font-size: 13px;
        font-weight: 500;
    }
}

.header .buttons {
    margin-left: auto;
    display: flex;
}

.site-notes {
    margin-top: 30px;

    @include bp-down($sm) {
        margin-top: 18px;
    }
}

.button {
    padding: 0;
    width: 80px;
    height: 33px;
    border-radius: 3px;
    border: solid 1px #cccdcf;
    background-color: #ffffff;
    font-size: 14px;
    font-family: var(--font-family-bold);
    letter-spacing: 0.08px;
    color: #2c3e50;
    margin-left: 7px;
    margin-bottom: 0;
    @include flex(center, center);

    @include bp-down($xs) {
        height: 25px;
        transform: translateY(8px);
    }
}

.photo-container {
    @include flex(flex-start);
    position: relative;
    margin-right: 55px;
    margin-bottom: 12px;
}

.photos {
    margin-top: 45px;
    display: flex;
    flex-wrap: wrap;

    @include bp-down($sm) {
        margin-top: 15px;
    }

    .title {
        flex-basis: 100%;
        font-weight: 500;
        margin-bottom: 20px;

        @include bp-down($sm) {
            color: #6a6d71;
        }

        @include bp-down($xs) {
            font-size: 14px;
            margin-bottom: 8px;
        }
    }

    ::v-deep img {
        border-radius: 3px;
        object-fit: contain;
        max-height: 200px;
        max-width: 200px;
        width: 200px;
        height: auto;
        object-fit: contain;

        @include bp-down($sm) {
            max-height: 82px;
            max-width: 90px;
            margin-right: 12px;
        }
    }
}

::v-deep .photos img-container img {
    margin-right: 10px;
    margin-bottom: 10px;
}

::v-deep .placeholder-container {
    flex-basis: 100%;
}

.notification {
    margin-top: 0;
}

.notification.success {
    margin-bottom: 20px;
    padding: 20px;
    border: 2px;
    border-radius: 4px;
}
.notification.success {
    background-color: #d4edda;
}
.notification.failed {
    background-color: #f8d7da;
}
</style>
