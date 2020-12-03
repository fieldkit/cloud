<template>
    <StandardLayout>
        <div class="container">
            <div class="menu">
                <router-link :to="{ name: 'adminUsers' }" class="link">Users</router-link>
                <router-link :to="{ name: 'adminStations' }" class="link">Stations</router-link>
            </div>

            <div class="status" v-if="status">
                <table>
                    <tbody>
                        <tr>
                            <th>Server</th>
                            <td>{{ status.serverName }}</td>
                        </tr>
                        <tr>
                            <th>Tag</th>
                            <td>{{ status.tag }}</td>
                        </tr>
                        <tr>
                            <th>Name</th>
                            <td>{{ status.name }}</td>
                        </tr>
                        <tr>
                            <th>GIT</th>
                            <td>{{ status.git.hash }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import CommonComponents from "@/views/shared";
import { PortalDeployStatus } from "@/api";

export default Vue.extend({
    name: "AdminMain",
    components: {
        StandardLayout,
        ...CommonComponents,
    },
    data(): {
        status: PortalDeployStatus | null;
    } {
        return {
            status: null,
        };
    },
    async mounted(): Promise<void> {
        await this.$services.api.getStatus().then((status) => {
            this.status = status;
        });
    },
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
