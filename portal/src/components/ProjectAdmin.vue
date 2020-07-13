<template>
    <div class="project-admin project-container" v-if="project">
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
                    <img alt="Fieldkit Project" v-if="project.mediaUrl" :src="getImageUrl(project)" class="project-image" />
                    <img alt="Default Fieldkit Project" v-else src="../assets/fieldkit_project.png" class="project-image" />
                </div>

                <div class="follow-btn"></div>
            </div>

            <div class="right">
                <div class="details-heading">
                    <div class="title">Project Details</div>
                    <div v-on:click="editProject" class="link">Edit Project</div>
                </div>
                <div class="details-top">
                    <div class="details-left">
                        <div class="project-detail" v-if="project.goal">Project Goal: {{ project.goal }}</div>
                        <div class="project-detail">{{ project.description }}</div>
                    </div>
                    <div class="details-right">
                        <div class="time-container" v-if="project.startTime">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            <template>Started: {{ project.startTime | prettyDate }}</template>
                        </div>
                        <div class="duration-container" v-if="displayProject.duration">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            <template>{{ displayProject.duration | prettyDuration }}</template>
                        </div>
                        <div class="location-container" v-if="project.location">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
                            <template>{{ project.location }}</template>
                        </div>
                        <div class="location-container" v-if="displayProject.places.native">
                            <img alt="Location" src="../assets/icon-location.png" class="icon" />
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

        <div class="row-section project-stations">
            <ProjectStations
                :project="project"
                :admin="false"
                :mapContainerSize="mapContainerSize"
                :listSize="listSize"
                :userStations="userStations"
            />
        </div>

        <div class="row-section data-readings">
            <div class="project-data">
                <ProjectDataFiles :projectStations="displayProject.stations" />
            </div>
            <div class="project-readings">
                <StationsReadings :project="project" />
            </div>
        </div>

        <div class="row-section manage-team-container">
            <div class="section-heading">Manage Team</div>
            <div class="users-container">
                <div class="user-row">
                    <div class="cell-heading">Members ({{ displayProject.users.length }})</div>
                    <div class="cell-heading">Role</div>
                    <div class="cell-heading"></div>
                    <div class="cell"></div>
                </div>
                <div class="user-row" v-for="projectUser in displayProject.users" v-bind:key="projectUser.user.email">
                    <div class="cell">
                        <UserPhoto :user="projectUser.user" />
                        <div class="invite-name">
                            <div v-if="projectUser.user.name != projectUser.user.email">{{ projectUser.user.name }}</div>
                            <div class="email">{{ projectUser.user.email }}</div>
                        </div>
                    </div>
                    <div class="cell">{{ projectUser.role }}</div>
                    <div class="cell invite-status">Invite {{ projectUser.membership.toLowerCase() }}</div>
                    <div class="cell">
                        <img
                            alt="Remove user"
                            src="../assets/close-icon.png"
                            class="remove-btn"
                            :data-user="projectUser.user.id"
                            v-on:click="(ev) => removeUser(ev, projectUser)"
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
                            v-model="form.inviteEmail"
                        />
                        <div class="validation-errors" v-if="$v.form.inviteEmail.$error || form.inviteDuplicate">
                            <span class="validation-error" v-if="!$v.form.inviteEmail.required">
                                Email is a required field.
                            </span>
                            <span class="validation-error" v-if="!$v.form.inviteEmail.email">
                                Must be a valid email address.
                            </span>
                            <span class="validation-error" v-if="form.inviteDuplicate">
                                This user is already invited.
                            </span>
                        </div>
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

        <ProjectActivity :displayProject="displayProject" v-if="false" />
    </div>
</template>

<script>
import { required, email } from "vuelidate/lib/validators";
import FKApi from "../api/api";
import * as ActionTypes from "@/store/actions";
import * as utils from "../utilities";
import ProjectStations from "../components/ProjectStations";
import ProjectActivity from "../components/ProjectActivity";
import ProjectDataFiles from "../components/ProjectDataFiles";
import StationsReadings from "../components/StationsReadings";
import CommonComponents from "@/views/shared";

