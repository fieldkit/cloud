<template>
    <div class="form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="@/assets/FieldKit_Logo_White.png" />
        <div v-if="!attempted">
            <form class="form" @submit.prevent="save">
                <h1 class="form-title">Recover Account</h1>
                <div class="form-subtitle">Enter your email address below, and password reset instructions will be sent.</div>
                <div class="form-group">
                    <TextField v-model="form.email" label="Email" keyboardType="email" />

                    <div class="form-errors" v-if="$v.form.email.$error">
                        <div v-if="!$v.form.email.required">Email is a required field.</div>
                        <div v-if="!$v.form.email.email">Must be a valid email address.</div>
                    </div>
                </div>
                <button class="form-submit" type="submit">Recover</button>
                <div>
                    <router-link :to="{ name: 'login' }" class="form-link">Back to Log In</router-link>
                </div>
            </form>
        </div>
        <div v-if="attempted" class="form success">
            <div v-if="!resending">
                <img alt="Success" src="@/assets/icon-success.svg" width="57px" class="form-header-icon" />
                <h1 class="form-title">Password Reset Email Sent</h1>
                <div class="form-subtitle">Check your inbox for the email with a link to reset your password.</div>
                <button class="form-submit" v-on:click="resend">Resend Email</button>
                <router-link :to="{ name: 'login' }" class="form-link">Back to Log In</router-link>
            </div>
            <div v-if="resending">
                <img alt="Resending" src="@/assets/Icon_Syncing2.png" width="57px" class="form-header-icon" />
                <p>Resending</p>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs } from "vuelidate/lib/validators";

export default Vue.extend({
    name: "RecoverAccountView",
    components: {
        ...CommonComponents,
    },
    data() {
        return {
            form: {
                email: "",
            },
            resending: false,
            attempted: false,
            busy: false,
        };
    },
    validations: {
        form: {
            email: {
                required,
                email,
            },
        },
    },
    methods: {
        save() {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }
            this.busy = true;
            return this.$services.api
                .sendResetPasswordEmail(this.form.email)
                .then(() => (this.attempted = true))
                .finally(() => (this.busy = true));
        },
        resend() {
            this.resending = true;
            return this.save().finally(() => (this.resending = false));
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms";

.reset-instructions {
    margin-bottom: 50px;
}

.form-submit {
    margin-top: 80px;

    @include bp-down($xs) {
        margin-top: 70px;
    }
}

.form-group {
    margin-top: 50px;
}

.form.success {
    .form-submit {
        margin-top: 50px;
    }
}

.form:not(.success) .form-subtitle {
    margin-top: -25px;
}
</style>
