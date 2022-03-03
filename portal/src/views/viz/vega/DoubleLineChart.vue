<template>
    <div class="viz doublelinechart"></div>
</template>

<script>
import _ from "lodash";
import { default as vegaEmbed } from "vega-embed";

import doubleLineSpec from "./doubleLine.v1.json";
import chartConfig from "./chartConfig.json";

import { TimeRange } from "../common";
import { TimeZoom } from "../viz";
import { applySensorMetaConfiguration } from "./customizations";

export default {
    name: "DoubleLineChart",
    props: {
        series: {
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
            const spec = _.cloneDeep(doubleLineSpec);
            spec.config = chartConfig;
            spec.layer[0].data = { name: "table0", values: this.series[0].data };
            spec.layer[0].encoding.y.title = this.series[0].vizInfo.label;
            spec.layer[1].data = { name: "table1", values: this.series[1].data };
            spec.layer[1].encoding.y.title = this.series[1].vizInfo.label;
            spec.width = "container";
            spec.height = "container";

            applySensorMetaConfiguration(spec, this.series);

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
