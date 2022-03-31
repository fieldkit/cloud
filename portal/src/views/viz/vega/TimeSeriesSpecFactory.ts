import _ from "lodash";
import { MapFunction, ChartSettings, SeriesData, getString, getSeriesThresholds } from "./SpecFactory";

export class TimeSeriesSpecFactory {
    constructor(private readonly allSeries, private readonly settings: ChartSettings = new ChartSettings(0, 0)) {}

    create() {
        const makeDataName = (i: number) => `table${i + 1}`;
        const makeValidDataName = (i: number) => `table${i + 1}Valid`;
        const makeStrokeName = (i: number) => `color${i ? "Right" : "Left"}`;
        const makeHoverName = (i: number) => `${i ? "RIGHT" : "LEFT"}`;
        const makeThresholdLevelAlias = (i: number, l: number) => `${i ? "right" : "left"}${l}`;
        const mapSeries = (mapFn: MapFunction<unknown>) => this.allSeries.map(mapFn);
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

        // This returns the domain for a single series. Primarily responsible
        // for constraining the axis domains to whatever minimums have been
        // configured for that sensor.
        const makeSeriesDomain = (series) => {
            const constrained = series.vizInfo.constrainedRanges;
            if (series.ds.graphing && constrained.length > 0) {
                const range = constrained[0];
                if (series.ds.shouldConstrainBy([range.minimum, range.maximum])) {
                    const d = [range.minimum, range.maximum];
                    console.log("viz: constrained", series.ds.graphing.dataRange, d);
                    return d;
                } else {
                    console.log(`viz: constrain-skip`);
                }
            } else {
                console.log(`viz: constrain-none`);
            }
            return series.queried.dataRange;
        };

        // Are the sensors being charted the same? If they are then we should
        // use the same axis domain for both, and pick one that covers both.
        const uniqueSensorKeys = _.uniq(this.allSeries.map((series) => series.vizInfo.key));
        const sameSensors = uniqueSensorKeys.length == 1 && this.allSeries.length > 1;
        const yDomainsAll = this.allSeries.map(makeSeriesDomain);
        const dataRangeAll = [_.min(yDomainsAll.map((dr: number[]) => dr[0])), _.max(yDomainsAll.map((dr: number[]) => dr[1]))];

        const makeDomainY = _.memoize((series) => {
            if (sameSensors) {
                console.log("viz: identical-y", dataRangeAll);
                return dataRangeAll;
            }
            return makeSeriesDomain(series);
        });

        const makeDomainX = () => {
            if (this.allSeries.length > 1) {
                const xDomainsAll = this.allSeries.map((series) => series.queried.timeRange);
                const timeRangeAll = [_.min(xDomainsAll.map((dr: number[]) => dr[0])), _.max(xDomainsAll.map((dr: number[]) => dr[1]))];
                console.log("viz: domain", xDomainsAll, timeRangeAll);
                return timeRangeAll;
            }
            return undefined;
        };

        const data = [
            {
                name: "brush_store",
            },
        ].concat(
            _.flatten(
                mapSeries((series, i) => {
                    const thresholds = getSeriesThresholds(series);
                    const hoverName = makeHoverName(i);
                    const transforms = thresholds
                        ? thresholds.levels.map((level, l: number) => {
                              return {
                                  type: "formula",
                                  expr: "datum.value <= " + level.value + " ? datum.value : null",
                                  as: makeThresholdLevelAlias(i, l),
                              };
                          })
                        : null;

                    series.queried.data.forEach((item) => {
                        item.name = hoverName;
                    });

                    return [
                        {
                            name: makeDataName(i),
                            values: series.queried.data,
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
                            ],
                        },
                    ];
                })
            )
        );

        const legends = _.flatten(
            mapSeries((series, i) => {
                const hoverName = makeHoverName(i);
                const hoverCheck = `hover.name == '${hoverName}' ? 1 : 0.1`;
                return {
                    titleColor: "#2c3e50",
                    labelColor: "#2c3e50",
                    title: series.vizInfo.label,
                    stroke: makeStrokeName(i),
                    orient: "top",
                    direction: "horizontal",
                    symbolType: "stroke",
                    padding: 10,
                    labelOpacity: {
                        signal: hoverCheck,
                    },
                    symbolOpacity: {
                        signal: hoverCheck,
                    },
                    titleOpacity: {
                        signal: hoverCheck,
                    },
                };
            })
        );

