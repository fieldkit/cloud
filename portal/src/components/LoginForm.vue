<template>
    <div>
        <div id="login-form" v-if="!accountCreated && !accountFailed && !showReset">
            <h1>{{ isLoggingIn ? "Log In to Your Account" : "Create an Account" }}</h1>
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
                <div class="input-container">
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
                <div class="input-container">
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
            </div>
            <div class="outer-input-container" v-if="!isLoggingIn">
                <div class="input-container">
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
            </div>
            <button v-on:click="submit" ref="submit">{{ isLoggingIn ? "Log In" : "Sign Up" }}</button>
            <div class="create-link" v-on:click="toggleForm">
                {{ isLoggingIn ? "Create an Account" : "Back to Log In" }}
            </div>
            <div class="forgot-link" v-if="isLoggingIn" v-on:click="showResetPassword">Forgot your password?</div>
        </div>

        <div id="forgot-link-container" v-if="showReset">
            <div class="outer-input-container" v-if="!resetSent">
                <div class="input-container">
                    <div class="reset-instructions">
                        Enter your email address below, and password reset instructions will be sent to you.
                    </div>
                    <input v-model="resetEmail" type="text" class="inputText" required="" />
                    <span class="floating-label">Email</span>
                </div>
            </div>
            <button class="save-btn" v-on:click="sendResetEmail" v-if="!resetSent">Submit</button>
            <div class="reset-sent" v-if="resetSent">Password reset email sent!</div>
        </div>

        <div id="notification" v-if="accountCreated || accountFailed">
            <div v-if="accountCreated">
                <p class="success">Success! Your account was created.</p>
                <p>Check your email to validate your account.</p>
            </div>
            <div v-if="accountFailed">
                <p class="error">Unfortunately we were unable to create your account.</p>
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
        },
        sendResetEmail() {
            this.api.sendResetPasswordEmail(this.resetEmail).then(() => {
                this.resetSent = true;
            });
        },
    },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#login-form,
#notification,
#forgot-link-container {
    width: 30%;
    background-color: white;
    display: inline-block;
    text-align: center;
    padding-bottom: 60px;
    padding-top: 20px;
}
h1 {
    font-weight: lighter;
}
button {
    margin-top: 20px;
    width: 70%;
    height: 50px;
    background-color: #ce596b;
    border: none;
    color: white;
    font-size: 18px;
    border-radius: 5px;
}
.outer-input-container {
    height: 65px;
    width: 70%;
    margin: auto;
}
.input-container {
    float: left;
    margin: 5px 0 0 0;
    width: 100%;
    text-align: left;
}
input {
    border: 0;
    border-bottom: 1px solid gray;
    outline: 0;
    font-size: 18px;
}
.floating-label {
    font-size: 18px;
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
.create-link {
    cursor: pointer;
    margin: 40px 0 0 0;
}
.success {
    margin: 15px 0 0 0;
    font-size: 18px;
    color: #3f8530;
}
.error {
    margin: 15px 0 0 0;
    font-size: 18px;
    color: #c42c44;
    padding: 10px;
}
.contact-link {
    cursor: pointer;
    text-decoration: underline;
}
.forgot-link {
    margin: 15px 0 0 0;
    cursor: pointer;
}
.reset-instructions {
    margin: 15px 0 18px 0;
}
.reset-sent {
    margin: 15px 0 0 0;
    font-size: 18px;
    color: #3f8530;
}
</style>
