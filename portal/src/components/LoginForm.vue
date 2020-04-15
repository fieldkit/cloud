<template>
    <div>
        <div id="login-form" v-if="!accountCreated && !accountFailed && !showReset">
            <h1>{{ isLoggingIn ? "Log In to Your Account" : "Create Your Account" }}</h1>
            <div class="outer-input-container" v-if="!isLoggingIn">
                <div class="input-container">
                    <input
                        id="name-field"
                        keyboardType="name"
                        autocorrect="false"
                        autocapitalizationType="none"
                        class="inputText"
                        required=""
                        v-model="name"
                        @keyup.enter="$refs.email.focus"
                        @blur="checkName"
                    />
                    <span class="floating-label">Name</span>
                </div>
                <div class="validation-error" id="no-name" v-if="noName">Name is a required field.</div>
                <div class="validation-error" id="name-too-long" v-if="nameTooLong">
                    Name must be less than 256 letters.
                </div>
                <div class="validation-error" id="name-has-space" v-if="nameHasSpace">Name must not contain spaces.</div>
            </div>
            <div class="outer-input-container">
                <div :class="'input-container ' + (isLoggingIn ? '' : ' middle-container')">
                    <img v-if="isLoggingIn" alt="Email" src="../assets/Icon_Email_login.png" class="email-img" />
                    <input
                        id="email-field"
                        ref="email"
                        class="inputText"
                        required=""
                        keyboardType="email"
                        autocorrect="false"
                        autocapitalizationType="none"
                        v-model="email"
                        @keyup.enter="$refs.password.focus"
                        @blur="checkEmail"
                    />
                    <span class="floating-label">Email</span>
                </div>
                <div class="validation-error" id="no-email" v-if="noEmail">Email is a required field.</div>
                <div class="validation-error" id="email-not-valid" v-if="emailNotValid">
                    Must be a valid email address.
                </div>
            </div>
            <div class="outer-input-container">
                <div class="input-container middle-container">
                    <img v-if="isLoggingIn" alt="Password" src="../assets/Icon_Password_login.png" class="password-img" />
                    <input
                        name="password"
                        class="inputText"
                        required=""
                        secure="true"
                        ref="password"
                        type="password"
                        v-model="password"
                        @keyup.enter="isLoggingIn ? submit() : $refs.confirmPassword.focus()"
                        @blur="checkPassword"
                    />
                    <span class="floating-label">Password</span>
                </div>
                <div class="validation-error" id="no-password" v-if="noPassword">Password is a required field.</div>
                <div class="validation-error" id="password-too-short" v-if="passwordTooShort">
                    Password must be at least 10 characters.
                </div>
                <div class="reset-link" v-if="isLoggingIn && !passwordTooShort && !noPassword" v-on:click="showResetPassword">
                    Reset password
                </div>
            </div>
            <div class="outer-input-container" v-if="!isLoggingIn">
                <div class="input-container middle-container">
                    <input
                        name="confirmPassword"
                        class="inputText"
                        required=""
                        secure="true"
                        ref="confirmPassword"
                        type="password"
                        v-model="confirmPassword"
                        @keyup.enter="submit"
                        @blur="checkConfirmPassword"
                    />
                    <span class="floating-label">Confirm Password</span>
                </div>
                <div class="validation-error" id="passwords-not-match" v-if="passwordsNotMatch">
                    Your passwords do not match.
                </div>
                <div class="policy-terms-container">
                    By creating an account you agree to our
                    <span class="bold">Privacy Policy</span>
                    and
                    <span class="bold">Terms of Use.</span>
                </div>
            </div>
            <button v-on:click="submit" ref="submit" class="submit-btn">{{ isLoggingIn ? "Log In" : "Create Account" }}</button>
            <div class="create-link" v-on:click="toggleForm">
                {{ isLoggingIn ? "Create an Account" : "Back to Log In" }}
            </div>
        </div>

        <div id="reset-container" v-if="showReset">
            <div v-if="!resetSent && !resendingReset">
                <div class="reset-heading">Password Reset</div>
                <div class="outer-input-container">
                    <div class="input-container">
                        <div class="reset-instructions">
                            Enter your email address below, and password reset instructions will be sent.
                        </div>
                        <input v-model="resetEmail" type="text" class="inputText" required="" />
                        <span class="floating-label">Email</span>
                    </div>
                </div>
                <button :class="sendingReset ? 'disabled' : 'submit-btn'" v-on:click="sendResetEmail" :disabled="sendingReset">
                    Submit
                </button>
                <div class="create-link" v-on:click="returnToLogin">
                    Go back to Log In
                </div>
            </div>
            <div v-if="resetSent && !resendingReset">
                <img alt="Success" src="../assets/Icon_Success.png" width="57px" />
                <p class="reset-heading">Password Reset Email Sent</p>
                <div class="notification-text">Check your inbox for the email with a link to reset your password.</div>
                <button v-on:click="resendReset" class="resubmit-btn">Resend Email</button>
                <div class="create-link" v-on:click="returnToLogin">
                    Go back to Log In
                </div>
            </div>
            <div class="resending" v-if="!resetSent && resendingReset">
                <img alt="Resending" src="../assets/Icon_Syncing2.png" width="57px" />
                <p class="reset-heading">Resending</p>
            </div>
        </div>

        <div id="notification" v-if="accountCreated || accountFailed">
            <div v-if="accountCreated">
                <img alt="Success" src="../assets/Icon_Success.png" width="57px" />
                <p class="success">Account Created</p>
                <div class="notification-text">We sent you an account validation email.</div>
                <div class="create-link" v-on:click="returnToLogin">
                    Go back to Log In
                </div>
            </div>
            <div v-if="accountFailed" class="notification-container">
                <img alt="Unsuccessful" src="../assets/Icon_Warning_error.png" width="57px" />
                <p class="error">A Problem Occurred</p>
                <p>Unfortunately we were unable to create your account at this time.</p>
                <p class="notification-text">Contact us to get assistance.</p>
                <p><a href="https://www.fieldkit.org/contact/" target="_blank" class="contact-link">Contact us</a></p>
            </div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";

