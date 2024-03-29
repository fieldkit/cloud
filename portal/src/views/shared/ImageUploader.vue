<template>
    <div v-if="!image || (!image.url && !preview)" class="placeholder-container">
        <img alt="Image" :src="placeholderImage" @click="imageInputClick" />
        <div class="upload-trigger">
            <label for="imageInput" class="upload-trigger">{{ $t("user.profile.photo.choose") }}</label>
            <span class="upload-trigger-msg">{{ $t("user.profile.photo.noFile") }}</span>
        </div>
        <input id="imageInput" ref="imageInput" type="file" accept="image/gif, image/jpeg, image/png" @change="upload" />
    </div>

    <div class="image-container" v-else>
        <img alt="Image" :src="photo" class="img" v-if="photo && !preview" />
        <img alt="Image" :src="preview" class="img" v-if="preview" />
        <label for="imageInput" class="upload-trigger-change">
            <template>{{ $t("user.profile.photo.change") }}</template>
        </label>
        <input id="imageInput" type="file" accept="image/gif, image/jpeg, image/png" @change="upload" />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import NewPhoto from "../../assets/profile-image-placeholder.svg";

export interface Image {
    url: string;
}

export interface UploadedImage {
    type: string;
    image: unknown;
    file: unknown;
}

export default Vue.extend({
    name: "ImageUploader",
    props: {
        image: {
            type: Object,
            required: false,
        },
        placeholder: {
            type: String,
            required: false,
        },
        allowPreview: {
            type: Boolean,
            default: true,
        },
    },
    data: () => {
        return {
            imageType: null,
            preview: null,
            acceptedTypes: ["jpg", "jpeg", "png", "gif"],
            photo: null,
        };
    },
    computed: {
        placeholderImage() {
            return this.placeholder || NewPhoto;
        },
    },
    mounted(this: any) {
        if (this.image && this.image.url) {
            return this.$services.api.loadMedia(this.image.url, { size: 800 }).then((photo) => {
                this.photo = photo;
            });
        }
    },
    methods: {
        acceptable(this: any, files: { type: string }[]): boolean {
            if (files.length != 1) {
                return false;
            }
            return this.acceptedTypes.filter((t) => files[0].type.indexOf(t) > -1).length > 0;
        },
        upload(this: any, ev) {
            this.preview = null;

            if (!this.acceptable(ev.target.files)) {
                return;
            }

            const image = ev.target.files[0];
            this.imageType = ev.target.files[0].type;

            if (this.allowPreview) {
                const reader = new FileReader();
                reader.readAsDataURL(image);
                reader.onload = (ev) => {
                    const r = ev?.target?.result;
                    if (r) {
                        this.preview = r;
                    }
                };
            }

            const uploaded: UploadedImage = {
                type: this.imageType,
                image: image, // TODO Remove
                file: image,
            };

            this.$emit("change", uploaded);
        },
        imageInputClick() {
            (this.$refs.imageInput as HTMLElement).click();
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.image-container {
    @include flex(baseline);
    flex-wrap: wrap;

    label {
        cursor: pointer;
    }

    img {
        margin-right: 12px;
    }
}

input {
    display: none;
}

.upload-trigger {
    margin: 14px 0 0;
    font-size: 14px;
    @include flex(flex-end);

    label {
        @include flex(center, center);
        padding: 0 17px;
        height: 30px;
        border-radius: 2px;
        border: solid 1px #cccdcf;
        background-color: #ffffff;
        font-size: 14px;
        font-weight: 600;
        margin: 0 15px 0 0;
    }

    &-msg {
        color: #6a6d71;
        font-weight: 500;
    }

    &-change {
        margin-top: 12px;
        font-size: 14px;
        color: #6a6d71;
        font-weight: 500;
    }
}
.img-placeholder {
    width: 220px;
    height: 200px;
    max-height: unset;
}
</style>
