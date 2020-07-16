<template>
    <div class="image-container">
        <img alt="Image" :src="placeholderImage" v-if="!image || (!image.url && !preview)" />
        <img alt="Image" :src="$config.baseUrl + image.url" class="image" v-if="image && image.url && !preview" />
        <img alt="Image" :src="preview" class="image" v-if="preview" />
        <br />
        <input type="file" accept="image/gif, image/jpeg, image/png" @change="upload" />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import NewPhoto from "../../assets/Profile_Image.png";

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
            required: true,
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
        };
    },
    computed: {
        placeholderImage() {
            return this.placeholder || NewPhoto;
        },
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

<style scoped></style>
