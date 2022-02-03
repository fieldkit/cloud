<template>
    <div class="viz doublelinechart"></div>
</template>

<script>
import _ from "lodash";
import { default as vegaEmbed } from "vega-embed";
import { expressionFunction } from "vega";

import doubleLineSpec from "./doubleLine.v1.json";
import chartConfig from "./chartConfig.json";

import { TimeRange } from "../common";
import { TimeZoom } from "../viz";

export default {
    name: "DoubleLineChart",
    props: {
        data: {
            type: Array,
            required: true,
        },
        labels: {
            type: Array,
            required: true,
        },
        valueSuffixes: {
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
            doubleLineSpec.config = chartConfig;
            doubleLineSpec.layer[0].data = { name: "table0", values: this.data[0].data };
            doubleLineSpec.layer[0].encoding.y.title = this.labels[0];
            doubleLineSpec.layer[1].data = { name: "table1", values: this.data[1].data };
            doubleLineSpec.layer[1].encoding.y.title = this.labels[1];
            doubleLineSpec.width = "container";
            doubleLineSpec.height = 300;

            expressionFunction("fkHumanReadable", (datum) => {
                if (_.isUndefined(datum)) {
                    return "N/A";
                }
                if (this.valueSuffixes) {
                    return `${datum.toFixed(3)} ${this.valueSuffixes[0]}`;
                }
                return `${datum.toFixed(3)}`;
            });

            await vegaEmbed(this.$el, doubleLineSpec, {
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
