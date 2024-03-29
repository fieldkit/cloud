<template>
    <div class="form-container">
        <Logo class="form-header-logo"></Logo>
        <form class="form" @submit.prevent="save">
            <template v-if="!success && !failed">
                <h1 class="form-title">Reset Password</h1>

                <div class="form-group">
                    <TextField v-model="form.password" label="Password" type="password" />

                    <div class="form-errors" v-if="$v.form.password.$error">
                        <div v-if="!$v.form.password.required">This is a required field.</div>
                        <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
                    </div>
                </div>

                <div class="form-group">
                    <TextField v-model="form.passwordConfirmation" label="Confirm Password" type="password" />

                    <div class="form-errors" v-if="$v.form.passwordConfirmation.$error">
                        <div v-if="!$v.form.passwordConfirmation.required">Confirmation is a required field.</div>
                        <div v-if="!$v.form.passwordConfirmation.sameAsPassword">Passwords must match.</div>
                    </div>
                </div>
                <button class="form-submit" v-on:click="save">Reset</button>
                <div>
                    <router-link :to="{ name: 'login' }" class="form-link">Back to Log In</router-link>
                </div>
            </template>
            <template v-if="success">
                <img src="@/assets/icon-success.svg" alt="Success" class="form-header-icon" width="57px" />
                <h1 class="form-title">Success!</h1>

                <router-link :to="{ name: 'login' }" class="form-link">Back to Log In</router-link>
            </template>
            <template v-if="failed">
                <img src="@/assets/icon-warning-error.svg" alt="Unsuccessful" class="form-header-icon" width="57px" />
                <h1 class="form-title">Password Not Reset</h1>
                <div class="form-subtitle">Unfortunately we were unable to reset your password.</div>
                <d>
                    Please
                    <a href="https://www.fieldkit.org/contact/" class="contact-link">contact us</a>
                    if you would like assistance.
                </d>
            </template>
        </form>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { required, email, minLength, sameAs } from "vuelidate/lib/validators";
import Logo from "@/views/shared/Logo.vue";

export default Vue.extend({
    name: "ResetPasswordView",
    components: {
        ...CommonComponents,
        Logo,
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
            return this.$services.api
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

<style scoped lang="scss">
@import "../../scss/forms.scss";

.contact-link {
    cursor: pointer;
    font-weight: 500;
    text-decoration: underline;

    @include bp-down($xs) {
        font-size: 14px;
    }
}
</style>
