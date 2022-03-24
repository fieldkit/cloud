<template>
    <div class="viz doublelinechart"></div>
</template>

<script>
import _ from "lodash";
import { default as vegaEmbed } from "vega-embed";

import { TimeRange } from "../common";
import { TimeZoom } from "../viz";
import { TimeSeriesSpecFactory } from "./TimeSeriesSpecFactory";

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
            const factory = new TimeSeriesSpecFactory(this.series);

            const spec = factory.create();

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
                view.view.addSignalListener("hover", (_, value) => {
                    console.log("vega:state(hover)", value);
                });
                view.view.addSignalListener("brush", (_, value) => {
                    scrubbed = value.time;
                    console.log("vega:state(brush)", value, this.vegaView.view.getState());
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
