<template>
    <StandardLayout>
        <div class="container">
            <router-link :to="{ name: 'adminMain' }" class="link">Back to Admin</router-link>

            <vue-confirm-dialog />

            <div class="delete-user">
                <div v-if="deletion.failed" class="notification failed">Oops, there was a problem.</div>

                <div v-if="deletion.success" class="notification success">They're toast.</div>

                <form id="delete-user-form" @submit.prevent="deleteUser">
                    <h2>Delete User</h2>
                    <div>
                        <TextField v-model="deletionForm.email" label="Buh-bye Email" />
                        <div class="validation-errors" v-if="$v.deletionForm.email.$error">
                            <div v-if="!$v.deletionForm.email.required">Email is a required field.</div>
                            <div v-if="!$v.deletionForm.email.email">Must be a valid email address.</div>
                        </div>
                    </div>
                    <div>
                        <TextField v-model="deletionForm.password" label="Your Password" type="password" />
                        <div class="validation-errors" v-if="$v.deletionForm.password.$error">
                            <div v-if="!$v.deletionForm.password.required">This is a required field.</div>
                            <div v-if="!$v.deletionForm.password.min">Password must be at least 10 characters.</div>
                        </div>
                    </div>
                    <button class="form-save-button" type="submit">Bye bye bye</button>
                </form>
            </div>

            <div class="clear-tnc">
                <div v-if="tnc.failed" class="notification failed">Oops, there was a problem.</div>

                <div v-if="tnc.success" class="notification success">Done, cleared that user's TNC status.</div>

                <form id="clear-tnc-form" @submit.prevent="clearTermsAndConditions">
                    <h2>Clear TNC</h2>
                    <div>
                        <TextField v-model="tncForm.email" label="Email" />
                        <div class="validation-errors" v-if="$v.tncForm.email.$error">
                            <div v-if="!$v.tncForm.email.required">Email is a required field.</div>
                            <div v-if="!$v.tncForm.email.email">Must be a valid email address.</div>
                        </div>
                    </div>
                    <button class="form-save-button" type="submit">Clear Terms and Conditions</button>
                </form>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs } from "vuelidate/lib/validators";

export default Vue.extend({
    name: "AdminMain",
    components: {
        StandardLayout,
        ...CommonComponents,
    },
    props: {},
    data: () => {
        return {
            deletionForm: {
                email: "",
                password: "",
            },
            tncForm: {
                email: "",
            },
            deletion: {
                success: false,
                failed: false,
            },
            tnc: {
                success: false,
                failed: false,
            },
        };
    },
    validations: {
        deletionForm: {
            email: {
                required,
                email,
            },
            password: {
                required,
                min: minLength(10),
            },
        },
        tncForm: {
            email: {
                required,
                email,
            },
        },
    },
    methods: {
        deleteUser(this: any) {
            this.$v.deletionForm.$touch();
            if (this.$v.deletionForm.$pending || this.$v.deletionForm.$error) {
                return;
            }

            return this.$confirm({
                message: `Are you sure? This operation cannot be undone.`,
                button: {
                    no: "No",
                    yes: "Yes",
                },
                callback: (confirm) => {
                    if (confirm) {
                        return this.$services.api.adminDeleteUser(this.form).then(
                            () => {
                                this.deletion.success = true;
                                this.deletion.failed = false;
                            },
                            () => {
                                this.deletion.failed = true;
                                this.deletion.success = false;
                            }
                        );
                    }
                },
            });
        },
        clearTermsAndConditions(this: any) {
            this.$v.tncForm.$touch();
            if (this.$v.tncForm.$pending || this.$v.tncForm.$error) {
                return;
            }

            return this.$confirm({
                message: `Are you sure? This operation cannot be undone.`,
                button: {
                    no: "No",
                    yes: "Yes",
                },
                callback: (confirm) => {
                    if (confirm) {
                        return this.$services.api.adminClearTermsAndConditions(this.tncForm).then(
                            () => {
                                this.tnc.success = true;
                                this.tnc.failed = false;
                            },
                            () => {
                                this.tnc.failed = true;
                                this.tnc.success = false;
                            }
                        );
                    }
                },
            });
        },
    },
});
</script>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    padding: 20px;
    text-align: left;
}
.delete-user {
    width: 400px;
    text-align: left;
}
.form-save-button {
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
.notification.success {
    margin-top: 20px;
    margin-bottom: 20px;
    padding: 20px;
    border: 2px;
    border-radius: 4px;
}
.notification.success {
    background-color: #d4edda;
}
.notification.failed {
    background-color: #f8d7da;
}
</style>
