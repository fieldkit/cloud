<template>
    <div class="project-public project-container" v-if="project">
        <div class="project-detail-card">
            <div class="photo-container">
                <ProjectPhoto :project="project" :image-size="150" />
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
            :layoutChanges="layoutChanges"
            :showStations="project.showStations"
            :mapBounds="mapBounds"
        />
        <StationSummary
            v-if="activeStation"
            :station="activeStation"
            :readings="false"
            :exploreContext="exploreContext"
            @close="onCloseSummary"
            v-bind:key="activeStation.id"
        />
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { mapGetters } from "vuex";
import * as utils from "../../utilities";
import { ProjectModule, DisplayStation, Project, DisplayProject, MappedStations, BoundingRectangle } from "@/store";
import StationsMap from "../shared/StationsMap.vue";
import StationSummary from "@/views/shared/StationSummary.vue";
import CommonComponents from "@/views/shared";

import { ExploreContext } from "@/views/viz/common";

export default Vue.extend({
    name: "ProjectBigMap",
    components: {
        ...CommonComponents,
        StationsMap,
        StationSummary,
    },
    data(): {
        layoutChanges: number;
        activeStationId: number | null;
    } {
        return {
            layoutChanges: 0,
            activeStationId: null,
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
        activeStation(): DisplayStation | null {
            if (this.activeStationId) {
                return this.$getters.stationsById[this.activeStationId];
            }
            return null;
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
        exploreContext(): ExploreContext {
            return new ExploreContext(this.project.id);
        },
    },
    methods: {
        getModuleImg(module: ProjectModule): string {
            return this.$loadAsset(utils.getModuleImg(module));
        },
        showSummary(station: DisplayStation): void {
            console.log("showSummay", station);
            this.activeStationId = station.id;
        },
        onCloseSummary(): void {
            this.activeStationId = null;
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
        width: 100%;
        position: fixed;
        border-top-right-radius: 10px;
        border-top-left-radius: 10px;
        bottom: 0px;
        text-align: center;
        padding-bottom: 10px;
        top: auto;
        right: auto;
        align-items: center;
        justify-content: center;
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
.detail-container {
    width: 75%;
}
.detail-description {
    font-family: $font-family-light;
    font-size: 14px;
}

.details {
    display: flex;
    border-radius: 2px;
    border: solid 1px var(--color-border);
    background-color: white;
}

.photo-container {
    width: 75px;
    height: 75px;
    margin: 10px;
    float: left;

    @include bp-down($sm) {
        display: none;
    }
}

.project-container {
    margin-top: -10px;
    width: 100%;
    height: 100%;
    margin: 0;
}
#station-summary {
    position: fixed;
    left: 100px;
    right: 100px;
}
</style>
