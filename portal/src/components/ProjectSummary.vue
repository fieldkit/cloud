<template>
    <div id="project-summary-container" v-if="viewingSummary">
        <div class="project-container" v-if="this.project">
            <div class="left">
                <img alt="Project image" src="../assets/fieldkit_project.png" />
            </div>
            <div class="left">
                <div id="project-name">{{ this.project.name }}</div>
                <div id="project-description">{{ this.project.description }}</div>
            </div>
            <div class="spacer"></div>
            <div class="section">
                <div class="section-control" v-on:click="toggleStations">
                    <div class="toggle-image-container">
                        <img
                            alt="Toggle FieldKit Stations section"
                            src="../assets/right.png"
                            v-if="!viewingStations"
                        />
                        <img
                            alt="Toggle FieldKit Stations section"
                            src="../assets/down.png"
                            v-if="viewingStations"
                        />
                    </div>
                    <div class="section-heading">FieldKit Stations ({{ this.stations.length }})</div>
                </div>
                <div class="section-content" v-if="viewingStations">
                    <div class="stations-container">
                        <div v-for="station in stations" v-bind:key="station.id">
                            <div class="station-cell">
                                <span class="station-name">{{ station.name }}</span> <br />
                                Last seen {{ getUpdatedDate(station) }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="spacer"></div>
            <div class="section">
                <div class="section-control" v-on:click="toggleActivity">
                    <div class="toggle-image-container">
                        <img
                            alt="Toggle FieldKit Stations section"
                            src="../assets/right.png"
                            v-if="!viewingActivity"
                        />
                        <img
                            alt="Toggle FieldKit Stations section"
                            src="../assets/down.png"
                            v-if="viewingActivity"
                        />
                    </div>
                    <div class="section-heading">Recent Activity</div>
                </div>
                <div class="section-content" v-if="viewingActivity">
                    Check back soon for a recent activity feed!
                </div>
            </div>
            <div class="spacer"></div>
            <div class="section">
                <div class="section-control" v-on:click="toggleTeam">
                    <div class="toggle-image-container">
                        <img
                            alt="Toggle FieldKit Stations section"
                            src="../assets/right.png"
                            v-if="!viewingTeam"
                        />
                        <img
                            alt="Toggle FieldKit Stations section"
                            src="../assets/down.png"
                            v-if="viewingTeam"
                        />
                    </div>
                    <div class="section-heading">Team</div>
                </div>
                <div class="section-content" v-if="viewingTeam">
                    <div class="user-container">
                        <div class="user-grid-heading">Members (1)</div>
                        <div class="user-grid-heading">Role</div>
                        <div class="user-grid-heading">Status</div>
                        <div class="row-border">
                            {{ this.user.name }} (you) <br />
                            <span class="email">{{ this.user.email }}</span>
                        </div>
                        <div class="row-border">Admin</div>
                        <div class="row-border">Active</div>
                    </div>
                </div>
            </div>
            <div class="spacer"></div>
        </div>
    </div>
</template>

<script>
export default {
    name: "ProjectSummary",
    data: () => {
        return {
            viewingSummary: false,
            viewingStations: false,
            viewingActivity: false,
            viewingTeam: false
        };
    },
    props: ["project", "stations", "user"],
    methods: {
        viewSummary() {
            this.viewingSummary = true;
        },
        toggleStations() {
            this.viewingStations = !this.viewingStations;
        },
        toggleActivity() {
            this.viewingActivity = !this.viewingActivity;
        },
        toggleTeam() {
            this.viewingTeam = !this.viewingTeam;
        },
        getUpdatedDate(station) {
            if (!station.status_json) {
                return "N/A";
            }
            const date = station.status_json.updated;
            const d = new Date(date);
            return d.toLocaleDateString("en-US");
        },
        closeSummary() {
            this.viewingSummary = false;
            if (this.$route.name != "projects") {
                this.$router.push({ name: "projects" });
            }
        }
    }
};
</script>

<style scoped>
#project-summary-container {
    background-color: #ffffff;
    width: 780px;
    position: absolute;
    top: 40px;
    left: 260px;
    padding: 0 15px 15px 15px;
    margin: 60px;
    z-index: 2;
}
.show-link {
    text-decoration: underline;
}
.project-container {
    padding: 10px;
    margin: 20px 0;
    font-size: 16px;
    font-weight: lighter;
    overflow: hidden;
}
.section {
    width: 100%;
    float: left;
}
.spacer {
    float: left;
    width: 100%;
    margin: 20px 0;
    border-bottom: 1px solid rgb(215, 220, 225);
    height: 1px;
}
.left {
    float: left;
}
.right {
    float: right;
}
.project-element {
    margin: 5px 5px 0 5px;
}
.toggle-image-container {
    display: inline-block;
    float: left;
    width: 40px;
}
.section-control {
    cursor: pointer;
}
.section-content {
    clear: both;
    float: left;
    padding-bottom: 20px;
    margin: 20px 0 10px 40px;
}
.section-heading {
    font-size: 20px;
    font-weight: normal;
    margin-top: 8px;
    display: inline-block;
}
#project-name {
    font-size: 24px;
    font-weight: bold;
    margin-left: 20px;
}
#project-description {
    margin-left: 20px;
}
.stations-container {
    display: grid;
    width: 720px;
    grid-template-columns: 1fr 1fr;
    grid-auto-rows: 40px;
    grid-gap: 40px;
}
.station-cell {
    font-size: 14px;
    padding: 10px;
    border: 1px solid rgb(215, 220, 225);
}
.station-name {
    font-weight: bold;
}
.user-container {
    display: grid;
    grid-template-columns: 280px 220px 220px;
    grid-auto-rows: 20px;
    grid-gap: 15px 0;
}
.row-border {
    border-top: 1px solid rgb(215, 220, 225);
    padding-top: 2px;
}
.user-grid-heading {
    font-weight: bold;
}
.email {
    font-size: 14px;
}
#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}
</style>
