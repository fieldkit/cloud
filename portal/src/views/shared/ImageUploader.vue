<template>
    <div>
        <img alt="Image" :src="placeholderImage" />
        <div class="upload-trigger">
            <label for="imageInput" class="upload-trigger"> Choose File </label>
            <span> No file chosen </span>
        </div>
    </div>

   <!-- <div class="image-container">

&lt;!&ndash;
        <template v-if="!image || (!image.url && !preview)">
&ndash;&gt;
&ndash;&gt;


        <img alt="Image" :src="$config.baseUrl + image.url" class="image" v-if="image && image.url && !preview" />
        <img alt="Image" :src="preview" class="image" v-if="preview" />
        <label for="imageInput">
            <template v-if="!preview"> Choose Image </template>
            <template v-if="preview"> Change Image </template>
        </label>
        <input id="imageInput" type="file" accept="image/gif, image/jpeg, image/png" @change="upload" />
    </div>-->

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

<style scoped lang="scss">
@import '../../scss/mixins';

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
    @include flex(center);

    label {
        @include flex(center, center);
        padding: 0 17px;
        height: 30px;
        border-radius: 2px;
        border: solid 1px #cccdcf;
        background-color: #ffffff;
        font-size: 14px;
        font-weight: 900;
        margin: 0 15px 0 0;
    }
}
</style>
