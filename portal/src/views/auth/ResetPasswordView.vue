<template>
    <div id="login-form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="../../assets/FieldKit_Logo_White.png" />
        <div>
            <div id="user-form-container">
                <div class="password-change" v-if="!success && !failed">
                    <div class="inner-password-change">
                        <div class="password-change-heading">Reset Password</div>

                        <div class="outer-input-container">
                            <div class="input-container">
                                <TextField v-model="form.password" label="Password" type="password" />

                                <div class="validation-errors" v-if="$v.form.password.$error">
                                    <div v-if="!$v.form.password.required">This is a required field.</div>
                                    <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
                                </div>
                            </div>
                        </div>

                        <div class="outer-input-container">
                            <div class="input-container middle-container">
                                <TextField v-model="form.passwordConfirmation" label="Confirm Password" type="password" />

                                <div class="validation-errors" v-if="$v.form.passwordConfirmation.$error">
                                    <div v-if="!$v.form.passwordConfirmation.required">Confirmation is a required field.</div>
                                    <div v-if="!$v.form.passwordConfirmation.sameAsPassword">Passwords must match.</div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <button class="save-btn" v-on:click="save">Reset</button>
                    <div>
                        <router-link :to="{ name: 'login' }" class="create-link">
                            Back to Log In
                        </router-link>
                    </div>
                </div>
                <div v-if="success">
                    <img alt="Success" src="../../assets/Icon_Success.png" width="57px" />
                    <p class="success">Success!</p>

                    <router-link :to="{ name: 'login' }" class="create-link">
                        Back to Log In
                    </router-link>
                </div>
                <div v-if="failed">
                    <img alt="Unsuccessful" src="../../assets/Icon_Warning_error.png" width="57px" />
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
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { required, email, minLength, sameAs } from "vuelidate/lib/validators";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "ResetPasswordView",
    components: {
        ...CommonComponents,
    },
    data: () => {
        return {
            form: {
                password: "",
                passwordConfirmation: "",
            },
            busy: false,
            success: false,
            failed: false,
        };
    },
    validations: {
        form: {
            password: { required, min: minLength(10) },
            passwordConfirmation: { required, min: minLength(10), sameAsPassword: sameAs("password") },
        },
    },
    methods: {
        save(this: any) {
            console.log("save");
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            this.success = false;
            this.failed = false;
            this.busy = true;
            const payload = {
                token: this.$route.query.token,
                password: this.form.password,
            };
            return new FKApi()
                .resetPassword(payload)
                .then(() => {
                    this.success = true;
                })
                .catch(() => {
                    this.failed = true;
                })
                .finally(() => {
                    this.busy = false;
                });
        },
    },
});
</script>

<style scoped>
#login-form-container {
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
.validation-errors {
    color: #c42c44;
    display: block;
    font-size: 14px;
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
.form-header-logo {
    width: 210px;
    margin-top: 150px;
    margin-bottom: 86px;
}
</style>
