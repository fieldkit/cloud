<template>
    <div class="follow-panel">
        <FollowControl :project="project">
            <template #default="{ following, followers, follow, unfollow }">
                <span v-if="following" v-on:click="unfollow" class="icon">
                    <img src="@/assets/icon-heart.svg" class="icon" />
                </span>
                <span v-if="!following" v-on:click="follow" class="icon">
                    <img src="@/assets/icon-heart-gray.svg" class="icon" />
                </span>
                <span v-show="followers > 1">{{ followers }} Follows</span>
                <span v-show="followers <= 1">{{ followers }} Follow</span>
            </template>
        </FollowControl>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { Project } from "@/api";
import FollowControl from "./FollowControl.vue";

export default Vue.extend({
    name: "FollowPanel",
    components: {
        FollowControl,
    },
    props: {
        project: {
            type: Object as PropType<Project>,
            required: true,
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.follow-panel {
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 18px;
    font-weight: 500;
    padding-top: 24px;
    border-radius: 1px;
    border-top: solid 1px var(--color-border);
    margin-top: 10px;

    @include bp-down($xs) {
        padding-top: 14px;
    }
}
.icon {
    cursor: pointer;
    margin-right: 10px;
    width: 22px;
    height: 20px;
}
</style>
