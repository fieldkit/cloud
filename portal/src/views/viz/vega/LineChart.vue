<template>
    <div>
        <div class="viz linechart"></div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import { default as vegaEmbed } from "vega-embed";

import { TimeRange } from "../common";
import { TimeZoom, SeriesData } from "../viz";
import { ChartSettings } from "./SpecFactory";
import { TimeSeriesSpecFactory } from "./TimeSeriesSpecFactory";

export default Vue.extend({
    name: "LineChart",
    props: {
        series: {
            type: Array as PropType<SeriesData[]>,
            required: true,
        },
    },
    data(): {
        vega: unknown | undefined;
    } {
        return {
            vega: undefined,
        };
    },
    async mounted(): Promise<void> {
        await this.refresh();
    },
    watch: {
        async series(): Promise<void> {
            await this.refresh();
        },
    },
    methods: {
        async refresh(): Promise<void> {
            const factory = new TimeSeriesSpecFactory(this.series, ChartSettings.Auto);

            const spec = factory.create();

            const vegaInfo = await vegaEmbed(this.$el, spec, {
                renderer: "svg",
                tooltip: { offsetX: -50, offsetY: 50 },
                actions: { source: false, editor: false, compiled: false },
            });

            this.vega = vegaInfo;

            let scrubbed = [];
            vegaInfo.view.addSignalListener("brush", (_, value) => {
                scrubbed = value.time;
            });
            vegaInfo.view.addEventListener("mouseup", () => {
                if (scrubbed.length == 2) {
                    this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                }
            });

            console.log("viz: vega:ready", vegaInfo.view.getState());
        },
        async downloadChart(fileFormat: string): Promise<void> {
            // From https://vega.github.io/vega/docs/api/view/#view_toImageURL
            await this.vega.view
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
});
</script>

<style scoped>
.viz {
    width: 100%;
}
</style>
