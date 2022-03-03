<template>
    <div class="project-container" v-if="project">
        <div class="header">
            <div class="left">
                <div class="project-name">{{ project.name }}</div>
                <div class="project-dashboard">Project Dashboard</div>
            </div>
            <div class="right activity">Activity</div>
        </div>
        <div class="details">
            <div class="left">
                <div class="photo">
                    <img alt="FieldKit Project" v-if="project.photo" :src="getImageUrl(project)" class="project-image" />
                    <img alt="Default FieldKit Project" v-else src="@/assets/fieldkit_project.png" class="project-image" />
                </div>

                <div class="follow-btn">
                    <span v-if="following" v-on:click="unfollowProject">
                        Following
                    </span>
                    <span v-else v-on:click="followProject">
                        <img alt="Follow" src="@/assets/heart_gray.png" class="icon" />
                        Follow
                    </span>
                </div>
            </div>

            <div class="right">
                <div class="details-heading">
                    Project Details
                </div>
                <div class="details-top">
                    <div class="details-left">
                        <div class="project-detail" v-if="project.goal">Project Goal: {{ project.goal }}</div>
                        <div class="project-detail">{{ project.description }}</div>
                    </div>
                    <div class="details-right">
                        <div class="time-container" v-if="project.startTime">
                            <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                            <template>Started: {{ project.startTime | prettyDate }}</template>
                        </div>
                        <div class="duration-container" v-if="displayProject.duration">
                            <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                            <template>{{ displayProject.duration | prettyDuration }}</template>
                        </div>
                        <div class="location-container" v-if="project.location">
                            <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                            <template>{{ project.location }}</template>
                        </div>
                        <div class="location-container" v-if="displayProject.places.native">
                            <img alt="Location" src="@/assets/icon-location.png" class="icon" />
                            <template>Native Lands: {{ displayProject.places.native }}</template>
                        </div>
                    </div>
                </div>

                <div class="details-bottom">
                    <div class="details-team">
                        <div class="title">Team</div>
                        <UserPhoto
                            v-for="projectUser in displayProject.users"
                            v-bind:key="projectUser.user.email"
                            :user="projectUser.user"
                        />
                    </div>
                    <div class="details-modules">
                        <div class="title">Modules</div>

                        <img
                            v-for="module in projectModules"
                            v-bind:key="module.name"
                            alt="Module icon"
                            class="module-icon"
                            :src="module.url"
                        />
                    </div>
                </div>
            </div>
        </div>

        <div class="project-stations">
            <ProjectStations
                :project="project"
                :admin="false"
                :mapContainerSize="mapContainerSize"
                :listSize="listSize"
                :userStations="userStations"
            />
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "@/store/strong-vue";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "ProjectLayout",
    components: {},
    props: {
        project: {
            type: Object,
            required: true,
        },
    },
    data: () => {
        return {};
    },
    computed: {},
    beforeMount() {
        console.log("ProjectLayout: beforeMount");
    },
    methods: {},
});
</script>

<style scoped></style>
