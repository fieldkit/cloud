import _ from "lodash";
import Vue from "vue";

import { Workspace, Group, Viz, HasSensorParams, FastTime, ChartType, TimeZoom, GeoZoom, Bookmark } from "./viz";
import { VizGroup } from "./VizGroup";
import { ForbiddenError } from "@/api";

export const VizWorkspace = Vue.extend({
    components: {
        VizGroup,
    },
    props: {
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    data() {
        return {};
    },
    async beforeMount(): Promise<void> {
        console.log(`viz: include-associated(1)`, this.workspace.allStationIds.length);
        const associated = await this.$services.api.getAssociatedStations(this.workspace.allStationIds[0]).catch(async (e) => {
            if (ForbiddenError.isInstance(e)) {
                await this.$router.push({
                    name: "login",
                    params: { errorMessage: String(this.$t("login.privateStation")) },
                    query: { after: this.$route.path },
                });
            }
        });
        if (associated) {
            const ids = associated.stations.map((s) => s.id);
            console.log(`viz: include-associated(1)`, associated);
            await this.workspace.addStationIds(ids);
        }

        return await this.workspace.query();
    },
    methods: {
        signal(ws: Workspace): Workspace {
            const bookmark = ws.bookmark();
            this.$emit("change", bookmark);
            return ws;
        },
        onGroupTimeZoomed(group: Group, zoom: TimeZoom) {
            group.log("zooming", zoom);
            return this.workspace
                .groupZoomed(group, zoom)
                .with((this as any).signal)
                .query();
        },
        onGraphTimeZoomed(group: Group, viz: Viz, zoom: TimeZoom) {
            viz.log("zooming", zoom);
            return this.workspace
                .graphTimeZoomed(viz, zoom)
                .with((this as any).signal)
                .query();
        },
        onGraphGeoZoomed(group: Group, viz: Viz, zoom: GeoZoom) {
            viz.log("zooming(geo)", zoom.bounds);
            return this.workspace
                .graphGeoZoomed(viz, zoom)
                .with((this as any).signal)
                .query();
        },
        onRemove(group: Group, viz: Viz) {
            viz.log("removing");
            return this.workspace
                .remove(viz)
                .with((this as any).signal)
                .query();
        },
        onCompare(group: Group, viz: Viz) {
            viz.log("compare");
            return this.workspace
                .compare(viz)
                .with((this as any).signal)
                .query();
        },
        onChangeSensors(group: Group, viz: Viz, params: HasSensorParams) {
            viz.log("change-sensors", params);
            return this.workspace
                .changeSensors(viz, params)
                .with((this as any).signal)
                .query();
        },
        onChangeChart(group: Group, viz: Viz, chartType: ChartType) {
            viz.log("chart", chartType);
            return this.workspace
                .changeChart(viz, chartType)
                .with((this as any).signal)
                .query();
        },
        onChangeLinkage(group: Group, viz: Viz, linking: boolean) {
            viz.log("linkage");
            return this.workspace
                .changeLinkage(viz, linking)
                .with((this as any).signal)
                .query();
        },
    },
    template: `
        <div class="groups-container">
            <template v-for="(group, index) in workspace.groups" :key="group.id">
                <VizGroup :group="group" :workspace="workspace" :topGroup="index == 0"
                    @group-time-zoomed="(...args) => onGroupTimeZoomed(group, ...args)"
                    @viz-time-zoomed="(...args) => onGraphTimeZoomed(group, ...args)"
                    @viz-geo-zoomed="(...args) => onGraphGeoZoomed(group, ...args)"
                    @viz-remove="(...args) => onRemove(group, ...args)"
                    @viz-compare="(...args) => onCompare(group, ...args)"
                    @viz-change-sensors="(...args) => onChangeSensors(group, ...args)"
                    @viz-change-chart="(...args) => onChangeChart(group, ...args)"
                    @viz-change-linkage="(...args) => onChangeLinkage(group, ...args)"
                />
            </template>
        </div>
	`,
});
