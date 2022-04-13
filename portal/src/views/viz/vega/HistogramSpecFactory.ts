import _ from "lodash";
import { MapFunction, ChartSettings, SeriesData, getString, getSeriesThresholds } from "./SpecFactory";

export class HistogramSpecFactory {
    constructor(private readonly allSeries, private readonly settings: ChartSettings = ChartSettings.Container) {}

    create() {
        const first = this.allSeries[0];

        return this.settings.apply({
            $schema: "https://vega.github.io/schema/vega-lite/v5.json",
            description: "FK Histogram Chart",
            config: {
                axis: {
                    labelFont: "Avenir Light",
                    labelFontSize: 12,
                    labelColor: "#6a6d71",
                    titleFont: "Avenir Light",
                    titleFontSize: 14,
                    titlePadding: 20,
                    tickSize: 10,
                    tickOpacity: 0,
                },
            },
            data: {
                name: "table",
                values: first.queried.data,
            },
            encoding: {
                x: {
                    bin: {
                        maxbins: 20,
                    },
                    field: "value",
                    axis: {
                        grid: false,
                    },
                    title: first.vizInfo.label,
                },
                y: {
                    aggregate: "count",
                    axis: {
                        grid: false,
                    },
                },
                color: {
                    bin: {
                        maxbins: 20,
                    },
                    field: "value",
                    scale: {
                        scheme: "inferno",
                        reverse: false,
                    },
                    legend: null,
                },
                tooltip: [
                    {
                        aggregate: "count",
                    },
                    {
                        field: "value",
                        bin: {
                            maxbins: 20,
                        },
                        title: "Range",
                    },
                ],
            },
            layer: [
                {
                    layer: [
                        {
                            mark: {
                                type: "bar",
                                strokeWidth: 2,
                            },
                        },
                    ],
                },
            ],
        });
    }
}