export default {
    name: "ProjectAdmin",
    components: {
        ...CommonComponents,
        ProjectStations,
        ProjectActivity,
        ProjectDataFiles,
        StationsReadings,
    },
    props: {
        displayProject: {
            required: true,
        },
        userStations: {
            required: true,
        },
    },
    data: () => {
        return {
            form: {
                inviteEmail: "",
                inviteDuplicate: false,
            },
            viewingActivityFeed: false,
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
    validations: {
        form: {
            inviteEmail: {
                required,
                email,
            },
        },
    },
    computed: {
        project() {
            return this.displayProject.project;
        },
        projectModules() {
            return this.displayProject.modules.map((m) => {
                return {
                    name: m.name,
                    url: this.getModuleImg(m),
                };
            });
        },
    },
    watch: {
        project() {
            this.form.inviteEmail = "";
        },
    },
    methods: {
        getProjectUserImage(projectUser) {
            if (projectUser.user.photo) {
                return this.$config.baseUrl + "/" + projectUser.user.photo.url;
            }
            return null;
        },
        editProject() {
            return this.$router.push({ name: "editProject", params: { id: this.project.id } });
        },
        addUpdate() {
            return this.$router.push({ name: "addProjectUpdate", params: { project: this.project } });
        },
        viewProfile() {
            return this.$emit("viewProfile");
        },
        closeActivityFeed() {
            this.viewingActivityFeed = false;
        },
        openActivityFeed() {
            this.viewingActivityFeed = true;
        },
        checkEmail() {
            this.$v.form.$touch();
            return !(this.$v.form.$pending || this.$v.form.$error);
        },
        sendInvite() {
            if (this.checkEmail()) {
                if (this.selectedRole == -1) {
                    this.selectedRole = 0;
                }
                const role = this.roleOptions.find((r) => {
                    return r.code == this.selectedRole;
                });
                const payload = {
                    projectId: this.project.id,
                    email: this.form.inviteEmail,
                    role: this.selectedRole,
                };
                return this.$store
                    .dispatch(ActionTypes.PROJECT_INVITE, payload)
                    .then(() => {
                        this.$v.form.$reset();
                        this.form.inviteEmail = "";
                        this.form.inviteDuplicate = false;
                    })
                    .catch(() => {
                        // TODO: Move this to vuelidate.
                        this.form.inviteDuplicate = true;
                    });
            }
        },
        removeUser(ev, projectUser) {
            const id = event.target.getAttribute("data-user");
            if (confirm("Are you sure you want to remove this team member?")) {
                const payload = {
                    projectId: this.project.id,
                    email: projectUser.user.email,
                };
                return this.$store.dispatch(ActionTypes.PROJECT_REMOVE, payload).then(() => {
                    this.$v.form.$reset();
                    this.form.inviteEmail = "";
                });
            }
        },
        getImageUrl(project) {
            return this.$config.baseUrl + "/" + project.photo;
        },
        getModuleImg(module) {
            return this.$loadAsset("modules-lg/" + utils.getModuleImg(module));
        },
    },
};
</script>

<style scoped>
.project-admin {
    display: flex;
    flex-direction: column;
    max-width: 1000px;
    padding-bottom: 100px;
}
.header {
    display: flex;
    flex-direction: row;
}
.header .left {
    margin-right: auto;
    display: flex;
    flex-direction: column;
}
.header .right {
    margin-left: auto;
}
.header .project-name {
    font-size: 24px;
    font-weight: bold;
    margin: 0 15px 0 0;
    display: inline-block;
}
.header .project-dashboard {
    font-size: 20px;
    font-weight: bold;
    margin: 0 15px 0 0;
    display: inline-block;
}

.details {
    display: flex;
    flex-direction: row;
}
.details > .left {
    flex: 1;
    border: 2px solid #d8dce0;
    border-radius: 2px;
    margin-right: 20px;
    background-color: white;
    padding: 20px;
}
.details > .right {
    flex: 2;
    border: 2px solid #d8dce0;
    border-radius: 2px;
    background-color: white;
    padding: 20px;
    display: flex;
    flex-direction: column;
}
.project-stations {
}

.row-section {
}

.details .details-heading {
    display: flex;
    flex-direction: row;
}
.details .details-heading .title {
    font-weight: bold;
    padding-bottom: 20px;
}
.details .details-heading .link {
    margin-left: auto;
    font-weight: bold;
    padding-bottom: 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
}
.details .details-top {
    display: flex;
    flex-direction: row;
    padding-bottom: 20px;
}
.details .details-left {
    padding-right: 20px;
    flex-grow: 1;
}
.details .details-right {
    flex-grow: 1;
}
.details .details-bottom {
    border-top: 1px solid #d8dce0;
    padding-top: 20px;
    display: flex;
    flex-direction: row;
}
.details-bottom .details-team {
    flex: 1;
}
.details-bottom .details-modules {
    flex: 1;
}
.details-bottom .title {
    font-weight: bold;
}
.row-section.data-readings {
    margin-top: 20px;
    display: flex;
    flex-direction: row;
}
.project-data {
    margin-right: 20px;
}
.project-data,
.project-readings {
    border: 2px solid #d8dce0;
    border-radius: 2px;
    background-color: white;
    padding: 20px;
    display: flex;
    flex-direction: column;
}
.data-readings .project-data {
    flex: 2;
}
.data-readings .project-readings {
    flex: 1;
}
.manage-team-container {
    margin-top: 20px;
    display: flex;
    flex-direction: column;
    border: 2px solid #d8dce0;
    background: white;
    padding: 20px;
}
.manage-team-container .user-row {
}

.user-row .cell {
    text-align: left;
}

.manage-team-container .section-heading {
    font-weight: bold;
    margin-top: 10px;
    margin-bottom: 20px;
}
/*
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
#activity-feed-container .heading {
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
.activity-btn {
    width: 110px;
    height: 25px;
    float: right;
    padding: 10px 10px 6px 10px;
    margin: -48px 58px 0 0;
    border: 1px solid #d8dce0;
    border-radius: 3px;
    text-align: center;
    cursor: pointer;
    font-size: 14px;
    font-weight: bold;
}
.activity-btn img {
    vertical-align: sub;
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
}
.goal-and-description {
    float: left;
    clear: both;
    width: 320px;
}
.project-detail {
    overflow-wrap: break-word;
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
    width: 32px;
    margin: 2px 5px;
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
*/
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
    align-items: center;
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
.validation-error {
    color: #c42c44;
    display: block;
}
.remove-btn {
    margin: 12px 0 0 0;
    float: right;
    cursor: pointer;
}
.invite-name {
    display: inline-block;
    vertical-align: top;
    margin-top: 20px;
}
.role-dropdown-container select {
    font-size: 13px;
    color: #6a6d71;
    border: none;
    padding: 2px 10px 2px 0;
    margin-left: -4px;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='17.538' height='9.865' viewBox='0 0 4.64 2.61'%3E%3Cpath d='M.493.087a.28.28 0 00-.406 0 .28.28 0 000 .406l2.03 2.03a.28.28 0 00.406 0l2.03-2.03a.28.28 0 000-.406.28.28 0 00-.406 0L2.32 1.914z' fill='%236a6d71'/%3E%3C/svg%3E%0A");
    background-repeat: no-repeat, repeat;
    background-position: right 0.7em top 50%, 0 0;
    background-size: 0.75em auto, 100%;
    box-sizing: border-box;
    -moz-appearance: none;
    -webkit-appearance: none;
    appearance: none;
}
</style>
