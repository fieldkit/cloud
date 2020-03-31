<template>
    <div>
        <SidebarNav viewing="projects" :projects="projects" :stations="stations" @showStation="showStation" />
        <HeaderBar :isAuthenticated="isAuthenticated" :user="user" ref="headerBar" />
        <div class="main-panel">
            <div id="user-form-container">
                <div class="password-change" v-if="!resetSuccess && !failedReset">
                    <div class="inner-password-change">
                        <div class="password-change-heading">Reset password</div>
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
                    <button class="save-btn" v-on:click="submitPasswordReset">Reset password</button>
                </div>
                <div class="success" v-if="resetSuccess">Password was reset!</div>
                <div class="reset-error" v-if="failedReset">Unfortunately we were unable to reset your password.</div>
            </div>
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
        SidebarNav,
    },
    props: [],
    data: () => {
        return {
            baseUrl: API_HOST,
            user: { username: "" },
            projects: [],
            stations: [],
            isAuthenticated: false,
            resetToken: "",
            newPassword: "",
            noPassword: false,
            passwordTooShort: false,
            confirmPassword: "",
            passwordsNotMatch: false,
            resetSuccess: false,
            failedReset: false,
        };
    },
    async beforeCreate() {
        this.api = new FKApi();
        this.api.getCurrentUser().then(user => {
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
        });
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
        },
        submitPasswordReset() {
            this.resetToken = this.$route.query.token;
            if (this.checkConfirmPassword) {
                const data = {
                    token: this.resetToken,
                    password: this.newPassword,
                };
                this.api
                    .resetPassword(data)
                    .then(result => {
                        this.resetSuccess = true;
                    })
                    .catch(e => {
                        this.failedReset = true;
                    });
            }
        },
    },
};
</script>

<style scoped>
#user-form-container {
    float: left;
    width: 700px;
    padding: 0 15px 15px 15px;
    margin: 60px;
    border: 1px solid rgb(215, 220, 225);
}

.input-container {
    float: left;
    margin: 10px 0 0 15px;
    width: 99%;
}
input:focus ~ .floating-label,
input:not(:focus):valid ~ .floating-label {
    top: -48px;
    font-size: 12px;
    opacity: 1;
}
input:invalid {
    box-shadow: none;
}
.inputText {
    color: #2c3e50;
    font-size: 14px;
    width: inherit;
    border: none;
    border-bottom: 2px solid rgb(235, 235, 235);
    font-size: 15px;
    padding-bottom: 4px;
}
.inputText:focus {
    border-bottom: 2px solid #52b5e4;
}
.floating-label {
    color: rgb(85, 85, 85);
    position: relative;
    top: -24px;
    pointer-events: none;
    transition: 0.2s ease all;
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
.success {
    margin: 15px 0 0 15px;
    font-size: 18px;
    color: #3f8530;
}
.reset-error {
    margin: 15px 0 0 15px;
    font-size: 18px;
    color: #c42c44;
}
</style>
