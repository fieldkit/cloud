<template>
    <div class="project-container" v-if="visible">
        <router-link :to="{ name: 'viewProject', params: { id: project.id } }">
            <div class="project-image-container">
                <ProjectPhoto :project="project" />
            </div>

            <img v-if="project.private" alt="Private project" src="@/assets/private.png" class="private-icon" />

            <div class="project-name">{{ project.name }}</div>
            <div class="project-description">{{ project.description }}</div>
        </router-link>
        <div class="invited-container" v-if="invited">
            <div class="accept" v-on:click.stop.prevent="onAccept">
                <img alt="Close" src="@/assets/icon-tick-blue.svg" />
                <span>{{ $t("project.invite.accept") }}</span>
            </div>
            <div class="reject" v-on:click.stop.prevent="onDecline">
                <img alt="Close" src="@/assets/icon-close-bold.svg" />
                <span>{{ $t("project.invite.decline") }}</span>
            </div>
        </div>
        <div class="social-container" v-else>
            <div class="social follows" v-if="project.following">
                <FollowControl :project="project">
                    <template #default="{ following, followers, follow, unfollow }">
                        <span v-if="following" v-on:click="unfollow" class="icon">
                            <img src="@/assets/icon-heart.svg" class="icon" />
                        </span>
                        <span v-if="!following" v-on:click="follow" class="icon">
                            <img src="@/assets/icon-heart-gray.svg" class="icon" />
                        </span>
                        {{ followers }}
                    </template>
                </FollowControl>
            </div>

            <!--<div class="social notifications" v-if="!project.notifications">
                <img alt="Notifications" src="@/assets/icon-notification.svg" width="15" />
                <span>2</span>
            </div>
            <div class="social comments" v-if="!project.comments" width="17">
                <img alt="Comments" src="@/assets/icon-message.svg" />
                <span>3</span>
            </div>-->
        </div>
    </div>
</template>
<script lang="ts">
import CommonComponents from "@/views/shared";
import FollowControl from "@/views/shared/FollowControl.vue";
import * as ActionTypes from "@/store/actions";

export default {
    name: "ProjectThumbnail",
    components: {
        ...CommonComponents,
        FollowControl,
    },
    props: {
        project: {
            type: Object,
            required: true,
        },
        invited: {
            type: Boolean,
            default: false,
        },
    },
    data() {
        return {
            visible: true,
            accepted: false,
        };
    },
    methods: {
        getImageUrl(this: any, project) {
            return this.$config.baseUrl + project.photo;
        },
        onAccept(this: any, ev: any) {
            console.log("accept", ev);
            return this.$store.dispatch(ActionTypes.ACCEPT_PROJECT, { projectId: this.project.id }).then(() => {
                this.accepted = true;
                this.visible = false;
            });
        },
        onDecline(this: any, ev: any) {
            console.log("decline", ev);
            return this.$store.dispatch(ActionTypes.DECLINE_PROJECT, { projectId: this.project.id }).then(() => {
                this.visible = false;
            });
        },
    },
};
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.project-container {
    flex-basis: 270px;
    border: 1px solid var(--color-border);
    margin: 0 12px 40px;
    padding-bottom: 17px;
    transition: all 0.33s;
    flex-direction: column;
    line-height: 1.4;
    @include flex();
    @include position(relative, 0 0 0 0);

    @include attention() {
        top: -3px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
    }

    @include bp-down($sm) {
        flex-basis: calc(50% - 24px);
    }

    @include bp-down($xs) {
        margin: 0 0 40px;
        flex-basis: 100%;
    }
}
.project-name {
    font-family: var(--font-family-bold);
    font-size: 16px;
    margin: 13px 15px 0 15px;
}
.project-description {
    overflow-wrap: break-word;
    font-family: $font-family-light;
    font-size: 14px;
    margin: 2px 15px 15px 15px;
    overflow-wrap: anywhere;
}
.project-image-container {
    height: 138px;
    text-align: center;
    border-bottom: 1px solid var(--color-border);
}
::v-deep .project-image {
    height: 100%;
    width: 100%;
    object-fit: cover;
}
.invited-icon {
    float: right;
    position: relative;
    z-index: 10;
    top: 0;
    right: 0;
}
.private-icon {
    float: right;
    margin: -14px 14px 0 0;
    position: relative;
    z-index: 10;
}
.social {
    @include flex(center);
    font-size: 14px;
    font-weight: 500;
    margin: 0 14px 0 15px;

    img {
        margin-right: 5px;
    }

    &-container {
        margin-top: auto;
        @include flex();
    }
}

.invited-container {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 7px 15px 0;
    margin-top: auto;

    * {
        font-family: var(--font-family-bold);
    }

    img {
        margin: 0 5px 0 3px;
        height: 100%;
    }
}
.invited-container .accept {
    padding: 5px;
    text-align: center;
    font-size: 14px;
    color: #2c3e50;
    border-radius: 3px;
    border: solid 1px #cccdcf;
    white-space: nowrap;
    flex-grow: 1;
    @include flex(center);
}
.invited-container .reject {
    padding: 5px;
    font-size: 14px;
    color: #2c3e50;
    white-space: nowrap;
    flex-grow: 1;
    @include flex(center);

    span {
        padding-top: 2px;
    }
}
</style>
