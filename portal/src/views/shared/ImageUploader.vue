<template>
    <div v-if="!image || (!image.url && !preview)" class="placeholder-container">
        <img alt="Image" :src="placeholderImage" />
        <div class="upload-trigger">
            <label for="imageInput" class="upload-trigger">Choose File</label>
            <span>No file chosen</span>
        </div>
        <input id="imageInput" type="file" accept="image/gif, image/jpeg, image/png" @change="upload" />
    </div>

    <div class="image-container" v-else>
        <img alt="Image" :src="photo" class="img" v-if="photo && !preview" />
        <img alt="Image" :src="preview" class="img" v-if="preview" />
        <label for="imageInput">
            <template>Change Image</template>
        </label>
        <input id="imageInput" type="file" accept="image/gif, image/jpeg, image/png" @change="upload" />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import NewPhoto from "../../assets/Profile_Image.png";
import FKApi from "@/api/api";

export interface Image {
    url: string;
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
            return new FKApi().loadMedia(this.image.url, { size: 800 }).then((photo) => {
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
                    this.preview = ev.target.result;
                };
            }

            this.$emit("change", {
                type: this.imageType,
                image: image, // TODO Remove
                file: image,
            });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.image-container {
    @include flex(baseline);
    flex-wrap: wrap;
}

label {
    transform: translateY(-7px);
    font-size: 14px;
    cursor: pointer;
    margin-left: 12px;

    @include bp-down($xs) {
        width: 100%;
        display: block;
        transform: translateY(5px);
        margin-left: 0;
    }
}

input {
    display: none;
}

.upload-trigger {
    margin: 14px 0 0;
    @include flex();

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
}
.img-placeholder {
    width: 220px;
    height: 200px;
    max-height: unset;
}
</style>
