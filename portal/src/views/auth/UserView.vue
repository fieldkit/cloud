<template>
    <StandardLayout>
        <div class="main-panel">
            <div class="form-edit" v-if="user">
                <h2>{{ $t("user.profile.title") }}</h2>
                <div class="notification success" v-if="notifySaved">{{ $t("user.profile.saved") }}</div>

                <ProfileForm :user="user" @save="saveForm" />

                <div class="notification success" v-if="notifyPasswordChanged">{{ $t("user.profile.password.changed") }}</div>
                <div class="notification failed" v-if="!passwordOk">{{ $t("user.profile.password.failed") }}</div>

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
    data(): {
        loading: boolean;
        notifySaved: boolean;
        notifyPasswordChanged: boolean;
        passwordOk: boolean;
    } {
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
        goBack(): void {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        async saveForm(form): Promise<void> {
            console.log("form", form);
            this.loading = true;
            if (form.image) {
                await this.$store.dispatch(ActionTypes.UPLOAD_USER_PHOTO, { type: form.image.type, file: form.image.file }).then(() => {
                    return this.$store.dispatch(ActionTypes.UPDATE_USER_PROFILE, { user: form }).then(() => {
                        this.loading = false;
                        this.notifySaved = true;
                    });
                });
            } else {
                await this.$store.dispatch(ActionTypes.UPDATE_USER_PROFILE, { user: form }).then(() => {
                    this.loading = false;
                    this.notifySaved = true;
                });
            }
        },
        async changePassword(form): Promise<void> {
            const user = this.user;
            if (!user) throw new Error(`no user`);

            console.log("form", form);

            this.passwordOk = true;

            const data = {
                userId: user.id,
                oldPassword: form.existing,
                newPassword: form.password,
            };
            await this.$services.api
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
@import "../../scss/mixins";
@import "../../scss/forms";
@import "../../scss/layout";

h2 {
    margin-bottom: 0;

    @include bp-down($xs) {
        margin: 40px 0 20px;
    }
}

::v-deep .user-image img {
    max-width: 400px;
    max-height: 400px;
}
::v-deep .heading {
    font-size: 16px;
    font-weight: 500;
}

#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
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
