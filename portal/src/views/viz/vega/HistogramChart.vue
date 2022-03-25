<template>
    <div class="viz histogram"></div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { default as vegaEmbed } from "vega-embed";

import { SeriesData } from "../viz";
import { ChartSettings } from "./SpecFactory";
import { HistogramSpecFactory } from "./HistogramSpecFactory";

export default Vue.extend({
    name: "Histogram",
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
        console.log("vega-mounted");
        await this.refresh();
    },
    watch: {
        async series(): Promise<void> {
            console.log("vega-watch-series");
            await this.refresh();
        },
    },
    methods: {
        async refresh() {
            const factory = new HistogramSpecFactory(this.series, ChartSettings.Container);

            const spec = factory.create();

            // histogramSpec.config = chartConfig;
            // histogramSpec.data = { name: "table", values: this.data.data };
            // histogramSpec.encoding.x.axis.title = this.label;
            // histogramSpec.width = "container";
            // histogramSpec.height = "container";

            const vegaInfo = await vegaEmbed(".histogram", spec, {
                renderer: "svg",
                tooltip: { offsetX: -50, offsetY: 50 },
                actions: { source: false, editor: false, compiled: false },
            });

            this.vega = vegaInfo;
        },
    },
});
</script>

<style scoped>
.viz {
    width: 100%;
}
</style>
