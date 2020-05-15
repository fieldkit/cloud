<template>
    <div id="project-summary-container">
        <div class="project-container" v-if="project">
            <div class="left">
                <div id="project-name">{{ project.name }}</div>
                <div class="dashboard-heading">Project Dashboard</div>
                <div class="project-image-actions-container">
                    <img alt="Fieldkit Project" v-if="project.media_url" :src="getImageUrl(project)" class="project-image" />
                    <img alt="Default Fieldkit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />

                    <div class="actions-icon-container">
                        <div class="action left">
                            <img src="../assets/update.png" />
                            <div class="label">Update</div>
                        </div>
                        <div class="action" v-on:click="viewProfile">
                            <img src="../assets/profile.png" />
                            <div class="label">View Profile</div>
                        </div>
                        <div class="action right">
                            <img src="../assets/share.png" />
                            <div class="label">Share</div>
                        </div>
                    </div>
                    <div class="stat follows">
                        <img alt="Follows" src="../assets/heart.png" class="follow-icon" />
                        <span>{{ numFollowers }} Follows</span>
                    </div>
                </div>
                <div class="project-details-container">
                    <div class="section-heading">Project Details</div>
                    <div v-on:click="editProject" class="edit-link">Edit Project</div>
                    <div class="goal-and-description">
                        <div class="project-detail" v-if="project.goal">Project Goal: {{ project.goal }}</div>
                        <div class="project-detail">{{ project.description }}</div>
                    </div>
                    <div class="time-and-location">
                        <div class="time" v-if="displayStartDate">
                            <img alt="Calendar" src="../assets/icon-calendar.png" class="icon" />
                            Started {{ displayStartDate }}
                        </div>
                        <div class="time" v-if="displayRunTime">
                            <img alt="Time" src="../assets/icon-time.png" class="icon" />
                            {{ displayRunTime }}
                        </div>
                        <div class="location" v-if="project.location">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            {{ project.location }}
                        </div>
                    </div>
                    <div class="space"></div>
                    <div class="team-icons">
                        <div class="icon-section-label">Team</div>
                        <img
                            v-for="user in projectUsers"
                            v-bind:key="user.user.id"
                            alt="User image"
                            :src="user.userImage"
                            class="user-icon"
                        />
                    </div>
                    <div class="module-icons">
                        <div class="icon-section-label">Modules</div>
                        <img v-for="module in modules" v-bind:key="module" alt="Module icon" class="module-icon" :src="module" />
                    </div>
                </div>
            </div>

            <ProjectStations
                :project="project"
                :mapContainerSize="mapContainerSize"
                :listSize="listSize"
                :userStations="userStations"
                @loaded="setModules"
            />

            <ProjectDataFiles />

            <StationsReadings :project="project" />

            <div class="manage-team-container">
                <div class="section-heading">Manage Team</div>
                <div class="users-container">
                    <div class="user-row">
                        <div class="cell-heading">Members ({{ projectUsers.length }})</div>
                        <div class="cell-heading">Role</div>
                        <div class="cell-heading"></div>
                        <div class="cell"></div>
                    </div>
                    <div class="user-row" v-for="user in projectUsers" v-bind:key="user.user.id">
                        <div class="cell">
                            <img alt="User image" :src="user.userImage" class="user-icon" />
                            {{ user.user.name }}
                            <br />
                            <span class="email">{{ user.user.email }}</span>
                        </div>
                        <div class="cell">{{ user.role }}</div>
                        <div class="cell invite-status">Invite {{ user.membership.toLowerCase() }}</div>
                        <div class="cell">
                            <img
                                alt="Remove user"
                                src="../assets/close-icon.png"
                                class="remove-btn"
                                :data-user="user.user.id"
                                v-on:click="removeUser"
                            />
                            <!-- <img
                                alt="Edit user"
                                src="../assets/edit.png"
                                class="edit-btn"
                                :data-user="user.user.id"
                                v-on:click="editUser"
                            /> -->
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
                        <div class="cell role-dropdown-container">
                            <select v-model="selectedRole">
                                <option v-for="role in roleOptions" v-bind:value="role.code" v-bind:key="role.code + role.name">
                                    {{ role.name }}
                                </option>
                            </select>
                        </div>
                        <div class="cell">
                            <button class="invite-btn" v-on:click="sendInvite">Invite</button>
                        </div>
                        <div class="cell"></div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import * as utils from "../utilities";
