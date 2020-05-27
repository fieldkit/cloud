<template>
    <!-- bottom margin has to be minus at least (height - 80) -->
    <div id="activity-feed-container">
        <div id="close-feed-btn" v-on:click="closeFeed">
            <img alt="Close" src="../assets/close.png" />
        </div>
        <div class="heading">Activity History</div>
        <div class="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <template v-if="!loading">
            <div class="inner-feed-container" v-if="activities.length > 0">
                <!-- just ingestions at the moment -->
                <div v-for="activity in activities" v-bind:key="activity.id" class="activity">
                    <div class="activity-icon">
                        <img :src="activity.icon" />
                    </div>
                    <div class="activity-text-container" v-if="activity.type == 'ingestion'">
                        <div class="activity-heading">Uploaded Data</div>
                        <div class="activity-date">{{ activity.time.toLocaleDateString() }}</div>
                        <div class="activity-text">
                            {{ activity.records.toLocaleString() }} readings uploaded from {{ activity.name }}, with
                            {{ activity.errors ? " an error." : "no errors." }}
                            <img src="../assets/Icon_Warning_error.png" v-if="activity.errors" class="activity-error" />
                        </div>
                    </div>
                    <div class="activity-text-container" v-if="activity.type == 'update'">
                        <div class="activity-heading">{{ activity.name }}</div>
                        <div class="activity-date">{{ activity.time.toLocaleDateString() }}</div>
                        <div class="activity-text">
                            {{ activity.text }}
                        </div>
                    </div>
                </div>
            </div>
            <div class="inner-feed-container" v-else>
                No recent activity.
            </div>
        </template>
    </div>
</template>
<script>
import FKApi from "../api/api";

export default {
    name: "ProjectActivity",
    data: () => {
        return {
            loading: true,
            activities: [],
        };
    },
    props: ["project", "viewing", "users"],
    watch: {
        viewing() {
            if (this.viewing) {
                this.loading = true;
                this.fetchActivity();
            }
        },
    },
    async beforeCreate() {
        this.api = new FKApi();
    },
    methods: {
        fetchActivity() {
            let imgPath = require.context("../assets/", false, /\.png$/);
            const img = "compass-icon.png";
            const compassImg = imgPath("./" + img);

            this.api.getProjectActivity(this.project.id).then(result => {
                this.activities = result.activities.map((a, i) => {
                    // starting with just two types of activity
                    if (a.type == "StationIngestion") {
                        return {
                            id: "ingestion-" + i,
                            type: "ingestion",
                            icon: compassImg,
                            name: a.station.name,
                            time: new Date(a.created_at),
                            records: a.meta.data.records,
                            errors: a.meta.errors,
                        };
                    }
                    if (a.type == "ProjectUpdate") {
                        const user = this.users.find(u => {
                            return u.user.id == a.meta.author.id;
                        });
                        return {
                            id: "update-" + i,
                            type: "update",
                            icon: user && user.userImage ? user.userImage : compassImg,
                            name: a.meta.author.name,
                            time: new Date(a.created_at),
                            text: a.meta.body,
                        };
                    }
                });
                this.loading = false;
            });
        },
        closeFeed() {
            this.$emit("closeActivity");
        },
    },
};
</script>

<style scoped>
.loading {
    text-align: center;
    width: 100%;
    float: left;
}
.loading img {
    width: 100px;
}
#activity-feed-container {
    position: relative;
    float: right;
    width: 375px;
    height: 100vh;
    margin: -80px -1px -2000px 0;
    border: 1px solid #d8dce0;
    background: white;
    overflow: scroll;
    z-index: 3;
}
.heading {
    font-size: 20px;
    font-weight: 600;
    float: left;
    margin: 25px;
}
#close-feed-btn {
    float: right;
    margin: 25px 25px 0 0;
    cursor: pointer;
}
.inner-feed-container {
    float: left;
    clear: both;
    margin: 0 25px 25px 25px;
    font-size: 20px;
}
.activity {
    float: left;
    margin: 15px 0;
}
.activity-icon {
    float: left;
    width: 40px;
    padding: 5px 0 0 0;
}
.activity-icon img {
    width: 30px;
}
.activity-text-container {
    float: right;
    width: 275px;
}
.activity-heading {
    font-size: 16px;
    font-weight: 600;
    float: left;
}
.activity-date {
    font-size: 12px;
    color: #6a6d71;
    font-weight: 600;
    float: left;
    margin: 4px 0 0 6px;
}
.activity-text {
    float: left;
    clear: both;
    font-size: 14px;
    margin: 4px 0 0 0;
}
.activity-error {
    width: 16px;
    vertical-align: sub;
}
</style>
