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
import chartStyles from "./chartStyles";
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
                downloadFileName: this.getFilename(this.series[0]),
                tooltip: {
                    offsetX: -50,
                    offsetY: 50,
                    formatTooltip: (value, sanitize) => {
                        return `<h3><span class="tooltip-color" style="color: ${sanitize(this.getTooltipColor(value.name))};">■</span>
                                        ${sanitize(value.title)}</h3>
                                        <p class="value">${sanitize(value.Value)}</p>
                                        <p class="time">${sanitize(value.time)}</p>`;
                    },
                },
                actions: { source: false, editor: false, compiled: false },
            });

            this.vega = vegaInfo;

            // Replace vega-embed save as icon with custom button
            const saveAs = document.createElement("svg");
            const summary = document.querySelector("summary");
            const svg = summary.querySelector("svg");
            svg.setAttribute("viewBox", "0 0 20 20");
            svg.innerHTML = '<g id="icon_SaveAs" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" stroke-linecap="round">' +
                                '<line x1="7.96030045" y1="1" x2="7.96030045" y2="11" id="Path-2" stroke="#2C3E50" stroke-width="1.5" stroke-linejoin="round"></line>' +
                                '<polyline id="Path-9" stroke="#2C3E50" stroke-width="1.5" stroke-linejoin="bevel" points="12.8961983 6.50366211 8.05585126 11 2.92245537 6.50366211"></polyline>' +
                                '<polyline id="Path-10" stroke="#2C3E50" stroke-width="1.5" stroke-linejoin="round" points="1 12.5363846 1 16.5 15.1181831 16.5 15.1181831 12.5363846"></polyline>' +
                            '</g>';
            const saveLabel = document.createElement("span");
            saveLabel.setAttribute("class", "save-label");
            saveLabel.innerHTML = "Save As";
            summary.appendChild(saveLabel);

            /*
            vegaInfo.view.addSignalListener("unit", (_, value) => {
                console.log("unit", value);
            });
            vegaInfo.view.addSignalListener("hover", (_, value) => {
                console.log("hover", value);
            });
            */

            let scrubbed = [];
            vegaInfo.view.addSignalListener("brush", (_, value) => {
                scrubbed = value.time;
            });
            vegaInfo.view.addEventListener("mouseup", () => {
                if (scrubbed.length == 2) {
                    this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                }
            });
            // Watch for brush drag outside the window
            vegaInfo.view.addEventListener("mousedown", (e) => {
                window.addEventListener("mouseup", (e) => {
                    if (scrubbed.length == 2 && e.target.nodeName !== "path") {
                        this.$emit("time-zoomed", new TimeZoom(null, new TimeRange(scrubbed[0], scrubbed[1])));
                    }
                });
            });

            console.log("viz: vega:ready", {
                state: vegaInfo.view.getState(),
                layouts: vegaInfo.view.data("all_layouts"),
            });
        },
        getFilename(series) {
            const stationName = series.vizInfo.station.name;
            const sensorName = series.vizInfo.name;

            return `${stationName}_${sensorName}`.replace("[^a-zA-Z0-9\\.\\-]", "_");
        },
        getTooltipColor(name) {
            if (name === "LEFT") {
                return chartStyles.primaryLine.stroke;
            }
            if (name === "RIGHT") {
                return chartStyles.secondaryLine.stroke;
            } else {
                return "#ccc";
            }
        },
    },
});
</script>

<style>
.viz {
    width: 100%;
}
.vega-embed summary{
    border-radius: 0px !important;
    width: 60px;
    height: 1em;
    display: flex;
    align-items: center;
    margin-right: 3.2em !important;
}
.vega-embed summary svg {
    width: 16px !important;
    height: 16px !important;
    display: inline-block;
}
.save-label {
    font-size: 10px;
    margin-left: 5px;
}
.vega-embed .vega-actions {
    right: 3em !important;
}

</style>