import FKApi from "../api/api";
import { API_HOST } from "../secrets";
import ProjectStations from "../components/ProjectStations";
import ProjectDataFiles from "../components/ProjectDataFiles";
import StationsReadings from "../components/StationsReadings";

export default {
    name: "ProjectAdmin",
    components: {
        ProjectStations,
        ProjectDataFiles,
        StationsReadings,
    },
    data: () => {
        return {
            baseUrl: API_HOST,
            displayStartDate: "",
            displayRunTime: "",
            projectUsers: [],
            inviteEmail: "",
            newUserImage: "",
            noEmail: false,
            emailNotValid: false,
            modules: [],
            numFollowers: 1,
            mapContainerSize: {
                width: "677px",
                height: "332px",
                outerWidth: "1022px",
            },
            listSize: {
                width: "345px",
                height: "332px",
                boxWidth: "274px",
            },
            selectedRole: -1,
            roleOptions: [
                {
                    code: -1,
                    name: "Select Role",
                },
                {
                    code: 0,
                    name: "Member",
                },
                {
                    code: 1,
                    name: "Administrator",
                },
            ],
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
                this.projectUsers = this.users;
            }
        },
    },
    async beforeCreate() {
        this.api = new FKApi();
    },
    mounted() {
        let imgPath = require.context("../assets/", false, /\.png$/);
        let img = "new_user.png";
        this.newUserImage = imgPath("./" + img);
    },
    methods: {
        reset() {
            this.fetchFollowers();
            this.updateDisplayDates();
            this.inviteEmail = "";
        },
        fetchFollowers() {
            this.api.getProjectFollows(this.project.id).then(result => {
                this.numFollowers = result.followers.length + 1;
            });
        },
        editProject() {
            this.$router.push({ name: "editProject", params: { id: this.project.id } });
        },
        viewProfile() {
            this.$emit("viewProfile");
        },
        setModules(modules) {
            this.modules = modules;
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
                if (this.selectedRole == -1) {
                    this.selectedRole = 0;
                }
                let role = this.roleOptions.find(r => {
                    return r.code == this.selectedRole;
                });
                const params = { email: this.inviteEmail, projectId: this.project.id, role: this.selectedRole };
                this.api.sendInvite(params).then(() => {
                    this.projectUsers.push({
                        user: {
                            id: "pending-" + Date.now(),
                            name: this.inviteEmail,
                            email: this.inviteEmail,
                        },
                        userImage: this.newUserImage,
                        role: role.name,
                        membership: "Pending",
                    });
                    this.inviteEmail = "";
                });
            }
        },
        removeUser(event) {
            let id = event.target.getAttribute("data-user");
            if (confirm("Are you sure you want to remove this team member?")) {
                const index = this.projectUsers.findIndex(u => {
                    return u.user.id == id;
                });
                const params = {
                    projectId: this.project.id,
                    email: this.projectUsers[index].user.email,
                };
                this.api.removeUserFromProject(params).then(() => {
                    // also remove from projectUsers
                    if (index > -1) {
                        this.projectUsers.splice(index, 1);
                    }
                });
            } else {
                // canceled
            }
        },
        editUser(event) {
            let id = event.target.getAttribute("data-user");
            console.log("edit user", id);
        },
        getImageUrl(project) {
            return this.baseUrl + "/projects/" + project.id + "/media";
        },
        updateDisplayDates() {
            this.displayRunTime = "";
            this.displayStartDate = "";
            if (this.project.start_time) {
                let d = new Date(this.project.start_time);
                this.displayStartDate = d.toLocaleDateString("en-US");
                this.displayRunTime = utils.getRunTime(this.project);
            }
        },
    },
};
</script>

