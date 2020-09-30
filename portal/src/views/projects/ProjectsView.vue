<template>
    <StandardLayout :viewingProjects="true">
        <div class="projects-view">
            <div class="container mine">
                <div class="header">
                    <h1 v-if="isAuthenticated">My Projects</h1>
                    <h1 v-if="!isAuthenticated">Projects</h1>
                    <div id="add-project" v-on:click="addProject" v-if="isAuthenticated">
                        <img alt="Add project" src="@/assets/icon-plus-round.svg" />
                        <span> Add Project </span>
                    </div>
                </div>

                <ProjectThumbnails :projects="userProjects" />
                <ProjectThumbnails :projects="invites.projects" :invited="true" v-if="invites" />
            </div>
            <div class="container community">
                <div class="header">
                    <h1>Community Projects</h1>
                </div>
                <ProjectThumbnails :projects="publicProjects" />
            </div>
        </div>
    </StandardLayout>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import StandardLayout from "../StandardLayout";
import ProjectThumbnails from "./ProjectThumbnails";
import * as ActionTypes from "@/store/actions";
import FKApi from "@/api/api";

export default {
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
    mounted() {
        return new FKApi().getInvitesByUser().then((invites) => {
            this.invites = invites;
        });
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
};
</script>

<style scoped lang="scss">
@import '../../scss/mixins';
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
    border-top: 1px solid #d8dce0;
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
    @include flex(flex-end);

    img {
        margin-right: 7px;
    }
}
#add-project img {
    vertical-align: bottom;
}
</style>
