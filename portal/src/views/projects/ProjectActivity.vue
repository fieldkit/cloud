<template>
    <div :class="'project-activity ' + containerClass">
        <div class="heading">
            <div class="title">
                Activity History
            </div>
            <div class="close-button" v-on:click="onClose">
                <img alt="Close" src="../../assets/close.png" />
            </div>
        </div>
        <div class="feed-container" v-if="visibleActivities.length > 0">
            <div v-for="activity in visibleActivities" v-bind:key="activity.id" class="activity">
                <div class="icon">
                    <img :src="activity.icon" />
                </div>
                <div class="panel" v-if="activity.type == 'ingestion'">
                    <div class="activity-heading">
                        <div class="title">Uploaded Data</div>
                        <div class="date">{{ activity.time.toLocaleDateString() }}</div>
                    </div>
                    <div class="activity-text">
                        {{ activity.records.toLocaleString() }} readings uploaded from {{ activity.name }}, with
                        {{ activity.errors ? " an error." : "no errors." }}
                        <img src="@/assets/Icon_Warning_error.png" v-if="activity.errors" class="activity-error" />
                    </div>
                </div>
                <div class="panel" v-if="activity.type == 'update'">
                    <div class="activity-heading">
                        <div class="title">{{ activity.name }}</div>
                        <div class="date">{{ activity.time.toLocaleDateString() }}</div>
                    </div>
                    <div class="activity-text">
                        {{ activity.text }}
                    </div>
                </div>
                <div class="panel" v-if="activity.type == 'deploy'">
                    <div class="activity-heading">
                        <div class="title">Deployed Station</div>
                        <div class="date">{{ activity.time.toLocaleDateString() }}</div>
                    </div>
                    <div class="activity-text">{{ activity.name }} was deployed.</div>
                    <div class="activity-text">
                        <img src="@/assets/icon-location.png" />
                        {{ activity.location[1] | prettyCoordinate }}, {{ activity.location[0] | prettyCoordinate }}
                    </div>
                </div>
            </div>
        </div>
        <div class="feed-container empty" v-else>
            No recent activity.
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "ProjectActivity",
    props: {
        containerClass: {
            type: String,
            default: "",
        },
        user: {
            required: true,
        },
        displayProject: {
            required: true,
        },
    },
    data: () => {
        return {
            activities: {
                activities: [],
            },
        };
    },
    computed: {
        loading(this: any) {
            return false;
        },
        visibleActivities(this: any) {
            const compass = this.$loadAsset("compass-icon.png");
            return _.take(
                this.activities.activities.map((a, i) => {
                    if (a.type == "StationIngestion") {
                        return {
                            id: "ingestion-" + i,
                            type: "ingestion",
                            icon: compass,
                            name: a.station.name,
                            time: new Date(a.createdAt),
                            records: a.meta.data.records,
                            errors: a.meta.errors,
                        };
                    } else if (a.type == "ProjectUpdate") {
                        return {
                            id: "update-" + i,
                            type: "update",
                            icon: compass,
                            name: a.meta.author.name,
                            time: new Date(a.createdAt),
                            text: a.meta.body,
                        };
                    } else if (a.type == "StationDeployed") {
                        return {
                            id: "deploy-" + i,
                            type: "deploy",
                            icon: compass,
                            name: a.station.name,
                            time: new Date(a.meta.deployedAt),
                            location: a.meta.location,
                        };
                    }
                }),
                50
            );
        },
    },
    mounted(this: any) {
        return this.refresh();
    },
    methods: {
        refresh(this: any) {
            return new FKApi().getProjectActivity(this.displayProject.id).then((activities) => {
                this.activities = activities;
                return activities;
            });
        },
        onClose() {
            this.$emit("close");
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
.heading {
    padding: 1em;
    display: flex;

    @include bp-down($sm) {
        font-size: 20px;
    }
}
.heading .title {
    font-size: 20px;
    font-weight: 500;
}
.heading .close-button {
    margin-left: auto;
    cursor: pointer;
}
.feed-container {
    display: flex;
    flex-direction: column;
}
.feed-container.empty {
    padding: 1em;
}
.feed-container .activity {
    display: flex;
    flex-direction: row;
    padding-right: 1em;
}

.feed-container .activity .icon {
    display: flex;
    flex-direction: column;
}
.activity .icon {
    width: 40px;
    padding-left: 1em;
    padding-right: 1em;
}
.activity .icon img {
    width: 30px;
}
.feed-container .activity .panel {
    display: flex;
    flex-direction: column;
    margin-bottom: 1em;
}
.activity-heading {
    display: flex;
}
.activity-heading .title {
    font-size: 16px;
    font-weight: 600;
}
.activity-heading .date {
    font-size: 12px;
    color: #6a6d71;
    font-weight: 600;
    margin: 4px 0 0 6px;
}
.activity-text {
    clear: both;
    font-size: 14px;
    margin: 4px 0 0 0;
}
.activity-error {
    width: 16px;
    vertical-align: sub;
}
</style>
