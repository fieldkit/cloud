<template>
    <div class="profile-change">
        <form @submit.prevent="saveForm">
            <h3 class="heading">Profile Picture</h3>
            <div class="user-image image-container">
                <ImageUploader :image="{ url: null }" v-if="!user.photo" @change="onImage" />
                <ImageUploader :image="{ url: user.photo.url }" v-if="user.photo" @change="onImage" />
            </div>
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
                <TextField v-model="form.bio" label="Short Description" />

                <div class="validation-errors" v-if="$v.form.name.$error">
                    <div v-if="!$v.form.bio.required">Short Description is a required field.</div>
                </div>
            </div>
            <div class="checkbox">
                <label>
                    Make my profile public
                    <input type="checkbox" id="checkbox" v-model="form.public" />
                    <span class="checkbox-btn"></span>
                </label>
            </div>
            <button class="button-solid" type="submit">Update</button>
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
            bio: {},
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

<style scoped lang="scss">
@import '../../scss/forms';
@import '../../scss/global';

.image-container {
    margin-bottom: 33px;

    @include bp-down($xs) {
        margin-bottom: 40px;
    }
}
::v-deep .user-image img {
    max-width: 198px!important;
    max-height: 198px!important;
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
form > div {
    margin-bottom: 20px;
}
.password-change-heading {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 40px;
}
.button-solid {
    margin-top: 25px;

    @include bp-down($xs) {
        width: 100%;
    }
}

.checkbox {
    margin-top: 25px;
}
</style>
