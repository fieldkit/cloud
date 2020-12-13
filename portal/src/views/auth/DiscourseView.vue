<template>
    <div class="form-container">
        <img class="form-header-logo" alt="FieldKit Logo" src="@/assets/FieldKit_Logo_White.png" />
        <LoginForm :forwardAfterQuery="forwardAfterQuery" :spoofing="false" @login="save" />
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";
import LoginForm from "./LoginForm.vue";

import FKApi, { LoginPayload } from "@/api/api";
import { ActionTypes, DiscourseParams, DiscourseLoginAction } from "@/store";
import { toSingleValue } from "@/utilities";

export default Vue.extend({
    components: {
        ...CommonComponents,
        LoginForm,
    },
    props: {},
    data(): {
        busy: boolean;
        failed: boolean;
    } {
        return {
            busy: false,
            failed: false,
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
                .dispatch(new DiscourseLoginAction(token, null, params))
                .then(
                    async () => {
                        // await this.leaveAfterAuth();
                    },
                    () => (this.failed = true)
                )
                .finally(() => {
                    this.busy = false;
                });
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
                .dispatch(new DiscourseLoginAction(null, payload, params))
                .then(
                    async () => {
                        // await this.leaveAfterAuth();
                    },
                    () => (this.failed = true)
                )
                .finally(() => {
                    this.busy = false;
                });
        },
        async leaveAfterAuth(): Promise<void> {
            const after = this.forwardAfterQuery;
            if (after.after) {
                await this.$router.push(after.after);
            } else {
                await this.$router.push({ name: "projects" });
            }
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms.scss";
</style>