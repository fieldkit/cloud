import _ from "lodash";
import { MapFunction, ChartSettings, SeriesData, getString, getSeriesThresholds } from "./SpecFactory";
import chartStyles from "./chartStyles";
import { DataRow } from "../api";

export class TimeSeriesSpecFactory {
    constructor(private readonly allSeries, private readonly settings: ChartSettings = ChartSettings.Container) {}

    create() {
        const mapSeries = (mapFn: MapFunction<unknown>) => this.allSeries.map(mapFn);
        const makeDataName = (i: number) => `table${i + 1}`;
        const makeValidDataName = (i: number) => `table${i + 1}Valid`;
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
        // Early hovering behavior.
        // `hover.name == '${makeHoverName(i)}' ? ${hovering} : ${otherwise}`;
        const ifHovering = alwaysShowHovering;

        const makeSeriesThresholds = (series: SeriesData) => {
            if (solidColors) {
                return undefined;
            }
            return getSeriesThresholds(series);
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
                    // console.log("viz: constrained", series.ds.graphing.dataRange, d);
                    return d;
                } else {
                    // console.log(`viz: constrain-skip`);
                }
            } else {
                // console.log(`viz: constrain-none`);
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

        const xDomainsAll = this.allSeries.map((series: SeriesData) => series.queried.timeRange);
        const timeRangeAll = [_.min(xDomainsAll.map((dr: number[]) => dr[0])), _.max(xDomainsAll.map((dr: number[]) => dr[1]))];

        // console.log("viz: time-domain", xDomainsAll, timeRangeAll);

        const makeDomainX = () => {
            // I can't think of a good reason to just always specify this.
            return timeRangeAll;
        };

        const data = [
            {
                name: "brush_store",
            },
        ]
            .concat(
                _.flatten(
                    mapSeries((series, i) => {
                        const hoverName = makeHoverName(i);
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

                        function calculateTimeDeltas(rows: DataRow[]) {
                            const initial: { prev: DataRow | null; deltas: number[] } = {
                                prev: null,
                                deltas: [] as number[],
                            };
                            return _.reduce(
                                rows,
                                (state, datum) => {
                                    if (state.prev != null) {
                                        state.deltas.push(datum.time - state.prev.time);
                                    }
                                    return _.extend(state, { prev: datum });
                                },
                                initial
                            ).deltas;
                        }

                        function calculateGaps(rows: DataRow[]) {
                            const initial: { prev: DataRow | null; gaps: number[] } = {
                                prev: null,
                                gaps: [0] as number[],
                            };
                            return _.reduce(
                                rows,
                                (state, datum) => {
                                    if (_.isNumber(datum.value)) {
                                        if (state.gaps[state.gaps.length - 1] > 0) {
                                            state.gaps.push(0);
                                        }
                                    } else {
                                        state.gaps[state.gaps.length - 1]++;
                                    }
                                    return state;
                                },
                                initial
                            ).gaps;
                        }

                        function rebin(rows: DataRow[], interval: number): DataRow[] {
                            const grouped = _(rows)
                                .groupBy((datum) => {
                                    return Math.floor(datum.time / interval) * interval;
                                })
                                .mapValues((bin, time) => {
                                    const valid = bin.filter((datum) => _.isNumber(datum.value));
                                    if (valid.length == 1) {
                                        return valid[0];
                                    }
                                    if (valid.length > 1) {
                                        // TODO Apply bin aggregate
                                        return valid[0];
                                    }
                                    return {
                                        time: Number(time),
                                        stationId: null,
                                        moduleId: null,
                                        sensorId: null,
                                        location: null,
                                        value: null,
                                    };
                                })
                                .value();

                            return _.values(grouped);
                        }

                        // TODO We can eventually remove hoverName here
                        const properties = { name: hoverName, vizInfo: series.vizInfo };
                        const original = series.queried.data.map((datum) => _.extend(datum, properties));

                        function sanitize(original: DataRow[]): DataRow[] {
                            const valid = original.filter((datum) => _.isNumber(datum.value));
                            const deltas = calculateTimeDeltas(valid);

                            console.log("viz: gaps", calculateGaps(original));
                            console.log("viz: deltas", _.mean(deltas) / 60000);

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

                            return original;
                        }

                        const sanitized = sanitize(original);

                        return [
                            {
                                name: makeDataName(i),
                                values: sanitized,
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
                    const hoverCheck = ifHovering(i, 1, 0.2);
                    const hoverCheckGrid = ifHovering(i, 0.5, 0.2);
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

        const seriesMarks = _.flatten(
            mapSeries((series, i) => {
                const hoverCheck = ifHovering(i, 1, 0.3);
                const scales = makeScales(i);
                const thresholds = makeSeriesThresholds(series);

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
                        data: makeValidDataName(i),
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
                            fill: {
                                value: "blue",
                            },
                            fillOpacity: {
                                value: 0.5,
                            },
                        },
                        update: {
                            fillOpacity: {
                                signal: `hover && hover.stationId == datum.stationId && hover.sensorId == datum.sensorId && hover.time == datum.time ? 1 : 0.0`,
                            },
                        },
                    },
                };

                if (thresholds) {
                    const hoverCheck = ifHovering(i, 1, 0.1);
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
                                        strokeWidth: { value: strokeWidth },
                                        defined: { signal: `isValid(datum.${alias})` },
                                    },
                                    update: {
                                        strokeOpacity: {
                                            signal: hoverCheck,
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
                    const hoverCheck = ifHovering(i, 1, 0.1);
                    const lineMark = {
                        type: "line",
                        style: i === 0 ? "primaryLine" : "secondaryLine",
                        from: { data: makeDataName(i) },
                        encode: {
                            enter: {
                                interpolate: { value: "cardinal" },
                                tension: { value: 0.9 },
                                strokeCap: { value: "round" },
                                x: { scale: scales.x, field: "time" },
                                y: { scale: scales.y, field: "value" },
                                strokeWidth: { value: 2 },
                                defined: { signal: "isValid(datum.value)" },
                            },
                            update: {
                                strokeOpacity: {
                                    signal: hoverCheck,
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
        );

        // TODO Unique by thresholds and perhaps y value?
        const ruleMarks = _.flatten(
            mapSeries((series, i) => {
                const domain = makeSeriesDomain(series);
                const scales = makeScales(i);
                const thresholds = getSeriesThresholds(series);
                if (thresholds) {
                    return thresholds.levels
                        .filter((level) => level.label != null)
                        .filter((level) => level.value != null && level.value >= domain[0] && level.value < domain[1])
                        .map((level) => {
                            return {
                                type: "rule",
                                from: { data: makeDataName(i) },
                                encode: {
                                    enter: {
                                        x: {
                                            scale: scales.x,
                                            value: timeRangeAll[0],
                                        },
                                        x2: {
                                            scale: scales.x,
                                            value: timeRangeAll[1],
                                        },
                                        y: {
                                            scale: scales.y,
                                            value: level.value,
                                        },
                                        y2: {
                                            scale: scales.y,
                                            value: level.value,
                                        },
                                        stroke: { value: level.color },
                                        strokeDash: { value: [1, 4] },
                                        strokeCap: { value: "round" },
                                        opacity: { value: 0.1 },
                                        strokeOpacity: { value: 0.1 },
                                        strokeWidth: {
                                            value: 1.5,
                                        },
                                    },
                                },
                            };
                        });
                }

                return [];
            })
        );

        const cellMarks = () => {
            return [
                {
                    name: "cell",
                    type: "path",
                    from: {
                        data: "all_layouts",
                    },
                    encode: {
                        enter: {
                            path: { field: "layout_path" },
                            fill: { value: "transparent" },
                            // strokeWidth: { value: 0.35 },
                            // stroke: { value: "red" },
                            tooltip: {
                                signal: `{
                                title: datum.vizInfo.label,
                                Value: join([round(datum.value*10)/10, datum.vizInfo.unitOfMeasure || ''], ' '),
                                time: timeFormat(datum.time, '%m/%d/%Y %H:%m'),
                            }`,
                            },
                        },
                    },
                },
            ];
        };

        const marks = [...cellMarks(), ...brushMarks, ...ruleMarks, ...seriesMarks];

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
