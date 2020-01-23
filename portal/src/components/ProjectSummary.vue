<template>
    <div id="project-summary-container" v-if="viewingSummary">
        <div class="project-container" v-if="this.project">
            <div v-if="project.media_url" class="left tall custom-project-image-container">
                <img alt="Fieldkit Project image" :src="getImageUrl(project)" class="custom-project-image" />
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
                    <div class="section-heading" v-if="projectStations.length != 0">
                        FieldKit Stations ({{ this.projectStations.length }})
                    </div>
                    <div class="section-heading" v-else>FieldKit Stations</div>
                </div>
                <div class="section-content" v-if="viewingStations">
                    <div class="station-dropdown">
                        Add a station:
                        <select v-model="stationOption" v-on:change="stationSelected">
                            <option
                                v-for="station in userStations"
                                v-bind:value="station.id"
                                v-bind:key="station.id"
                            >
                                {{ station.name }}
                            </option>
                        </select>
                    </div>
                    <div class="stations-container">
                        <div v-for="station in projectStations" v-bind:key="station.id">
                            <div class="station-cell">
                                <div class="delete-link">
                                    <img
                                        alt="Info"
                                        src="../assets/Delete.png"
                                        :data-id="station.id"
                                        v-on:click="deleteStation"
                                    />
                                </div>
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
                    <div class="users-container">
                        <div class="user-row">
                            <div class="cell-heading">Members ({{ this.projectUsers.length }})</div>
                            <div class="cell-heading">Role</div>
                            <div class="cell-heading">Status</div>
                            <div class="cell"></div>
                        </div>
                        <div class="user-row" v-for="user in projectUsers" v-bind:key="user.id">
                            <div class="cell">
                                {{ user.name }} <br />
                                <span class="email">{{ user.email }}</span>
                            </div>
                            <div class="cell">{{ user.role }}</div>
                            <div class="cell">{{ user.status }}</div>
                            <div class="cell">
                                <img
                                    alt="Remove user"
                                    src="../assets/close.png"
                                    class="remove-btn"
                                    :data-user="user.id"
                                    v-on:click="removeUser"
                                />
                            </div>
                        </div>
                        <div class="user-row">
                            <div class="cell">
                                <input
                                    class="text-input"
                                    placeholder="New member email"
                                    keyboardType="email"
                                    autocorrect="false"
                                    autocapitalizationType="none"
                                    v-model="inviteEmail"
                                    @blur="checkEmail"
                                />
                                <span class="validation-error" id="no-email" v-if="noEmail">
                                    Email is a required field.
                                </span>
                                <span class="validation-error" id="email-not-valid" v-if="emailNotValid">
                                    Must be a valid email address.
                                </span>
                            </div>
                            <div class="cell"></div>
                            <div class="cell">
                                <button class="save-btn" v-on:click="sendInvite">Invite</button>
                            </div>
                            <div class="cell"></div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="spacer"></div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
import { API_HOST } from "../secrets";

export default {
    name: "ProjectSummary",
    data: () => {
        return {
            baseUrl: API_HOST,
            projectStations: [],
            viewingSummary: false,
            viewingStations: false,
            viewingActivity: false,
            viewingTeam: false,
            displayStartDate: "",
            displayEndDate: "",
            projectUsers: [],
            inviteEmail: "",
            noEmail: false,
            emailNotValid: false,
            stationOption: ""
        };
    },
    props: ["project", "userStations", "users"],
    watch: {
        project() {
            if (this.project) {
                this.reset();
            }
        },
        users() {
            if (this.users) {
                this.projectUsers = this.users.map(u => {
                    // temp: put in role and status
                    u.role = "Admin";
                    u.status = "Active";
                    return u;
                });
            }
        }
    },
    methods: {
        viewSummary() {
            this.viewingSummary = true;
        },
        reset() {
            this.stationOption = "";
            this.projectStations = [];
            this.fetchStations();
            this.updateDisplayDates();
            this.viewingStations = false;
            this.viewingActivity = false;
            this.viewingTeam = false;
            this.inviteEmail = "";
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
        fetchStations() {
            const api = new FKApi();
            api.getStationsByProject(this.project.id).then(result => {
                this.projectStations = result.stations;
            });
        },
        stationSelected() {
            const api = new FKApi();
            const params = {
                projectId: this.project.id,
                stationId: this.stationOption
            };
            api.addStationToProject(params).then(() => {
                this.fetchStations();
            });
        },
        deleteStation(event) {
            const stationId = event.target.getAttribute("data-id");
            if (window.confirm("Are you sure you want to remove this station?")) {
                const api = new FKApi();
                const params = {
                    projectId: this.project.id,
                    stationId: stationId
                };
                api.removeStationFromProject(params).then(() => {
                    this.fetchStations();
                });
            }
        },
        checkEmail() {
            this.noEmail = false;
            this.emailNotValid = false;
            this.noEmail = !this.inviteEmail || this.inviteEmail.length == 0;
            if (this.noEmail) {
                return false;
            }
            // eslint-disable-next-line
            let emailPattern = /^([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$/;
            this.emailNotValid = !emailPattern.test(this.inviteEmail);
            return !this.emailNotValid;
        },
        sendInvite() {
            let valid = this.checkEmail();
            if (valid) {
                this.$emit("inviteUser", { email: this.inviteEmail, projectId: this.project.id });
                this.projectUsers.push({
                    id: "pending-" + Date.now(),
                    name: "New member",
                    email: this.inviteEmail,
                    role: "",
                    status: "Pending"
                });
                this.inviteEmail = "";
            }
        },
        removeUser(event) {
            let id = event.target.getAttribute("data-user");
            if (confirm("Are you sure you want to remove this team member?")) {
                const index = this.projectUsers.findIndex(u => {
                    return u.id == id;
                });
                // remove pending-, if newly added member
                if (id.indexOf("pending-") == 0) {
                    id = id.split("pending-")[1];
                    // and make in range for type integer
                    id = parseInt(id / 1000000);
                }
                const params = {
                    projectId: this.project.id,
                    userId: id,
                    email: this.projectUsers[index].email
                };
                this.$emit("removeUser", params);
                // also remove from projectUsers
                if (index > -1) {
                    this.projectUsers.splice(index, 1);
                }
            } else {
                // canceled
            }
        },
        getImageUrl(project) {
            return this.baseUrl + "/projects/" + project.id + "/media/?t=" + Date.now();
        },
        updateDisplayDates() {
            this.displayStartDate = "";
            this.displayEndDate = "";
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
    padding: 0 15px 15px 15px;
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
    margin: 20px 0 10px 20px;
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
.station-dropdown {
    float: left;
    clear: both;
    margin-bottom: 20px;
}
.station-dropdown select {
    font-size: 16px;
    border: 1px solid lightgray;
    border-radius: 4px;
    padding: 2px 4px;
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
.delete-link {
    float: right;
    /*opacity: 0;*/
}
.delete-link:hover {
    opacity: 1;
}
.user-row {
    display: grid;
    grid-template-columns: 280px 200px 200px 40px;
    border-bottom: 1px solid rgb(215, 220, 225);
    padding: 10px 0;
}
.cell-heading {
    font-weight: bold;
}
.email {
    font-size: 14px;
}
.text-input {
    border: none;
    border-radius: 5px;
    background: rgb(240, 240, 240);
    font-size: 15px;
    padding: 4px 0 4px 8px;
}
.save-btn {
    padding: 2px 8px;
    font-size: 15px;
    color: white;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}
.validation-error {
    color: #c42c44;
    display: block;
}
.remove-btn {
    cursor: pointer;
}
#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}
</style>
