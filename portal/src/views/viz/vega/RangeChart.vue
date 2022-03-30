<template>
    <div class="viz rangechart"></div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { default as vegaEmbed } from "vega-embed";

import { SeriesData } from "../viz";
import { ChartSettings } from "./SpecFactory";
import { RangeSpecFactory } from "./RangeSpecFactory";

export default Vue.extend({
    name: "RangeChart",
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
        async refresh() {
            const factory = new RangeSpecFactory(this.series, ChartSettings.Container);

            const spec = factory.create();

            const vegaInfo = await vegaEmbed(".rangechart", spec, {
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
