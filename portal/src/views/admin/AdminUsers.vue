<template>
    <StandardLayout>
        <div class="container">
            <router-link :to="{ name: 'adminMain' }" class="link">
                Back to Admin
            </router-link>

            <vue-confirm-dialog />

            <div class="delete-user">
                <div v-if="deletion.failed" class="notification failed">
                    Oops, there was a problem.
                </div>

                <div v-if="deletion.success" class="notification success">
                    They're toast.
                </div>

                <form id="delete-user-form" @submit.prevent="deleteUser">
                    <h2>Delete User</h2>
                    <div>
                        <TextField v-model="form.email" label="Buh-bye Email" />
                        <div class="validation-errors" v-if="$v.form.email.$error">
                            <div v-if="!$v.form.email.required">Email is a required field.</div>
                            <div v-if="!$v.form.email.email">Must be a valid email address.</div>
                        </div>
                    </div>
                    <div>
                        <TextField v-model="form.password" label="Your Password" type="password" />
                        <div class="validation-errors" v-if="$v.form.password.$error">
                            <div v-if="!$v.form.password.required">This is a required field.</div>
                            <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
                        </div>
                    </div>
                    <button class="form-save-button" type="submit">Bye bye bye</button>
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

import FKApi from "@/api/api";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "AdminMain",
    components: {
        StandardLayout,
        ...CommonComponents,
    },
    props: {},
    data: () => {
        return {
            form: {
                email: "",
                password: "",
            },
            deletion: {
                success: false,
                failed: false,
            },
        };
    },
    validations: {
        form: {
            email: {
                required,
                email,
            },
            password: {
                required,
                min: minLength(10),
            },
        },
    },
    methods: {
        deleteUser(this: any) {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
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
                        return new FKApi().adminDeleteUser(this.form).then(
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