export default {
    name: "LoginForm",
    data: () => {
        return {
            email: "",
            password: "",
            name: "",
            emailNotValid: false,
            noEmail: false,
            noPassword: false,
            passwordTooShort: false,
            isLoggingIn: true,
            noName: false,
            nameTooLong: false,
            nameHasSpace: false,
            confirmPassword: "",
            passwordsNotMatch: false,
            accountCreated: false,
            accountFailed: false,
            showReset: false,
            resetSent: false,
            resetEmail: "",
            sendingReset: false,
            resendingReset: false,
        };
    },
    async beforeCreate() {
        this.api = new FKApi();
    },
    methods: {
        async login() {
            try {
                const auth = await this.api.login(this.email.toLowerCase(), this.password);
                const isAuthenticated = await this.api.authenticated();
                if (isAuthenticated) {
                    this.userToken = auth;
                    this.$router.push({ name: "projects" });
                } else {
                    alert("Unfortunately we were unable to log you in. Please check your credentials and try again.");
                }
            } catch (error) {
                alert("Unfortunately we were unable to log you in. Please check your credentials and try again.");
            }
        },

        async register() {
            if (this.checkConfirmPassword) {
                const user = {
                    name: this.name,
                    email: this.email,
                    password: this.password,
                    confirmPassword: this.confirmPassword,
                };

                this.api
                    .register(user)
                    .then(() => {
                        this.accountCreated = true;
                    })
                    .catch(() => {
                        this.accountFailed = true;
                    });
            }
        },

        toggleForm() {
            this.isLoggingIn = !this.isLoggingIn;
        },

        returnToLogin() {
            this.showReset = false;
            this.resetSent = false;
            this.accountCreated = false;
            this.accountFailed = false;
            this.isLoggingIn = true;
        },

        checkEmail() {
            // reset these first
            this.noEmail = false;
            this.emailNotValid = false;
            // then check
            this.noEmail = !this.email || this.email.length == 0;
            if (this.noEmail) {
                return;
            }
            // eslint-disable-next-line
            let emailPattern = /^([a-zA-Z0-9_+\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$/;
            this.emailNotValid = !emailPattern.test(this.email);
        },

        checkPassword() {
            this.noPassword = false;
            this.passwordTooShort = false;
            this.noPassword = !this.password || this.password.length == 0;
            if (this.noPassword) {
                return;
            }
            this.passwordTooShort = this.password.length < 10;
        },

        checkConfirmPassword() {
            this.passwordsNotMatch = this.password != this.confirmPassword;
            return this.passwordsNotMatch;
        },

        checkName() {
            this.noName = !this.name || this.name.length == 0;
            if (this.noName) {
                return;
            }
            let matches = this.name.match(/\s/g);
            this.nameHasSpace = matches && matches.length > 0;
            this.nameTooLong = this.name.length > 255;
        },

        submit(event) {
            if (event) {
                event.preventDefault();
            }
            if (!this.email || !this.password) {
                this.alert("Please provide both an email address and password.");
                return;
            }
            if (this.isLoggingIn) {
                this.login();
            } else {
                this.register();
            }
        },

        showResetPassword() {
            this.showReset = true;
            this.resetSent = false;
            this.sendingReset = false;
            this.resendingReset = false;
        },

        sendResetEmail() {
            this.sendingReset = true;
            this.api.sendResetPasswordEmail(this.resetEmail).then(() => {
                this.resetSent = true;
                this.sendingReset = false;
                this.resendingReset = false;
            });
        },

        resendReset() {
            this.resetSent = false;
            this.resendingReset = true;
            setTimeout(this.sendResetEmail, 500);
        },
    },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#login-form,
#notification,
#reset-container {
    width: 460px;
    background-color: white;
    display: inline-block;
    text-align: center;
    padding-bottom: 60px;
    padding-top: 20px;
}
#notification {
    padding-top: 58px;
}
#reset-container {
    padding-top: 78px;
}
h1 {
    font-weight: 500;
    font-size: 24px;
    margin-bottom: 55px;
}
.bold {
    font-weight: bold;
}
.submit-btn,
.resubmit-btn {
    margin-top: 50px;
    width: 300px;
    height: 45px;
    background-color: #ce596b;
    border: none;
    color: white;
    font-size: 18px;
    font-weight: 600;
    border-radius: 5px;
}
.resubmit-btn {
    margin-top: 0;
}
.disabled {
    margin-top: 20px;
    width: 300px;
    height: 50px;
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
.email-img,
.password-img {
    width: 20px;
    float: right;
    margin-bottom: -18px;
}
ul {
    list-style-type: none;
    padding: 0;
}
li {
    display: inline-block;
    margin: 0 10px;
}
.validation-error {
    float: left;
    color: #c42c44;
    display: block;
    font-size: 14px;
    margin: -20px 0 0 0;
}
.policy-terms-container {
    font-size: 13px;
    line-height: 16px;
    text-align: left;
}
.create-link {
    cursor: pointer;
    margin: 30px 0 0 0;
    font-size: 14px;
    font-weight: 500;
}
.success,
.error {
    margin: 25px 0 25px 0;
    font-size: 24px;
}
.notification-container {
    width: 300px;
    margin: auto;
}
.notification-text {
    width: 300px;
    margin: 0 auto 80px auto;
}
.contact-link {
    cursor: pointer;
    font-size: 14px;
}
.reset-heading {
    font-size: 24px;
    font-weight: 500;
}
.reset-link {
    float: right;
    margin: -16px 0 0 0;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
}
.reset-instructions {
    font-size: 16px;
    text-align: center;
    margin: 15px 0 54px 0;
}
</style>
