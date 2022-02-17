<template>
    <div>
        <button v-on:click="downloadChart('png')" v-if="false">Download chart png</button>
        <button v-on:click="downloadChart('svg')" v-if="false">Download chart svg</button>
        <div class="viz linechart"></div>
    </div>
</template>

<script>
import _ from "lodash";
import { default as vegaEmbed } from "vega-embed";
import { expressionFunction } from "vega";
import lineSpec from "./line.v1.json";
import chartConfig from "./chartConfig.json";

import { TimeRange } from "../common";
import { TimeZoom } from "../viz";

export default {
    name: "LineChart",
    props: {
        data: {
            type: Object,
            required: true,
        },
        label: {
            type: String,
            required: true,
        },
        valueSuffix: {
            type: String,
            required: true,
        },
        thresholds: {
            type: Array,
            required: true,
        },
        constrainDataAxis: {
            type: Array,
            required: true,
        },
    },
    mounted: function() {
        console.log("vega-mounted");
        this.refresh();
    },
    watch: {
        label() {
            console.log("vega-watch-label");
            this.refresh();
        },
        data() {
            console.log("vega-watch-data");
            this.refresh();
        },
    },
    methods: {
        async refresh() {
            const spec = _.cloneDeep(lineSpec);
            spec.config = chartConfig;
            spec.data = { name: "table", values: this.data.data };
            spec.layer[0].encoding.y.axis.title = this.label;
            spec.width = "container";
            spec.height = "container";

            if (this.thresholds.length > 0) {
                console.log("viz-thresholds", this.thresholds);
                spec.layer[0].layer = this.thresholds;
            }

            if (this.constrainDataAxis && this.constrainDataAxis.length > 0) {
                const range = this.constrainDataAxis[0];
                spec.layer[0].encoding.y.scale.domain = [range.minimum, range.maximum];
            }

            expressionFunction("fkHumanReadable", (datum) => {
                if (_.isUndefined(datum)) {
                    return "N/A";
                }
                if (this.valueSuffix) {
                    return `${datum.toFixed(3)} ${this.valueSuffix}`;
                }
                return `${datum.toFixed(3)}`;
            });

            await vegaEmbed(this.$el, spec, {
                renderer: "svg",
                tooltip: { offsetX: -50, offsetY: 50 },
                actions: { source: false, editor: false, compiled: false },
            }).then((view) => {
                this.vegaView = view;
                let scrubbed = [];
                view.view.addSignalListener("brush", (_, value) => {
                    scrubbed = value.time;
                });
                this.vegaView.view.addEventListener("mouseup", () => {
                    console.log("vega-line-brush", scrubbed);
                    if (scrubbed.length == 2) {
                        this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                    }
                });
            });
        },
        // From https://vega.github.io/vega/docs/api/view/#view_toImageURL
        async downloadChart(fileFormat) {
            await this.vegaView.view
                .toImageURL(fileFormat, 2)
                .then(function(url) {
                    const link = document.createElement("a");
                    link.setAttribute("href", url);
                    link.setAttribute("target", "_blank");
                    link.setAttribute("download", "vega-export." + fileFormat);
                    link.dispatchEvent(new MouseEvent("click"));
                })
                .catch(function(error) {
                    console.log(error);
                });
        },
    },
};
</script>

<style scoped>
.viz {
    width: 100%;
}
</style>