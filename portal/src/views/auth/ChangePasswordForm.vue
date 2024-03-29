<template>
    <form @submit.prevent="saveForm">
        <h3 class="heading">Change Password</h3>
        <div>
            <TextField v-model="form.existing" :label="$t('user.profile.form.password.existing.label')" type="password" />

            <div class="validation-errors" v-if="$v.form.existing.$error">
                <div v-if="!$v.form.existing.required">{{ $t("user.profile.form.password.existing.required") }}</div>
                <div v-if="!$v.form.existing.min">{{ $t("user.profile.form.password.existing.min") }}</div>
            </div>
        </div>
        <div>
            <TextField v-model="form.password" :label="$t('user.profile.form.password.password.label')" type="password" />

            <div class="validation-errors" v-if="$v.form.password.$error">
                <div v-if="!$v.form.password.required">{{ $t("user.profile.form.password.password.required") }}</div>
                <div v-if="!$v.form.password.min">{{ $t("user.profile.form.password.password.min") }}</div>
            </div>
        </div>
        <div>
            <TextField v-model="form.passwordConfirmation" :label="$t('user.profile.form.password.confirmation.label')" type="password" />

            <div class="validation-errors" v-if="$v.form.passwordConfirmation.$error">
                <div v-if="!$v.form.passwordConfirmation.required">{{ $t("user.profile.form.password.confirmation.required") }}</div>
                <div v-if="!$v.form.passwordConfirmation.sameAsPassword">
                    {{ $t("user.profile.form.password.confirmation.sameAsPassword") }}
                </div>
            </div>
        </div>
        <button class="button-solid" type="submit">{{ $t("user.profile.form.password.update") }}</button>
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

<style scoped lang="scss">
@import "../../scss/global";
@import "../../scss/mixins";

.main-panel {
    display: flex;
    flex-direction: column;
    max-width: 700px;
    padding: 20px;
}
.heading {
    font-family: var(--font-family-bold);
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

.notification {
    margin: 20px;
    padding: 20px;
    background-color: #d4edda;
    border: 2px;
    border-radius: 4px;
}

form {
    margin-top: 45px;
    padding-bottom: 80px;
}

.button-solid {
    margin-top: 15px;
    margin-bottom: 20px;

    @include bp-down($xs) {
        width: 100%;
    }
}
</style>
