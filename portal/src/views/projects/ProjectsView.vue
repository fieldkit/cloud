<template>
    <StandardLayout :viewingProjects="true">
        <div class="projects-view">
            <div class="container mine" v-if="userProjects.length > 0">
                <div class="header">
                    <h1 v-if="isAuthenticated">{{ $t("projects.title.mine") }}</h1>
                    <h1 v-if="!isAuthenticated">{{ $t("projects.title.anonymous") }}</h1>
                    <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                        <i class="icon icon-plus-round"></i>
                        <span>{{ $t("projects.add") }}</span>
                    </div>
                </div>

                <ProjectThumbnails :projects="userProjects" />
                <ProjectThumbnails :projects="invites.projects" :invited="true" v-if="invites" />
            </div>
            <div class="container community">
                <div class="header">
                    <h1>{{ $t("projects.title.community") }}</h1>
                </div>
                <ProjectThumbnails :projects="publicProjects" />
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import { mapState, mapGetters } from "vuex";
import StandardLayout from "../StandardLayout";
import ProjectThumbnails from "./ProjectThumbnails";
import { getPartnerForcedLandingPage } from "@/views/shared/partners";

export default Vue.extend({
    name: "ProjectsView",
    components: {
        StandardLayout,
        ProjectThumbnails,
    },
    data() {
        return {
            invites: null,
        };
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({
            userProjects: (s) => Object.values(s.stations.user.projects),
            publicProjects: (s) => Object.values(s.stations.community.projects),
        }),
    },
    async mounted(): Promise<void> {
        const forcedLandingPage = getPartnerForcedLandingPage();
        if (forcedLandingPage != null) {
            await this.$router.push(forcedLandingPage);
            return;
        }
        if (this.isAuthenticated) {
            this.invites = await this.$services.api.getInvitesByUser();
        }
    },
    methods: {
        goBack() {
            if (window.history.length > 1) {
                this.$router.go(-1);
            } else {
                this.$router.push("/");
            }
        },
        addProject() {
            this.$router.push({ name: "addProject" });
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
.projects-view {
    display: flex;
    flex-direction: column;
    padding: 10px 72px 60px;
    text-align: left;

    @include bp-down($lg) {
        padding: 10px 45px 60px;
    }

    @include bp-down($sm) {
        padding: 0 20px 30px;
    }

    @include bp-down($xs) {
        padding: 0 10px 30px;
    }
}

.container.community {
    border-top: 1px solid var(--color-border);
}
.container .header {
    display: flex;
    flex-direction: row;
    align-items: center;
    margin-bottom: 30px;
    margin-top: 40px;

    @include bp-down($lg) {
        margin-bottom: 20px;
        margin-top: 30px;
    }

    @include bp-down($xs) {
        margin-bottom: 25px;
        margin-top: 20px;
    }

    h1 {
        font-size: 36px;
        margin: 0;

        @include bp-down($lg) {
            font-size: 32px;
        }

        @include bp-down($xs) {
            font-size: 24px;
        }
    }
}
#loading {
    width: 100%;
    height: 100%;
    background-color: rgba(255, 255, 255, 0.65);
    text-align: center;
}
.no-user-message {
    float: left;
    font-size: 20px;
    margin: 40px 0 0 40px;
}
.show-link {
    text-decoration: underline;
}
#add-project {
    margin-left: auto;
    cursor: pointer;
    font-size: 16px;
    @include flex(center);

    i {
        margin-right: 7px;
        margin-top: -3px;
    }
}
</style>
