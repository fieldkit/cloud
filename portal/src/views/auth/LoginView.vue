<template>
    <div class="form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="@/assets/FieldKit_Logo_White.png" />
        <form class="form" @submit.prevent="save">
            <h1 class="form-title">Log In to Your Account</h1>
            <div class="form-group" v-if="spoofing">
                <TextField v-model="form.spoofEmail" label="Spoof Email" />
                <div class="form-errors" v-if="$v.form.spoofEmail.$error">
                    <div v-if="!$v.form.spoofEmail.required">Spoof Email is a required field.</div>
                    <div v-if="!$v.form.spoofEmail.email">Must be a valid email address.</div>
                </div>
            </div>
            <div class="form-group">
                <TextField v-model="form.email" label="Email" />
                <div class="form-errors" v-if="$v.form.email.$error">
                    <div v-if="!$v.form.email.required">Email is a required field.</div>
                    <div v-if="!$v.form.email.email">Must be a valid email address.</div>
                </div>
            </div>
            <div class="form-group">
                <TextField v-model="form.password" label="Password" type="password" />
                <div class="form-errors" v-if="$v.form.password.$error">
                    <div v-if="!$v.form.password.required">This is a required field.</div>
                    <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
                </div>
            </div>
            <div class="form-link-recover">
                <router-link :to="{ name: 'recover', query: forwardAfterQuery() }" class="form-link">Reset password</router-link>
            </div>
            <div v-if="failed" class="login-failed">
                Unfortunately we were unable to log you in. Please check your credentials and try again.
            </div>
            <button class="form-submit" type="submit">Log In</button>
            <div>
                <router-link :to="{ name: 'register', query: forwardAfterQuery() }" class="form-link">Create an account</router-link>
            </div>
        </form>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs, requiredIf } from "vuelidate/lib/validators";

import FKApi, { LoginPayload } from "@/api/api";
import { ActionTypes } from "@/store";

export default Vue.extend({
    components: {
        ...CommonComponents,
    },
    props: {
        spoofing: {
            type: Boolean,
            default: false,
        },
    },
    data(): {
        form: {
            spoofEmail: string;
            email: string;
            password: string;
        };
        busy: boolean;
        failed: boolean;
    } {
        return {
            form: {
                spoofEmail: "",
                email: "",
                password: "",
            },
            busy: false,
            failed: false,
        };
    },
    validations() {
        return {
            form: {
                spoofEmail: {
                    // eslint-disable-next-line
                    required: requiredIf(function (this: any) {
                        return this.spoofing;
                    }),
                    email,
                },
                email: {
                    required,
                    email,
                },
                password: { required, min: minLength(10) },
            },
        };
    },
    methods: {
        forwardAfterQuery(): { after?: string } {
            const after = this.$route.query.after;
            if (after) {
                if (_.isArray(after) && after.length > 0 && after[0]) {
                    return { after: after[0] };
                }
                return { after: after as string };
            }
            return {};
        },
        createPayload(): LoginPayload {
            if (this.spoofing) {
                return new LoginPayload(this.form.spoofEmail, this.form.email + " " + this.form.password);
            }
            return new LoginPayload(this.form.email, this.form.password);
        },
        async save(): Promise<void> {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            const payload = this.createPayload();

            this.busy = true;
            this.failed = false;

            await this.$store
                .dispatch(ActionTypes.LOGIN, payload)
                .then(
                    async () => {
                        await this.leaveAfterAuth();
                    },
                    () => (this.failed = true)
                )
                .finally(() => {
                    this.busy = false;
                });
        },
        async leaveAfterAuth(): Promise<void> {
            const after = this.forwardAfterQuery();
            if (after.after) {
                await this.$router.push(after.after);
            } else {
                await this.$router.push({ name: "projects" });
            }
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms.scss";
</style>
