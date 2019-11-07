<template>
    <div>
        <SidebarNav viewing="projects" :projects="projects" :stations="stations" @showStation="showStation" />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" ref="headerBar" />
        <div class="main-panel" v-show="!loading && isAuthenticated">
            <div class="view-user" v-if="!isEditing">
                <div id="user-name">{{ this.user.name }}</div>
                <div id="edit-user">
                    <img alt="Edit user" src="../assets/edit.png" v-on:click="editUser" />
                </div>
                <div class="user-element">{{ this.user.email }}</div>
                <div class="user-element">{{ this.user.bio }}</div>
            </div>

            <div id="user-form-container" v-if="isEditing">
                <div id="close-form-btn" v-on:click="closeForm">
                    <img alt="Close" src="../assets/close.png" />
                </div>
                <div class="input-label">Name:</div>
                <input v-model="user.name" placeholder="Name" class="text-input wide-text-input" />
                <div class="input-label">Email:</div>
                <input v-model="user.email" placeholder="Email" class="text-input wide-text-input" />
                <div class="input-label">Bio:</div>
                <input v-model="user.bio" placeholder="Bio" class="text-input wide-text-input" />
                <div id="public-checkbox-container">
                    <input type="checkbox" id="checkbox" v-model="publicProfile" />
                    <label for="checkbox">Make my profile public</label>
                    <img alt="Info" src="../assets/info.png" />
                </div>
                <div class="image-container">
                    <img
                        alt="User image"
                        :src="baseUrl + '/user/' + user.id + '/media'"
                        v-if="user.media_url && !previewImage"
                    />
                    <img :src="previewImage" class="uploading-image" v-if="!user.media_url || previewImage" />
                    <br />
                    {{
                        this.user.media_url ? "Update your profile image: " : "Add an image to your profile: "
                    }}
                    <input type="file" accept="image/gif, image/jpeg, image/png" @change="uploadImage" />
                </div>
                <button class="save-btn" v-on:click="submitUpdate">Update</button>
            </div>
        </div>
        <div id="loading" v-if="loading">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div v-if="failedAuth" class="no-auth-message">
            <p>
                Please
                <router-link :to="{ name: 'login' }" class="show-link">
                    log in
                </router-link>
                to view account.
            </p>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
import HeaderBar from "../components/HeaderBar";
import SidebarNav from "../components/SidebarNav";
import { API_HOST } from "../secrets";

export default {
    name: "UserView",
    components: {
        HeaderBar,
        SidebarNav
    },
    props: ["id"],
    data: () => {
        return {
            baseUrl: API_HOST,
            user: { username: "" },
            publicProfile: true,
            previewImage: null,
            projects: [],
            stations: [],
            acceptedImageTypes: ["jpg", "jpeg", "png", "gif"],
            isEditing: false,
            isAuthenticated: false,
            failedAuth: false,
            loading: false
        };
    },
    async beforeCreate() {
        this.api = new FKApi();
        this.api
            .getCurrentUser()
            .then(user => {
                this.user = user;
                this.isAuthenticated = true;
                this.api.getProjects().then(projects => {
                    if (projects && projects.projects.length > 0) {
                        this.projects = projects.projects;
                    }
                });
                this.api.getStations().then(s => {
                    this.stations = s.stations;
                });
            })
            .catch(() => {
                this.loading = false;
                this.failedAuth = true;
            });
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        editUser() {
            this.isEditing = true;
        },
        submitUpdate() {
            this.loading = true;
            if (this.sendingImage) {
                this.api.uploadUserImage({ type: this.imageType, image: this.sendingImage }).then(() => {
                    this.$refs.headerBar.refreshImage(this.previewImage);
                    this.api.updateUser(this.user).then(() => {
                        this.isEditing = false;
                        this.loading = false;
                    });
                });
            } else {
                this.api.updateUser(this.user).then(() => {
                    this.isEditing = false;
                    this.loading = false;
                });
            }
        },
        uploadImage(event) {
            this.previewImage = null;
            this.sendingImage = null;
            let valid = false;
            if (event.target.files.length > 0) {
                this.acceptedImageTypes.forEach(t => {
                    if (event.target.files[0].type.indexOf(t) > -1) {
                        valid = true;
                    }
                });
            }
            if (!valid) {
                return;
            }
            this.imageType = event.target.files[0].type;
            const image = event.target.files[0];
            this.sendingImage = image;
            const reader = new FileReader();
            reader.readAsDataURL(image);
            reader.onload = event => {
                this.previewImage = event.target.result;
            };
        },
        closeForm() {
            this.isEditing = false;
        }
    }
};
</script>

<style scoped>
.view-user {
    margin: 20px 40px;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.no-auth-message {
    float: left;
    font-size: 20px;
    margin: 40px 0 0 40px;
}
.show-link {
    text-decoration: underline;
}
#user-name {
    font-size: 24px;
    font-weight: bold;
    margin: 30px 15px 0 20px;
    display: inline-block;
}
#edit-user {
    display: inline-block;
    cursor: pointer;
}
.user-element {
    margin-left: 20px;
}

#user-form-container {
    width: 700px;
    padding: 0 15px 15px 15px;
    margin: 60px;
    border: 1px solid rgb(215, 220, 225);
}
.input-label {
    width: 60px;
    text-align: right;
    float: left;
    clear: both;
    margin: 0 12px 30px 0;
}
.wide-text-input {
    float: left;
    width: 600px;
    border: none;
    border-bottom: 2px solid rgb(235, 235, 235);
    font-size: 15px;
    padding-bottom: 4px;
}
#public-checkbox-container {
    float: left;
    margin: 0 0 0 30px;
    width: 98%;
}
#public-checkbox-container input {
    float: left;
    margin: 5px;
}
#public-checkbox-container label {
    float: left;
    margin: 2px 5px;
}
#public-checkbox-container img {
    float: left;
    margin: 2px 5px;
}
.image-container {
    width: 98%;
    margin: 15px 0 0 30px;
    float: left;
}

.save-btn {
    width: 300px;
    height: 50px;
    font-size: 18px;
    color: white;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    margin: 15px 0 20px 30px;
    cursor: pointer;
}

#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}
</style>
