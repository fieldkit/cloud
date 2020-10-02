<template>
    <StandardLayout>
        <div class="container">
            <router-link :to="{ name: 'adminMain' }" class="link">Back to Admin</router-link>

            <vue-confirm-dialog />

            <table class="stations">
                <thead>
                    <tr class="header">
                        <th>ID</th>
                        <th>Name</th>
                        <th>Device ID</th>
                        <th>Owner</th>
                        <th>Recording</th>
                        <th>Last Ingestion</th>
                        <th>Created</th>
                        <th>Updated</th>
                        <th>Firmware</th>
                        <th>Location</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="station in stations" v-bind:key="station.id">
                        <td>{{ station.id }}</td>
                        <td class="name">{{ station.name }}</td>
                        <td class="device-id">{{ station.deviceId }}</td>
                        <td class="owner">{{ station.owner.name }}</td>
                        <td class="date recording">{{ station.recordingStartedAt | prettyTime }}</td>
                        <td class="date ingestion">{{ station.lastIngestionAt | prettyTime }}</td>
                        <td class="date created">{{ station.createdAt | prettyDate }}</td>
                        <td class="date updated">{{ station.updatedAt | prettyTime }}</td>
                        <td class="firmware">
                            <template v-if="station.firmwareNumber">
                                {{ station.firmwareNumber }}
                            </template>
                        </td>
                        <td class="location">
                            <div v-if="station.location">
                                {{ station.location.latitude | prettyCoordinate }}, {{ station.location.longitude | prettyCoordinate }}
                            </div>
                        </td>
                        <td>
                            <div class="button" v-on:click="deleteStation(station)">Delete</div>
                        </td>
                    </tr>
                </tbody>
            </table>

            <div>
                <PaginationControls :page="page" :totalPages="totalPages" @new-page="onNewPage" />
            </div>
        </div>
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import CommonComponents from "@/views/shared";
import PaginationControls from "@/views/shared/PaginationControls.vue";

import FKApi, { EssentialStation } from "@/api/api";

export default Vue.extend({
    name: "AdminStations",
    components: {
        ...CommonComponents,
        StandardLayout,
        PaginationControls,
    },
    props: {},
    data: () => {
        return {
            stations: [],
            page: 0,
            pageSize: 50,
            totalPages: 0,
        };
    },
    mounted(this: any) {
        return this.refresh();
    },
    methods: {
        refresh() {
            return this.$services.api.getAllStations(this.page, this.pageSize).then((page) => {
                this.totalPages = Math.ceil(page.total / this.pageSize);
                this.stations = page.stations;
            });
        },
        onNewPage(this: any, page: number) {
            this.page = page;
            return this.refresh();
        },
        deleteStation(this: any, station: EssentialStation) {
            return this.$confirm({
                message: `Are you sure? This operation cannot be undone.`,
                button: {
                    no: "No",
                    yes: "Yes",
                },
                callback: (confirm) => {
                    if (confirm) {
                        return this.$services.api.deleteStation(station.id).then(() => {
                            return this.refresh();
                        });
                    }
                },
            });
        },
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
.form-save-button {
    margin-top: 50px;
    width: 300px;
    height: 45px;
    background-color: #ce596b;
    border: none;
    color: white;
    font-size: 18px;
    font-weight: 600;
    border-radius: 5px;
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
.stations {
    font-size: 15px;
    margin-top: 1em;
}
.stations thead tr th {
    background-color: #fcfcfc;
    border-bottom: 2px solid #d8dce0;
}
.stations tbody tr {
    padding: 5px;
}
.stations td {
    padding: 5px;
}
.stations .device-id {
    font-family: monospace;
    font-size: 14px;
}
.stations .date {
    text-align: right;
}
.stations .button {
    font-size: 12px;
    padding: 5px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
    text-align: center;
}
</style>
