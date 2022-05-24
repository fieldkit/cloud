<template>
    <div>
        <div class="viz scrubber"></div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import { default as vegaEmbed } from "vega-embed";

import { TimeRange } from "../common";
import { TimeZoom, SeriesData } from "../viz";
import { ScrubberSpecFactory, ChartSettings } from "./ScrubberSpecFactory";
import { DiscussionState } from "@/store/modules/discussion";
import { ActionTypes } from "@/store";

export default {
    name: "Scrubber",
    props: {
        series: {
            type: Array as PropType<SeriesData[]>,
            required: true,
        },
        visible: {
            type: Array,
            required: true,
        },
    },
    data(): { vega: unknown | null } {
        return { vega: null };
    },
    async mounted(): Promise<void> {
        await this.refresh();
    },
    watch: {
        async series(): Promise<void> {
            await this.refresh();
        },
        async visible() {
            this.pickRange(this.visible);
        },
    },
    computed: {
        dataEvents () {
            return this.$getters.dataEvents;
        }
    },
    methods: {
        async refresh(): Promise<void> {
            const factory = new ScrubberSpecFactory(this.series, new ChartSettings(TimeRange.mergeArrays([this.visible])), this.dataEvents);

            const spec = factory.create();

            const vegaInfo = await vegaEmbed(this.$el, spec, {
                renderer: "svg",
                actions: { source: false, editor: false, compiled: false },
            });

            this.vega = vegaInfo;

            let scrubbed = [];
            vegaInfo.view.addSignalListener("brush", (_, value) => {
                if (value.time) {
                    scrubbed = value.time;
                } else if (this.series[0].data) {
                    scrubbed = this.series[0].data.timeRange;
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
            vegaInfo.view.addEventListener("mouseup", () => {
                console.log("viz: vega:scrubber-brush", scrubbed);
                if (scrubbed.length == 2) {
                    this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                }
            });

            console.log("viz: scrubber", {
                state: vegaInfo.view.getState(),
                data: vegaInfo.view.data("data_1")
            });

            this.pickRange(this.visible);
        },
        async brush(times: number[]): Promise<void> {
            if (!this.vega || !this.series[0].queried) {
                console.log("viz: vega:scrubber:brush-ignore");
                return;
            }
            const x = times.map((v) => this.vega.view.scale("x")(v));
            console.log("viz: vega:scrubber:brush", times, x);
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
        async pickRange(timeRange: number[]): Promise<void> {
            const first = this.series[0];
            if (first.ds) {
                const xDomainsAll = this.series.map((series: SeriesData) => series.queried.timeRange);
                const allRanges = [...xDomainsAll, this.visible];
                const timeRangeAll = TimeRange.mergeArraysIgnoreExtreme(allRanges).toArray();
                if (_.isEqual(timeRangeAll, timeRange)) {
                    await this.brush([]);
                } else {
                    await this.brush(timeRange);
                }
            }
        },
        //TODO move to store
        // async getDataEvents() {
        //     await this.$services.api
        //         .getDataEvents(JSON.stringify(this.parentData))
        //         .then( (response) => {
        //             console.log(response)
        //         })
        //         .catch((e) => {
        //             console.error(e);
        //             //this.errorMessage = CommentsErrorsEnum.postComment;
        //         });

        // },
        getDataEvents() {
            this.$store.dispatch(ActionTypes.NEED_DATA_EVENTS, { bookmark: JSON.stringify(this.parentData) });
        }
    },
};
</script>

<style lang="scss">
@import "src/scss/variables";
.viz {
    width: 100%;
}
</style>
