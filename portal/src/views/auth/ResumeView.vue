<template>
    <div></div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import { ResumeAction } from "@/store";

export default Vue.extend({
    components: {},
    props: {},
    data(): {} {
        return {};
    },
    async mounted(): Promise<void> {
        if (this.$route.params.token) {
            await this.$store.dispatch(new ResumeAction(this.$route.params.token));
            await this.leaveAfterAuth();
        }
    },
    methods: {
        forwardAfterQuery(): { after?: string } {
            const after = this.$route.query.after;
            if (after) {
                if (_.isArray(after) && after.length > 0 && after[0]) {
                    return { after: after[0] };
                }
                return { after: after as string };
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
