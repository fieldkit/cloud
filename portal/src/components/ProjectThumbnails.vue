<template>
    <div>
        <div v-for="project in displayProjects" v-bind:key="project.id" class="project-container">
            <router-link :to="{ name: 'viewProject', params: { id: project.id } }">
                <div class="project-image-container">
                    <img alt="Fieldkit Project" v-if="project.media_url" :src="getImageUrl(project)" class="project-image" />
                    <img alt="Default Fieldkit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />
                </div>
                <img v-if="project.private" alt="Private project" src="../assets/private.png" class="private-icon" />

                <div class="project-name">{{ project.name }}</div>
                <div class="project-description">{{ project.description }}</div>
                <div class="stats-icon-container">
                    <div class="stat follows" v-if="project.numFollowers">
                        <img alt="Follows" src="../assets/heart.png" class="follow-icon" />
                        <span>{{ project.numFollowers }}</span>
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
import { API_HOST } from "../secrets";
import FKApi from "../api/api";

export default {
    name: "ProjectThumbnails",
    data: () => {
        return {
            baseUrl: API_HOST,
            displayProjects: [],
        };
    },
    props: ["projects"],
    mounted() {},
    async beforeCreate() {
        this.api = new FKApi();
    },
    watch: {
        projects() {
            if (this.projects) {
                this.projects.forEach(p => {
                    this.api
                        .getProjectFollows(p.id)
                        .then(result => {
                            // fudging this so every project has at least one follower for now
                            // owner should be following, anyway
                            p.numFollowers = result.followers.length + 1;
                            this.displayProjects.push(p);
                        })
                        .catch(() => {
                            this.displayProjects.push(p);
                        });
                });
            }
        },
    },
    methods: {
        getImageUrl(project) {
            return this.baseUrl + "/projects/" + project.id + "/media/";
        },
    },
};
</script>

<style scoped>
.project-container {
    position: relative;
    float: left;
    width: 270px;
    height: 265px;
    margin: 0 12px 25px 12px;
    border: 1px solid rgb(235, 235, 235);
}
.project-name {
    font-weight: bold;
    font-size: 16px;
    margin: 10px 15px 0 15px;
}
.project-description {
    font-weight: lighter;
    font-size: 14px;
    margin: 0 15px 10px 15px;
}
.project-image-container {
    height: 138px;
    text-align: center;
    border-bottom: 1px solid rgb(235, 235, 235);
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
