<template>
    <div id="sidebar-nav">
        <div class="header">
            <img alt="Fieldkit Logo" id="header-logo" src="../assets/FieldKit_Logo_White.png" />
        </div>
        <div id="inner-nav">
            <router-link :to="{ name: 'projects' }">
                <div class="nav-label">
                    <img alt="Projects" src="../assets/projects.png" />
                    <span>
                        Projects
                        <div class="selected" v-if="viewingProjects"></div>
                    </span>
                </div>
            </router-link>
            <div v-for="project in projects" v-bind:key="project.id">
                <router-link :to="{ name: 'projectById', params: { id: project.id } }" class="project-link">
                    {{ project.name }}
                </router-link>
            </div>
            <router-link :to="{ name: 'stations' }">
                <div class="nav-label">
                    <img alt="Stations" src="../assets/stations.png" />
                    <span>
                        Stations
                        <div class="selected" v-if="viewingStations"></div>
                    </span>
                </div>
            </router-link>
            <div v-for="station in stations" v-bind:key="station.id">
                <div class="station-link" v-on:click="showStation" :data-id="station.id">
                    {{ station.name }}
                </div>
            </div>

            <!-- <router-link :to="{ name: 'data' }">
                <div class="nav-label">
                    <img alt="Data" src="../assets/data.png" />
                    <span>
                        Data
                        <div class="selected" v-if="viewingData"></div>
                    </span>
                </div>
            </router-link> -->
        </div>
    </div>
</template>

<script>
export default {
    name: "SidebarNav",
    data: () => {
        return {
            viewingProjects: false,
            viewingStations: false,
            viewingData: false
        };
    },
    props: ["viewing", "stations", "projects"],
    mounted() {
        switch (this.viewing) {
            case "projects":
                this.viewingProjects = true;
                break;
            case "stations":
                this.viewingStations = true;
                break;
            case "data":
                this.viewingData = true;
                break;
            default:
                break;
        }
    },
    methods: {
        showStation(event) {
            const id = event.target.getAttribute("data-id");
            const station = this.stations.find(s => {
                return s.id == id;
            });
            this.$emit("showStation", station);
        }
    }
};
</script>

<style scoped>
#sidebar-nav {
    position: absolute;
    top: 0;
    background-color: #1b80c9;
    width: 240px;
    height: 100%;
    min-height: 600px;
    color: white;
}
#sidebar-nav .header {
    border-bottom: 1px solid rgba(255, 255, 255, 0.5);
}
#header-logo {
    width: 140px;
    margin: 16px auto;
}

#inner-nav {
    width: inherit;
    text-align: left;
    padding-top: 20px;
    padding-left: 15px;
}
.nav-label {
    font-weight: bold;
    font-size: 16px;
    margin: 12px 0;
    cursor: pointer;
}
.nav-label img {
    vertical-align: bottom;
}
.selected {
    width: 0;
    height: 0;
    border-top: 15px solid transparent;
    border-bottom: 15px solid transparent;
    border-right: 15px solid white;
    float: right;
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
    margin: 0 0 0 37px;
    display: inline-block;
}
</style>
