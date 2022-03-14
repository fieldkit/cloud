<template>
    <div class="project-public project-container" v-if="project">

        <div class="project-detail-card">
            <div class="photo-container">
                <ProjectPhoto :project="project" :image-size="150"/>
            </div>
            <div class="detail-container">
                <h3 class="detail-title">{{ project.name }}</h3>
                <div class="detail-description">{{ project.description }}</div>
                <router-link :to="{ name: 'viewProject' }" class="link">Project Dashboard ></router-link>
            </div> 
        </div>

        <StationsMap
            @show-summary="showSummary"
            :mapped="mappedProject"
            :layoutChanges="layoutChanges   "
            :showStations="project.showStations"
            :mapBounds="mapBounds"
        />


    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { mapGetters } from "vuex";
import * as utils from "../../utilities";
import { ProjectModule, DisplayStation, Project, DisplayProject, MappedStations, BoundingRectangle } from "@/store";
import StationsMap from "../shared/StationsMap.vue";
import CommonComponents from "@/views/shared";
import { twitterCardMeta } from "@/social";

export default Vue.extend({
    name: "ProjectBigMap",
    components: {
        ...CommonComponents,
        StationsMap,
    },
    data(): {
        layoutChanges: number;
    } {
        return {
            layoutChanges: 0,
        };
    },
    metaInfo() {
        return {
            meta: twitterCardMeta(this.displayProject),
            afterNavigation() {
                console.log("hello: after-navigation");
            },
        };
    },
    props: {
        user: {
            required: true,
        },
        userStations: {
            type: Array as PropType<DisplayStation[]>,
            required: true,
        },
        displayProject: {
            type: Object as PropType<DisplayProject>,
            required: true,
        },
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated", isBusy: "isBusy" }),
        project(): Project {
            return this.displayProject.project;
        },
        projectStations(): DisplayStation[] {
            return this.$getters.projectsById[this.displayProject.id].stations;
        },
        mappedProject(): MappedStations | null {
            return this.$getters.projectsById[this.project.id].mapped;
        },
        projectModules(): { name: string; url: string }[] {
            return this.$getters.projectsById[this.displayProject.id].modules.map((m) => {
                return {
                    name: m.name,
                    url: this.getModuleImg(m),
                };
            });
        },
        mapBounds(): BoundingRectangle {
            if (this.project.bounds?.min && this.project.bounds?.max) {
                return new BoundingRectangle(this.project.bounds?.min, this.project.bounds?.max);
            }

            return MappedStations.defaultBounds();
        },
    },
    mounted () {
        //console.log("IMAGE SIZE", this.imageSize)
    },
    methods: {
        getModuleImg(module: ProjectModule): string {
            return this.$loadAsset(utils.getModuleImg(module));
        },
        getTeamHeading(): string {
            // TODO i18n
            const members = this.displayProject.users.length == 1 ? "member" : "members";
            return "Project Team (" + this.displayProject.users.length + " " + members + ")";
        },
        showSummary() {
            console.log("SHOW SUMMARY");
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/project";
@import "../../scss/global";

.project-detail-card {
    display: flex;
    border: 1px solid var(--color-border);
    padding: 1px;
    border-radius: 3px;
    position: relative;
    z-index: 50;
    width: 349px;
    position: absolute;
    top: 95px;
    right: 28px;
    box-sizing: border-box;
    background-color: #ffffff;

    @include bp-down($sm) {
        flex: 0 0 calc(50% - 18px);
    }

    @include bp-down($xs) {
        flex: 0 0 100%;
        margin: 0 0 10px;
    }
    .link {
        color: $color-fieldkit-primary;
        font-size: 12px;
    }
}
.detail-title {
    font-family: $font-family-bold;
    font-size: 18px;
    margin-top: 15px;
    margin-bottom: 5px;
}
.detail-containiner {
    width: 240px;
}
.detail-description {
    font-family: $font-family-light;
    font-size: 14px;
}
.photo-container{
    width: 95px;
}

.project-public {
    display: flex;
    flex-direction: column;

    @include bp-down($sm) {
        padding-bottom: 20px;
    }
}
.header {
    display: flex;
    flex-direction: row;
}
.header .left {
    margin-right: auto;
    display: flex;
    flex-direction: column;
}
.header .right {
    margin-left: auto;
}
.header .project-name {
    font-size: 24px;
    font-family: var(--font-family-bold);
    margin: 0 15px 0 0;
    display: inline-block;
}
.header .project-dashboard {
    font-size: 20px;
    font-family: var(--font-family-bold);
    margin: 0 15px 0 0;
    display: inline-block;
    margin-top: 10px;
    margin-bottom: 20px;
}

.details {
    display: flex;
    border-radius: 2px;
    border: solid 1px var(--color-border);
    background-color: white;

    @include bp-down($sm) {
        flex-wrap: wrap;
    }
}

.details .project-detail {
    font-family: $font-family-light;
    margin-bottom: 9px;
    line-height: 1.5;
    overflow-wrap: anywhere;
}

.photo-container {
    width: 75px;
    height: 75px;
    margin: 10px;

    @include bp-down($xs) {
        width: 100%;
    }
}

::v-deep .project-image {
    width: 100%;
    height: auto;
}

.project-container {
    margin-top: -10px;
    width: 100% !important;
    margin: 0;

    @include bp-down($sm) {
        margin-top: -28px;
    }
}

.recent-activity {
    border-radius: 2px;
    border: solid 1px var(--color-border);
    background-color: #ffffff;
    padding: 17px 23px;
    flex: 1;

    @include bp-down($sm) {
        padding: 19px 10px;
    }

    h1 {
        margin: 0 0 23px;
        font-size: 20px;
        font-weight: 500;

        @include bp-down($sm) {
            font-size: 18px;
        }
    }

    h2 {
        margin: 0;
        font-size: 16px;
        font-weight: 500;

        + span {
            margin-left: auto;
            padding-left: 10px;
            font-weight: 300;
            font-size: 14px;
            color: #6a6d71;
        }
    }

    li {
        @include flex(flex-start);
        font-size: 16px;
        margin-bottom: 30px;

        > div {
            display: flex;
            flex-direction: column;
            flex: 1;
            line-height: 1.3;

            > div {
                display: flex;
            }
        }

        img {
            margin-right: 16px;
            margin-top: 2px;
        }

        p {
            margin-top: 5px;
            font-size: 14px;
        }
    }
}



</style>
