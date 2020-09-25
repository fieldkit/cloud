<template v-if="icon">
    <img v-if="userImage" :src="userImage" class="default-user-icon" />
    <img v-else src="@/assets/Profile_Image.png" class="default-user-icon" />
</template>

<script lang="ts">
import Vue, { PropType } from "vue";

interface UserPhoto {
    photo: { url: string | null };
}

export default Vue.extend({
    name: "UserPhoto",
    props: {
        user: {
            type: Object as PropType<UserPhoto>,
            required: true,
        },
        icon: {
            type: Boolean,
            default: true,
        },
    },
    computed: {
        userImage(): string | null {
            if (this.user && this.user.photo) {
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
}
</style>
