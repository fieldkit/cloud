<template>
    <div>
        <slot :following="following" :followers="followers" :follow="followProject" :unfollow="unfollowProject"></slot>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { Project } from "@/api";

export class ProjectFollowing {
    constructor(public readonly following: boolean, public readonly total: number) {}
}

export default Vue.extend({
    name: "FollowControl",
    props: {
        project: {
            type: Object as PropType<Project>,
            required: true,
        },
    },
    data(): {
        action: ProjectFollowing | null;
    } {
        return {
            action: null,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        following(): boolean {
            return this.action ? this.action.following : this.project.following.following;
        },
        followers(): number {
            return this.action ? this.action.total : this.project.following.total;
        },
    },
    methods: {
        async followProject(): Promise<void> {
            if (this.isAuthenticated) {
                await this.$store.dispatch(ActionTypes.PROJECT_FOLLOW, { projectId: this.project.id });
                this.action = new ProjectFollowing(true, this.followersExceptMyself() + 1);
            } else {
                console.log("follow-control: anonymous");
            }
        },
        async unfollowProject(): Promise<void> {
            if (this.isAuthenticated) {
                await this.$store.dispatch(ActionTypes.PROJECT_UNFOLLOW, { projectId: this.project.id });
                this.action = new ProjectFollowing(false, this.followersExceptMyself());
            } else {
                console.log("follow-control: anonymous");
            }
        },
        followersExceptMyself(): number {
            if (this.project.following.following) {
                return this.project.following.total - 1;
            }
            return this.project.following.total;
        },
    },
});
</script>

<style scoped lang="scss"></style>
