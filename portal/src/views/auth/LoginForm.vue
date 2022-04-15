<template>
    <form class="form" @submit.prevent="save">
        <h1 class="form-title">Log In to Your Account</h1>
        <div class="form-group" v-if="spoofing">
            <TextField v-model="form.spoofEmail" :label="$t('login.form.spoofEmail.label')" />
            <div class="form-errors" v-if="$v.form.spoofEmail.$error">
                <div v-if="!$v.form.spoofEmail.required">{{ $t("login.form.spoofEmail.required") }}</div>
                <div v-if="!$v.form.spoofEmail.email">{{ $t("login.form.spoofEmail.valid") }}</div>
            </div>
        </div>
        <div class="form-group">
            <TextField v-model="form.email" :label="$t('login.form.email.label')" />
            <div class="form-errors" v-if="$v.form.email.$error">
                <div v-if="!$v.form.email.required">{{ $t("login.form.email.required") }}</div>
                <div v-if="!$v.form.email.email">{{ $t("login.form.email.valid") }}</div>
            </div>
        </div>
        <div class="form-group">
            <TextField v-model="form.password" :label="$t('login.form.password.label')" type="password" />
            <div class="form-errors" v-if="$v.form.password.$error">
                <div v-if="!$v.form.password.required">{{ $t("login.form.password.required") }}</div>
                <div v-if="!$v.form.password.min">{{ $t("login.form.password.valid") }}</div>
            </div>
        </div>
        <div class="form-link-recover">
            <router-link :to="{ name: 'recover', query: forwardAfterQuery }" class="form-link">{{ $t("login.resetLink") }}</router-link>
        </div>
        <div v-if="failed" class="login-failed">
            {{ $t("login.loginFailed") }}
        </div>
        <button class="form-submit" type="submit">
            <div class="loading-spinner-wrap">
                <Spinner class="loading-spinner" v-if="busy" />
            </div>
            <template v-if="!busy">
                {{ $t("login.loginButton") }}
            </template>
        </button>
        <div>
            <router-link :to="{ name: 'register', query: forwardAfterQuery }" class="form-link">
                {{ $t("login.createAccountLink") }}
            </router-link>
        </div>
    </form>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs, requiredIf } from "vuelidate/lib/validators";

import FKApi, { LoginPayload } from "@/api/api";
import { ActionTypes } from "@/store";

export default Vue.extend({
    name: "LoginForm",
    components: {
        ...CommonComponents,
    },
    props: {
        spoofing: {
            type: Boolean,
            default: false,
        },
        failed: {
            type: Boolean,
            default: false,
        },
        forwardAfterQuery: {
            type: Object as PropType<{ after?: string }>,
            default: false,
        },
        busy: {
            type: Boolean,
            default: false,
        },
    },
    data(): {
        form: {
            spoofEmail: string;
            email: string;
            password: string;
        };
    } {
        return {
            form: {
                spoofEmail: "",
                email: "",
                password: "",
            },
        };
    },
    validations() {
        return {
            form: {
                spoofEmail: {
                    // eslint-disable-next-line
                    required: requiredIf(function(this: any) {
                        return this.spoofing;
                    }),
                    email,
                },
                email: {
                    required,
                    email,
                },
                password: { required, min: minLength(10) },
            },
        };
    },
    methods: {
        createPayload(): LoginPayload {
            if (this.spoofing) {
                return new LoginPayload(this.form.spoofEmail, this.form.email + " " + this.form.password);
            }
            return new LoginPayload(this.form.email, this.form.password);
        },
        async save(): Promise<void> {
            this.$v.form.$touch();
            if (this.busy || this.$v.form.$pending || this.$v.form.$error) {
                return;
            }
            const payload = this.createPayload();
            this.$emit("login", payload);
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms.scss";

.form-submit {
    position: relative;
}

.loading-spinner {
    width: 20px; 
    height: 20px;
    margin: 0 10px 0;
    background: linear-gradient(to right, #fff 10%, rgba(0, 0, 0, 0) 42%);

    &-wrap {
        @include position(absolute, 50% null null 50%);
        @include flex(center, center);
        transform: translate(-50%, -50%);
    }

    &:before {
        background-color: #fff;
    }

    &:after {
        background-color: var(--color-secondary);
    }
}
</style>
