import _ from "lodash";
import { MapFunction, ChartSettings, SeriesData, getString, getSeriesThresholds } from "./SpecFactory";

export class RangeSpecFactory {
    constructor(private readonly allSeries, private readonly settings) {}

    create() {
        const first = this.allSeries[0];

        return this.settings.apply({
            $schema: "https://vega.github.io/schema/vega-lite/v5.json",
            description: "FK Range Chart",
            config: {
                customFormatTypes: true,
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
                axisX: {
                    title: null,
                },
                axisY: {},
                view: {
                    fill: "#ffffff",
                    stroke: "transparent",
                },
            },
            data: {
                name: "table",
                values: first.queried.data,
            },
            transform: [
                {
                    bin: {
                        maxbins: 20,
                    },
                    field: "time",
                    as: "bin_time",
                },
                {
                    aggregate: [
                        {
                            op: "min",
                            field: "value",
                            as: "minimum",
                        },
                        {
                            op: "max",
                            field: "value",
                            as: "maximum",
                        },
                    ],
                    groupby: ["bin_time", "bin_time_end"],
                },
                {
                    window: [
                        {
                            op: "lead",
                            field: "bin_time_end",
                            as: "lead_bin_time_end",
                        },
                    ],
                },
            ],
            encoding: {
                x: {
                    field: "bin_time",
                    type: "temporal",
                    axis: {
                        formatType: "time",
                        tickCount: 8,
                        grid: false,
                        domain: true,
                        domainColor: "#cccccc",
                    },
                    title: null,
                },
                x2: {
                    field: "bin_time_end",
                    type: "temporal",
                },
                y: {
                    field: "minimum",
                    type: "quantitative",
                    scale: {
                        zero: false,
                    },
                    axis: {
                        grid: false,
                    },
                },
                y2: {
                    field: "maximum",
                },
                tooltip: [
                    {
                        field: "maximum",
                        formatType: "number",
                        format: ".3",
                    },
                    {
                        field: "minimum",
                        formatType: "number",
                        format: ".3",
                    },
                ],
            },
            layer: [
                {
                    layer: [
                        {
                            mark: {
                                type: "area",
                                tooltip: false,
                                interpolate: "step-after",
                                color: {
                                    x1: 1,
                                    x2: 1,
                                    y1: 1,
                                    y2: 0,
                                    gradient: "linear",
                                    stops: [
                                        {
                                            offset: 0,
                                            color: "#000004",
                                        },
                                        {
                                            offset: 0.1,
                                            color: "#170C3A",
                                        },
                                        {
                                            offset: 0.2,
                                            color: "#420A68",
                                        },
                                        {
                                            offset: 0.3,
                                            color: "#6B186E",
                                        },
                                        {
                                            offset: 0.4,
                                            color: "#932667",
                                        },
                                        {
                                            offset: 0.5,
                                            color: "#BB3754",
                                        },
                                        {
                                            offset: 0.6,
                                            color: "#DD513A",
                                        },
                                        {
                                            offset: 0.7,
                                            color: "#F3771A",
                                        },
                                        {
                                            offset: 0.8,
                                            color: "#FCA50A",
                                        },
                                        {
                                            offset: 0.9,
                                            color: "#F6D645",
                                        },
                                        {
                                            offset: 1,
                                            color: "#FCFFA4",
                                        },
                                    ],
                                },
                                strokeWidth: 1,
                                stroke: "#ffffff",
                            },
                        },
                    ],
                },
                {
                    layer: [
                        {
                            encoding: {
                                x: {
                                    field: "bin_time_end",
                                    type: "temporal",
                                },
                                x2: {
                                    field: "lead_bin_time_end",
                                    type: "temporal",
                                },
                            },
                            mark: {
                                type: "area",
                                tooltip: false,
                                interpolate: "step-before",
                                color: {
                                    x1: 1,
                                    x2: 1,
                                    y1: 1,
                                    y2: 0,
                                    gradient: "linear",
                                    stops: [
                                        {
                                            offset: 0,
                                            color: "#000004",
                                        },
                                        {
                                            offset: 0.1,
                                            color: "#170C3A",
                                        },
                                        {
                                            offset: 0.2,
                                            color: "#420A68",
                                        },
                                        {
                                            offset: 0.3,
                                            color: "#6B186E",
                                        },
                                        {
                                            offset: 0.4,
                                            color: "#932667",
                                        },
                                        {
                                            offset: 0.5,
                                            color: "#BB3754",
                                        },
                                        {
                                            offset: 0.6,
                                            color: "#DD513A",
                                        },
                                        {
                                            offset: 0.7,
                                            color: "#F3771A",
                                        },
                                        {
                                            offset: 0.8,
                                            color: "#FCA50A",
                                        },
                                        {
                                            offset: 0.9,
                                            color: "#F6D645",
                                        },
                                        {
                                            offset: 1,
                                            color: "#FCFFA4",
                                        },
                                    ],
                                },
                                strokeWidth: 1,
                                stroke: "#ffffff",
                            },
                        },
                    ],
                },
                {
                    layer: [
                        {
                            encoding: {
                                tooltip: [
                                    {
                                        field: "maximum",
                                        formatType: "number",
                                        format: ".3",
                                    },
                                    {
                                        field: "minimum",
                                        formatType: "number",
                                        format: ".3",
                                    },
                                ],
                            },
                            mark: {
                                type: "bar",
                                tooltip: true,
                                fillOpacity: 0,
                                strokeWidth: 1,
                                stroke: "#ffffff",
                            },
                        },
                    ],
                },
            ],
        });
    }
}
