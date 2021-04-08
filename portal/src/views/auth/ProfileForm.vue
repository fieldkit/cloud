<template>
    <div class="profile-change">
        <form @submit.prevent="saveForm">
            <h3 class="heading">{{ $t("user.profile.photo.title") }}</h3>
            <div class="user-image image-container">
                <ImageUploader :image="{ url: null }" v-if="!user.photo" @change="onImage" />
                <ImageUploader :image="{ url: user.photo.url }" v-if="user.photo" @change="onImage" />
            </div>
            <div>
                <TextField v-model="form.name" :label="$t('user.profile.form.name.label')" />

                <div class="validation-errors" v-if="$v.form.name.$error">
                    <div v-if="!$v.form.name.required">{{ $t("user.profile.form.name.required") }}</div>
                </div>
            </div>
            <div>
                <TextField v-model="form.email" :label="$t('user.profile.form.email.label')" />

                <div class="validation-errors" v-if="$v.form.email.$error">
                    <div v-if="!$v.form.email.required">{{ $t("user.profile.form.email.required") }}</div>
                    <div v-if="!$v.form.email.email">{{ $t("user.profile.form.email.required") }}</div>
                    <div v-if="!$v.form.email.taken">
                        {{ $t("user.profile.form.email.taken") }}
                        <router-link :to="{ name: 'recover' }" class="recover-link">
                            {{ $t("user.profile.form.email.recover") }}
                        </router-link>
                    </div>
                </div>
            </div>
            <div>
                <TextField v-model="form.bio" :label="$t('user.profile.form.bio.label')" />

                <div class="validation-errors" v-if="$v.form.name.$error">
                    <div v-if="!$v.form.bio.required">{{ $t("user.profile.form.bio.required") }}</div>
                </div>
            </div>
            <div class="checkbox" v-show="false">
                <label>
                    {{ $t("user.profile.form.makePublic") }}
                    <input type="checkbox" id="checkbox" v-model="form.publicProfile" />
                    <span class="checkbox-btn"></span>
                </label>
            </div>
            <button class="button-solid" type="submit">{{ $t("user.profile.form.update") }}</button>
        </form>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { required, email, minLength, sameAs } from "vuelidate/lib/validators";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { DisplayStation } from "@/store";

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
    data(): {
        loading: boolean;
        form: {
            id: number;
            name: string;
            email: string;
            bio: string;
            image: any;
            publicProfile: boolean;
        };
        changePasswordForm: {
            existing: string;
            password: string;
            passwordConfirmation: string;
        };
        notifySaved: boolean;
        notifyPasswordChanged: boolean;
        emailAvailable: boolean;
    } {
        return {
            loading: false,
            form: {
                id: this.user.id,
                name: this.user.name,
                email: this.user.email,
                bio: this.user.bio,
                image: null,
                publicProfile: false,
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
        goBack(): void {
            window.history.length > 1 ? this.$router.go(-1) : this.$router.push("/");
        },
        showStation(station: DisplayStation): Promise<any> {
            return this.$router.push({ name: "mapStation", params: { id: String(station.id) } });
        },
        saveForm(): void {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                console.log("save form, validation error");
                return;
            }

            this.$emit("save", this.form);
        },
        changePassword(): void {
            this.$v.changePasswordForm.$touch();
            if (this.$v.changePasswordForm.$pending || this.$v.changePasswordForm.$error) {
                console.log("password form, validation error");
                return;
            }

            this.$emit("change-password", this.changePasswordForm);
        },
        onImage(image): void {
            this.form.image = image;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms";
@import "../../scss/global";

.image-container {
    margin-bottom: 33px;

    @include bp-down($xs) {
        margin-bottom: 40px;
    }
}
::v-deep .user-image img {
    max-width: 198px !important;
    max-height: 198px !important;
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
