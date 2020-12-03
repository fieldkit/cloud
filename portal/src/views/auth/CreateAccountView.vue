<template>
    <div class="form-container">
        <div class="form-wrap">
            <img class="form-header-logo" alt="FieldKit Logo" src="@/assets/FieldKit_Logo_White.png" />
            <form v-if="!created" class="form" @submit.prevent="save">
                <h1 class="form-title">Create Your Account</h1>
                <div class="form-group">
                    <TextField v-model="form.name" label="Name" />

                    <div class="form-errors" v-if="$v.form.name.$error">
                        <div v-if="!$v.form.name.required">Name is a required field.</div>
                    </div>
                </div>
                <div class="form-group">
                    <TextField v-model="form.email" label="Email" keyboardType="email" />

                    <div class="form-errors" v-if="$v.form.email.$error">
                        <div v-if="!$v.form.email.required">Email is a required field.</div>
                        <div v-if="!$v.form.email.email">Must be a valid email address.</div>
                        <div v-if="!$v.form.email.taken">
                            This address appears to already be registered.
                            <router-link :to="{ name: 'recover' }" class="form-link recover">Recover Account</router-link>
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <TextField v-model="form.password" label="Password" type="password" />

                    <div class="form-errors" v-if="$v.form.password.$error">
                        <div v-if="!$v.form.password.required">This is a required field.</div>
                        <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
                    </div>
                </div>
                <div class="form-group">
                    <TextField v-model="form.passwordConfirmation" label="Confirm password" type="password" />

                    <div class="form-errors" v-if="$v.form.passwordConfirmation.$error">
                        <div v-if="!$v.form.passwordConfirmation.required">Confirmation is a required field.</div>
                        <div v-if="!$v.form.passwordConfirmation.sameAsPassword">Passwords must match.</div>
                    </div>
                    <div class="form-policy">
                        By creating an account you agree to our
                        <span class="bold">Privacy Policy</span>
                        and
                        <span class="bold">Terms of Use.</span>
                    </div>
                </div>
                <button class="form-submit" type="submit">Create Account</button>
                <div>
                    <router-link :to="{ name: 'login' }" class="form-link">Back to Log In</router-link>
                </div>
            </form>
            <form class="form" v-if="created">
                <template v-if="!resending">
                    <img src="@/assets/icon-success.svg" alt="Success" class="form-header-icon" width="57px" />
                    <h1 class="form-title">Account Created</h1>
                    <h2 class="form-subtitle">We sent you an account validation email.</h2>
                    <button class="form-submit" v-on:click="resend">Resend Email</button>
                    <router-link :to="{ name: 'login' }" class="form-link">Back to Log In</router-link>
                </template>
                <template v-if="resending">
                    <img src="@/assets/Icon_Syncing2.png" alt="Resending" class="form-header-icon" width="57px" />
                    <p class="form-link">Resending</p>
                </template>
            </form>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs } from "vuelidate/lib/validators";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "CreateAccountView",
    components: {
        ...CommonComponents,
    },
    data(): {
        form: {
            name: string;
            email: string;
            password: string;
            passwordConfirmation: string;
        };
        available: boolean;
        creating: boolean;
        resending: boolean;
        created: { id: number } | null;
    } {
        return {
            form: {
                name: "",
                email: "",
                password: "",
                passwordConfirmation: "",
            },
            available: true,
            creating: true,
            resending: false,
            created: null,
        };
    },
    validations: {
        form: {
            name: {
                required,
            },
            email: {
                required,
                email,
                taken: function (this: any) {
                    return this.available;
                },
            },
            password: { required, min: minLength(10) },
            passwordConfirmation: { required, min: minLength(10), sameAsPassword: sameAs("password") },
        },
    },
    methods: {
        async save(): Promise<void> {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            this.creating = true;

            await this.$services.api
                .register(this.form)
                .then((created) => {
                    this.created = created;
                })
                .catch((error) => {
                    if (error.status === 400) {
                        this.available = false;
                    } else {
                        return this.$seriousError(error);
                    }
                })
                .finally(() => {
                    this.creating = false;
                });
        },
        async resend(): Promise<void> {
            if (!this.created) throw new Error(`nothing to resend`);
            this.resending = true;
            await this.$services.api.resendCreateAccount(this.created.id).then(() => {
                this.resending = false;
            });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms.scss";
</style>
