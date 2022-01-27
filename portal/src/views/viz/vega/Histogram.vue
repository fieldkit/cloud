<template>
    <div class="viz histogram"></div>
</template>

<script>
import { default as vegaEmbed } from "vega-embed";
import histogramSpec from "./histo.vl.json";
import chartConfig from "./chartConfig.json";

export default {
    name: "Histogram",
    props: {
        data: {
            type: Object,
            required: true,
        },
        label: {
            type: String,
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
            histogramSpec.config = chartConfig;
            histogramSpec.data = { name: "table", values: this.data.data };
            histogramSpec.encoding.x.axis.title = this.label;

            const vegaView = await vegaEmbed(".histogram", histogramSpec, {
                renderer: "svg",
                tooltip: { offsetX: -50, offsetY: 50 },
                actions: { source: false, editor: false, compiled: false },
            });

            this.vegaView = vegaView;
        },
    },
};
</script>

<style scoped>
.viz {
    width: 100%;
}
</style>
