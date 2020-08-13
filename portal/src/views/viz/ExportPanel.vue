<template>
    <div :class="'export-panel ' + containerClass">
        <div class="heading">
            <div class="title">
                Export
            </div>
            <div class="close-button" v-on:click="onClose">
                <img alt="Close" src="../../assets/close.png" />
            </div>
        </div>
        <div class="export-options">
            <div class="button" @click="onExportCSV">CSV</div>
        </div>
        <div class="user-exports" v-if="history">
            <div class="previous-heading">Previous Exports</div>
            <div v-for="de in history" class="export" v-bind:key="de.id">
                <div class="kind">{{ prettyKind(de.kind) }}</div>
                <div class="created">{{ de.createdAt | prettyTime }}</div>
                <div class="busy" v-if="!de.downloadUrl">Generating</div>
                <div class="size" v-if="de.downloadUrl">{{ de.size | prettyBytes }}</div>
                <div class="download" v-if="de.downloadUrl"><a :href="$config.baseUrl + de.downloadUrl">Download</a></div>
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
import { ExportDataAction } from "@/store/typed-actions";
import { GlobalState } from "@/store/modules/global";
import { ExportStatus } from "@/api/api";

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
    data: () => {
        return {};
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
        prettyKind(kind: string): string {
            return "CSV";
        },
        onClose() {
            this.$emit("close");
        },
        onExportCSV() {
            console.log("viz: exporting");
            return this.$store.dispatch(new ExportDataAction(this.bookmark));
        },
        onDownload(ev: any, de: ExportStatus) {
            console.log("viz: download", de);
        },
    },
});
</script>

<style>
.export-panel {
}
.export-panel .heading {
    padding: 1em;
    display: flex;
}
.export-panel .heading .title {
    font-size: 20px;
    font-weight: 500;
}
.export-panel .export {
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
    font-size: 12px;
    padding: 10px;
    background-color: #ffffff;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
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
}
</style>
