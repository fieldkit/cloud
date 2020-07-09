import _ from "lodash";
import Vue from "vue";
import * as d3 from "d3";

import { Time, TimeRange, Margins, ChartLayout } from "./common";
import { Graph, QueriedData, Workspace } from "./viz";

export const D3Map = Vue.extend({
    name: "D3Map",
    components: {},
    data() {
        return {};
    },
    props: {
        viz: {
            type: Graph,
            required: true,
        },
        workspace: {
            type: Workspace,
            required: true,
        },
    },
    computed: {
        data(): QueriedData {
            if (this.viz.data && !this.viz.data.empty) {
                return this.viz.data;
            }
            return null;
        },
    },
    watch: {
        data(newValue, oldValue) {
            this.viz.log("graphing (data)");
            this.refresh();
        },
    },
    mounted() {
        this.viz.log("mounted");
        this.refresh();
    },
    updated() {
        this.viz.log("updated");
    },
    methods: {
        refresh() {
            if (!this.data) {
                return;
            }
            const layout = new ChartLayout(1050, 340, new Margins({ top: 5, bottom: 50, left: 50, right: 0 }));
            const data = this.data;
            const timeRange = data.timeRange;
            const dataRange = data.dataRange;
            const charts = [
                {
                    layout: layout,
                },
            ];

            const svg = d3
                .select(this.$el)
                .select(".chart")
                .selectAll("svg")
                .data(charts)
                .join((enter) => {
                    const adding = enter
                        .append("svg")
                        .attr("preserveAspectRatio", "xMidYMid meet")
                        .attr("width", (c) => c.layout.width)
                        .attr("height", (c) => c.layout.height);

                    return adding;
                });
        },
    },
    template: `<div class="viz histogram"><div class="chart"></div></div>`,
});
