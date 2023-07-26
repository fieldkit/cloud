<template>
    <div>
        <div class="viz scrubber"></div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import { default as vegaEmbed } from "vega-embed";

import { isMobile } from "@/utilities";
import { TimeRange } from "../common";
import { TimeZoom, SeriesData } from "../viz";
import { ScrubberSpecFactory, ChartSettings } from "./ScrubberSpecFactory";
import { DiscussionState } from "@/store/modules/discussion";
import { ActionTypes } from "@/store";

export default Vue.extend({
    name: "Scrubber",
    props: {
        series: {
            type: Array as PropType<SeriesData[]>,
            required: true,
        },
        visible: {
            type: Object as PropType<TimeRange>,
            required: true,
        },
        dragging: {
            type: Boolean,
            required: true,
        },
    },
    data(): {
        vega: any | null;
    } {
        return { vega: null };
    },
    async mounted(): Promise<void> {
        await this.refresh();
    },
    watch: {
        async series(): Promise<void> {
            console.log("viz:", "scrubber: refresh(ignored, series)");
            // await this.refresh();
        },
        async dragging(dragging): Promise<void> {
            console.log("viz:", "scrubber: dragging", dragging, this.visible);
            this.pickRange(this.visible);
        },
        async visible(): Promise<void> {
            this.pickRange(this.visible);
        },
        async dataEvents(): Promise<void> {
            await this.refresh();
        },
    },
    computed: {
        dataEvents() {
            return this.$getters.dataEvents;
        },
    },
    methods: {
        async refresh(): Promise<void> {
            console.log("viz:", "scrubber: refresh");

            const factory = new ScrubberSpecFactory(
                this.series,
                new ChartSettings(TimeRange.mergeArrays([this.visible]), undefined, { w: 0, h: 0 }, false, false, isMobile()),
                this.dataEvents.filter(event => (event.start > this.visible.start && event.end < this.visible.end)),
            );

            const spec = factory.create();

            const vegaInfo = await vegaEmbed(this.$el as HTMLElement, spec, {
                renderer: "svg",
                actions: { source: false, editor: false, compiled: false },
            });

            this.vega = vegaInfo;

            // eslint-disable-next-line
            let scrubbed: number[] = [];
            vegaInfo.view.addSignalListener("brush", (_, value) => {
                if (value.time) {
                    scrubbed = value.time;
                } else if (this.series[0].queried) {
                    scrubbed = this.series[0].queried.timeRange;
                }
            });
            // vegaInfo.view.addSignalListener("scrub_handle_left", (_, value) => {
            //     console.log("SCRUB HANDLE RIGHT", value)
            // });
            // vegaInfo.view.addSignalListener("scrub_handle_right", (_, value) => {
            //     console.log("SCRUB HANDLE RIGHT", value)
            // });
            // vegaInfo.view.addEventListener("mousedown", (evt, value) => {
            //     console.log(evt, value);
            // });
            vegaInfo.view.addSignalListener("event_click", (_, value) => {
              this.$emit("event-clicked", value);
            });
            vegaInfo.view.addEventListener("mouseup", () => {
                if (scrubbed.length == 2) {
                    console.log("viz: vega:scrubber:brush-zoomed", scrubbed);
                    this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                }
            });

            console.log("viz: scrubber", {
                state: vegaInfo.view.getState(),
                data: vegaInfo.view.data("data_1"),
            });

            this.pickRange(this.visible);
        },
        inflatedBrushTimes(times: TimeRange): number[] {
            const x = times.toArray().map((v) => this.vega.view.scale("x")(v));
            const minimumWidth = 5;
            const halfMinimumWidth = minimumWidth / 2;
            const width = x[1] - x[0];
            if (width > minimumWidth) {
                return x;
            }
            const middle = x[0] + width / 2;
            return [middle - halfMinimumWidth, middle + halfMinimumWidth];
        },
        async brush(times: TimeRange): Promise<void> {
            if (!this.vega || !this.series[0].queried) {
                console.log("viz: vega:scrubber:brushing-ignore");
                return;
            }

            // console.log("viz: vega:scrubber:brushing", times, x);
            const x = this.inflatedBrushTimes(times);

            try {
                await this.vega.view
                    .signal("brush_x", x)
                    .signal("brush_tuple", {
                        fields: [
                            {
                                field: "time",
                                channel: "x",
                                type: "R",
                            },
                        ],
                        values: times,
                    })
                    .runAsync();
            } catch (error) {
                console.log("viz: error", error);
            }
        },
        async pickRange(timeRange: TimeRange): Promise<void> {
            const first = this.series[0];
            if (first.ds) {
                await this.brush(timeRange);
            }
        },
    },
});
</script>

<style lang="scss">
@import "src/scss/variables";
.viz {
    width: 100%;
}
g.left_scrub:hover > path {
    cursor: ew-resize;
}
g.right_scrub:hover > path {
    cursor: ew-resize;
}
g.brush_brush:hover {
    cursor: grab;
}
g.brush_brush:hover:active {
    cursor: grabbing;
}
</style>
