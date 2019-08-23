<template>
    <div>
        <div id="white-header" class="header">
            <div class="user-name">{{ this.user.name }}</div>
            <div class="log-out" v-on:click="logout">Log out</div>
        </div>

        <div id="sidebar-nav">
            <div class="header">
                <img alt="Fieldkit Logo" id="header-logo" src="../assets/FieldKit_Logo_White.png" />
            </div>
            <div id="inner-nav">
                <div class="nav-label">
                    <div class="selected"></div>
                    <img alt="Projects" src="../assets/projects.png" />
                    <span>Projects</span>
                    <span class="small-nav-text">Default FieldKit Project</span>
                </div>
                <div class="nav-label">
                    <img alt="Stations" src="../assets/stations.png" />
                    Stations
                </div>
                <div class="nav-label">
                    <img alt="Data" src="../assets/data.png" />
                    Data
                </div>
            </div>
        </div>

        <div class="dashboard">
            <div class="section">
                <h1>My Projects</h1>
                <img alt="Default Fieldkit Project" src="../assets/default.png" />
            </div>
            <div class="section">
                <h1>Community</h1>
            </div>
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";

export default {
    name: "dashboard",
    data: () => {
        return { user: {}, stations: {} };
    },
    async beforeCreate() {
        const api = new FKApi();
        if (api.authenticated()) {
            this.user = await api.getCurrentUser();
            this.stations = await api.getStations();
            // console.log("this is the user info", this.user);
            // console.log("this is the station info", this.stations);
        }
    },
    methods: {
        logout() {
            const api = new FKApi();
            api.logout();
            this.$router.push({ name: "login" });
        }
    }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#white-header {
    color: gray;
    text-align: right;
}
.user-name {
    padding: 12px 20px 0 0;
}
.log-out {
    padding: 0 20px 0 0;
    cursor: pointer;
}
.header {
    width: 100%;
    height: 70px;
    border-bottom: 2px solid rgb(235, 235, 235);
}
#sidebar-nav {
    position: absolute;
    top: 0;
    background-color: #1b80c9;
    width: 240px;
    height: 100%;
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
    text-align: left;
    padding-top: 20px;
    padding-left: 15px;
}
.nav-label {
    font-weight: bold;
    font-size: 16px;
    margin: 12px 0;
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

.dashboard {
    position: absolute;
    left: 240px;
    color: rgb(41, 61, 81);
    margin-left: 90px;
    text-align: left;
}
.dashboard h1 {
    font-size: 38px;
    margin-top: 40px;
}
.dashboard .section {
    width: 800px;
    padding-bottom: 40px;
    border-bottom: 2px solid rgb(235, 235, 235);
}
</style>
