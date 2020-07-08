<template>
    <div id="login-form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="../../assets/FieldKit_Logo_White.png" />
        <div>
            <form id="form" @submit.prevent="save">
                <h1>Log In to Your Account</h1>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input
                            keyboardType="email"
                            autocorrect="false"
                            autocapitalizationType="none"
                            class="inputText"
                            v-model="form.email"
                        />
                        <span class="floating-label">Email</span>

                        <div class="validation-errors" v-if="$v.form.email.$error">
                            <div v-if="!$v.form.email.required">Email is a required field.</div>
                            <div v-if="!$v.form.email.email">Must be a valid email address.</div>
                        </div>
                    </div>
                </div>
                <div class="outer-input-container">
                    <div class="input-container middle-container">
                        <input name="password" secure="true" type="password" class="inputText" v-model="form.password" />
                        <span class="floating-label">Password</span>

                        <div class="validation-errors" v-if="$v.form.password.$error">
                            <div v-if="!$v.form.password.required">This is a required field.</div>
                            <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
                        </div>
                    </div>
                </div>
                <div v-if="failed" class="login-failed">
                    Unfortunately we were unable to log you in. Please check your credentials and try again.
                </div>
                <button class="form-save-btn" type="submit">Log In</button>
                <div>
                    <router-link :to="{ name: 'recover' }" class="create-link">
                        Reset My Password
                    </router-link>
                </div>
                <div>
                    <router-link :to="{ name: 'register' }" class="create-link">
                        Create an Account
                    </router-link>
                </div>
            </form>
        </div>
    </div>
</template>

<script>
import { required, email, minLength, sameAs } from "vuelidate/lib/validators";
import FKApi, { LoginPayload } from "@/api/api";
import * as ActionTypes from "@/store/actions";

export default {
    data() {
        return {
            form: {
                email: "",
                password: "",
            },
            busy: false,
            failed: false,
        };
    },
    validations: {
        form: {
            email: {
                required,
                email,
            },
            password: { required, min: minLength(10) },
        },
    },
    methods: {
        save() {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            const payload = new LoginPayload(this.form.email, this.form.password);

            this.busy = true;
            this.failed = false;

            return this.$store
                .dispatch(ActionTypes.AUTHENTICATE, payload)
                .then(
                    () => this.$router.push(this.$route.query.after || { name: "projects" }),
                    () => (this.failed = true)
                )
                .finally(() => {
                    this.busy = false;
                });
        },
    },
};
</script>

<style scoped>
#login-form-container {
    width: 100%;
    min-height: 100%;
    background-image: linear-gradient(#52b5e4, #1b80c9);
}
#form,
#notifications {
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
.form-save-btn {
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
ul {
    list-style-type: none;
    padding: 0;
}
li {
    display: inline-block;
    margin: 0 10px;
}
.policy-terms-container {
    font-size: 13px;
    line-height: 16px;
    text-align: left;
}
.create-link {
    cursor: pointer;
    display: block;
    margin-top: 30px;
    font-size: 14px;
    font-weight: bold;
}
.form-header-logo {
    width: 210px;
    margin-top: 150px;
    margin-bottom: 86px;
}
.validation-errors {
    color: #c42c44;
    display: block;
    font-size: 14px;
}
.recover-link {
    color: black;
    font-weight: bold;
    margin-top: 10px;
    margin-bottom: 30px;
    display: block;
}
.login-failed {
    color: #c42c44;
    font-weight: bold;
    width: 50%;
    margin: auto;
}
</style>
