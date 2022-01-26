<template>
    <div class="viz rangechart"></div>
</template>

<script>
import { default as vegaEmbed } from "vega-embed";
import rangeSpec from "./range.vl.json";
import chartConfig from "./chartConfig.json";

export default {
    name: "RangeChart",
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
            console.log("vega-updated");

            rangeSpec.config = chartConfig;
            rangeSpec.data = { name: "table", values: this.data.data };
            // rangeSpec.layer[0].encoding.y.title = this.label;

            const vegaView = await vegaEmbed(".rangechart", rangeSpec, {
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
