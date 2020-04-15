<template>
    <div id="white-header" class="header">
        <div class="header-inner-section">
            <div class="menu-icon-container">
                <img alt="Menu icon" src="../assets/menu.png" v-on:click="toggleSidebar" />
            </div>
            <div class="text-elements">
                <div class="user-name">
                    <router-link :to="{ name: 'user' }" class="account-link" v-if="isAuthenticated">
                        {{ this.user.name }}
                    </router-link>
                </div>
                <div class="log-out" v-if="isAuthenticated" v-on:click="logout">Log out</div>
                <router-link :to="{ name: 'login' }" class="log-in" v-if="!isAuthenticated">
                    Log in
                </router-link>
            </div>
            <img alt="User image" :src="userImage" class="user-image" v-if="isAuthenticated && user.media_url" />
        </div>
    </div>
</template>

<script>
import FKApi from "../api/api";
import { API_HOST } from "../secrets";

export default {
    name: "HeaderBar",
    props: ["isAuthenticated", "user"],
    data: () => {
        return {
            baseUrl: API_HOST,
            userImage: null,
        };
    },
    watch: {
        user() {
            this.userImage = this.baseUrl + "/user/" + this.user.id + "/media";
        },
    },
    methods: {
        logout() {
            const api = new FKApi();
            api.logout();
            this.$router.push({ name: "login" });
        },
        refreshImage(image) {
            this.userImage = image;
        },
        toggleSidebar() {
            let wide = document.getElementById("sidebar-nav-wide");
            let narrow = document.getElementById("sidebar-nav-narrow");
            let wideDisplay = getComputedStyle(wide, null).display;
            if (wideDisplay == "block") {
                wide.style.display = "none";
                narrow.style.display = "block";
            } else {
                wide.style.display = "block";
                narrow.style.display = "none";
            }
            this.$emit("sidebarToggled");
        },
    },
};
</script>

<style scoped>
#white-header {
    width: auto;
    color: gray;
    text-align: right;
    overflow: hidden;
}
.header-inner-section {
    width: 100%;
    height: 69px;
    float: left;
    border-bottom: 2px solid rgb(235, 235, 235);
}
.user-image {
    float: right;
    margin: 10px 12px 0 0;
    max-width: 50px;
    max-height: 50px;
}
.menu-icon-container {
    float: left;
    margin: 20px 0 0 15px;
    cursor: pointer;
}
.text-elements {
    float: right;
}
.user-name {
    padding: 12px 20px 0 0;
}
.account-link {
    text-decoration: underline;
    cursor: pointer;
}
.log-out,
.log-in {
    padding: 0 20px 0 0;
    cursor: pointer;
}
</style>
