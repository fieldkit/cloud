import _ from "lodash";
import { ChartSettings, SeriesData } from "./SpecFactory";
import { TimeRange } from "../common";

export { ChartSettings };

export class ScrubberSpecFactory {
    constructor(private readonly allSeries, private readonly settings: ChartSettings = ChartSettings.Container) {}

    create() {
        const first = this.allSeries[0]; // TODO
        const xDomainsAll = this.allSeries.map((series: SeriesData) => series.queried.timeRange);
        // We ignore extreme ranges here because of this.settings.timeRange
        const allRanges = [...xDomainsAll, this.settings.timeRange.toArray()];
        const timeRangeAll = TimeRange.mergeArraysIgnoreExtreme(allRanges).toArray();

        return {
            $schema: "https://vega.github.io/schema/vega-lite/v5.json",
            description: "FK Scrubber Spec",
            width: "container",
            height: 50,
            config: {
                axis: {
                    labelFont: "Avenir Light",
                    labelFontSize: 12,
                    labelColor: "#2c3e50",
                    titleColor: "#2c3e50",
                    titleFont: "Avenir Light",
                    titleFontSize: 14,
                    titlePadding: 10,
                    tickSize: 10,
                    tickOpacity: 0,
                    domain: false,
                },
                axisX: {
                    title: null,
                    tickSize: 20,
                },
                view: {
                    fill: "#f4f5f7",
                    stroke: "transparent",
                },
            },
            data: {
                name: "table",
                values: first.queried.data,
            },
            encoding: {
                x: {
                    title: "Time",
                    field: "time",
                    type: "temporal",
                    scale: {
                        domain: timeRangeAll,
                    },
                    axis: {
                        formatType: "time",
                        labelPadding: -14,
                        titlePadding: 5,
                        tickCount: 8,
                        grid: false,
                    },
                },
                y: {
                    field: "value",
                    type: "quantitative",
                    axis: {
                        labelOpacity: 0,
                        titleOpacity: 0,
                        grid: false,
                    },
                },
            },
            layer: [
                {
                    params: [
                        {
                            name: "brush",
                            select: {
                                type: "interval",
                                encodings: ["x"],
                            },
                        },
                    ],
                    mark: {
                        type: "area",
                        color: "#DCDEDF",
                    },
                },
                {
                    transform: [
                        {
                            filter: {
                                param: "brush",
                            },
                        },
                    ],
                    mark: {
                        type: "area",
                        color: "#52b5e4",
                    },
                },
                /*
                {
                    mark: {
                        type: "image",
                        width: 30,
                        height: 30,
                    },
                    encoding: {
                        x: {
                            field: "time",
                            type: "temporal",
                        },
                        url: {
                            field: "img",
                            type: "nominal",
                        },
                    },
                    data: {
                        values: [], // Annotations
                    },
                },
                */
            ],
            resolve: {
                axis: {
                    x: "shared",
                    y: this.allSeries.length > 1 ? "independent" : "shared",
                },
            },
        };
    }
}
