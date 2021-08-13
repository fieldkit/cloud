<template>
    <div class="form-container">
        <div class="form-wrap">
            <img class="form-header-logo" :alt="$t('layout.logo.alt')" src="@/assets/FieldKit_Logo_White.png" />
            <form v-if="!created" class="form" @submit.prevent="save">
                <h1 class="form-title">{{ $t("createAccount.form.title") }}</h1>
                <div class="form-group">
                    <TextField v-model="form.name" :label="$t('createAccount.form.name.label')" />

                    <div class="form-errors" v-if="$v.form.name.$error">
                        <div v-if="!$v.form.name.required">{{ $t("createAccount.form.name.required") }}</div>
                    </div>
                </div>
                <div class="form-group">
                    <TextField v-model="form.email" :label="$t('createAccount.form.email.label')" keyboardType="email" />

                    <div class="form-errors" v-if="$v.form.email.$error">
                        <div v-if="!$v.form.email.required">{{ $t("createAccount.form.email.required") }}</div>
                        <div v-if="!$v.form.email.email">{{ $t("createAccount.form.email.valid") }}</div>
                        <div v-if="!$v.form.email.taken">
                            {{ $t("createAccount.form.email.taken") }}
                            <router-link :to="{ name: 'recover' }" class="form-link recover">
                                {{ $t("createAccount.form.recover") }}
                            </router-link>
                        </div>
                    </div>
                </div>
                <div class="form-group">
                    <TextField v-model="form.password" :label="$t('createAccount.form.password.label')" type="password" />

                    <div class="form-errors" v-if="$v.form.password.$error">
                        <div v-if="!$v.form.password.required">{{ $t("createAccount.form.password.required") }}</div>
                        <div v-if="!$v.form.password.min">{{ $t("createAccount.form.password.valid") }}</div>
                    </div>
                </div>
                <div class="form-group">
                    <TextField
                        v-model="form.passwordConfirmation"
                        :label="$t('createAccount.form.passwordConfirmation.label')"
                        type="password"
                    />

                    <div class="form-errors" v-if="$v.form.passwordConfirmation.$error">
                        <div v-if="!$v.form.passwordConfirmation.required">
                            {{ $t("createAccount.form.passwordConfirmation.required") }}
                        </div>
                        <div v-if="!$v.form.passwordConfirmation.sameAsPassword">
                            {{ $t("createAccount.form.passwordConfirmation.valid") }}
                        </div>
                    </div>
                </div>
                <div class="checkbox">
                    <div>
                        <label>
                            <span>
                                {{ $t("createAccount.form.terms.label") }}
                                <router-link to="terms" class="bold">{{ $t("createAccount.form.terms.labelLink") }}</router-link>
                            </span>
                            <input type="checkbox" id="checkbox" v-model="form.tncAccept" />
                            <span class="checkbox-btn"></span>
                        </label>
                    </div>
                </div>
                <div class="form-errors" v-if="$v.form.tncAccept.$error">
                    <div v-if="!$v.form.tncAccept.required">{{ $t("createAccount.form.terms.required") }}</div>
                </div>
                <button class="form-submit" type="submit">{{ $t("createAccount.form.createAccountButton") }}</button>
                <div>
                    <router-link :to="{ name: 'login' }" class="form-link">{{ $t("createAccount.form.backButton") }}</router-link>
                </div>
            </form>
            <form class="form" v-if="created">
                <template v-if="!resending">
                    <img src="@/assets/icon-success.svg" alt="Success" class="form-header-icon" width="57px" />
                    <h1 class="form-title">{{ $t("createAccount.form.created.title") }}</h1>
                    <h2 class="form-subtitle">{{ $t("createAccount.form.created.subtitle") }}</h2>
                    <button class="form-submit" v-on:click="resend">{{ $t("createAccount.form.created.resendButton") }}</button>
                    <router-link :to="{ name: 'login' }" class="form-link">{{ $t("createAccount.form.backButton") }}</router-link>
                </template>
                <template v-if="resending">
                    <img
                        src="@/assets/Icon_Syncing2.png"
                        :alt="$t('createAccount.form.resendMessage')"
                        class="form-header-icon"
                        width="57px"
                    />
                    <p class="form-link">{{ $t("createAccount.form.resendMessage") }}</p>
                </template>
            </form>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { required, email, minLength, sameAs } from "vuelidate/lib/validators";
import FKApi from "@/api/api";

export default Vue.extend({
    name: "CreateAccountView",
    components: {
        ...CommonComponents,
    },
    data(): {
        form: {
            name: string;
            email: string;
            password: string;
            passwordConfirmation: string;
            tncAccept: boolean;
        };
        available: boolean;
        creating: boolean;
        resending: boolean;
        created: { id: number } | null;
    } {
        return {
            form: {
                name: "",
                email: "",
                password: "",
                passwordConfirmation: "",
                tncAccept: false,
            },
            available: true,
            creating: true,
            resending: false,
            created: null,
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
                taken: function (this: any) {
                    return this.available;
                },
            },
            password: { required, min: minLength(10) },
            passwordConfirmation: { required, min: minLength(10), sameAsPassword: sameAs("password") },
            tncAccept: { checked: (value) => value === true },
        },
    },
    methods: {
        async save(): Promise<void> {
            this.$v.form.$touch();
            if (this.$v.form.$pending || this.$v.form.$error) {
                return;
            }

            this.creating = true;

            await this.$services.api
                .register(this.form)
                .then((created) => {
                    this.created = created;
                })
                .catch((error) => {
                    if (error.status === 400) {
                        this.available = false;
                    } else {
                        return this.$seriousError(error);
                    }
                })
                .finally(() => {
                    this.creating = false;
                });
        },
        async resend(): Promise<void> {
            if (!this.created) throw new Error(`nothing to resend`);
            this.resending = true;
            await this.$services.api.resendCreateAccount(this.created.id).then(() => {
                this.resending = false;
            });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms";
@import "../../scss/global";

.checkbox {
    margin-top: 25px;
    span {
        text-align: left;
        float: left;
    }
}
</style>
