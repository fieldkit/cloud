<template>
    <div class="image-container">
        <img alt="User image" src="../../assets/Profile_Image.png" v-if="!image.url && !preview" />
        <img alt="User image" :src="$config.baseUrl + image.url" class="image" v-if="image.url && !preview" />
        <img alt="User image" :src="preview" class="image" v-if="!image.url || preview" />
        <br />
        <input type="file" accept="image/gif, image/jpeg, image/png" @change="upload" />
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export interface Image {
    url: string;
}

export default Vue.extend({
    name: "ImageUploader",
    components: {},
    props: {
        image: {
            type: Object,
            required: true,
        },
    },
    data: () => {
        return {
            imageType: null,
            preview: null,
            acceptedTypes: ["jpg", "jpeg", "png", "gif"],
        };
    },
    methods: {
        acceptable(files: { type: string }[]): boolean {
            if (files.length != 1) {
                return false;
            }
            return this.acceptedTypes.filter((t) => files[0].type.indexOf(t) > -1).length > 0;
        },
        upload(ev) {
            this.preview = null;

            if (!this.acceptable(ev.target.files)) {
                return;
            }

            const image = ev.target.files[0];
            this.imageType = ev.target.files[0].type;

            const reader = new FileReader();
            reader.readAsDataURL(image);
            reader.onload = (ev) => {
                this.preview = ev.target.result;
            };

            this.$emit("change", {
                type: this.imageType,
                image: image,
            });
        },
    },
});
</script>

<style scoped></style>
