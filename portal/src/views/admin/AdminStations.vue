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
                <tbody v-for="station in stations" v-bind:key="station.id">
                    <tr v-on:click="selected(station)">
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
                            <div class="button" @click.stop="deleteStation(station)">Delete</div>
                        </td>
                    </tr>
                    <tr v-if="focused && focused.id == station.id">
                        <td colspan="11" class="focused">
                            <div class="row">
                                <div class="uploads" v-if="focused.uploads.length > 0">
                                    <h3>Uploads</h3>
                                    <table>
                                        <thead>
                                            <tr>
                                                <th>ID</th>
                                                <th>Date</th>
                                                <th>Type</th>
                                                <th>Size</th>
                                                <th>Blocks</th>
                                                <th>URL</th>
                                                <th></th>
                                            </tr>
                                        </thead>
                                        <tbody v-for="upload in focused.uploads" v-bind:key="upload.id">
                                            <tr>
                                                <td>{{ upload.id }}</td>
                                                <td>{{ upload.time | prettyTime }}</td>
                                                <td>{{ upload.type }}</td>
                                                <td>{{ upload.size }}</td>
                                                <td>{{ upload.blocks }}</td>
                                                <td>{{ upload.url }}</td>
                                                <td><button v-on:click="onProcessUpload(upload)" class="button">Process</button></td>
                                            </tr>
                                        </tbody>
                                    </table>
                                </div>
                                <div class="tools">
                                    <h3>Information</h3>
                                    <div v-if="focused.data">
                                        Data:
                                        <div>
                                            <b>{{ focused.data.start | prettyTime }}</b>
                                        </div>
                                        <div>
                                            <b>{{ focused.data.end | prettyTime }}</b>
                                        </div>
                                        <div>
                                            <b>{{ focused.data.numberOfSamples }}</b>
                                            Samples
                                        </div>
                                    </div>
                                    <div>
                                        Server Logs:
                                        <a :href="urlForServerLogs(station)" target="_blank">10 days</a>
                                    </div>
                                    <div>
                                        <button v-on:click="onExplore(station)" class="button">Explore Data</button>
                                    </div>
                                    <div>
                                        <button v-on:click="onProcessRecords(station)" class="button">Process Recs</button>
                                    </div>
                                    <div>
                                        <button v-on:click="onProcessData(station)" class="button">Process Data</button>
                                    </div>
                                    <h3>Transfer</h3>
                                    <div>
                                        <TransferStation :station="station" @transferred="(user) => onTransferred(station, user)" />
                                    </div>
                                </div>
                            </div>
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
import FKApi, { Station, SimpleUser, EssentialStation } from "@/api/api";
import { BookmarkFactory } from "@/views/viz/viz";
import TransferStation from "./TransferStation.vue";

export default Vue.extend({
    name: "AdminStations",
    components: {
        ...CommonComponents,
        StandardLayout,
        PaginationControls,
        TransferStation,
    },
    props: {},
    data(): {
        stations: EssentialStation[];
        page: number;
        pageSize: number;
        totalPages: number;
        focused: Station | null;
    } {
        return {
            stations: [],
            page: 0,
            pageSize: 50,
            totalPages: 0,
            focused: null,
        };
    },
    mounted(this: any) {
        return this.refresh();
    },
    methods: {
        async selected(station: EssentialStation): Promise<void> {
            if (this.focused && this.focused.id == station.id) {
                return;
            }
            this.focused = null;
            this.focused = await this.$services.api.getStation(station.id);
        },
        async refresh(): Promise<void> {
            await this.$services.api.getAllStations(this.page, this.pageSize).then((page) => {
                this.totalPages = Math.ceil(page.total / this.pageSize);
                this.stations = page.stations;
            });
        },
        onNewPage(page: number): Promise<void> {
            this.page = page;
            return this.refresh();
        },
        async deleteStation(station: EssentialStation): Promise<void> {
            await this.$confirm({
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
        urlForServerLogs(station: EssentialStation): string {
            const deviceIdBase64 = Buffer.from(station.deviceId, "hex").toString("base64");
            const deviceIdHex = station.deviceId;
            const query = `device_id:("${deviceIdBase64}" or "${deviceIdHex}")`;
            return "https://code.conservify.org/logs-viewer/?range=864000&query=" + encodeURIComponent(query);
        },
        async onProcessUpload(upload: { id: number }): Promise<void> {
            await this.$confirm({
                message: `Are you sure? This could take a while for certain files.`,
                button: {
                    no: "No",
                    yes: "Yes",
                },
                callback: async (confirm) => {
                    if (confirm) {
                        await this.$services.api.adminProcessUpload(upload.id);
                    }
                },
            });
        },
        async onProcessData(station: EssentialStation): Promise<void> {
            await this.$confirm({
                message: `Are you sure? This could take a while for certain stations.`,
                button: {
                    no: "No",
                    yes: "Yes",
                },
                callback: async (confirm) => {
                    if (confirm) {
                        await this.$services.api.adminProcessStationData(station.id);
                    }
                },
            });
        },
        async onProcessRecords(station: EssentialStation): Promise<void> {
            await this.$confirm({
                message: `Are you sure? This could take a while for certain stations.`,
                button: {
                    no: "No",
                    yes: "Yes",
                },
                callback: async (confirm) => {
                    if (confirm) {
                        await this.$services.api.adminProcessStation(station.id, true);
                    }
                },
            });
        },
        onTransferred(station: EssentialStation, user: SimpleUser): void {
            station.owner = user;
        },
        async onExplore(station: EssentialStation): Promise<void> {
            const bm = BookmarkFactory.forStation(station.id);
            await this.$router.push({ name: "exploreBookmark", params: { bookmark: JSON.stringify(bm) } });
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

td.focused {
    border: 1px solid #cfcfcf;
    border-radius: 4px;
    padding: 2em;
}
.row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
}

.row .tools {
    width: 400px;
}
.row .uploads {
}
</style>
