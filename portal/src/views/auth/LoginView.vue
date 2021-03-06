<template>
    <div class="form-container" v-if="!sso">
        <img class="form-header-logo" alt="FieldKit Logo" src="@/assets/FieldKit_Logo_White.png" />
        <LoginForm :forwardAfterQuery="forwardAfterQuery" :spoofing="spoofing" :failed="failed" @login="save" />
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";
import LoginForm from "./LoginForm.vue";

import FKApi, { LoginPayload } from "@/api/api";
import { ActionTypes } from "@/store";
import { toSingleValue } from "@/utilities";

export default Vue.extend({
    components: {
        ...CommonComponents,
        LoginForm,
    },
    props: {
        spoofing: {
            type: Boolean,
            default: false,
        },
    },
    data(): {
        busy: boolean;
        failed: boolean;
        sso: boolean;
    } {
        return {
            busy: false,
            failed: false,
            sso: this.$config.sso,
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
    },
    async mounted(): Promise<void> {
        if (!this.spoofing) {
            if (this.$config.sso) {
                try {
                    const url = await this.$services.api.loginUrl(null);
                    window.location.replace(url);
                } catch (error) {
                    console.log(`sso-url error: ${error}`);
                    this.sso = false;
                }
            }
        } else {
            this.sso = false;
        }
    },
    methods: {
        async save(payload: LoginPayload): Promise<void> {
            this.busy = true;
            this.failed = false;

            await this.$store
                .dispatch(ActionTypes.LOGIN, payload)
                .then(
                    async () => {
                        await this.leaveAfterAuth();
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
