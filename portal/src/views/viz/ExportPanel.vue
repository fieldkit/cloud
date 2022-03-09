<template>
    <div :class="'export-panel ' + containerClass">
        <div class="heading">
            <div class="title">Export</div>
            <div class="close-button icon icon-close" v-on:click="onClose"></div>
        </div>
        <div class="export-options">
            <div class="button" @click="onExportCSV">CSV</div>
            <div class="button" @click="onExportJSONLines">JSON Lines</div>
        </div>
        <div class="user-exports" v-if="history">
            <div class="previous-heading">Previous Exports</div>
            <div v-if="history.length == 0">No previous exports.</div>
            <div v-for="de in history" class="export" v-bind:key="de.id">
                <div class="kind">{{ prettyKind(de.format) }}</div>
                <div class="created">{{ de.createdAt | prettyTime }}</div>
                <div class="busy" v-if="!de.downloadUrl">{{ de.progress | prettyPercentage }}</div>
                <div class="size" v-if="de.downloadUrl">{{ de.size | prettyBytes }}</div>
                <a
                    class="download"
                    v-if="de.downloadUrl"
                    v-bind:class="{ downloaded: false }"
                    :href="$config.baseUrl + de.downloadUrl"
                    @click="(ev) => onDownload(de)"
                >
                    <template v-if="isDownloaded(de.id)">Downloaded</template>
                    <template v-else>Download</template>
                </a>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import CommonComponents from "@/views/shared";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { ExportDataAction, ExportParams } from "@/store/typed-actions";
import { GlobalState } from "@/store/modules/global";
import { ExportStatus } from "@/api/api";

interface Data {
    downloaded: { [index: number]: boolean };
}

export default Vue.extend({
    name: "ExportPanel",
    components: {
        ...CommonComponents,
    },
    props: {
        bookmark: {
            type: Object,
            required: true,
        },
        containerClass: {
            type: String,
            required: true,
        },
    },
    data: (): Data => {
        return {
            downloaded: {},
        };
    },
    computed: {
        ...mapState({
            history: (gs: GlobalState) => {
                if (gs.exporting.history) {
                    return _.reverse(_.sortBy(Object.values(gs.exporting.history), (de) => de.createdAt));
                }
                return [];
            },
        }),
    },
    mounted() {
        return this.$store.dispatch(ActionTypes.NEED_EXPORTS);
    },
    methods: {
        isDownloaded(this: any, id: number): boolean {
            return this.downloaded[id] === true;
        },
        prettyKind(kind: string): string {
            const map = {
                csv: "CSV",
                "json-lines": "JSON Lines",
            };
            return map[kind] || "Unknown";
        },
        async onClose(): Promise<void> {
            this.$emit("close");
        },
        async onExportCSV(): Promise<void> {
            console.log("viz: exporting");
            await this.$store.dispatch(new ExportDataAction(this.bookmark, new ExportParams({ csv: true })));
        },
        async onExportJSONLines(): Promise<void> {
            console.log("viz: exporting");
            await this.$store.dispatch(new ExportDataAction(this.bookmark, new ExportParams({ jsonLines: true })));
        },
        onDownload(de: ExportStatus): void {
            console.log("viz: download", de);
            Vue.set(this.downloaded, de.id, true);
        },
    },
});
</script>

<style>
.export-panel .heading {
    padding: 1em;
    display: flex;
    align-items: center;
}
.export-panel .heading .title {
    font-size: 20px;
    font-weight: 500;
}
.export-panel .heading .close-button {
    margin-left: auto;
    cursor: pointer;
}
.export-panel .user-exports {
    display: flex;
    flex-direction: column;
}

.user-exports {
    padding: 20px;
}

.user-exports .export {
    text-align: left;
    margin-bottom: 1em;
    display: grid;
    font-size: 13px;
    grid-template-columns: 1fr 2fr 1fr 1fr;
    align-items: baseline;
}

.download {
    display: block;
    padding: 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
    text-align: center;
}

.download.downloaded {
    color: #efefef;
}

.download a {
    text-align: center;
}

.export-options {
    padding: 20px;
}

.previous-heading {
    font-size: 20px;
    font-weight: 500;
    text-align: left;
    margin-bottom: 1em;
}

.button {
    font-size: 12px;
    padding: 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
    margin-bottom: 1em;
}
</style>
