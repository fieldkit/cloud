<template>
    <div class="form-container" v-if="credentialsNecessary">
        <Logo class="form-header-logo"></Logo>
        <LoginForm :forwardAfterQuery="forwardAfterQuery" :spoofing="false" :busy="busy" @login="save" />
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";
import LoginForm from "./LoginForm.vue";

import FKApi, { LoginPayload } from "@/api/api";
import { ActionTypes, DiscourseParams, LoginDiscourseAction } from "@/store";
import { toSingleValue } from "@/utilities";
import Logo from "@/views/shared/Logo.vue";

export default Vue.extend({
    components: {
        ...CommonComponents,
        LoginForm,
        Logo,
    },
    props: {},
    data(): {
        busy: boolean;
        failed: boolean;
        credentialsNecessary: boolean;
    } {
        return {
            busy: false,
            failed: false,
            credentialsNecessary: false,
        };
    },
    computed: {
        forwardAfterQuery(): { after?: string } {
            const after = toSingleValue(this.$route.query.after);
            if (after) {
                return { after: after };
            }
            return {};
        },
        discourseParams(): DiscourseParams | null {
            const sso = toSingleValue(this.$route.query.sso);
            const sig = toSingleValue(this.$route.query.sig);
            if (sso && sig) {
                return new DiscourseParams(sso, sig);
            }
            return null;
        },
    },
    async mounted(): Promise<void> {
        const params = this.discourseParams;
        const token = this.$state.user.token;
        if (params && token) {
            await this.$store
                .dispatch(new LoginDiscourseAction(token, null, params))
                .then(
                    () => {
                        // NOTE The action handler will navigate away. So nothign should happen after.
                    },
                    () => {
                        this.failed = true;
                        this.credentialsNecessary = true;
                    }
                )
                .finally(() => {
                    this.busy = false;
                });
        } else {
            this.credentialsNecessary = true;
        }
    },
    methods: {
        async save(payload: LoginPayload): Promise<void> {
            const params = this.discourseParams;
            if (!params) {
                this.failed = true;
                return;
            }

            this.busy = true;
            this.failed = false;

            await this.$store
                .dispatch(new LoginDiscourseAction(null, payload, params))
                .then(
                    () => {
                        // NOTE The action handler will navigate away. So nothign should happen after.
                    },
                    () => {
                        this.failed = true;
                    }
                )
                .finally(() => {
                    this.busy = false;
                });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms.scss";
</style>
