<template>
    <div></div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import { LoginOidcAction } from "@/store";
import { toSingleValue } from "@/utilities";

export default Vue.extend({
    components: {},
    props: {},
    data(): {} {
        return {};
    },
    async mounted(): Promise<void> {
        console.log(`route:`, this.$route);
        console.log(`query:`, this.$route.query);
        const p = this.$route.query;
        const state = toSingleValue(this.$route.query.state);
        const sessionState = toSingleValue(this.$route.query.session_state);
        const code = toSingleValue(this.$route.query.code);
        if (state && sessionState && code) {
            await this.$store.dispatch(new LoginOidcAction(this.$route.params.token, state, sessionState, code));
            // await this.leaveAfterAuth();
        }
    },
    methods: {},
});
</script>

<style scoped lang="scss"></style>