        const axes = [
            {
                orient: "bottom",
                scale: "x",
                domain: makeDomainX(),
                tickCount: 8,
                labelPadding: -24,
                tickSize: 30,
                tickDash: [2, 2],
            },
        ].concat(
            _.flatten(
                mapSeries((series, i) => {
                    if (i > 2) throw new Error(`viz: Too many axes`);
                    const hoverName = makeHoverName(i);
                    const hoverCheck = `hover.name == '${hoverName}' ? 1 : 0.2`;
                    const makeOrientation = (i: number) => (i == 0 ? "left" : "right");
                    const makeAxisScale = (i: number) => (i == 0 ? "y" : "y2");

                    return {
                        title: series.vizInfo.label,
                        orient: makeOrientation(i),
                        scale: makeAxisScale(i),
                        domain: makeDomainY(series),
                        tickCount: 5,
                        titlePadding: 10,
                        domainOpacity: 0,
                        titleOpacity: {
                            signal: hoverCheck,
                        },
                        gridOpacity: {
                            signal: hoverCheck,
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
                const yDomain = makeDomainY(series);
                const xDomain = makeDomainX();

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

        const marks = brushMarks.concat(
            _.flatten(
                mapSeries((series, i) => {
                    const hoverName = makeHoverName(i);
                    const hoverCheck = `hover.name == '${hoverName}' ? 1 : 0.3`;
                    const scales = makeScales(i);
                    const title = series.vizInfo.label;
                    const suffix = series.vizInfo.unitOfMeasure || "";
                    const thresholds = getSeriesThresholds(series);

                    const firstLineMark = {
                        type: "line",
                        from: {
                            data: makeValidDataName(i),
                        },
                        encode: {
                            enter: {
                                interpolate: {
                                    value: "cardinal",
                                },
                                tension: {
                                    value: 0.9,
                                },
                                x: {
                                    scale: scales.x,
                                    field: "time",
                                },
                                y: {
                                    scale: scales.y,
                                    field: "value",
                                },
                                stroke: {
                                    value: "#cccccc",
                                },
                                strokeWidth: {
                                    value: 1,
                                },
                                strokeDash: {
                                    value: [4, 4],
                                },
                            },
                            update: {
                                strokeOpacity: {
                                    signal: hoverCheck,
                                },
                            },
                        },
                    };

                    const symbolMark = {
                        type: "symbol",
                        from: {
                            data: makeDataName(i),
                        },
                        encode: {
                            enter: {
                                x: {
                                    scale: scales.x,
                                    field: "time",
                                },
                                y: {
                                    scale: scales.y,
                                    field: "value",
                                },
                                stroke: {
                                    value: null,
                                },
                                strokeWidth: {
                                    value: 2,
                                },
                                size: {
                                    value: 100,
                                },
                                tooltip: {
                                    signal: `{
                                        title: '${title}',
                                        Value: join([round(datum.value*10)/10, '${suffix}'], ' '),
                                        time: timeFormat(datum.time, '%m/%d/%Y %H:%m'),
                                    }`,
                                },
                                fill: {
                                    value: "blue",
                                },
                                fillOpacity: {
                                    value: 0,
                                },
                            },
                            hover: {
                                fillOpacity: {
                                    value: 0.5,
                                },
                            },
                            update: {
                                fillOpacity: {
                                    value: 0,
                                },
                            },
                        },
                    };

                    if (thresholds) {
                        const thresholdsMarks = thresholds.levels
                            .map((level, l: number) => {
                                const alias = makeThresholdLevelAlias(i, l);
                                const strokeWidth = 2 + (thresholds.levels.length - l) / thresholds.levels.length;

                                return {
                                    type: "line",
                                    from: { data: makeDataName(i) },
                                    encode: {
                                        enter: {
                                            interpolate: { value: "cardinal" },
                                            tension: { value: 0.9 },
                                            strokeCap: { value: "round" },
                                            x: { scale: scales.x, field: "time" },
                                            y: { scale: scales.y, field: alias },
                                            stroke: { value: level.color },
                                            strokeWidth: { value: strokeWidth },
                                            defined: { signal: `isValid(datum.${alias})` },
                                        },
                                        update: {
                                            strokeOpacity: {
                                                signal: `hover.name == '${hoverName}' ? 1 : 0.1`,
                                            },
                                        },
                                    },
                                };
                            })
                            .reverse();

                        return [
                            {
                                type: "group",
                                marks: _.concat([firstLineMark], thresholdsMarks as never[], [symbolMark] as never[]),
                            },
                        ];
                    } else {
                        const lineMark = {
                            type: "line",
                            from: { data: makeDataName(i) },
                            encode: {
                                enter: {
                                    interpolate: { value: "cardinal" },
                                    tension: { value: 0.9 },
                                    strokeCap: { value: "round" },
                                    x: { scale: scales.x, field: "time" },
                                    y: { scale: scales.y, field: "value" },
                                    stroke: this.defaultStroke(),
                                    strokeWidth: { value: 2 },
                                    defined: { signal: "isValid(datum.value)" },
                                },
                                update: {
                                    strokeOpacity: {
                                        signal: `hover.name == '${hoverName}' ? 1 : 0.1`,
                                    },
                                },
                            },
                        };
                        return [
                            {
                                type: "group",
                                marks: [firstLineMark, lineMark, symbolMark],
                            },
                        ];
                    }
                })
            )
        );

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
                value: {
                    name: makeHoverName(0),
                },
                on: [
                    {
                        events: "symbol:mouseover",
                        update: "datum",
                    },
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
                        update: "clampRange(panLinear(brush_translate_anchor.extent_x, brush_translate_delta.x / span(brush_translate_anchor.extent_x)), 0, width)",
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
                        update: '(!isArray(brush_time) || (+invert("x", brush_x)[0] === +brush_time[0] && +invert("x", brush_x)[1] === +brush_time[1])) ? brush_scale_trigger : {}',
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
                    gridDash: [2, 2],
                },
            },
            signals: signals,
            data: data,
            scales: scales,
            legends: legends,
            axes: axes,
            marks: marks,
        });
    }

    private defaultStroke() {
        return {
            value: {
                x1: 1,
                y1: 1,
                x2: 1,
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
        };
    }
}