<style scoped>
#project-summary-container {
    width: 1080px;
    margin: 0 0 0 30px;
    background-color: #ffffff;
    z-index: 2;
}
#project-name {
    font-size: 24px;
    font-weight: bold;
    margin: 0 15px 0 0;
    display: inline-block;
}
.dashboard-heading {
    font-size: 18px;
}
.show-link {
    text-decoration: underline;
}
.project-container {
    font-size: 16px;
    font-weight: lighter;
    overflow: hidden;
}
.project-image-actions-container {
    width: 288px;
    height: 295px;
    float: left;
    margin: 17px 22px 0 0;
    padding: 25px;
    border: 1px solid #d8dce0;
    text-align: center;
}
.project-image {
    max-width: 288px;
    max-height: 139px;
}
.actions-icon-container {
    width: 288px;
    margin: 30px 0 0 0;
    border-bottom: 1px solid #d8dce0;
}
.actions-icon-container .label {
    font-size: 14px;
    font-weight: bold;
}
.action {
    display: inline-block;
    margin: 0 0 17px 0;
    cursor: pointer;
}
.action.left {
    float: left;
    margin-left: 6px;
}
.action.right {
    float: right;
    margin-right: 6px;
}
.stat {
    font-size: 18px;
    font-weight: 600;
    margin: 24px 0;
}
.stat img {
    margin-right: 12px;
}
.project-details-container {
    width: 610px;
    height: 295px;
    float: left;
    margin: 17px 22px 0 0;
    padding: 25px;
    border: 1px solid #d8dce0;
}
.section-heading {
    font-size: 20px;
    font-weight: 600;
    float: left;
    margin: 0 0 35px 0;
}
.edit-link {
    float: right;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
}
.goal-and-description {
    float: left;
    clear: both;
    width: 320px;
}
.project-detail {
    font-size: 16px;
    line-height: 24px;
}
.time-and-location {
    float: right;
    width: 230px;
    font-size: 14px;
}
.time,
.location {
    margin: 0 0 8px 0;
}
.space {
    width: 100%;
    float: left;
    margin: 30px 0 0 0;
    border-bottom: solid 1px #d8dce0;
}
.team-icons,
.module-icons {
    margin: 22px 0 0 0;
    width: 225px;
    float: left;
}
.module-icons {
    margin: 22px 0 0 50px;
    float: left;
}
.icon-section-label {
    font-size: 14px;
    font-weight: 600;
}
.user-icon,
.module-icon {
    width: 35px;
    margin: 0 5px;
}
.manage-team-container {
    width: 1022px;
    float: left;
    margin: 22px 0 0 0;
    border: 1px solid #d8dce0;
}
.manage-team-container {
    width: 972px;
    padding: 25px;
}
.users-container {
    width: 100%;
    float: left;
    clear: both;
}
.invite-btn {
    width: 80px;
    height: 28px;
    border-radius: 3px;
    border: 1px solid #cccdcf;
    font-size: 14px;
    font-weight: bold;
    background-color: #ffffff;
    cursor: pointer;
}
.user-row {
    display: grid;
    font-size: 13px;
    grid-template-columns: 322px 193px 389px 40px;
    border-bottom: 1px solid rgb(215, 220, 225);
    padding: 10px 0;
}
.invite-status {
    color: #0a67aa;
    font-weight: 600;
}
.cell-heading {
    font-size: 14px;
    font-weight: bold;
}
.users-container .user-icon {
    float: left;
}
.text-input {
    border: none;
    border-radius: 5px;
    font-size: 15px;
    padding: 4px 0 4px 8px;
}
.role-dropdown-container select {
    font-size: 16px;
    border: 1px solid lightgray;
    border-radius: 4px;
    padding: 2px 4px;
}
.validation-error {
    color: #c42c44;
    display: block;
}

.edit-btn,
.remove-btn {
    margin: 12px 0 0 0;
    float: right;
    cursor: pointer;
}
.edit-btn {
    margin-right: 12px;
}
</style>
