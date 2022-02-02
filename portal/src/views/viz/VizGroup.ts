import _ from "lodash";
import Vue from "vue";

import { Workspace, Viz, Group, TimeZoom } from "./viz";
import { VizGraph } from "./VizGraph";
import { VegaScrubber as Scrubber } from "./VegaScrubber";

export const VizGroup = Vue.extend({
    components: {
        VizGraph,
        Scrubber,
    },
    props: {
        group: {
            type: Group,
            required: true,
        },
        workspace: {
            type: Workspace,
            required: true,
        },
        topGroup: {
            type: Boolean,
            required: true,
        },
    },
    mounted() {
        this.group.log("mounted", this.group, this.topGroup);
    },
    methods: {
        raiseGroupZoomed(zoom: TimeZoom, ...args) {
            return this.$emit("group-time-zoomed", zoom, ...args);
        },
        raiseVizTimeZoomed(...args) {
            return this.$emit("viz-time-zoomed", ...args);
        },
        raiseVizGeoZoomed(...args) {
            return this.$emit("viz-geo-zoomed", ...args);
        },
        raiseRemove(...args) {
            return this.$emit("viz-remove", ...args);
        },
        raiseCompare(...args) {
            return this.$emit("viz-compare", ...args);
        },
        raiseChangeSensors(...args) {
            return this.$emit("viz-change-sensors", ...args);
        },
        raiseChangeChart(...args) {
            return this.$emit("viz-change-chart", ...args);
        },
        raiseChangeLinkage(...args) {
            return this.$emit("viz-change-linkage", ...args);
        },
        isLinked(index: number): boolean {
            if (index == 0) {
                return false;
            }
            return this.group.vizes.length > 1;
        },
    },
    template: `
		<div class="">
			<div class="group-container">
				<template v-for="(viz, index) in group.vizes" :key="viz.id">
                    <div style="display: none;">Group:{{group.id}} Viz:{{viz.id}} {{isLinked(index)}} {{index}} {{topGroup}}</div>
					<div class="icons-container" v-if="!topGroup || index > 0" v-bind:class="{ 'linked': isLinked(index), 'unlinked': !isLinked(index) }">
						<div class="invisible-spacing-icon"></div>
						<div class="icon" v-on:click="(ev) => raiseChangeLinkage(viz, !isLinked(index))" v-bind:class="{ 'icon-open-link': !isLinked(index), 'icon-link': isLinked(index) }"></div>
						<div class="icon remove-icon" v-on:click="(ev) => raiseRemove(viz)"></div>
					</div>

					<VizGraph :viz="viz" :workspace="workspace"
						@viz-time-zoomed="(...args) => raiseVizTimeZoomed(viz, ...args)"
						@viz-geo-zoomed="(...args) => raiseVizGeoZoomed(viz, ...args)"
						@viz-remove="(...args) => raiseRemove(viz, ...args)"
						@viz-compare="(...args) => raiseCompare(viz, ...args)"
						@viz-change-sensors="(...args) => raiseChangeSensors(viz, ...args)"
						@viz-change-chart="(...args) => raiseChangeChart(viz, ...args)"
						/>
					</component>
				</template>
				<div v-if="group.scrubbers && !group.scrubbers.empty">
					<Scrubber :scrubbers="group.scrubbers" @viz-time-zoomed="(...args) => raiseGroupZoomed(...args)" />
				</div>
				<div v-else class="loading-scrubber">
					Loading Scrubber
				</div>
			</div>
		</div>
	`,
});
