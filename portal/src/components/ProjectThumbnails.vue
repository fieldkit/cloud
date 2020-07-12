<template>
    <div class="project-thumbnails">
        <div v-for="project in projects" v-bind:key="project.id" class="project-container">
            <router-link :to="{ name: 'viewProject', params: { id: project.id } }">
                <div class="project-image-container">
                    <img alt="FieldKit Project" v-if="project.photo" :src="getImageUrl(project)" class="project-image" />
                    <img alt="FieldKit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />
                </div>

                <img v-if="project.private" alt="Private project" src="../assets/private.png" class="private-icon" />

                <div class="project-name">{{ project.name }}</div>
                <div class="project-description">{{ project.description }}</div>
                <div class="stats-icon-container">
                    <div class="stat follows" v-if="project.numberOfFollowers">
                        <img alt="Follows" src="../assets/heart.png" class="follow-icon" />
                        <span>{{ project.numberOfFollowers }}</span>
                    </div>
                    <div class="stat notifications" v-if="project.notifications">
                        <img alt="Notifications" src="../assets/notification.png" class="notify-icon" />
                        <span>2</span>
                    </div>
                    <div class="stat comments" v-if="project.comments">
                        <img alt="Comments" src="../assets/comment.png" class="comment-icon" />
                        <span>3</span>
                    </div>
                </div>
            </router-link>
        </div>
    </div>
</template>

<script>
import FKApi from "@/api/api";

export default {
    name: "ProjectThumbnails",
    data: () => {
        return {};
    },
    props: { projects: { required: true } },
    methods: {
        getImageUrl(project) {
            return this.$config.baseUrl + "/" + project.photo;
        },
    },
};
</script>

<style scoped>
.project-thumbnails {
    display: flex;
    flex-wrap: wrap;
    text-align: left;
}
.project-container {
    width: 270px;
    height: 265px;
    border: 1px solid #d8dce0;
    margin-right: 20px;
    margin-bottom: 20px;
}
.project-name {
    font-weight: bold;
    font-size: 16px;
    margin: 10px 15px 0 15px;
}
.project-description {
    overflow-wrap: break-word;
    font-weight: lighter;
    font-size: 14px;
    margin: 0 15px 10px 15px;
}
.project-image-container {
    height: 138px;
    text-align: center;
    border-bottom: 1px solid #d8dce0;
}
.project-image {
    max-width: 270px;
    max-height: 138px;
}
.private-icon {
    float: right;
    margin: -14px 14px 0 0;
    position: relative;
    z-index: 10;
}
.stats-icon-container {
    position: absolute;
    bottom: 20px;
}
.stat {
    display: inline-block;
    font-size: 14px;
    font-weight: 600;
    margin: 0 14px 0 15px;
}
.stat img {
    float: left;
    margin: 2px 4px 0 0;
}
</style>
