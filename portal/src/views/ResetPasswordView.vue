<template>
    <div id="login-container">
        <img
            v-bind:style="{
                width: '210px',
                marginTop: '150px',
                marginBottom: '86px',
            }"
            alt="Fieldkit Logo"
            src="../assets/FieldKit_Logo_White.png"
        />
        <br />
        <div id="user-form-container">
            <div class="password-change" v-if="!resetSuccess && !failedReset">
                <div class="inner-password-change">
                    <div class="password-change-heading">Reset password</div>

                    <div class="outer-input-container">
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
                    </div>

                    <div class="outer-input-container">
                        <div class="input-container middle-container">
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
                </div>
                <button class="save-btn" v-on:click="submitPasswordReset">Reset password</button>
            </div>
            <div v-if="resetSuccess">
                <img alt="Success" src="../assets/Icon_Success.png" width="57px" />
                <p class="success">Password Reset</p>

                <router-link :to="{ name: 'login' }" class="create-link">
                    Go to Log In
                </router-link>
            </div>
            <div v-if="failedReset">
                <img alt="Unsuccessful" src="../assets/Icon_Warning_error.png" width="57px" />
                <p class="error">Password Not Reset</p>
                <div class="notification-text">Unfortunately we were unable to reset your password.</div>
                <p>
                    Please
                    <a href="https://www.fieldkit.org/contact/" class="contact-link">contact us</a>
                    if you would like assistance.
                </p>
            </div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";

export default {
    name: "ResetPasswordView",
    components: {},
    props: [],
    data: () => {
        return {
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
    },
    methods: {
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
                    .then(() => {
                        this.resetSuccess = true;
                    })
                    .catch(() => {
                        this.failedReset = true;
                    });
            }
        },
    },
};
</script>

<style scoped>
#login-container {
    width: 100%;
    min-height: 100%;
    background-image: linear-gradient(#52b5e4, #1b80c9);
}
#user-form-container {
    width: 460px;
    background-color: white;
    display: inline-block;
    text-align: center;
    padding-bottom: 60px;
    padding-top: 78px;
}
.outer-input-container {
    width: 300px;
    margin: auto;
}
.input-container {
    margin: auto;
    width: 100%;
    text-align: left;
}
.middle-container {
    margin-top: 22px;
}
input {
    background: none;
    border: 0;
    border-bottom: 2px solid #d8dce0;
    outline: 0;
    font-size: 18px;
    padding-bottom: 2px;
}
.floating-label {
    font-size: 16px;
    color: #6a6d71;
}

.password-change-heading {
    font-size: 24px;
    font-weight: 500;
    margin-bottom: 50px;
}

.save-btn {
    width: 300px;
    height: 45px;
    color: white;
    font-size: 18px;
    font-weight: 600;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    margin: 40px 0 20px 0;
    cursor: pointer;
}

.validation-error {
    float: left;
    color: #c42c44;
    display: block;
    font-size: 14px;
    margin: -20px 0 0 0;
}
.success,
.error {
    margin: 25px 0 25px 0;
    font-size: 24px;
}
.notification-text {
    width: 300px;
    margin: 0 auto 40px auto;
}
.contact-link {
    cursor: pointer;
    text-decoration: underline;
}
</style>
