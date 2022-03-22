<template>
    <div class="viz doublelinechart"></div>
</template>

<script>
import _ from "lodash";
import { default as vegaEmbed } from "vega-embed";
import * as vegaLite from "vega-lite";

import newDoubleLineSpec from "./doublelinevega.json";
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
    data() {
        return {
            testing: true,
        };
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
            const spec = _.cloneDeep(this.testing ? newDoubleLineSpec : doubleLineSpec);

            if (this.testing) {
                // spec.width = "container";
                // spec.height = "container";
            } else {
                // const compiled = vegaLite.compile(spec);
                // console.log("COMPILED", compiled);
                // console.log("COMPILED", JSON.stringify(compiled.spec));

                spec.width = "container";
                spec.height = "container";
                spec.config = chartConfig;
                spec.layer[0].data = { name: "table0", values: this.series[0].data };
                spec.layer[0].encoding.y.title = this.series[0].vizInfo.label;
                spec.layer[1].data = { name: "table1", values: this.series[1].data };
                spec.layer[1].encoding.y.title = this.series[1].vizInfo.label;
            }

            applySensorMetaConfiguration(spec, this.series);

            await vegaEmbed(this.$el, spec, {
                renderer: "svg",
                tooltip: { offsetX: -50, offsetY: 50 },
                actions: { source: false, editor: false, compiled: false },
            }).then((view) => {
                this.vegaView = view;
                let scrubbed = [];
                view.view.addSignalListener("unit", (_, value) => {
                    console.log("vega:state(unit)", value, this.vegaView.view.getState());
                });
                view.view.addSignalListener("brush2_x", (_, value) => {
                    // console.log("vega:state(brush2_x)", value);
                });
                view.view.addSignalListener("brush2_tuple", (_, value) => {
                    // console.log("vega:state(brush2_tuple)", value);
                });
                view.view.addSignalListener("brush2_time", (_, value) => {
                    // console.log("vega:state(brush2_time)", value);
                });
                view.view.addSignalListener("hover", (_, value) => {
                    console.log("vega:state(hover)", value);
                });
                view.view.addSignalListener("brush2", (_, value) => {
                    scrubbed = value.time;
                    console.log("vega:state(brush2)", value, this.vegaView.view.getState());
                });
                this.vegaView.view.addEventListener("mouseup", () => {
                    console.log("vega-line-brush", scrubbed);
                    if (scrubbed && scrubbed.length == 2) {
                        this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                    }
                });

                console.log("vega:state", this.vegaView.view.getState());
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
