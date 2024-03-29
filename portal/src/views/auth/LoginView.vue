<template>
    <div class="form-container" v-if="!sso">
        <ForbiddenBanner v-if="errorMessage" :title="errorMessage" :subtitle="$t('login.loginError')"></ForbiddenBanner>
        <Logo class="form-header-logo"></Logo>
        <LoginForm :forwardAfterQuery="forwardAfterQuery" :spoofing="spoofing" :failed="failed" :busy="busy" @login="save" />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import LoginForm from "./LoginForm.vue";
import Logo from "../shared/Logo.vue";

import { LoginPayload } from "@/api/api";
import { ActionTypes } from "@/store";
import { toSingleValue } from "@/utilities";
import ForbiddenBanner from "@/views/shared/ForbiddenBanner.vue";
import { getPartnerCustomizationWithDefault } from "@/views/shared/partners";

export default Vue.extend({
    components: {
        ...CommonComponents,
        LoginForm,
        Logo,
        ForbiddenBanner,
    },
    props: {
        spoofing: {
            type: Boolean,
            default: false,
        },
        errorMessage: {
            type: String,
            required: false,
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
            let params;

            if (this.$route.query.params) {
                params = JSON.parse(this.$route.query.params as string);
            }

            if (after.after) {
                await this.$router.push({ path: after.after, query: params }).catch((e) => {
                    console.log(e);
                });
            } else {
                await this.$router.push({ name: this.partnerCustomization().routeAfterLogin }).catch((e) => {
                    console.log(e);
                });
            }
        },
        partnerCustomization() {
            return getPartnerCustomizationWithDefault();
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/forms.scss";
</style>
