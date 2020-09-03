<template>
    <div class="form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="@/assets/FieldKit_Logo_White.png" />
        <form class="form" @submit.prevent="save">
            <h1 class="form-title"> Log In to Your Account </h1>
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
                <router-link :to="{ name: 'recover', query: forwardAfterQuery() }" class="form-link">
                    Reset password
                </router-link>
            </div>
            <div v-if="failed" class="login-failed">
                Unfortunately we were unable to log you in. Please check your credentials and try again.
            </div>
            <button class="form-submit" type="submit">Log In</button>
            <div>
                <router-link :to="{ name: 'register', query: forwardAfterQuery() }" class="form-link">
                    Create an account
                </router-link>
            </div>
        </form>
    </div>
</template>

<script lang="ts">
    import _ from "lodash";
    import Vue from "vue";
    import CommonComponents from "@/views/shared";

    import { required, email, minLength, sameAs } from "vuelidate/lib/validators";

    import FKApi, { LoginPayload } from "@/api/api";
    import * as ActionTypes from "@/store/actions";

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
        data() {
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
            if (this.spoofing) {
                return {
                    form: {
                        spoofEmail: {
                            required,
                            email,
                        },
                        email: {
                            required,
                            email,
                        },
                        password: { required, min: minLength(10) },
                    },
                };
            }
            return {
                form: {
                    email: {
                        required,
                        email,
                    },
                    password: { required, min: minLength(10) },
                },
            };
        },
        methods: {
            forwardAfterQuery() {
                if (this.$route.query.after) {
                    return {
                        after: this.$route.query.after,
                    };
                }
                return {};
            },
            createPayload() {
                if (this.spoofing) {
                    return new LoginPayload(this.form.spoofEmail, this.form.email + " " + this.form.password);
                }
                return new LoginPayload(this.form.email, this.form.password);
            },
            save() {
                this.$v.form.$touch();
                if (this.$v.form.$pending || this.$v.form.$error) {
                    return;
                }

                const payload = this.createPayload();

                this.busy = true;
                this.failed = false;

                return this.$store
                    .dispatch(ActionTypes.LOGIN, payload)
                    .then(
                        () => {
                            if (this.$route.query.after) {
                                if (_.isArray(this.$route.query.after)) {
                                    return this.$router.push(this.$route.query.after[0]);
                                }
                                return this.$router.push(this.$route.query.after as string);
                            }
                            return this.$router.push({ name: "projects" });
                        },
                        () => (this.failed = true)
                    )
                    .finally(() => {
                        this.busy = false;
                    });
            },
        },
    });
</script>

<style scoped lang="scss">

    @import '../../scss/forms.scss';

</style>
