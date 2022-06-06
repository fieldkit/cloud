<template>
    <StandardLayout>
        <div class="container-wrap">
            <DoubleHeader
                backRoute="viewStation"
                :title="$tc('station.allPhotos')"
                :backTitle="$tc('layout.backStationDashboard')"
                :backRouteParams="{ id: $route.params.stationId }"
            >
                <template v-slot:default>
                    <div class="button">
                        <i class="icon icon-share"></i>
                        Add photos
                    </div>
                </template>
            </DoubleHeader>

            <div class="flex flex-wrap flex-space-between" v-if="media">
                <div class="photo-wrap" v-for="photo in photos" v-bind:key="photo.key">
                    <button class="photo-options">
                        <ListItemOptions @listItemOptionClick="onPhotoOptionClick($event, photo)" :options="photoOptions"></ListItemOptions>
                    </button>
                    <AuthenticatedPhoto :url="photo.url" />
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
import { NoteMedia, Notes, PortalNoteMedia } from "@/views/notes/model";
import { ActionTypes } from "@/store";
import ListItemOptions from "@/views/shared/ListItemOptions.vue";
import NewPhoto from "@/assets/image-placeholder.svg";

export default Vue.extend({
    name: "StationPhotosView",
    components: {
        StandardLayout,
        DoubleHeader,
        AuthenticatedPhoto,
        ListItemOptions,
    },
    computed: {
        photos(): NoteMedia[] {
            return NoteMedia.onlyPhotos(this.$state.notes.media);
        },
        media(): PortalNoteMedia[] {
            return this.$state.notes.media;
        },
    },
    data: (): {
        photoOptions: {
            label: string;
            event: string;
        }[];
    } => {
        return {
            photoOptions: [],
        };
    },
    methods: {
        onPhotoOptionClick() {
            console.log("da");
        },
    },
    mounted() {
        this.photoOptions = [
            {
                label: this.$tc("station.replacePhoto"),
                event: "replace-image",
            },
            {
                label: this.$tc("station.setAsStationImage"),
                event: "use-as-station-image",
            },
            {
                label: this.$tc("station.setAsStationImage"),
                event: "delete-media",
            },
        ];
    },
    beforeMount(): Promise<any> {
        return this.$store.dispatch(ActionTypes.NEED_NOTES, { id: this.$route.params.stationId });
    },
});
</script>

<style scoped lang="scss">
@import "src/scss/mixins";

.photo-wrap {
    margin-top: 10px;
    flex: 0 0 calc(50% - 5px);
    position: relative;

    &:nth-of-type(1) {
        flex: 0 0 100%;
        margin-top: 15px;
    }

    img {
        width: 100%;
        height: 100%;
        object-fit: cover;
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
        border: solid 1px #d8dce0;
        font-size: 14px;
    }
}
</style>
