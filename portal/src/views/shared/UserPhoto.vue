<template v-if="icon">
    <div v-if="userImage">
        <img :src="userImage" class="default-user-icon" />
        <slot></slot>
    </div>
    <div v-else>
        <img src="@/assets/profile-image-placeholder.svg" class="default-user-icon" />
        <slot></slot>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";

interface UserPhoto {
    updatedAt?: number;
    photo: {
        url: string | null;
    };
}

export default Vue.extend({
    name: "UserPhoto",
    props: {
        user: {
            type: Object as PropType<UserPhoto>,
            required: false,
        },
        icon: {
            type: Boolean,
            default: true,
        },
    },
    computed: {
        userImage(): string | null {
            if (this.user && this.user.photo) {
                if (this.user.updatedAt) {
                    return this.$config.baseUrl + this.user.photo.url + "?ts=" + this.user.updatedAt;
                }
                return this.$config.baseUrl + this.user.photo.url;
            }
            return null;
        },
    },
});
</script>

<style scoped>
.default-user-icon {
    margin: 10px 12px 0 0;
    width: 44px;
    height: 44px;
    border-radius: 50%;
    object-fit: cover;
    image-rendering: -webkit-optimize-contrast;
}
</style>
