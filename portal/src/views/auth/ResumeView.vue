<template>
    <div></div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import { ResumeAction, LoginOidcAction } from "@/store";
import { toSingleValue } from "@/utilities";

export default Vue.extend({
    components: {},
    props: {},
    data(): {} {
        return {};
    },
    async mounted(): Promise<void> {
        const p = this.$route.query;
        if (this.$route.params.token) {
            await this.$store.dispatch(new ResumeAction(this.$route.params.token));
            await this.leaveAfterAuth();
        }

        const state = toSingleValue(this.$route.query.state);
        const sessionState = toSingleValue(this.$route.query.session_state);
        const code = toSingleValue(this.$route.query.code);
        if (state && sessionState && code) {
            await this.$store.dispatch(new LoginOidcAction(this.$route.params.token, state, sessionState, code));
            await this.leaveAfterAuth();
        }
    },
    methods: {
        forwardAfterQuery(): { after?: string } {
            const after = toSingleValue(this.$route.query.after);
            if (after) {
                return { after: after };
            }
            return {};
        },
        async leaveAfterAuth(): Promise<void> {
            const after = this.forwardAfterQuery();
            if (after.after) {
                await this.$router.push(after.after);
            } else {
                await this.$router.push({ name: "projects" });
            }
        },
    },
});
</script>

<style scoped lang="scss"></style>
