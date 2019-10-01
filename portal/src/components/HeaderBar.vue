<template>
    <div id="white-header" class="header">
        <div class="user-name">{{ this.user.name }}</div>
        <div class="log-out" v-if="isAuthenticated" v-on:click="logout">Log out</div>
        <router-link :to="{ name: 'login' }" class="log-in" v-if="!isAuthenticated">
            Log in
        </router-link>
    </div>
</template>

<script>
import FKApi from "../api/api";

export default {
    name: "HeaderBar",
    props: ["isAuthenticated", "user"],
    data: () => {
        return {};
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

<style scoped>
#white-header {
    color: gray;
    text-align: right;
}
.user-name {
    padding: 12px 20px 0 0;
}
.log-out,
.log-in {
    padding: 0 20px 0 0;
    cursor: pointer;
}
</style>
