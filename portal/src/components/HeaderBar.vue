<template>
    <div id="white-header" class="header">
        <div class="header-inner-section">
            <div class="menu-icon-container">
                <img alt="Menu icon" src="../assets/menu.png" v-on:click="toggleSidebar" />
            </div>
            <div class="text-elements">
                <div class="user-name">
                    <router-link v-if="user" :to="{ name: 'editUser' }" class="account-link">
                        {{ user.name }}
                    </router-link>
                </div>
                <div class="log-out" v-if="isAuthenticated" v-on:click="logout">Log out</div>
                <router-link :to="{ name: 'login', query: { redirect: $route.fullPath } }" class="log-in" v-if="!isAuthenticated">
                    Log in
                </router-link>
            </div>
            <img v-if="user && userImage" alt="User image" :src="userImage" class="user-image" />
        </div>
    </div>
</template>

<script>
import FKApi from "@/api/api";
import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";

export default {
    name: "HeaderBar",
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({ user: (s) => s.user.user }),
        userImage() {
            if (this.$store.state.user.user.photo) {
                return this.$config.baseUrl + this.$store.state.user.user.photo.url;
            }
            return null;
        },
    },
    methods: {
        logout() {
            return this.$store.dispatch(ActionTypes.LOGOUT).then(() => {
                return this.$router.push({ name: "login" });
            });
        },
        toggleSidebar() {
            /*
            const wide = document.getElementById("sidebar-nav-wide");
            const narrow = document.getElementById("sidebar-nav-narrow");
            const wideDisplay = getComputedStyle(wide, null).display;
            if (wideDisplay == "block") {
                wide.style.display = "none";
                narrow.style.display = "block";
            } else {
                wide.style.display = "block";
                narrow.style.display = "none";
            }
			*/
            this.$emit("toggled");
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
    border-bottom: 2px solid #d8dce0;
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
