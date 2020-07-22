<template>
    <StandardLayout>
        <div class="container">
            <router-link :to="{ name: 'adminUsers' }" class="link">
                Users
            </router-link>
            <router-link :to="{ name: 'adminStations' }" class="link">
                Stations
            </router-link>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import CommonComponents from "@/views/shared";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "AdminMain",
    components: {
        StandardLayout,
        ...CommonComponents,
    },
    props: {},
    data: () => {
        return {};
    },
    methods: {},
});
</script>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    padding: 20px;
    text-align: left;
}
.notification.success {
    margin-top: 20px;
    margin-bottom: 20px;
    padding: 20px;
    border: 2px;
    border-radius: 4px;
}
.notification.success {
    background-color: #d4edda;
}
.notification.failed {
    background-color: #f8d7da;
}
.link {
    display: block;
    margin-bottom: 1em;
    font-size: 18px;
    font-weight: bold;
}
</style>
