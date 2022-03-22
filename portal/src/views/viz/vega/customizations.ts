import _ from "lodash";
import { expressionFunction } from "vega";

expressionFunction("fkHumanReadable", (datum, suffix) => {
    if (_.isUndefined(datum)) {
        return "N/A";
    }
    if (suffix) {
        return `${datum.toFixed(3)} ${suffix}`;
    }
    return `${datum.toFixed(3)}`;
});

type VizStrings = { [index: string]: string };

interface ThresholdLevel {
    label: VizStrings | null;
    value: number | null;
    offset: number | null;
    color: string;
}

interface Thresholds {
    label: VizStrings | null;
    levels: ThresholdLevel[];
}

function getString(d: VizStrings | null): string | null {
    if (d) {
        return d["enUS"] || d["en-US"]; // HACK Portal compatibility.
    }
    return null;
}

function getSeriesThresholds(series): Thresholds | null {
    if (series.vizInfo.viz.length > 0) {
        const vizConfig = series.vizInfo.viz[0];
        const thresholds = vizConfig.thresholds;
        if (thresholds && thresholds.levels && thresholds.levels.length > 0) {
            return thresholds;
        }
    }

    return null;
}

interface VegaDataNode {
    values: { name: string }[];
    transform: unknown;
}
interface VegaMarkNode {
    encode: {
        enter: {
            tooltip: {
                signal: string;
            };
        };
    };
}

