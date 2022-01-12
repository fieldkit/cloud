<template>
    <StandardLayout>
        <div class="container">
            <div class="header">
                <h1 v-if="isAuthenticated">{{ $t("projects.title.mine") }}</h1>
                <h1 v-if="!isAuthenticated">{{ $t("projects.title.anonymous") }}</h1>
                <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                    <img alt="Add project" src="@/assets/icon-plus-round.svg" />
                    <span>{{ $t("projects.add") }}</span>
                </div>
            </div>

            <ProjectThumbnails :projects="userProjects" />
            <ProjectThumbnails :projects="invites.projects" :invited="true" v-if="invites" />
        </div>
    </StandardLayout>
</template>

<script>
import Vue from "vue";
import { mapState, mapGetters } from "vuex";
import StandardLayout from "./StandardLayout";
import * as ActionTypes from "@/store/actions";

export default Vue.extend({
    name: "ProjectsView",
    components: {
        StandardLayout,
    },
    data() {
        return {};
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({
            userProjects: (s) => Object.values(s.stations.user.projects),
            publicProjects: (s) => Object.values(s.stations.community.projects),
        }),
    },
    mounted() {
        return this.$services.api.getInvitesByUser().then((invites) => {
            this.invites = invites;
        });
    },
    methods: {},
});
</script>

<style scoped lang="scss">
@import "../scss/mixins";

.container {
  text-align: left;
}

h1 {
    font-size: 36px;
}
</style>
