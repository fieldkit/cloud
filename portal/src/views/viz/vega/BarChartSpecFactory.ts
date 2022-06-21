import _ from "lodash";
import { MapFunction, ChartSettings, SeriesData, getString, getSeriesThresholds } from "./SpecFactory";
import chartStyles from "./chartStyles";
import { DataRow } from "../api";
import { makeRange } from "../common";

export class BarChartSpecFactory {
    constructor(private readonly allSeries, private readonly settings: ChartSettings = ChartSettings.Container) {}

    create() {
        const mapSeries = (mapFn: MapFunction<unknown>) => this.allSeries.map(mapFn);
        const makeDataName = (i: number) => `table${i + 1}`;
        const makeValidDataName = (i: number) => `table${i + 1}Valid`;
        const makeBarDataName = (i: number) => `table${i + 1}Bar`;
        const makeStrokeName = (i: number) => `color${i ? "Right" : "Left"}`;
        const makeHoverName = (i: number) => `${i ? "RIGHT" : "LEFT"}`;
        const makeThresholdLevelAlias = (i: number, l: number) => `${i ? "right" : "left"}${l}`;
        const makeScales = (i: number) => {
            if (i == 0) {
                return {
                    x: "x",
                    y: "y",
                };
            } else {
                return {
                    x: "x2",
                    y: "y2",
                };
            }
        };

        const solidColors = true;

        // Always showing hovering state.
        const alwaysShowHovering = (i: number, hovering: any, otherwise: any) => `${hovering}`;
        const ifHovering = alwaysShowHovering;

        const makeSeriesThresholds = (series: SeriesData) => {
            if (solidColors) {
                return undefined;
            }
            return getSeriesThresholds(series);
        };

        const filteredData = mapSeries((series, i) => {
            // TODO We can eventually remove hoverName here
            const hoverName = makeHoverName(i);
            const properties = { name: hoverName, vizInfo: series.vizInfo };
            const original = series.queried.data;

            function sanitize(original: DataRow[]): DataRow[] {
                // This was the first approach we tried to prevent
                // the infilled missing values from creating
                // distracting graphs of too many valid data
                // islands.
                /*
                    const valid = original.filter((datum) => _.isNumber(datum.value));
                    const deltas = calculateTimeDeltas(valid);

                    // console.log("viz: gaps", calculateGaps(original));
                    // console.log("viz: deltas", _.mean(deltas) / 60000);

                    // Very simple heuristic for enabling re-bin. We
                    // basically rebin to the average interval if more
                    // of the data is invalid than valid. I think we can
                    // do better.
                    if (original.length - valid.length > valid.length) {
                        const meanBetweenValid = _.mean(deltas);
                        const interval = Math.ceil(meanBetweenValid / 60000) * 60000;
                        console.log("viz: rebin", interval);
                        return rebin(original, interval);
                    } else {
                        console.log("viz: rebin-skip", original.length - valid.length, valid.length);
                    }
                */

                // This is the another approach we're trying,
                // basically remove missing data if the queried
                // resolution is finer than the expected sensor
                // data's resolution.
                if (series.queried.shouldIgnoreMissing) {
                    return original.filter((datum) => _.isNumber(datum.value));
                }

                return original;
            }

            return sanitize(/*removeOutliersHack(original)*/ original).map((datum) =>
                _.extend(datum, properties, { timeDate: new Date(datum.time) })
            );
        });

        // This returns the domain for a single series. Primarily responsible
        // for constraining the axis domains to whatever minimums have been
        // configured for that sensor.
        const makeSeriesDomain = (series, i: number) => {
            const data = filteredData[i].filter((datum) => _.isNumber(datum.value)).map((datum) => datum.value);
            const filteredRange = data.length > 0 ? makeRange(data) : [0, 0];
            const constrained = series.vizInfo.constrainedRanges;
            if (series.ds.graphing && constrained.length > 0) {
                const range = constrained[0];
                if (series.ds.shouldConstrainBy(filteredRange, [range.minimum, range.maximum])) {
                    const d = [range.minimum, range.maximum];
                    // console.log("viz: constrained", series.ds.graphing.dataRange, d);
                    return d;
                } else {
                    // console.log(`viz: constrain-skip`);
                }
            } else {
                // console.log(`viz: constrain-none`);
            }
            return filteredRange;
        };

        // Are the sensors being charted the same? If they are then we should
        // use the same axis domain for both, and pick one that covers both.
        const uniqueSensorKeys = _.uniq(this.allSeries.map((series) => series.vizInfo.key));
        const sameSensors = uniqueSensorKeys.length == 1 && this.allSeries.length > 1;
        const yDomainsAll = this.allSeries.map((series, i: number) => makeSeriesDomain(series, i));
        const dataRangeAll = [_.min(yDomainsAll.map((dr: number[]) => dr[0])), _.max(yDomainsAll.map((dr: number[]) => dr[1]))];

        const makeDomainY = _.memoize((i: number, series) => {
            if (sameSensors) {
                console.log("viz: identical-y", dataRangeAll);
                return dataRangeAll;
            }
            return makeSeriesDomain(series, i);
        });

        const xDomainsAll = this.allSeries.map((series: SeriesData) => series.queried.timeRange);
        const timeRangeAll = [_.min(xDomainsAll.map((dr: number[]) => dr[0])), _.max(xDomainsAll.map((dr: number[]) => dr[1]))];

        // console.log("viz: time-domain", xDomainsAll, timeRangeAll);

        const makeDomainX = () => {
            // I can't think of a good reason to just always specify this.
            return timeRangeAll;
        };

        const xDomain = makeDomainX();

        const data = [
            {
                name: "brush_store",
            },
        ]
            .concat(
                _.flatten(
                    mapSeries((series, i) => {
                        const thresholds = makeSeriesThresholds(series);
                        const scales = makeScales(i);
                        const transforms = thresholds
                            ? thresholds.levels.map((level, l: number) => {
                                  return {
                                      type: "formula",
                                      expr: "datum.value <= " + level.value + " ? datum.value : null",
                                      as: makeThresholdLevelAlias(i, l),
                                  };
                              })
                            : null;

                        return [
                            {
                                name: makeDataName(i),
                                values: filteredData[i],
                                transform: transforms,
                            },
                            {
                                name: makeValidDataName(i),
                                source: makeDataName(i),
                                transform: [
                                    {
                                        type: "filter",
                                        expr: "isValid(datum.value)",
                                    },
                                    {
                                        type: "formula",
                                        expr: `scale('${scales.x}', datum.time)`,
                                        as: "layout_x",
                                    },
                                    {
                                        type: "formula",
                                        expr: `scale('${scales.y}', datum.value)`,
                                        as: "layout_y",
                                    },
                                ],
                            },
                            {
                                name: makeBarDataName(i),
                                source: makeValidDataName(i),
                                transform: [
                                    {
                                        field: "timeDate",
                                        type: "timeunit",
                                        units: ["year", "month", "date", "hours"],
                                        as: ["barStartDate", "barEndDate"],
                                    },
                                    {
                                        type: "formula",
                                        expr: "time(datum.barEndDate)",
                                        as: "barEnd",
                                    },
                                    {
                                        type: "formula",
                                        expr: "time(datum.barStartDate)",
                                        as: "barStart",
                                    },
                                ],
                            },
                        ];
                    })
                )
            )
            .concat([
                {
                    name: "all_layouts",
                    source: this.allSeries.map((series, i) => makeValidDataName(i)),
                    transform: [
                        {
                            type: "voronoi",
                            x: "layout_x",
                            y: "layout_y",
                            size: [{ signal: "width" }, { signal: "height" }],
                            as: "layout_path",
                        },
                    ],
                },
            ] as any[]);

        const legends = _.flatten(
            mapSeries((series, i) => {
                if (sameSensors && i > 0) {
                    return [];
                }
                return [
                    {
                        titleColor: "#2c3e50",
                        labelColor: "#2c3e50",
                        title: series.vizInfo.label,
                        stroke: makeStrokeName(i),
                        orient: "top",
                        direction: "horizontal",
                        symbolType: "stroke",
                        padding: 10,
                    },
                ];
            })
        );

        const axes = [
            {
                orient: "bottom",
                scale: "x",
                domain: xDomain,
                tickCount: 6,
                labelPadding: -24,
                tickSize: 30,
                tickDash: [2, 2],
                title: "Time",
                format: {
                    year: "%m/%d/%Y",
                    quarter: "%m/%d/%Y",
                    month: "%m/%d/%Y",
                    week: "%m/%d/%Y",
                    date: "%m/%d/%Y",
                    hours: "%m/%d/%Y %H:%M",
                    minutes: "%m/%d/%Y %H:%M",
                    seconds: "%m/%d/%Y %H:%M",
                    milliseconds: "%m/%d/%Y %H:%M",
                },
            },
        ].concat(
            _.flatten(
                mapSeries((series, i) => {
                    if (i > 2) throw new Error(`viz: Too many axes`);
                    const hoverCheck = ifHovering(i, 1, 0.2);
                    const hoverCheckGrid = ifHovering(i, 0.5, 0.2);
                    const makeOrientation = (i: number) => (i == 0 ? "left" : "right");
                    const makeAxisScale = (i: number) => (i == 0 ? "y" : "y2");

                    return {
                        title: series.vizInfo.label,
                        orient: makeOrientation(i),
                        scale: makeAxisScale(i),
                        domain: makeDomainY(i, series),
                        tickCount: 5,
                        titlePadding: 10,
                        domainOpacity: 0,
                        gridDash: [],
                        titleOpacity: {
                            signal: hoverCheck,
                        },
                        gridOpacity: {
                            signal: hoverCheckGrid,
                        },
                        labelOpacity: {
                            signal: hoverCheck,
                        },
                    };
                })
            )
        );

        const scales = _.flatten(
            mapSeries((series, i) => {
                const yDomain = makeDomainY(i, series);

                return [
                    {
                        name: makeScales(i).x,
                        type: "time",
                        range: "width",
                        domain: xDomain
                            ? xDomain
                            : {
                                  data: makeDataName(i),
                                  field: "time",
                              },
                    },
                    {
                        name: makeScales(i).y,
                        type: "linear",
                        range: "height",
                        nice: true,
                        zero: true,
                        domain: yDomain
                            ? yDomain
                            : {
                                  data: makeDataName(i),
                                  field: "value",
                              },
                    },
                ];
            })
        ).concat(
            mapSeries((series, i) => {
                const strokeName = makeStrokeName(i);
                const thresholds = getSeriesThresholds(series);
                if (thresholds) {
                    return {
                        name: strokeName,
                        type: "ordinal",
                        domain: thresholds.levels.filter((d) => d.label).map((d) => getString(d.label)),
                        range: thresholds.levels.map((d) => d.color),
                    };
                } else {
                    return {
                        name: strokeName,
                        type: "ordinal",
                    };
                }
            })
        );

        const brushMarks = [
            {
                name: "brush_brush_bg",
                type: "rect",
                clip: true,
                encode: {
                    enter: {
                        fill: {
                            value: "#333",
                        },
                        fillOpacity: {
                            value: 0.125,
                        },
                    },
                    update: {
                        x: [
                            {
                                test: 'data("brush_store").length',
                                signal: "brush_x[0]",
                            },
                            {
                                value: 0,
                            },
                        ],
                        y: [
                            {
                                test: 'data("brush_store").length',
                                value: 0,
                            },
                            {
                                value: 0,
                            },
                        ],
                        x2: [
                            {
                                test: 'data("brush_store").length',
                                signal: "brush_x[1]",
                            },
                            {
                                value: 0,
                            },
                        ],
                        y2: [
                            {
                                test: 'data("brush_store").length',
                                field: {
                                    group: "height",
                                },
                            },
                            {
                                value: 0,
                            },
                        ],
                    },
                },
            },
            {
                name: "brush_brush",
                type: "rect",
                clip: true,
                encode: {
                    enter: {
                        fill: {
                            value: "transparent",
                        },
                    },
                    update: {
                        x: [
                            {
                                test: 'data("brush_store").length',
                                signal: "brush_x[0]",
                            },
                            {
                                value: 0,
                            },
                        ],
                        y: [
                            {
                                test: 'data("brush_store").length',
                                value: 0,
                            },
                            {
                                value: 0,
                            },
                        ],
                        x2: [
                            {
                                test: 'data("brush_store").length',
                                signal: "brush_x[1]",
                            },
                            {
                                value: 0,
                            },
                        ],
                        y2: [
                            {
                                test: 'data("brush_store").length',
                                field: {
                                    group: "height",
                                },
                            },
                            {
                                value: 0,
                            },
                        ],
                        stroke: [
                            {
                                test: "brush_x[0] !== brush_x[1]",
                                value: "white",
                            },
                            {
                                value: null,
                            },
                        ],
                    },
                },
            },
        ];

        const seriesMarks = _.flatten(
            mapSeries((series, i) => {
                const hoverCheck = ifHovering(i, 1, 0.3);
                const scales = makeScales(i);
                return [
                    {
                        type: "group",
                        marks: [
                            {
                                type: "rect",
                                style: i === 0 ? "primaryBar" : "secondaryBar",
                                from: { data: makeBarDataName(i) },
                                encode: {
                                    enter: {
                                        x2: { scale: scales.x, field: "barStart", offset: 1 },
                                        x: { scale: scales.x, field: "barEnd" },
                                        y: { scale: scales.y, field: "value" },
                                        y2: { scale: scales.y, value: 0 },
                                    },
                                    update: {
                                        strokeOpacity: {
                                            signal: hoverCheck,
                                        },
                                    },
                                },
                            },
                        ],
                    },
                ];
            })
        );

        const marks = [...brushMarks, ...seriesMarks];

        const interactiveSignals = [
            {
                name: "width",
                init: "containerSize()[0]",
                on: [
                    {
                        events: "window:resize",
                        update: "containerSize()[0]",
                    },
                ],
            },
            {
                name: "height",
                init: "containerSize()[1]",
                on: [
                    {
                        events: "window:resize",
                        update: "containerSize()[1]",
                    },
                ],
            },
            {
                name: "hover",
                on: [
                    {
                        events: "@cell:mouseover",
                        update: "datum",
                    },
                    /*
                    {
                        events: "@cell:mouseout",
                        update: "null",
                    },
                    */
                ],
            },
            {
                name: "unit",
                value: {},
                on: [
                    {
                        events: "mousemove",
                        update: "isTuple(group()) ? group() : unit",
                    },
                ],
            },
            {
                name: "brush",
                update: 'vlSelectionResolve("brush_store", "union")',
            },
            {
                name: "brush_x",
                value: [],
                on: [
                    {
                        events: {
                            source: "view",
                            type: "mousedown",
                            filter: ['!event.item || event.item.mark.name !== "brush_brush"'],
                        },
                        update: "[x(unit), x(unit)]",
                    },
                    {
                        events: {
                            source: "window",
                            type: "mousemove",
                            consume: true,
                            between: [
                                {
                                    source: "view",
                                    type: "mousedown",
                                    filter: ['!event.item || event.item.mark.name !== "brush_brush"'],
                                },
                                {
                                    source: "window",
                                    type: "mouseup",
                                },
                            ],
                        },
                        update: "[brush_x[0], clamp(x(unit), 0, width)]",
                    },
                    {
                        events: {
                            signal: "brush_scale_trigger",
                        },
                        update: '[scale("x", brush_time[0]), scale("x", brush_time[1])]',
                    },
                    {
                        events: [
                            {
                                source: "view",
                                type: "dblclick",
                            },
                        ],
                        update: "[0, 0]",
                    },
                    {
                        events: {
                            signal: "brush_translate_delta",
                        },
                        update:
                            "clampRange(panLinear(brush_translate_anchor.extent_x, brush_translate_delta.x / span(brush_translate_anchor.extent_x)), 0, width)",
                    },
                    {
                        events: {
                            signal: "brush_zoom_delta",
                        },
                        update: "clampRange(zoomLinear(brush_x, brush_zoom_anchor.x, brush_zoom_delta), 0, width)",
                    },
                ],
            },
            {
                name: "brush_time",
                on: [
                    {
                        events: {
                            signal: "brush_x",
                        },
                        update: 'brush_x[0] === brush_x[1] ? null : invert("x", brush_x)',
                    },
                ],
            },
            {
                name: "brush_scale_trigger",
                value: {},
                on: [
                    {
                        events: [
                            {
                                scale: "x",
                            },
                        ],
                        update:
                            '(!isArray(brush_time) || (+invert("x", brush_x)[0] === +brush_time[0] && +invert("x", brush_x)[1] === +brush_time[1])) ? brush_scale_trigger : {}',
                    },
                ],
            },
            {
                name: "brush_tuple",
                on: [
                    {
                        events: [
                            {
                                signal: "brush_time",
                            },
                        ],
                        update: "brush_time ? {fields: brush_tuple_fields, values: [brush_time]} : null",
                    },
                ],
            },
            {
                name: "brush_tuple_fields",
                value: [
                    {
                        field: "time",
                        channel: "x",
                        type: "R",
                    },
                ],
            },
            {
                name: "brush_translate_anchor",
                value: {},
                on: [
                    {
                        events: [
                            {
                                source: "view",
                                type: "mousedown",
                                markname: "brush_brush",
                            },
                        ],
                        update: "{x: x(unit), y: y(unit), extent_x: slice(brush_x)}",
                    },
                ],
            },
            {
                name: "brush_translate_delta",
                value: {},
                on: [
                    {
                        events: [
                            {
                                source: "window",
                                type: "mousemove",
                                consume: true,
                                between: [
                                    {
                                        source: "view",
                                        type: "mousedown",
                                        markname: "brush_brush",
                                    },
                                    {
                                        source: "window",
                                        type: "mouseup",
                                    },
                                ],
                            },
                        ],
                        update: "{x: brush_translate_anchor.x - x(unit), y: brush_translate_anchor.y - y(unit)}",
                    },
                ],
            },
            {
                name: "brush_zoom_anchor",
                on: [
                    {
                        events: [
                            {
                                source: "view",
                                type: "wheel",
                                consume: true,
                                markname: "brush_brush",
                            },
                        ],
                        update: "{x: x(unit), y: y(unit)}",
                    },
                ],
            },
            {
                name: "brush_zoom_delta",
                on: [
                    {
                        events: [
                            {
                                source: "view",
                                type: "wheel",
                                consume: true,
                                markname: "brush_brush",
                            },
                        ],
                        force: true,
                        update: "pow(1.001, event.deltaY * pow(16, event.deltaMode))",
                    },
                ],
            },
            {
                name: "brush_modify",
                on: [
                    {
                        events: {
                            signal: "brush_tuple",
                        },
                        update: 'modify("brush_store", brush_tuple, true)',
                    },
                ],
            },
        ];

        const staticSignals = [
            {
                name: "hover",
                value: {
                    name: makeHoverName(0),
                },
            },
            {
                name: "brush_x",
                value: [],
            },
        ];

        const allSignals = this.settings.w == 0 ? interactiveSignals : staticSignals;

        return this.buildSpec(
            allSignals as never[],
            data as never[],
            scales as never[],
            axes as never[],
            legends as never[],
            marks as never[],
            false
        );
    }

    private buildSpec(signals: never[], data: never[], scales: never[], axes: never[], legends: never[], marks: never[], verbose = false) {
        if (verbose) {
            console.log("viz: signals", signals);
            console.log("viz: data", data);
            console.log("viz: scales", scales);
            console.log("viz: axes", axes);
            console.log("viz: legends", legends);
            console.log("viz: marks", marks);
        }

        return this.settings.apply({
            $schema: "https://vega.github.io/schema/vega/v5.json",
            description: "FK Time Series",
            padding: 5,
            config: {
                background: "white",
                axis: {
                    labelFont: "Avenir Light",
                    labelFontSize: 12,
                    labelColor: "#2c3e50",
                    titleColor: "#2c3e50",
                    titleFont: "Avenir Light",
                    titleFontSize: 14,
                    titlePadding: 20,
                    tickSize: 10,
                    tickOpacity: 0,
                    domain: false,
                    grid: true,
                    gridDash: [],
                },
                style: chartStyles,
            },
            signals: signals,
            data: data,
            scales: scales,
            legends: legends,
            axes: axes,
            marks: marks,
        });
    }
}
