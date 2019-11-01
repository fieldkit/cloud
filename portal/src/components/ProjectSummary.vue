<template>
    <div id="project-summary-container" v-if="viewingSummary">
        <div class="project-container" v-if="this.project">
            <div v-if="project.media_url" class="left tall custom-project-image-container">
                <img
                    alt="Fieldkit Project image"
                    :src="baseUrl + '/projects/' + project.id + '/media'"
                    class="custom-project-image"
                />
            </div>
            <div v-else class="left tall">
                <img alt="Fieldkit Project image" src="../assets/fieldkit_project.png" />
            </div>
            <div class="left">
                <div id="project-name">{{ this.project.name }}</div>
                <div id="edit-project">
                    <img alt="Edit project" src="../assets/edit.png" v-on:click="editProject" />
                </div>
                <div class="project-element">{{ this.project.description }}</div>
                <div class="project-element">{{ this.project.location }}</div>
                <div class="left project-element">
                    <div class="left" v-if="displayStartDate">
                        <span class="date-label">Start Date</span> <br />
                        {{ this.displayStartDate }}
                    </div>
                    <div class="left end-date" v-if="displayEndDate">
                        <span class="date-label">End Date</span> <br />
                        {{ this.displayEndDate }}
                    </div>
                </div>
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
import { API_HOST } from "../secrets";

export default {
    name: "ProjectSummary",
    data: () => {
        return {
            baseUrl: API_HOST,
            viewingSummary: false,
            viewingStations: false,
            viewingActivity: false,
            viewingTeam: false,
            displayStartDate: "",
            displayEndDate: ""
        };
    },
    props: ["project", "stations", "user"],
    watch: {
        project() {
            if (this.project) {
                this.updateDisplayDates();
            }
        }
    },
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
        editProject() {
            this.$router.push({ name: "editProject", params: { id: this.project.id } });
        },
        updateDisplayDates() {
            if (this.project.start_time) {
                let d = new Date(this.project.start_time);
                this.displayStartDate = d.toLocaleDateString("en-US");
            }
            if (this.project.end_time) {
                let d = new Date(this.project.end_time);
                this.displayEndDate = d.toLocaleDateString("en-US");
            }
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
.tall {
    height: 100%;
}
.left {
    float: left;
}
.right {
    float: right;
}
.project-element {
    max-width: 425px;
    font-size: 14px;
    margin: 10px 20px 0 20px;
}
.date-label {
    font-size: 12px;
    color: rgb(134, 134, 134);
}
.end-date {
    margin-left: 20px;
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
    margin: 0 15px 0 20px;
    display: inline-block;
}
#edit-project {
    display: inline-block;
    cursor: pointer;
}
.custom-project-image-container {
    text-align: center;
}
.custom-project-image {
    max-width: 275px;
    max-height: 135px;
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
