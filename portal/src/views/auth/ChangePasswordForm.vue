<template>
    <form @submit.prevent="saveForm">
        <h3>Change Password</h3>
        <div>
            <TextField v-model="form.existing" label="Current password" type="password" />

            <div class="validation-errors" v-if="$v.form.existing.$error">
                <div v-if="!$v.form.existing.required">This is a required field.</div>
                <div v-if="!$v.form.existing.min">Password must be at least 10 characters.</div>
            </div>
        </div>
        <div>
            <TextField v-model="form.password" label="New password" type="password" />

            <div class="validation-errors" v-if="$v.form.password.$error">
                <div v-if="!$v.form.password.required">This is a required field.</div>
                <div v-if="!$v.form.password.min">Password must be at least 10 characters.</div>
            </div>
        </div>
        <div>
            <TextField v-model="form.passwordConfirmation" label="Confirm new password" type="password" />

            <div class="validation-errors" v-if="$v.form.passwordConfirmation.$error">
                <div v-if="!$v.form.passwordConfirmation.required">Confirmation is a required field.</div>
                <div v-if="!$v.form.passwordConfirmation.sameAsPassword">Passwords must match.</div>
            </div>
        </div>
        <button class="save" type="submit">Change password</button>
    </form>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs } from "vuelidate/lib/validators";

import Promise from "bluebird";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default Vue.extend({
    name: "ChangePasswordFOrm",
    components: {
        ...CommonComponents,
    },
    props: {
        user: {
            type: Object,
            required: true,
        },
    },
    data: function (this: any) {
        return {
            form: {
                existing: "",
                password: "",
                passwordConfirmation: "",
            },
        };
    },
    validations: {
        form: {
            existing: {
                required,
                min: minLength(10),
            },
            password: { required, min: minLength(10) },
            passwordConfirmation: { required, min: minLength(10), sameAsPassword: sameAs("password") },
        },
    },
    methods: {
        saveForm() {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                console.log("save form, validation error");
                return;
            }

            return this.$emit("save", this.form);
        },
    },
});
</script>

<style scoped>
.main-panel {
    display: flex;
    flex-direction: column;
    max-width: 700px;
    padding: 20px;
}
.heading {
    font-weight: bold;
    font-size: 24px;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.input-container {
    display: flex;
    flex-direction: column;
    margin: 10px 0 0 0px;
}
.password-change-heading {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 40px;
}
form > div {
    margin-bottom: 20px;
}

button.save {
    margin-top: 20px;
    margin-bottom: 20px;
    width: 300px;
    height: 50px;
    color: white;
    font-size: 18px;
    font-weight: bold;
    background-color: #ce596b;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}
.validation-errors {
    color: #c42c44;
    display: block;
    font-size: 14px;
    margin-bottom: 25px;
}
.notification {
    margin: 20px;
    padding: 20px;
    background-color: #d4edda;
    border: 2px;
    border-radius: 4px;
}

form {
    margin-top: 60px;
}

</style>
