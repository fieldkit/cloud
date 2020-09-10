<template>
    <StandardLayout>
        <div class="main-panel">
            <div class="form-edit" v-if="user">
                <h2> My Account </h2>
                <div class="notification success" v-if="notifySaved">
                    Profile saved.
                </div>

                <ProfileForm :user="user" @save="saveForm" />

                <div class="notification success" v-if="notifyPasswordChanged">
                    Password changed.
                </div>
                <div class="notification failed" v-if="!passwordOk">
                    Please check your password and try again.
                </div>

                <ChangePasswordForm :user="user" @save="changePassword" />
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import ProfileForm from "./ProfileForm.vue";
import ChangePasswordForm from "./ChangePasswordForm.vue";

import Promise from "bluebird";
import FKApi from "@/api/api";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "UserView",
    components: {
        StandardLayout,
        ProfileForm,
        ChangePasswordForm,
    },
    props: {
        id: {
            type: Number,
            required: false,
        },
    },
    data: () => {
        return {
            loading: false,
            notifySaved: false,
            notifyPasswordChanged: false,
            passwordOk: true,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
            stations: (s: GlobalState) => s.stations.user.stations,
            userProjects: (s: GlobalState) => s.stations.user.projects,
        }),
    },
    methods: {
        goBack() {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        saveForm(form) {
            console.log("form", form);
            this.loading = true;
            if (form.image) {
                return this.$store.dispatch(ActionTypes.UPLOAD_USER_PHOTO, { type: form.image.type, file: form.image.file }).then(() => {
                    return this.$store.dispatch(ActionTypes.UPDATE_USER_PROFILE, { user: form }).then(() => {
                        this.loading = false;
                        this.notifySaved = true;
                    });
                });
            } else {
                return this.$store.dispatch(ActionTypes.UPDATE_USER_PROFILE, { user: form }).then(() => {
                    this.loading = false;
                    this.notifySaved = true;
                });
            }
        },
        changePassword(form) {
            console.log("form", form);

            this.passwordOk = true;

            const data = {
                userId: this.user.id,
                oldPassword: form.existing,
                newPassword: form.password,
            };
            return new FKApi()
                .updatePassword(data)
                .then(() => {
                    this.loading = false;
                    this.notifyPasswordChanged = true;
                    return Promise.delay(2000).then(() => {
                        this.notifyPasswordChanged = false;
                    });
                })
                .catch(() => {
                    this.passwordOk = false;
                })
                .finally(() => {
                    this.loading = false;
                });
        },
    },
});
</script>

<style scoped lang="scss">
@import '../../scss/mixins';
@import '../../scss/forms';
@import '../../scss/layout';

.heading {
    font-weight: bold;
    font-size: 24px;
}
.image-container {
    margin-bottom: 40px;
}
::v-deep .user-image img {
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

.validation-errors {
    color: #c42c44;
    display: block;
    font-size: 14px;
    margin-bottom: 25px;
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