function applyDoubleSensorConfiguration(spec, allSeries) {
    // This could be cleaned up quite a bit. I'm waiting until I've got a better
    // understanding of Vega's internals. One big opportunity here is that we
    // have two different routes through this depending on if the charted series
    // has a custom viz color scheme. These usually rely on actual value
    // thresholds rather than stops/offsets like the linear gradient approach
    // uses for the default viz.
    //
    // TODO Generalize N-series logic (remove left vs right and iterate)
    // TODO Generalize color scheme logic.
    // TODO Combine single and dual line.

    const defaultStroke = {
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

    const leftData: VegaDataNode | undefined = _.first(spec.data.filter((d) => d.name == "table1"));
    const rightData: VegaDataNode | undefined = _.first(spec.data.filter((d) => d.name == "table2"));
    if (!leftData || !rightData) throw new Error("viz: Unable to find data node(s)");

    leftData.values = allSeries[0].data;
    rightData.values = allSeries[1].data;

    leftData.values.forEach((item) => {
        item.name = "LEFT";
    });

    rightData.values.forEach((item) => {
        item.name = "RIGHT";
    });

    const markGroups = spec.marks.filter((m) => m.type == "group");

    const leftSuffix = allSeries[0].vizInfo.unitOfMeasure || "";
    const leftMarksGroup = markGroups[0];
    const leftMarks: VegaMarkNode | undefined = _.first(leftMarksGroup.marks.filter((m) => m.encode?.enter?.tooltip?.signal));
    if (leftMarks) {
        leftMarks.encode.enter.tooltip.signal = `{
            title: '${allSeries[0].vizInfo.label}',
            Value: join([round(datum.value*10)/10, '${leftSuffix}'], ' '),
            time: timeFormat(datum.time, '%m/%d/%Y %H:%m'),
        }`;
    }

    const rightSuffix = allSeries[1].vizInfo.unitOfMeasure || "";
    const rightMarksGroup = markGroups[1];
    const rightMarks: VegaMarkNode | undefined = _.first(rightMarksGroup.marks.filter((m) => m.encode?.enter?.tooltip?.signal));
    if (rightMarks) {
        rightMarks.encode.enter.tooltip.signal = `{
            title: '${allSeries[1].vizInfo.label}',
            Value: join([round(datum.value*10)/10, '${rightSuffix}'], ' '),
            time: timeFormat(datum.time, '%m/%d/%Y %H:%m'),
        }`;
    }

    spec.axes[1].title = allSeries[0].vizInfo.label;
    spec.axes[2].title = allSeries[1].vizInfo.label;

    spec.legends[0].title = allSeries[0].vizInfo.label;
    spec.legends[1].title = allSeries[1].vizInfo.label;

    const thresholdsLeft = getSeriesThresholds(allSeries[0]);
    if (thresholdsLeft) {
        const thresholdsLeftTransforms = thresholdsLeft.levels.map((d, i: number) => {
            return {
                type: "formula",
                expr: "datum.value <= " + d.value + " ? datum.value : null",
                as: "left" + i,
            };
        });

        leftData.transform = thresholdsLeftTransforms;

        const thresholdsLeftMarks = thresholdsLeft.levels.map((d, i: number) => {
            const strokewidth = 2 + (thresholdsLeft.levels.length - i) / thresholdsLeft.levels.length;
            return {
                type: "line",
                from: { data: "table1" },
                encode: {
                    enter: {
                        interpolate: { value: "cardinal" },
                        tension: { value: 0.9 },
                        strokeCap: { value: "round" },
                        x: { scale: "x", field: "time" },
                        y: { scale: "y", field: "left" + i },
                        stroke: { value: d.color },
                        strokeWidth: { value: strokewidth },
                        defined: { signal: "isValid(datum.left" + i + ")" },
                    },
                    update: {
                        strokeOpacity: {
                            signal: "hover.name == 'LEFT' ? 1 : 0.1",
                        },
                    },
                },
            };
        });

        leftMarksGroup.marks = [leftMarksGroup.marks[0]].concat(thresholdsLeftMarks.reverse()).concat(leftMarksGroup.marks[1]);

        spec.scales = spec.scales.concat([
            {
                name: "colorLeft",
                type: "ordinal",
                domain: thresholdsLeft.levels.filter((d) => d.label).map((d) => getString(d.label)),
                range: thresholdsLeft.levels.map((d) => d.color),
            },
        ]);
    } else {
        const thresholdsLeftMarks = {
            type: "line",
            from: { data: "table1" },
            encode: {
                enter: {
                    interpolate: { value: "cardinal" },
                    tension: { value: 0.9 },
                    strokeCap: { value: "round" },
                    x: { scale: "x", field: "time" },
                    y: { scale: "y", field: "value" },
                    stroke: defaultStroke,
                    strokeWidth: { value: 2 },
                    defined: { signal: "isValid(datum.value)" },
                },
                update: {
                    strokeOpacity: {
                        signal: "hover.name == 'LEFT' ? 1 : 0.1",
                    },
                },
            },
        };

        leftMarksGroup.marks = [leftMarksGroup.marks[0]].concat([thresholdsLeftMarks]).concat(leftMarksGroup.marks[1]);

        spec.scales = spec.scales.concat([
            {
                name: "colorLeft",
                type: "ordinal",
            },
        ]);
    }

    const thresholdsRight = getSeriesThresholds(allSeries[1]);
    if (thresholdsRight) {
        const thresholdsRightTransforms = thresholdsRight.levels.map((d, i: number) => {
            return {
                type: "formula",
                expr: "datum.value <= " + d.value + " ? datum.value : null",
                as: "right" + i,
            };
        });

        rightData.transform = thresholdsRightTransforms;

        const thresholdsRightMarks = thresholdsRight.levels.map((d, i: number) => {
            const strokewidth = 2 + (thresholdsRight.levels.length - i) / thresholdsRight.levels.length;
            return {
                type: "line",
                from: { data: "table2" },
                encode: {
                    enter: {
                        interpolate: { value: "cardinal" },
                        tension: { value: 0.9 },
                        strokeCap: { value: "round" },
                        x: { scale: "x2", field: "time" },
                        y: { scale: "y2", field: "right" + i },
                        stroke: { value: d.color },
                        strokeWidth: { value: strokewidth },
                        defined: { signal: "isValid(datum.right" + i + ")" },
                    },
                    update: {
                        strokeOpacity: {
                            signal: "hover.name == 'RIGHT' ? 1 : 0.1",
                        },
                    },
                },
            };
        });

        rightMarksGroup.marks = [rightMarksGroup.marks[0]].concat(thresholdsRightMarks.reverse()).concat(rightMarksGroup.marks[1]);

        spec.scales = spec.scales.concat([
            {
                name: "colorRight",
                type: "ordinal",
                domain: thresholdsRight.levels.filter((d) => d.label).map((d) => getString(d.label)),
                range: thresholdsRight.levels.map((d) => d.color),
            },
        ]);
    } else {
        const thresholdsRightMarks = {
            type: "line",
            from: { data: "table2" },
            encode: {
                enter: {
                    interpolate: { value: "cardinal" },
                    tension: { value: 0.9 },
                    strokeCap: { value: "round" },
                    x: { scale: "x2", field: "time" },
                    y: { scale: "y2", field: "value" },
                    stroke: defaultStroke,
                    strokeWidth: { value: 2 },
                    defined: { signal: "isValid(datum.value)" },
                },
                update: {
                    strokeOpacity: {
                        signal: "hover.name == 'RIGHT' ? 1 : 0.1",
                    },
                },
            },
        };

        rightMarksGroup.marks = [rightMarksGroup.marks[0]].concat([thresholdsRightMarks]).concat(rightMarksGroup.marks[1]);

        spec.scales = spec.scales.concat([
            {
                name: "colorRight",
                type: "ordinal",
            },
        ]);
    }
}

function applySingleSensorConfiguration(spec, allSeries) {
    for (let i = 0; i < allSeries.length; ++i) {
        const series = allSeries[i];

        const thresholds = getSeriesThresholds(series);
        if (thresholds) {
            console.log("viz:thresholds", thresholds);

            const levels = thresholds.levels;
            const thresholdLayers = levels
                .map((level, i) => {
                    const strokeWidth = 2 + (levels.length - i) / levels.length;

                    if (level.label) {
                        return {
                            transform: [
                                {
                                    calculate: "datum.value <= " + level.value + " ? datum.value : null",
                                    as: "layerValue" + i,
                                },
                                {
                                    calculate: "datum.layerValue" + i + " <= " + level.value + " ? '" + getString(level.label) + "' : null",
                                    as: getString(thresholds.label),
                                },
                            ],
                            encoding: {
                                y: { field: "layerValue" + i },
                                stroke: {
                                    field: getString(thresholds.label),
                                    legend: {
                                        orient: "top",
                                    },
                                    scale: {
                                        domain: levels.filter((l) => l.label).map((d) => getString(d.label)),
                                        range: levels.map((d) => d.color),
                                    },
                                },
                            },
                            mark: {
                                type: "line",
                                interpolate: "cardinal",
                                tension: 0.9,
                                strokeCap: "round",
                                strokeWidth: strokeWidth,
                            },
                        };
                    }

                    return {
                        transform: [
                            {
                                calculate: "datum.value <= " + level.value + " ? datum.value : null",
                                as: "layerValue" + i,
                            },
                        ],
                        encoding: {
                            y: { field: "layerValue" + i },
                            stroke: {
                                scale: {
                                    domain: levels.filter((l) => l.label != null).map((d) => getString(d.label)),
                                    range: levels.map((d) => d.color),
                                },
                            },
                        },
                        mark: {
                            type: "line",
                            interpolate: "cardinal",
                            tension: 0.9,
                            strokeCap: "round",
                            strokeWidth: strokeWidth,
                        },
                    };
                })
                .reverse();

            spec.layer[i].layer[1].layer = thresholdLayers;

            spec.layer[i].layer.splice(2, 1);
        }

        if (series.length == 1) {
            const suffix = series.vizInfo.unitOfMeasure;
            spec.layer[1].encoding.tooltip[0].format = suffix;
        }

        const constrained = series.vizInfo.constrainedRanges;
        if (constrained.length > 0) {
            const range = constrained[0];
            if (series.ds.shouldConstrainBy([range.minimum, range.maximum])) {
                console.log("viz:constrained", range, series.ds.graphing.dataRange);
                spec.layer[i].encoding.y.scale.domain = [range.minimum, range.maximum];
            } else {
                console.log(`viz:constrain-skip`);
            }
        } else {
            console.log(`viz:constrain-none`);
        }
    }
}

export function applySensorMetaConfiguration(spec, series) {
    if (series.length == 2) {
        applyDoubleSensorConfiguration(spec, series);
    } else {
        applySingleSensorConfiguration(spec, series);
    }
}
