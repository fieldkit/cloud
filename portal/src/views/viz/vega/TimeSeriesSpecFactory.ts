import _ from "lodash";
import { ChartSettings, DataRow, SeriesData, getSeriesThresholds, getAxisLabel } from "./SpecFactory";
import chartStyles from "./chartStyles";
import { makeRange, truncateTime, addDays, addSeconds } from "../common";

export interface TimeSeriesDataRow extends DataRow {
    gap: number;
    minimumGap: number;
}

export class TimeSeriesData {
    constructor(private readonly allSeries: SeriesData[]) {}
}

export class TimeSeriesSpecFactory {
    constructor(
        private readonly allSeries: SeriesData[],
        private readonly settings: ChartSettings = ChartSettings.Container,
        private readonly brushable = true,
        private readonly draggable = true
    ) {}

    create() {
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

        const isBarChart = (series: SeriesData): boolean => {
            return series.vizInfo.viz.length == 1;
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

        const xDomainsAll: number[][] = this.allSeries.map((series: SeriesData) => series.visible.toArray());

        const eachDataEnd = this.allSeries.map((series) => series.queried.dataEnd).filter((v) => v !== null);
        const dataEnd = eachDataEnd.length > 0 ? _.max(eachDataEnd) : new Date().getTime();

        const timeRangeAll =
            xDomainsAll.length == 0
                ? null
                : ([_.min(xDomainsAll.map((dr: number[]) => dr[0])), _.max(xDomainsAll.map((dr: number[]) => dr[1]))] as number[]);

        if (timeRangeAll) {
            // console.log("viz: time-domain", xDomainsAll, timeRangeAll, timeRangeAll[1] - timeRangeAll[0]);
        }

        const makeDomainX = (): number[] | null => {
            // I can't think of a good reason to just always specify this.
            return timeRangeAll;
        };

        const xDomain: number[] | null = makeDomainX();

        const filteredData = this.allSeries.map((series, i): TimeSeriesDataRow[] => {
            // TODO We can eventually remove hoverName here
            const hoverName = makeHoverName(i);
            const properties = { name: hoverName, vizInfo: series.vizInfo, series: i };
            const original = series.queried.data;

            // If a sensor has a custom filter, that information will be in the vizInfo object.
            const afterCustomFilter = series.vizInfo.applyCustomFilter(original);

            // Decorate the rows with information about the series and the vizInfo details necessary for rendering/tooltips.
            const afterProperties = afterCustomFilter.map((datum) => _.extend(datum, properties));

            // Add gap information so we can determine where missing data lies.
            const addGaps = (rows) => {
                for (let i = 0; i < rows.length; ++i) {
                    if (i == rows.length - 1) {
                        rows[i].gap = 0;
                    } else {
                        rows[i].gap = (rows[i + 1].time - rows[i].time) / 1000;
                    }
                }
                return rows;
            };

            const afterGapsAdded = addGaps(afterProperties);

            // console.log("viz: info", series.vizInfo, "gap", maybeMinimumGap, "bucket-size", series.queried.bucketSize);

            const maybeMinimumGap = series.vizInfo.minimumGap;

            const addMinimumGap = (rows) => {
                if (!maybeMinimumGap || !series.queried.bucketSize) {
                    return rows;
                }

                if (series.queried.bucketSize > maybeMinimumGap) {
                    return rows.map((datum) => _.extend(datum, { minimumGap: series.queried.bucketSize }));
                }

                return rows.map((datum) => _.extend(datum, { minimumGap: maybeMinimumGap }));
            };

            const afterMostMinimumGapAdded = addMinimumGap(afterGapsAdded);

            return afterMostMinimumGapAdded;
        });

        // This returns the domain for a single series. Primarily responsible
        // for constraining the axis domains to whatever minimums have been
        // configured for that sensor.
        const makeSeriesDomain = (series, i: number) => {
            const data = filteredData[i].map((datum) => datum.value).filter((value): value is number => _.isNumber(value));
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
        const uniqueSensorUnits = _.uniq(this.allSeries.map((series) => series.vizInfo.unitOfMeasure));
        const sameSensorUnits = uniqueSensorUnits.length == 1 && this.allSeries.length > 1;
        const yDomainsAll = this.allSeries.map((series, i: number) => makeSeriesDomain(series, i));
        const dataRangeAll = [_.min(yDomainsAll.map((dr: number[]) => dr[0])), _.max(yDomainsAll.map((dr: number[]) => dr[1]))];

        const timeLabel = "Time (" + Intl.DateTimeFormat().resolvedOptions().timeZone + ")";

        const makeDomainY = _.memoize((i: number, series) => {
            if (sameSensorUnits) {
                // console.log("viz: identical-y", dataRangeAll);
                return dataRangeAll;
            }
            return makeSeriesDomain(series, i);
        });

        const getBarConfiguration = (i: number, timeRange: number[] | null): { units: string[]; step: number | undefined } => {
            const bucketSize = this.allSeries[i].queried.bucketSize;
            const step = bucketSize > 300 ? bucketSize / 60 : 5;
            return {
                units: ["year", "month", "date", "hours", "minutes"],
                step: step,
            };
        };

        const barChartData = () => {
            return _.flatten(
                this.allSeries.map((series, i) => {
                    if (!isBarChart(series)) {
                        return [];
                    }

                    return [
                        {
                            name: makeBarDataName(i),
                            source: makeValidDataName(i),
                            transform: [
                                _.extend(
                                    {
                                        field: "time",
                                        type: "timeunit",
                                        as: ["barStartDate", "barEndDate"],
                                    },
                                    getBarConfiguration(i, timeRangeAll)
                                ),
                                {
                                    type: "formula",
                                    expr: timeRangeAll
                                        ? `time(clamp(datum.barStartDate , ${timeRangeAll[0]}, ${timeRangeAll[1]}))`
                                        : `time(datum.barStartDate)`,
                                    as: "barStart",
                                },
                                {
                                    type: "formula",
                                    expr: timeRangeAll
                                        ? `time(clamp(datum.barEndDate, ${timeRangeAll[0]}, ${timeRangeAll[1]}))`
                                        : `time(datum.barEndDate)`,
                                    as: "barEnd",
                                },
                                {
                                    type: "aggregate",
                                    groupby: ["barStart", "barEnd"],
                                    ops: ["mean"],
                                    fields: ["value"],
                                    as: ["value"],
                                },
                                {
                                    type: "formula",
                                    expr: "time(datum.barStart) + ((time(datum.barEnd) - time(datum.barStart)) / 2)",
                                    as: "barMiddle",
                                },
                                { type: "formula", expr: i, as: "series" },
                            ],
                        },
                    ];
                })
            );
        };

        const regularData = () => {
            return _.flatten(
                this.allSeries.map((series, i) => {
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
                        : [];

                    return [
                        {
                            name: makeDataName(i),
                            values: filteredData[i],
                            transform: [
                                /*
                                    This breaks dragging, for some reason.
                                    Instead of this we're using the clip
                                    functionality, which didn't work w/o custom
                                    SVG paths. The clip path otherwise was just
                                    an empty rect.
                                    {
                                        type: "filter",
                                        expr: "inrange(datum.time, visible_times)",
                                    },
                                    */
                                ...transforms,
                            ],
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
            );
        };

        const layoutData = () => {
            if (this.settings.tiny || this.settings.mobile) {
                return [];
            }
            return [
                {
                    name: "all_layouts",
                    source: this.allSeries.map((series: SeriesData, i: number) => makeValidDataName(i)),
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
            ];
        };

        const data = [
            {
                name: "brush_store",
            },
            ...regularData(),
            ...barChartData(),
            ...layoutData(),
        ];

        const legends = _.flatten(
            this.allSeries.map((series, i) => {
                if (this.settings.tiny) {
                    return [];
                }
                if (sameSensors && i > 0) {
                    return [];
                }
                return [
                    {
                        titleColor: "#2c3e50",
                        labelColor: "#2c3e50",
                        title: series.vizInfo.chartLabel,
                        stroke: makeStrokeName(i),
                        orient: "top",
                        direction: "horizontal",
                        symbolType: "stroke",
                        padding: 10,
                    },
                ];
            })
        );

        const tinyXAxis = () => {
            return [
                {
                    orient: "bottom",
                    scale: "x",
                    domain: {
                        signal: "visible_times",
                    },
                    tickCount: undefined,
                    /*
                    values: {
                        signal: "visible_times",
                    },
                    */
                    values: xDomain,
                    titleFontSize: 12,
                    titlePadding: 4,
                    titleFontWeight: "normal",
                    labelPadding: 0,
                    format: "%m/%d/%Y %H:%M",
                    encode: {
                        labels: {
                            update: {
                                align: {
                                    signal: "item === item.mark.items[0] ? 'left' : 'right'",
                                },
                            },
                        },
                    },
                },
            ];
        };

        const optionalXAxis = () => {
            if (!xDomain) {
                return [];
            }

            if (this.settings.mobile) {
                return tinyXAxis();
            }

            const getTimeTicks = () => {
                if (this.settings.estimated && xDomain) {
                    const number = Math.ceil(this.settings.estimated.w / 250);
                    const range = xDomain[1] - xDomain[0];
                    const step = range / number;
                    return _.range(0, number + 1).map((n) => xDomain[0] + n * step);
                }
                return xDomain;
            };
            const values = getTimeTicks();

            return [
                {
                    orient: "bottom",
                    scale: "x",
                    domain: {
                        signal: "visible_times",
                    },
                    tickCount: undefined,
                    labelPadding: -24,
                    tickSize: 30,
                    tickDash: [2, 2],
                    title: timeLabel,
                    values: values,
                    format: "%m/%d/%Y %H:%M",
                    encode: {
                        labels: {
                            update: {
                                align: {
                                    signal:
                                        "item === item.mark.items[0] ? 'left' : item === item.mark.items[item.mark.items.length - 1] ? 'right' : 'center'",
                                },
                            },
                        },
                    },
                },
            ];
        };

        const axes = [
            ...(this.settings.tiny ? tinyXAxis() : optionalXAxis()),
            ..._.flatten(
                this.allSeries.map((series, i) => {
                    if (i > 2) throw new Error(`viz: Too many axes`);
                    const hoverCheck = ifHovering(i, 1, 0.2);
                    const hoverCheckGrid = ifHovering(i, 0.5, 0.2);
                    const makeOrientation = (i: number) => (i == 0 ? "left" : "right");
                    const makeAxisScale = (i: number) => (i == 0 ? "y" : "y2");

                    if (this.settings.tiny) {
                        return {
                            title: series.vizInfo.axisLabel,
                            orient: makeOrientation(i),
                            scale: makeAxisScale(i),
                            domain: makeDomainY(i, series),
                            titleFontSize: 10,
                            gridDash: [],
                            tickCount: 5,
                            domainOpacity: 0,
                            titlePadding: 10,
                            titleOpacity: 1,
                            gridOpacity: 0,
                        };
                    }

                    return {
                        title: series.vizInfo.axisLabel,
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
            ),
        ];

        const standardScales = this.allSeries.map((series, i) => {
            const yDomain = makeDomainY(i, series);

            return [
                {
                    name: makeScales(i).x,
                    type: "time",
                    range: "width",
                    domain: xDomain
                        ? {
                              signal: "visible_times",
                          }
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
        });

        const thresholdScales = this.allSeries.map((series, i) => {
            const strokeName = makeStrokeName(i);
            const thresholds = getSeriesThresholds(series);
            if (thresholds) {
                return {
                    name: strokeName,
                    type: "ordinal",
                    domain: thresholds.levels.filter((d) => d.label).map((d) => getAxisLabel(d)),
                    range: thresholds.levels.map((d) => d.color),
                };
            } else {
                return {
                    name: strokeName,
                    type: "ordinal",
                    domain: [""],
                    range: [i == 0 ? "#003F5C" : "#6ef0da"], // primaryLine, secondaryLine
                };
            }
        });

        const scales = [..._.flatten(standardScales), ..._.flatten(thresholdScales)];

        const brushMarks = () => {
            if (!this.brushable) {
                return [];
            }

            return [
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
        };

        const seriesMarks = _.flatten(
            this.allSeries
                .map((series, i): any => /* TODO: Any */ {
                    const hoverCheck = ifHovering(i, 1, 0.3);
                    const scales = makeScales(i);
                    const thresholds = makeSeriesThresholds(series);

                    if (isBarChart(series)) {
                        return [
                            {
                                type: "symbol",
                                from: {
                                    data: makeBarDataName(i),
                                },
                                encode: {
                                    enter: {
                                        stroke: { value: null },
                                        strokeWidth: { value: 2 },
                                        size: { value: 100 },
                                        fill: { value: "blue" },
                                        fillOpacity: { value: 0.5 },
                                        x: { scale: scales.x, field: "barMiddle" },
                                        y: { scale: scales.y, field: "value" },
                                    },
                                    update: {
                                        fillOpacity: {
                                            signal: `hover && hover.series == datum.series && hover.time >= datum.barStart && hover.time <= datum.barEnd ? 1 : 0`,
                                        },
                                        x: { scale: scales.x, field: "barMiddle" },
                                        y: { scale: scales.y, field: "value" },
                                    },
                                },
                            },
                            {
                                type: "group",
                                marks: [
                                    {
                                        type: "rect",
                                        style: i === 0 ? "primaryBar" : "secondaryBar",
                                        from: { data: makeBarDataName(i) },
                                        encode: {
                                            enter: {
                                                x: { scale: scales.x, field: "barStart" },
                                                x2: { scale: scales.x, field: "barEnd" },
                                                y: { scale: scales.y, field: "value" },
                                                y2: { scale: scales.y, value: 0 },
                                            },
                                            update: {
                                                strokeOpacity: {
                                                    signal: hoverCheck,
                                                },
                                                x: { scale: scales.x, field: "barStart" },
                                                x2: { scale: scales.x, field: "barEnd" },
                                                y: { scale: scales.y, field: "value" },
                                                y2: { scale: scales.y, value: 0 },
                                            },
                                        },
                                    },
                                ],
                            },
                        ];
                    }

                    const symbolMarks = () => {
                        if (this.settings.mobile || this.settings.tiny) {
                            return [];
                        }
                        return [
                            {
                                type: "symbol",
                                from: {
                                    data: makeValidDataName(i),
                                },
                                encode: {
                                    enter: {
                                        stroke: { value: null },
                                        strokeWidth: { value: 2 },
                                        size: { value: 100 },
                                        fill: { value: "blue" },
                                        fillOpacity: { value: 0.5 },
                                        x: { scale: scales.x, field: "time" },
                                        y: { scale: scales.y, field: "value" },
                                    },
                                    update: {
                                        fillOpacity: {
                                            signal: `hover && hover.stationId == datum.stationId && hover.sensorId == datum.sensorId && hover.time == datum.time ? 1 : 0.0`,
                                        },
                                        x: { scale: scales.x, field: "time" },
                                        y: { scale: scales.y, field: "value" },
                                    },
                                },
                            },
                        ];
                    };

                    const dashedLineMark = {
                        type: "line",
                        from: {
                            data: makeValidDataName(i),
                        },
                        encode: {
                            enter: {
                                interpolate: { value: "cardinal" },
                                tension: { value: 0.9 },
                                stroke: { value: "#cccccc" },
                                strokeWidth: { value: 1 },
                                strokeDash: { value: [4, 4] },
                                x: { scale: scales.x, field: "time" },
                                y: { scale: scales.y, field: "value" },
                            },
                            update: {
                                strokeOpacity: {
                                    signal: hoverCheck,
                                },
                                x: { scale: scales.x, field: "time" },
                                y: { scale: scales.y, field: "value" },
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
                                    from: { data: makeValidDataName(i) },
                                    encode: {
                                        enter: {
                                            interpolate: { value: "cardinal" },
                                            tension: { value: 0.9 },
                                            strokeCap: { value: "round" },
                                            strokeWidth: { value: strokeWidth },
                                            defined: {
                                                signal: `!datum.minimumGap || datum.${alias} <= datum.minimumGap`,
                                            },
                                            x: { scale: scales.x, field: "time" },
                                            y: { scale: scales.y, field: alias },
                                        },
                                        update: {
                                            strokeOpacity: {
                                                signal: hoverCheck,
                                            },
                                            x: { scale: scales.x, field: "time" },
                                            y: { scale: scales.y, field: alias },
                                        },
                                    },
                                };
                            })
                            .reverse();

                        return [
                            {
                                type: "group",
                                clip: {
                                    path: { signal: "chart_clip" },
                                },
                                marks: [dashedLineMark, ...thresholdsMarks, ...symbolMarks()],
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
                                    strokeWidth: { value: 2 },
                                    defined: { signal: `!datum.minimumGap || datum.gap <= datum.minimumGap` },
                                    strokeOpacity: {
                                        signal: hoverCheck,
                                    },
                                    x: { scale: scales.x, field: "time" },
                                    y: { scale: scales.y, field: "value" },
                                },
                                update: {
                                    x: { scale: scales.x, field: "time" },
                                    y: { scale: scales.y, field: "value" },
                                },
                            },
                        };
                        return [
                            {
                                type: "group",
                                clip: {
                                    path: { signal: "chart_clip" },
                                },
                                marks: [dashedLineMark, lineMark, ...symbolMarks()],
                            },
                        ];
                    }
                })
                .reverse()
        );

        // define mark stacking priority
        const getMarkSortIndex = (d) => {
            if (d["marks"]) {
                if (d["marks"][0].type === "line") return 1;
                if (d["marks"][0].type === "rect") return 2;
                else return 0;
            } else return 0;
        };

        // TODO Unique by thresholds and perhaps y value?
        const ruleMarks = () => {
            return _.flatten(
                this.allSeries.map((series, i) => {
                    const domain = makeSeriesDomain(series, i);
                    const scales = makeScales(i);
                    const thresholds = getSeriesThresholds(series);
                    if (thresholds && timeRangeAll) {
                        return thresholds.levels
                            .filter((level) => level.label != null && !level.hidden)
                            .filter((level) => level.start != null && level.start >= domain[0] && level.start < domain[1])
                            .map((level) => {
                                return {
                                    type: "rule",
                                    encode: {
                                        enter: {
                                            x: { scale: scales.x, value: timeRangeAll[0] },
                                            x2: { scale: scales.x, value: timeRangeAll[1] },
                                            y: { scale: scales.y, value: level.start },
                                            y2: { scale: scales.y, value: level.start },
                                            stroke: { value: level.color },
                                            strokeDash: { value: [1, 4] },
                                            strokeCap: { value: "round" },
                                            opacity: { value: 0.5 },
                                            strokeOpacity: { value: 0.5 },
                                            strokeWidth: {
                                                value: 1.5,
                                            },
                                        },
                                        update: {
                                            x: { scale: scales.x, value: timeRangeAll[0] },
                                            x2: { scale: scales.x, value: timeRangeAll[1] },
                                            y: { scale: scales.y, value: level.start },
                                            y2: { scale: scales.y, value: level.start },
                                        },
                                    },
                                };
                            });
                    }

                    return [];
                })
            );
        };

        const twelveHourMarks = () => {
            if (!this.settings.mobile) {
                return [];
            }

            if (!timeRangeAll) {
                return [];
            }

            const timeRange = (timeRangeAll[1] - timeRangeAll[0]) / 1000;
            if (timeRange > 86400 * 4) {
                return [];
            }

            const marksStart = addDays(truncateTime(new Date(timeRangeAll[0])), -2);
            const marksEnd = addDays(truncateTime(new Date(timeRangeAll[1])), 2);
            const elapsed = marksEnd.getTime() - marksStart.getTime();
            const halves = (elapsed / 1000 / 86400) * 2;

            const scales = makeScales(0);

            return [
                {
                    type: "group",
                    clip: {
                        path: { signal: "chart_clip" },
                    },
                    marks: _.range(0, halves)
                        .map((i) => {
                            return addSeconds(new Date(marksStart), i * (86400 / 2));
                        })
                        .map((mark) => {
                            // console.log("viz:", mark, mark.getTime());

                            return {
                                type: "rule",
                                encode: {
                                    enter: {
                                        x: { scale: scales.x, value: mark.getTime() },
                                        x2: { scale: scales.x, value: mark.getTime() },
                                        y: { value: 0 },
                                        y2: { signal: "height" },
                                        stroke: { value: "#ddd" },
                                        strokeWidth: { value: 1 },
                                    },
                                    update: {
                                        x: { scale: scales.x, value: mark.getTime() },
                                        x2: { scale: scales.x, value: mark.getTime() },
                                        y: { value: 0 },
                                        y2: { signal: "height" },
                                    },
                                },
                            };
                        }),
                },
            ];
        };

        const cellMarks = () => {
            if (this.settings.tiny || this.settings.mobile) {
                return [];
            }

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
                                unitOfMeasure: datum.vizInfo.unitOfMeasure,
                                value: datum.value,
                                time: timeFormat(datum.time, '%m/%d/%Y %H:%M'),
                                name: datum.name
                            }`,
                            },
                        },
                        update: {
                            path: { field: "layout_path" },
                        },
                    },
                },
            ];
        };

        const marks = [
            ...cellMarks(),
            ...brushMarks(),
            ...ruleMarks(),
            ...twelveHourMarks(),
            ...seriesMarks.sort((a, b) => {
                return getMarkSortIndex(b) - getMarkSortIndex(a);
            }),
        ];

        const brushSignals = [
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
            {
                name: "visible_times",
                value: xDomain || [0, 1],
            },
            {
                name: "chart_background",
                value: "transparent",
            },
        ];

        const calculateDragLimits = () => {
            const timeNow = new Date().getTime();
            const endOffset = timeRangeAll ? timeRangeAll[1] - timeRangeAll[0] : 0;
            if (timeRangeAll && dataEnd) {
                if (timeRangeAll[0] > dataEnd || timeRangeAll[1] > dataEnd) {
                    return { start: timeNow - endOffset, end: timeNow };
                }
            }
            const start = dataEnd && timeRangeAll ? dataEnd - endOffset : timeNow;
            const end = dataEnd || timeNow;
            return { start, end };
        };

        const { start: dragMaximumStart, end: dragMaximumEnd } = calculateDragLimits();

        const dragSignals = [
            {
                name: "down",
                value: null,
                on: [
                    { events: "touchend", update: "null" },
                    { events: "mousedown, touchstart", update: "xy()" },
                ],
            },
            {
                name: "drag_maximum_start",
                value: dragMaximumStart,
            },
            {
                name: "drag_maximum_end",
                value: dragMaximumEnd,
            },
            {
                name: "drag_delta",
                value: [0, 0],
                on: [
                    {
                        events: [
                            {
                                type: "touchmove",
                                consume: true,
                                filter: "event.touches.length === 1",
                            },
                        ],
                        update: "down ? [down[0] - x(), y() - down[1]] : [0, 0]",
                    },
                ],
            },
            {
                name: "drag_x",
                value: [],
                on: [
                    {
                        events: {
                            signal: "drag_delta",
                        },
                        update: "[down[0] + drag_delta[0], down[1] + drag_delta[1]]",
                    },
                ],
            },
            {
                name: "xcur",
                value: xDomain || [0, 1],
                on: [
                    /*
                    {
                        events: "mousedown, touchstart, touchend",
                        update: "slice(xdom)",
                    },
                    */
                ],
            },
            {
                name: "visible_times_calc",
                value: xDomain || [0, 1],
                on: [
                    {
                        events: { signal: "drag_delta" },
                        update:
                            "[ceil(min(xcur[0] + span(xcur) * drag_delta[0] / width, drag_maximum_start)), ceil(min(xcur[1] + span(xcur) * drag_delta[0] / width, drag_maximum_end))]",
                    },
                ],
            },
            {
                name: "visible_times",
                value: xDomain || [0, 1],
                on: [
                    {
                        events: { signal: "visible_times_calc" },
                        update: "visible_times_calc",
                    },
                ],
            },
            {
                name: "drag_time",
                on: [
                    {
                        events: {
                            signal: "down",
                        },
                        update: "down === null ? visible_times : null",
                    },
                ],
            },
            {
                name: "chart_background",
                value: "transparent",
                on: [
                    {
                        events: {
                            signal: "down",
                        },
                        update: "down === null ? 'transparent' : '#efefef'",
                    },
                ],
            },
        ];

        const containerSignals = () => {
            return [
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
                    name: "chart_clip_size",
                    init: '[scale("x", visible_times[1]), containerSize()[1]]',
                    on: [
                        {
                            events: { scale: "x" },
                            update: '[scale("x", visible_times[1]), containerSize()[1]]',
                        },
                    ],
                },
                {
                    name: "chart_clip",
                    init: "'M 0 0 H ' + (chart_clip_size[0]) + ' V ' + (chart_clip_size[1]) + ' H 0 Z'",
                    on: [
                        {
                            events: { signal: "chart_clip_size" },
                            update: "'M 0 0 H ' + (chart_clip_size[0]) + ' V ' + (chart_clip_size[1]) + ' H 0 Z'",
                        },
                    ],
                },
            ];
        };

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
            {
                name: "visible_times_calc",
                value: xDomain || [0, 1],
            },
            {
                name: "visible_times",
                value: xDomain || [0, 1],
            },
            {
                name: "chart_background",
                value: "transparent",
            },
            {
                name: "drag_time",
                value: null,
            },
        ];

        const interactiveSignals = () => {
            const standard = [
                {
                    name: "hover",
                    on: [
                        {
                            events: "@cell:mouseover",
                            update: "datum",
                        },
                        {
                            events: "@cell:mouseout",
                            update: "null",
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
            ];

            if (this.brushable) {
                return [...containerSignals(), ...standard, ...brushSignals];
            }

            if (this.draggable) {
                return [...containerSignals(), ...standard, ...dragSignals];
            }

            throw new Error();
        };

        const allSignals = () => {
            if (this.settings.tiny) {
                const tinySignals = [...containerSignals(), ...staticSignals];
                return tinySignals;
            }

            if (this.settings.size.w == 0) {
                if (this.brushable || this.draggable) {
                    return interactiveSignals();
                }
            }

            return [...containerSignals(), ...staticSignals];
        };

        return this.buildSpec(
            allSignals() as never[],
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
                background: "#ffffff",
                group: {
                    fill: {
                        signal: "chart_background",
                    },
                },
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
                legend: {
                    layout: {
                        top: {
                            anchor: "left",
                            direction: this.settings.mobile ? "vertical" : "horizontal",
                            margin: 0,
                        },
                    },
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
}
