import _ from "lodash";
import { expressionFunction } from "vega";

// import fieldkitHumidityData from "./fieldkitHumidityData.json";
// import fieldkitTemperatureData from "./fieldkitTemperatureData.json";

expressionFunction("fkHumanReadable", (datum, suffix) => {
    if (_.isUndefined(datum)) {
        return "N/A";
    }
    if (suffix) {
        return `${datum.toFixed(3)} ${suffix}`;
    }
    return `${datum.toFixed(3)}`;
});

function getString(d) {
    return d["enUS"] || d["en-US"]; // HACK Portal compatibility.
}

function getSeriesThresholds(series) {
    if (series.vizInfo.viz.length > 0) {
        const vizConfig = series.vizInfo.viz[0];
        const thresholds = vizConfig.thresholds;
        if (thresholds && thresholds.levels && thresholds.levels.length > 0) {
            return thresholds;
        }
    }

    return {
        title: null,
        levels: [
            {
                label: null,
                value: 10000,
                color: "#afefaf",
            },
        ],
    };
}

function applyDoubleSensorConfiguration(spec, allSeries) {
    // This could be cleaned up quite a bit. I'm waiting until I've got a better
    // understanding of Vega's internals.

    spec.data[0].values = allSeries[0].data;
    spec.data[1].values = allSeries[1].data;

    spec.data[0].values.forEach((item) => {
        item.name = "LEFT";
    });

    spec.data[1].values.forEach((item) => {
        item.name = "RIGHT";
    });

    const leftSuffix = allSeries[0].vizInfo.unitOfMeasure || "";
    spec.marks[0].marks[1].encode.enter.tooltip.signal = `{
        title: '${allSeries[0].vizInfo.label}',
        Value: join([round(datum.value*10)/10, '${leftSuffix}'], ' '),
        time: timeFormat(datum.time, '%m/%d/%Y %H:%m'),
    }`;

    const rightSuffix = allSeries[1].vizInfo.unitOfMeasure || "";
    spec.marks[1].marks[1].encode.enter.tooltip.signal = `{
        title: '${allSeries[1].vizInfo.label}',
        Value: join([round(datum.value*10)/10, '${rightSuffix}'], ' '),
        time: timeFormat(datum.time, '%m/%d/%Y %H:%m'),
    }`;

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

        spec.data[0].transform = thresholdsLeftTransforms;

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

        spec.marks[0].marks = [spec.marks[0].marks[0]].concat(thresholdsLeftMarks.reverse()).concat(spec.marks[0].marks[1]);

        spec.scales = spec.scales.concat([
            {
                name: "colorLeft",
                type: "ordinal",
                domain: thresholdsLeft.levels.filter((d) => d.label).map((d) => getString(d.label)),
                range: thresholdsLeft.levels.map((d) => d.color),
            },
        ]);
    } else {
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

        spec.data[1].transform = thresholdsRightTransforms;

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
                        x: { scale: "x", field: "time" },
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

        spec.marks[1].marks = [spec.marks[1].marks[0]].concat(thresholdsRightMarks.reverse()).concat(spec.marks[1].marks[1]);

        spec.scales = spec.scales.concat([
            {
                name: "colorRight",
                type: "ordinal",
                domain: thresholdsRight.levels.filter((d) => d.label).map((d) => getString(d.label)),
                range: thresholdsRight.levels.map((d) => d.color),
            },
        ]);
    } else {
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
