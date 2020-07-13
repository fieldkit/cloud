<template>
    <StandardLayout>
        <div class="main-panel" v-show="!isBusy && isAuthenticated" v-if="user">
            <div id="user-form-container">
                <div id="account-heading">My Account</div>
                <div class="notification" v-if="savedNotification">
                    Profile saved.
                </div>
                <div class="image-container">
                    <div id="image-heading">Profile picture</div>
                    <img src="../assets/Profile_Image.png" v-if="!user.photo && !previewImage" />
                    <img alt="User image" :src="$config.baseUrl + user.photo.url" v-if="user.photo && !previewImage" class="user-image" />
                    <img :src="previewImage" class="user-image" v-if="!user.photo || previewImage" />
                    <br />
                    <input type="file" accept="image/gif, image/jpeg, image/png" @change="uploadImage" />
                </div>
                <div class="input-container">
                    <input v-model="user.name" class="inputText" type="text" required="" />
                    <span class="floating-label">Name</span>
                </div>
                <div class="input-container">
                    <input v-model="user.email" type="text" class="inputText" required="" />
                    <span class="floating-label">Email</span>
                </div>
                <div class="input-container">
                    <input v-model="user.bio" type="text" class="inputText" required="" />
                    <span class="floating-label">Short Description</span>
                </div>
                <div id="public-checkbox-container">
                    <input type="checkbox" id="checkbox" v-model="publicProfile" />
                    <label for="checkbox">Make my profile public</label>
                </div>
                <button class="save-btn" v-on:click="submitUpdate">Update</button>

                <div class="password-change">
                    <div class="inner-password-change">
                        <div class="password-change-heading">Change password</div>
                        <div class="input-container">
                            <input v-model="oldPassword" secure="true" type="password" class="inputText" required="" />
                            <span class="floating-label">Current password</span>
                        </div>
                        <div class="input-container">
                            <input
                                v-model="newPassword"
                                secure="true"
                                type="password"
                                class="inputText"
                                required=""
                                @blur="checkPassword"
                            />
                            <span class="floating-label">New password</span>
                        </div>
                        <span class="validation-error" id="no-password" v-if="noPassword">Password is a required field.</span>
                        <span class="validation-error" id="password-too-short" v-if="passwordTooShort">
                            Password must be at least 10 characters.
                        </span>
                        <div class="input-container">
                            <input
                                v-model="confirmPassword"
                                secure="true"
                                type="password"
                                class="inputText"
                                required=""
                                @blur="checkConfirmPassword"
                            />
                            <span class="floating-label">Confirm new password</span>
                        </div>
                        <span class="validation-error" v-if="passwordsNotMatch">
                            Passwords do not match.
                        </span>
                    </div>
                    <button class="save-btn" v-on:click="submitPasswordChange">Change password</button>

                    <div class="forgot-link-container">
                        <div class="forgot-link" v-if="!showReset && !resetSent" v-on:click="showResetPassword">Forgot your password?</div>
                        <div class="input-container" v-if="showReset">
                            <div class="reset-instructions">
                                Enter your email address and password reset instructions will be sent to you.
                            </div>
                            <input v-model="resetEmail" type="text" class="inputText" required="" />
                            <span class="floating-label">Email</span>
                        </div>
                        <button class="save-btn" v-if="showReset" v-on:click="sendResetEmail">Submit</button>
                        <div class="reset-sent" v-if="resetSent">Password reset email sent!</div>
                    </div>
                </div>
            </div>
        </div>
        <div id="loading" v-if="isBusy">
            <img alt="" src="../assets/progress.gif" />
        </div>
        <div v-if="!isAuthenticated" class="no-user-message">
            <p>
                Please
                <router-link :to="{ name: 'login', query: { redirect: { name: 'user' } } }" class="show-link">
                    log in
                </router-link>
                to view account.
            </p>
        </div>
    </StandardLayout>
</template>

