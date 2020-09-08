<template>
    <router-link :to="{ name: 'viewProject', params: { id: project.id } }" class="project-container" v-if="visible">
        <div class="project-image-container">
            <!--
            <div v-if="invited" class="invited-icon">
                Project Invite
            </div>
            -->
            <ProjectPhoto :project="project" />
        </div>

        <img v-if="project.private" alt="Private project" src="@/assets/private.png" class="private-icon" />

        <div class="project-name">{{ project.name }}</div>
        <div class="project-description">{{ project.description }}</div>
        <div class="invited-container" v-if="invited">
            <div class="accept" v-on:click.stop.prevent="onAccept">
                <img alt="Close" src="@/assets/tick.png" />
                Accept Invite
            </div>
            <div class="reject" v-on:click.stop.prevent="onDecline">
                <img alt="Close" src="@/assets/Delete.png" />
                Decline
            </div>
        </div>
        <div class="social-container" v-else>
            <div class="social follows" v-if="project.following">
                <img alt="Follows" src="@/assets/heart.png" class="follow-icon" />
                <span>{{ project.following.total }}</span>
            </div>
            <!--<div class="social notifications" v-if="!project.notifications">
                <img alt="Notifications" src="@/assets/notification.png" class="notify-icon" />
                <span>2</span>
            </div>
            <div class="social comments" v-if="!project.comments">
                <img alt="Comments" src="@/assets/comment.png" class="comment-icon" />
                <span>3</span>
            </div>-->
        </div>
    </router-link>
</template>
<script lang="ts">
import CommonComponents from "@/views/shared";
import * as ActionTypes from "@/store/actions";

export default {
    name: "ProjectThumbnail",
    components: {
        ...CommonComponents,
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
@import '../../scss/mixins';

.project-container {
    flex-basis: calc(33.33% - 24px);
    border: 1px solid #d8dce0;
    margin: 0 12px 40px;
    padding-bottom: 17px;
    transition: all 0.33s;
    flex-direction: column;
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
        flex-basis: calc(100% - 24px);
    }
}
.project-name {
    font-weight: bold;
    font-size: 16px;
    margin: 13px 15px 0 15px;
}
.project-description {
    overflow-wrap: break-word;
    font-weight: lighter;
    font-size: 14px;
    margin: 2px 15px 15px 15px;
}
.project-image-container {
    height: 138px;
    text-align: center;
    border-bottom: 1px solid #d8dce0;

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
    display: inline-block;
    font-size: 14px;
    font-weight: 600;
    margin: 0 14px 0 15px;

    &-container {
        margin-top: auto;
    }
}
.social img {
    float: left;
    margin: 2px 4px 0 0;
}
.invited-container {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 7px 15px 0;
    margin-top: auto;

    img {
        margin: 0 5px 0 3px;
    }
}
.invited-container .accept {
    padding: 5px;
    text-align: center;
    font-size: 14px;
    font-weight: 900;
    color: #2c3e50;
    border-radius: 3px;
    border: solid 1px #cccdcf;
    white-space: nowrap;
    @include flex(center);
}
.invited-container .reject {
    padding: 5px;
    text-align: center;
    font-size: 14px;
    font-weight: 900;
    color: #2c3e50;
    white-space: nowrap;
    @include flex(center);
}
</style>
