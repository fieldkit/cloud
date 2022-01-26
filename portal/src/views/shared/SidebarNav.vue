<template>
    <div class="container-side" v-bind:class="{ active: !sidebar.narrow }">
        <div class="sidebar-header">
            <router-link :to="{ name: 'projects' }">
                <Logo />
            </router-link>
        </div>
        <a class="sidebar-trigger" v-on:click="toggleSidebar">
            <img alt="Menu icon" src="@/assets/icon-menu.svg" width="32" height="22" />
        </a>
        <div id="inner-nav">
            <div class="nav-section">
                <router-link :to="{ name: 'projects' }">
                    <div class="nav-label">
                        <img alt="Projects" src="@/assets/icon-projects.svg" />
                        <span v-bind:class="{ selected: viewingProjects }">{{ $t("layout.side.projects.title") }}</span>
                    </div>
                </router-link>
                <div v-for="project in projects" v-bind:key="project.id">
                    <router-link
                        :to="{ name: 'viewProject', params: { id: project.id } }"
                        class="nav-link"
                        v-bind:class="{ selected: viewingProject && viewingProject.id === project.id }"
                        @click.native="closeMenuOnMobile()"
                    >
                        {{ project.name }}
                    </router-link>
                </div>
            </div>

            <div class="nav-section">
                <router-link :to="{ name: 'mapAllStations' }">
                    <div class="nav-label">
                        <img alt="Stations" src="@/assets/icon-stations.svg" />
                        <span v-bind:class="{ selected: viewingStations }">{{ $t("layout.side.stations.title") }}</span>
                    </div>
                </router-link>
                <div v-for="station in stations" v-bind:key="station.id">
                    <span
                        class="nav-link"
                        v-on:click="showStation(station)"
                        v-bind:class="{ selected: viewingStations && viewingStation && viewingStation.id === station.id }"
                    >
                        {{ station.name }}
                    </span>
                </div>
                <div v-if="isAuthenticated && stations.length == 0" class="nav-link">{{ $t("layout.side.stations.empty") }}</div>
            </div>
        </div>
        <div class="sidebar-header sidebar-compass">
            <router-link :to="{ name: 'projects' }">
                <img :alt="$t('layout.logo.compass.alt')" src="@/assets/logo-compass.svg" width="45" height="45" />
            </router-link>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import Logo from "@/views/shared/Logo.vue";

export default Vue.extend({
    name: "SidebarNav",
    components: {
        Logo,
    },
    props: {
        viewingProject: { type: Object, default: null },
        viewingStation: { type: Object, default: null },
        viewingProjects: { type: Boolean, default: false },
        viewingStations: { type: Boolean, default: false },
        isAuthenticated: { type: Boolean, required: true },
        stations: {
            required: true,
        },
        projects: {
            required: true,
        },
        narrow: {
            type: Boolean,
            default: false,
        },
    },
    mounted(): void {
        const desktopBreakpoint = 1040;
        const windowAny: any = window;
        const resizeObserver = new windowAny.ResizeObserver((entries) => {
            if (entries[0].contentRect.width < desktopBreakpoint) {
                if (!this.sidebar.narrow) {
                    this.sidebar.narrow = true;
                }
            }
        });

        resizeObserver.observe(document.querySelector("body"));
    },
    data(): {
        sidebar: {
            narrow: boolean;
        };
    } {
        return {
            sidebar: {
                narrow: window.screen.availWidth <= 1040,
            },
        };
    },
    methods: {
        showStation(station: unknown): void {
            this.$emit("show-station", station);
            this.closeMenuOnMobile();
        },
        closeMenuOnMobile(): void {
            if (window.screen.availWidth < 1040) {
                this.sidebar.narrow = true;
            }
        },
        toggleSidebar(): void {
            this.sidebar.narrow = !this.sidebar.narrow;
            this.$emit("sidebar-toggle");
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";

.container-side {
    background: #fff;
    width: 65px;
    flex: 0 0 65px;
    transition: all 0.25s;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.28);
    z-index: $z-index-menu;

    &.active {
        width: 240px;
        flex: 0 0 240px;
    }

    @include bp-down($md) {
        width: 0;
        background: #fff;
        height: 100%;
        @include position(fixed, 0 null null 0);
    }
}

#sidebar-nav-narrow img {
    margin-top: 10px;
}
.sidebar-header {
    height: 66px;
    border-bottom: 1px solid rgba(235, 235, 235, 1);
    opacity: 0;
    transition: 0.25s all;
    @include flex(center, center);

    @at-root .container-side.active & {
        opacity: 1;
    }

    @include bp-down($sm) {
        justify-content: flex-start;
        padding: 0 20px;
        height: 54px;
    }

    > a {
        height: 100%;
        display: flex;
    }
}

#inner-nav {
    float: left;
    text-align: left;
    padding: 20px 15px 0;
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.33s ease-in;

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
    font-family: var(--font-family-bold);
    font-size: 16px;
    margin: 12px 0;
    cursor: pointer;
}
.nav-label img {
    vertical-align: sub;
    margin: 0 10px 0 5px;
}
.selected {
    border-bottom: 2px solid $color-primary;
    height: 100%;
    display: inline-block;

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
    }
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

.nav-link {
    cursor: pointer;
    font-family: $font-family-light;
    font-size: 14px;
    margin: 0 0 0 30px;
    display: inline-block;
    line-height: 1.2;

    &.selected {
        padding-bottom: 2px;
    }
}
.sidebar-compass {
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 1;
    transform: translateX(0);
    width: 65px;
    height: 66px;
    @include position(absolute, 0 null null 0);

    @at-root .container-side.active & {
        transition: all 0.33s;
        opacity: 0;
        visibility: hidden;
        transform: translateX(100px);
    }

    @include bp-down($md) {
        display: none;
    }
}
.sidebar-compass img {
    align-self: center;
}
.sidebar-trigger {
    transition: all 0.25s;
    cursor: pointer;
    @include position(absolute, 23px null null 77px);

    @include bp-down($md) {
        left: 10px;
        top: 16px;
    }

    .container-side.active & {
        left: 251px;

        @include bp-down($md) {
            left: 188px;
        }
    }
}
</style>
