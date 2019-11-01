<template>
    <div id="white-header" class="header">
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
            userImage: null
        };
    },
    watch: {
        user() {
            this.userImage = this.baseUrl + "/user/" + this.user.id + "/media";
        }
    },
    methods: {
        logout() {
            const api = new FKApi();
            api.logout();
            this.$router.push({ name: "login" });
        },
        refreshImage(image) {
            this.userImage = image;
        }
    }
};
</script>

<style scoped>
#white-header {
    color: gray;
    text-align: right;
}
.user-image {
    float: right;
    margin: 10px 12px 0 0;
    max-width: 50px;
    max-height: 50px;
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
