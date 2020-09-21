<template>
    <div class="container-side" v-bind:class="{ active: !narrow }">
        <div class="sidebar-header">
            <router-link :to="{ name: 'projects' }">
                <img alt="FieldKit Logo" id="header-logo" src="@/assets/logo-fieldkit.svg" />
            </router-link>
        </div>
        <div id="inner-nav">
            <div class="nav-section">
                <router-link :to="{ name: 'projects' }">
                    <div class="nav-label">
                        <img alt="Projects" src="@/assets/icon-projects.svg" />
                        <span v-bind:class="{ selected: viewingProjects }">
                            Projects
                        </span>
                    </div>
                </router-link>
                <div v-for="project in projects" v-bind:key="project.id">
                    <router-link
                        :to="{ name: 'viewProject', params: { id: project.id } }"
                        class="project-link"
                        v-bind:class="{ selected: viewingProject && viewingProject.id === project.id }"
                        @click.native="closeMenuOnMobile()"
                    >
                        {{ project.name }}
                    </router-link>
                </div>
            </div>

            <div class="nav-section">
                <router-link :to="{ name: 'stations' }">
                    <div class="nav-label">
                        <img alt="Stations" src="@/assets/icon-stations.svg" />
                        <span v-bind:class="{ selected: viewingStations }">
                            Stations
                        </span>
                    </div>
                </router-link>
                <div v-for="station in stations" v-bind:key="station.id">
                    <div class="station-link" v-on:click="showStation(station)">
                        {{ station.name }}
                    </div>
                </div>
                <div v-if="isAuthenticated && stations.length == 0" class="station-link">
                    No stations added
                </div>
            </div>
        </div>
        <div class="sidebar-header sidebar-compass">
            <img alt="FieldKit Compass Logo" src="@/assets/logo-compass.svg" width="45" height="45" />
        </div>
    </div>
</template>

<script>
export default {
    name: "SidebarNav",
    props: {
        viewingProject: { type: Object, default: null },
        viewingProjects: { default: false },
        viewingStations: { default: false },
        isAuthenticated: { required: true },
        stations: { required: true },
        projects: { required: true },
        narrow: {
            type: Boolean,
            default: false,
        },
    },
    methods: {
        showStation(station) {
            this.$emit("show-station", station);
        },
        closeMenuOnMobile() {
            if (window.screen.availWidth < 1040) {
                this.$emit("toggle-menu", true);
            }
        }
    },
};
</script>

<style scoped lang="scss">
@import '../../scss/mixins';

.container-side {
    width: 65px;
    flex: 0 0 65px;
    transition: all 0.33s;
    overflow: hidden;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.28);
    z-index: $z-index-menu;

    &.active {
        width: 240px;
        flex: 0 0 240px;
    }

    @include bp-down($md) {
        width: 0;
        z-index: $z-index-top;
        background: #fff;
        height: 100%;
        @include position(fixed, 0 null null 0);
    }
}

#sidebar-nav-narrow img {
    margin-top: 10px;
}
.sidebar-header {
    height: 70px;
    border-bottom: 1px solid rgba(235, 235, 235, 1);
    opacity: 0;
    @include flex(center, center);

    @at-root .container-side.active & {
       opacity: 1;
    }

    @include bp-down($md) {
        justify-content: flex-start;
        padding: 0 20px;
    }
}
#header-logo {
    width: 140px;
    margin: 16px auto;
}

#inner-nav {
    float: left;
    text-align: left;
    padding: 20px 15px 0;
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.33s;

    @at-root .container-side.active & {
        opacity: 1;
        visibility: visible;
        width: 200px;
    }
}
.nav-section {
    margin-bottom: 40px;

    > div {
        padding: 4px 0;
    }
}
.nav-label {
    @include flex(center);
    font-weight: bold;
    font-size: 16px;
    margin: 12px 0;
    cursor: pointer;
}
.nav-label img {
    vertical-align: sub;
    margin: 0 10px 0 5px;
}
.selected {
    border-bottom: 2px solid #1b80c9;
    height: 100%;
    display: inline-block;
}
.unselected {
    display: inline-block;
}
.small-nav-text {
    font-weight: normal;
    font-size: 13px;
    margin: 20px 0 0 37px;
    display: inline-block;
}
.project-link {
    line-height: 1.2;
}
.project-link,
.station-link {
    cursor: pointer;
    font-weight: normal;
    font-size: 14px;
    margin: 0 0 0 30px;
    display: inline-block;
}
.sidebar-compass {
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 1;
    transform: translateX(0);
    width: 65px;
    @include position(absolute, 0 null null 0);

    @at-root .container-side.active & {
        transition: all 0.33s;
        opacity: 0;
        transform: translateX(100px);
    }
}
.sidebar-compass img {
    align-self: center;
}
</style>
