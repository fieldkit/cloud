<template>
    <div>
        <div id="sidebar-nav-wide">
            <div class="sidebar-header">
                <router-link :to="{ name: 'projects' }" class="project-link">
                    <img alt="Fieldkit Logo" id="header-logo" src="../assets/FieldKit_Logo_Blue.png" />
                </router-link>
            </div>
            <div id="inner-nav">
                <div class="nav-section">
                    <router-link :to="{ name: 'projects' }">
                        <div class="nav-label">
                            <img alt="Projects" src="../assets/Icon_Projects_blue.png" />
                            <span :class="viewingProjects ? 'selected' : 'unselected'">
                                Projects
                            </span>
                        </div>
                    </router-link>
                    <div v-for="project in projects" v-bind:key="project.id">
                        <router-link :to="{ name: 'viewProject', params: { id: project.id } }" class="project-link">
                            {{ project.name }}
                        </router-link>
                    </div>
                </div>

                <div class="nav-section">
                    <router-link :to="{ name: 'stations' }">
                        <div class="nav-label">
                            <img alt="Stations" src="../assets/Icon_Station_blue.png" />
                            <span :class="viewingStations ? 'selected' : 'unselected'">
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
        </div>
        <div id="sidebar-nav-narrow">
            <div class="sidebar-header">
                <img alt="Fieldkit Compass Logo" src="../assets/compass.png" />
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "SidebarNav",
    props: {
        viewingProjects: { default: false },
        viewingStations: { default: false },
        viewingData: { default: false },
        isAuthenticated: { required: true },
        stations: { required: true },
        projects: { required: true },
    },
    methods: {
        showStation(station) {
            this.$emit("showStation", station);
        },
    },
};
</script>

<style scoped>
#sidebar-nav-wide {
    display: block;
    float: left;
    width: 240px;
    height: 100%;
    min-height: 600px;
    border-right: 1px solid rgba(235, 235, 235, 1);
}
#sidebar-nav-narrow {
    display: none;
    float: left;
    width: 65px;
    height: 100%;
    min-height: 600px;
    border-right: 1px solid rgba(235, 235, 235, 1);
}
#sidebar-nav-narrow img {
    margin-top: 10px;
}
.sidebar-header {
    width: 100%;
    height: 70px;
    float: left;
    border-bottom: 1px solid rgba(235, 235, 235, 1);
}
#header-logo {
    width: 140px;
    margin: 16px auto;
}

#inner-nav {
    width: 92%;
    float: left;
    text-align: left;
    padding-top: 20px;
    padding-left: 15px;
}
.nav-section {
    margin-bottom: 40px;
}
.nav-label {
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
.project-link,
.station-link {
    cursor: pointer;
    font-weight: normal;
    font-size: 14px;
    margin: 0 0 0 30px;
    display: inline-block;
}
</style>