<script>
import Promise from "bluebird";
import StandardLayout from "./StandardLayout";
import FKApi from "../api/api";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default {
    name: "UserView",
    components: {
        StandardLayout,
    },
    props: {
        id: {
            type: Number,
            required: false,
        },
    },
    data: () => {
        return {
            savedNotification: false,
            publicProfile: true,
            previewImage: "",
            acceptedImageTypes: ["jpg", "jpeg", "png", "gif"],
            oldPassword: "",
            newPassword: "",
            noPassword: false,
            passwordTooShort: false,
            confirmPassword: "",
            passwordsNotMatch: false,
            showReset: false,
            resetSent: false,
            resetEmail: "",
        };
    },
    async beforeCreate() {
        this.api = new FKApi();
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s) => s.user.user,
            stations: (s) => s.stations.stations.user,
            userProjects: (s) => s.stations.projects.user,
        }),
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        showStation(station) {
            this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        checkPassword() {
            this.noPassword = false;
            this.passwordTooShort = false;
            this.noPassword = !this.newPassword || this.newPassword.length == 0;
            if (this.noPassword) {
                return;
            }
            this.passwordTooShort = this.newPassword.length < 10;
        },
        checkConfirmPassword() {
            this.passwordsNotMatch = this.newPassword != this.confirmPassword;
            return this.passwordsNotMatch;
        },
        submitPasswordChange() {
            if (this.oldPassword && this.checkConfirmPassword) {
                this.loading = true;
                const data = {
                    userId: this.user.id,
                    oldPassword: this.oldPassword,
                    newPassword: this.newPassword,
                };
                return this.api.updatePassword(data).then(() => {
                    // TODO: indicate success
                    this.loading = false;
                });
            }
        },
        showResetPassword() {
            this.showReset = true;
        },
        sendResetEmail() {
            return this.api.sendResetPasswordEmail(this.resetEmail).then(() => {
                this.showReset = false;
                this.resetSent = true;
            });
        },
        submitUpdate() {
            this.loading = true;
            if (this.sendingImage) {
                return this.$store.dispatch(ActionTypes.UPLOAD_USER_PHOTO, { type: this.imageType, image: this.sendingImage }).then(() => {
                    return this.$store.dispatch(ActionTypes.UPDATE_USER_PROFILE, { user: this.user }).then(() => {
                        this.loading = false;
                        this.savedNotification = true;
                        // NOTE Considering a flash library for this.
                        return Promise.delay(2000).then(() => {
                            this.savedNotification = false;
                        });
                    });
                });
            } else {
                return this.$store.dispatch(ActionTypes.UPDATE_USER_PROFILE, { user: this.user }).then(() => {
                    this.loading = false;
                    this.savedNotification = true;
                    // NOTE Considering a flash library for this.
                    return Promise.delay(2000).then(() => {
                        this.savedNotification = false;
                    });
                });
            }
        },
        uploadImage(event) {
            this.previewImage = null;
            this.sendingImage = null;
            let valid = false;
            if (event.target.files.length > 0) {
                this.acceptedImageTypes.forEach((t) => {
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
            reader.onload = (event) => {
                this.previewImage = event.target.result;
            };
        },
    },
};
</script>

<style scoped>
#account-heading {
    font-weight: bold;
    font-size: 24px;
    margin: 15px 0 0 15px;
}
#close-form-btn {
    float: right;
    margin-top: 15px;
    cursor: pointer;
}
.image-container {
    width: 98%;
    margin: 15px 0 30px 15px;
    float: left;
}
.user-image {
    max-width: 400px;
    max-height: 400px;
}
.view-user {
    margin: 20px 40px;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.no-user-message {
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
    float: left;
    width: 700px;
    padding: 0 15px 15px 15px;
}

.input-container {
    float: left;
    margin: 10px 0 0 15px;
    width: 99%;
}

#public-checkbox-container {
    float: left;
    margin: 0 0 15px 15px;
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

.password-change-heading {
    font-size: 16px;
    font-weight: 500;
    margin: 0 0 20px 15px;
}
.password-change {
    margin: 40px 0;
}
.inner-password-change {
    float: left;
}

.save-btn {
    width: 300px;
    height: 50px;
    color: white;
    font-size: 18px;
    font-weight: bold;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    margin: 10px 0 20px 15px;
    cursor: pointer;
}

.validation-error {
    float: left;
    color: #c42c44;
    display: block;
    font-size: 14px;
    margin: -20px 0 5px 15px;
}
.forgot-link {
    margin: 15px 0 0 15px;
    text-decoration: underline;
    cursor: pointer;
}
.reset-instructions {
    margin: 15px 0 18px 0;
}
.reset-sent {
    margin: 15px 0 0 15px;
    font-size: 18px;
    color: #3f8530;
}
.notification {
    margin: 20px;
    padding: 20px;
    background-color: #d4edda;
    border: 2px;
    border-radius: 4px;
}
</style>
