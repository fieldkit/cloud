<template>
    <StandardLayout>
        <div class="container">
            <div class="menu">
                <router-link :to="{ name: 'adminUsers' }" class="link">Users</router-link>
                <router-link :to="{ name: 'adminStations' }" class="link">Stations</router-link>
            </div>

            <div class="status" v-if="status">
                <table>
                    <thead>
                        <tr>
                            <th>Server Logs</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>
                                <a
                                    href="https://code.conservify.org/logs-viewer/?range=3600&query=tag:fkprd%20OR%20tag:fkdev"
                                    target="_blank"
                                >
                                    fkprd and fkdev: all
                                </a>
                            </td>
                        </tr>
                        <tr>
                            <td>
                                <a href="https://code.conservify.org/logs-viewer/?range=86400&query=zaplevel:error" target="_blank">
                                    fkprd and fkdev: errors
                                </a>
                            </td>
                        </tr>
                        <tr>
                            <td>
                                <a
                                    href="https://code.conservify.org/logs-viewer/?range=86400&query=_exists_:data_processed&include=device_id,user_id,meta_errors,data_errors,meta_processed,data_processed,blocks,station_name"
                                    target="_blank"
                                >
                                    fkprd and fkdev: ingestions
                                </a>
                            </td>
                        </tr>
                        <tr>
                            <td>
                                <a
                                    href="https://code.conservify.org/logs-viewer/?range=86400&query=message:%22station%20conflict%22"
                                    target="_blank"
                                >
                                    fkprd and fkdev: conflicts
                                </a>
                            </td>
                        </tr>
                    </tbody>
                </table>
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
