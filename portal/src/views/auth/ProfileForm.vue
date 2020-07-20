<template>
    <div class="profile-change">
        <form @submit.prevent="saveForm">
            <h3>Profile Picture</h3>
            <div class="user-image image-container">
                <ImageUploader :image="{ url: null }" v-if="!user.photo" @change="onImage" />
                <ImageUploader :image="{ url: user.photo.url }" v-if="user.photo" @change="onImage" />
            </div>
            <h3>User Profile</h3>
            <div>
                <TextField v-model="form.name" label="Name" />

                <div class="validation-errors" v-if="$v.form.name.$error">
                    <div v-if="!$v.form.name.required">Name is a required field.</div>
                </div>
            </div>
            <div>
                <TextField v-model="form.email" label="Email" />

                <div class="validation-errors" v-if="$v.form.email.$error">
                    <div v-if="!$v.form.email.required">Email is a required field.</div>
                    <div v-if="!$v.form.email.email">Must be a valid email address.</div>
                    <div v-if="!$v.form.email.taken">
                        This address appears to already be registered.
                        <router-link :to="{ name: 'recover' }" class="recover-link">
                            Recover Account
                        </router-link>
                    </div>
                </div>
            </div>
            <div>
                <TextField v-model="form.bio" label="Bio" />

                <div class="validation-errors" v-if="$v.form.name.$error">
                    <div v-if="!$v.form.bio.required">Bio is a required field.</div>
                </div>
            </div>
            <button class="save" type="submit">Update</button>
        </form>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs } from "vuelidate/lib/validators";

import Promise from "bluebird";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default Vue.extend({
    name: "ProfileForm",
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
            loading: false,
            form: {
                id: this.user.id,
                name: this.user.name,
                email: this.user.email,
                bio: this.user.bio,
                image: null,
            },
            changePasswordForm: {
                existing: "",
                password: "",
                passwordConfirmation: "",
            },
            notifySaved: false,
            notifyPasswordChanged: false,
            emailAvailable: true,
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
            },
            bio: {
                required,
            },
        },
        changePasswordForm: {
            existing: {
                required,
                min: minLength(10),
            },
            password: { required, min: minLength(10) },
            passwordConfirmation: { required, min: minLength(10), sameAsPassword: sameAs("password") },
        },
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        showStation(station) {
            return this.$router.push({ name: "viewStation", params: { id: station.id } });
        },
        saveForm() {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                console.log("save form, validation error");
                return;
            }

            return this.$emit("save", this.form);
        },
        changePassword() {
            this.$v.changePasswordForm.$touch();
            if (this.$v.changePasswordForm.$pending || this.$v.changePasswordForm.$error) {
                console.log("password form, validation error");
                return;
            }

            return this.$emit("change-password", this.changePasswordForm);
        },
        onImage(image) {
            this.form.image = image;
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
.image-container {
    margin-bottom: 40px;
}
/deep/ .user-image img {
    max-width: 400px;
    max-height: 400px;
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
form > div {
    margin-bottom: 20px;
}
#public-checkbox-container {
}
#public-checkbox-container input {
}
#public-checkbox-container label {
}
#public-checkbox-container img {
}
.password-change-heading {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 40px;
}
.password-change {
    margin-top: 20px;
}
button.save {
    margin-top: 20px;
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
</style>
