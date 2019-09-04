<template>
    <div id="LoginForm">
        <h1>Log In to Your Account</h1>
        <input
            id="email-field"
            placeholder="Email"
            keyboardType="email"
            autocorrect="false"
            autocapitalizationType="none"
            v-model="email"
            @keyup.enter="$refs.password.focus"
            @blur="checkEmail"
        />
        <span class="validation-error" id="no-email" v-if="noEmail">Email is a required field.</span>
        <span class="validation-error" id="email-not-valid" v-if="emailNotValid"
            >Must be a valid email address.</span
        >
        <br />
        <br />
        <input
            name="password"
            placeholder="Password"
            secure="true"
            ref="password"
            type="password"
            v-model="password"
            @keyup.enter="login"
            @blur="checkPassword"
        />
        <span class="validation-error" id="no-password" v-if="noPassword">Password is a required field.</span>
        <span class="validation-error" id="password-too-short" v-if="passwordTooShort"
            >Password must be at least 10 characters.</span
        >
        <br />
        <button v-on:click="login" ref="submit">Log In</button>
    </div>
</template>

<script>
import FKApi from "../../api/api";

export default {
    name: "LoginForm",
    data: () => {
        return {
            email: "",
            password: "",
            user: {},
            emailNotValid: false,
            noEmail: false,
            noPassword: false,
            passwordTooShort: false
        };
    },
    methods: {
        async login(event) {
            event.preventDefault();
            try {
                const api = new FKApi();
                const auth = await api.login(this.email, this.password);
                const isAuthenticated = await api.authenticated();
                if (isAuthenticated) {
                    this.userToken = auth;
                    this.$router.push({ name: "dashboard" });
                } else {
                    alert(
                        "Unfortunately we were unable to log you in. Please check your credentials and try again."
                    );
                }
            } catch (error) {
                alert(
                    "Unfortunately we were unable to log you in. Please check your credentials and try again."
                );
            }
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
            let emailPattern = /^([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$/;
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
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1 {
    font-weight: lighter;
}
button {
    margin-top: 50px;
    width: 70%;
    height: 50px;
    background-color: #ce596b;
    border: none;
    color: white;
    font-size: 18px;
    border-radius: 5px;
}
input {
    border: 0;
    border-bottom: 1px solid gray;
    outline: 0;
    height: 30px;
    width: 70%;
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
div {
    width: 30%;
    background-color: white;
    display: inline-block;
    text-align: center;
    padding-bottom: 60px;
    padding-top: 20px;
}
.validation-error {
    color: #c42c44;
    display: block;
}
</style>
