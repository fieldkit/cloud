<template>
    <StandardLayout>
        <vue-confirm-dialog />
        <div class="container-wrap">
            <DoubleHeader
                :backRoute="projectId ? 'viewStation' : 'viewStationFromMap'"
                :title="$tc('station.allPhotos')"
                :backTitle="$tc('layout.backStationDashboard')"
                :backRouteParams="{ projectId, stationId }"
            >
                <template v-slot:default>
                    <button class="button-social">
                        <label for="imageInput">
                            <i class="icon icon-photo"></i>
                            {{ $tc("station.btn.addPhotos") }}
                        </label>
                    </button>
                    <input id="imageInput" type="file" accept="image/jpeg, image/png" @change="upload" />
                </template>
            </DoubleHeader>

            <div class="flex flex-wrap flex-space-between" v-if="media">
                <div class="photo-wrap" v-for="photo in photos" v-bind:key="photo.key">
                    <button class="photo-options">
                        <ListItemOptions :options="photoOptions" @listItemOptionClick="onPhotoOptionClick($event, photo)"></ListItemOptions>
                    </button>
                    <AuthenticatedPhoto :url="photo.url" :loading="photo.id === loadingPhotoId" />
                </div>
                <div v-if="photos.length === 0" class="empty-photos">
                    {{ $t("station.emptyPhotos") }}
                </div>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "@/views/StandardLayout.vue";
import DoubleHeader from "@/views/shared/DoubleHeader.vue";
import AuthenticatedPhoto from "@/views/shared/AuthenticatedPhoto.vue";
import { AddedPhoto, NoteMedia, PortalNoteMedia } from "@/views/notes/model";
import { ActionTypes } from "@/store";
import ListItemOptions, { ListItemOption } from "@/views/shared/ListItemOptions.vue";

export default Vue.extend({
    name: "StationPhotosView",
    components: {
        StandardLayout,
        DoubleHeader,
        AuthenticatedPhoto,
        ListItemOptions,
    },
    computed: {
        projectId(): number {
            return parseInt(this.$route.params.projectId, 10);
        },
        stationId(): number {
            return parseInt(this.$route.params.stationId, 10);
        },
        photos(): NoteMedia[] {
            return NoteMedia.onlyPhotos(this.$state.notes.media);
        },
        media(): PortalNoteMedia[] {
            return this.$state.notes.media;
        },
    },
    data: (): {
        photoOptions: ListItemOption[];
        loadingPhotoId: number | null;
    } => {
        return {
            photoOptions: [],
            loadingPhotoId: null,
        };
    },
    methods: {
        async onPhotoOptionClick(event: string, photo: PortalNoteMedia) {
            if (event === "delete-photo") {
                await this.$confirm({
                    message: this.$tc("confirmDeletePhoto"),
                    button: {
                        no: this.$tc("cancel"),
                        yes: this.$tc("delete"),
                    },
                    callback: async (confirm) => {
                        if (confirm) {
                            this.loadingPhotoId = photo.id;
                            this.$services.api
                                .deleteMedia(photo.id)
                                .then(async () => {
                                    await this.$store.dispatch(ActionTypes.NEED_NOTES, { id: this.$route.params.stationId });
                                })
                                .finally(() => {
                                    this.onFinishedPhotoAction();
                                });
                        }
                    },
                });
            }
            if (event === "set-as-station-image") {
                this.loadingPhotoId = photo.id;
                this.$services.api
                    .setStationImage(this.stationId, photo.id)
                    .then(async () => {
                        await this.$store.dispatch(ActionTypes.NEED_NOTES, { id: this.$route.params.stationId });
                    })
                    .finally(() => {
                        this.onFinishedPhotoAction();
                    });
            }
        },
        onFinishedPhotoAction(): void {
            window.scrollTo({ top: 0, left: 0, behavior: "smooth" });
            this.loadingPhotoId = null;
        },
        upload(this: any, ev) {
            const image = ev.target.files[0];

            const reader = new FileReader();
            reader.onload = (ev) => {
                if (ev?.target?.result) {
                    const photo = new AddedPhoto(image.type, image, ev.target.result);
                    this.$services.api.uploadStationMedia(this.stationId, photo.key, photo.file).then(() => {
                        this.$store.dispatch(ActionTypes.NEED_NOTES, { id: this.stationId });
                        return [];
                    });
                }
            };

            reader.readAsDataURL(image);
        },
    },
    mounted() {
        this.photoOptions = [
            {
                label: this.$tc("station.btn.setAsStationImage"),
                event: "set-as-station-image",
                icon: "icon-home",
            },
            {
                label: this.$tc("station.btn.deletePhoto"),
                event: "delete-photo",
                icon: "icon-trash",
            },
        ];
    },
    beforeMount(): Promise<any> {
        return this.$store.dispatch(ActionTypes.NEED_NOTES, { id: this.stationId });
    },
});
</script>

<style scoped lang="scss">
@import "src/scss/mixins";

.photo-wrap {
    margin-top: 10px;
    flex: 0 0 calc(50% - 5px);
    position: relative;
    min-height: 200px;
    background-color: #e2e4e6;

    &:nth-of-type(1) {
        flex: 0 0 100%;
        margin-top: 15px;
    }

    @include bp-down($xs) {
        &:nth-of-type(even) {
            ::v-deep .options-btns {
                right: auto;
            }
        }
    }

    ::v-deep img {
        width: 100%;
        min-height: 200px;
        height: 100%;
        object-fit: cover;
        border-radius: 2px;
    }
}

.photo-options {
    @include position(absolute, 20px 20px null null);
    width: 35px;
    height: 35px;
    background-color: #fff;
    border: solid 1px #cccdcf;
    padding: 0;
    border-radius: 50px;
    z-index: $z-index-top;

    ::v-deep .options-trigger {
        padding: 0;
        display: flex;
        justify-content: center;
        font-weight: bold;
        transform: translate(-1px, -4px);
    }

    ::v-deep .options-btns {
        right: 0;
        top: 40px;
        border: solid 1px var(--color-border);
        border-radius: 1px;
        padding-right: 25px;

        > * {
            font-size: 14px;
        }
    }
}

input[type="file"] {
    display: none;
}

.button-social {
    padding: 0 20px;

    label {
        display: flex;
        align-items: center;
        cursor: pointer;
    }
}

.empty-photos {
    margin-top: 10px;
}
</style>
