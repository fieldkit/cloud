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

export function getString(d: VizStrings | null): string | null {
    if (d) {
        return d["enUS"] || d["en-US"]; // HACK Portal compatibility.
    }
    return null;
}

export function getSeriesThresholds(series): Thresholds | null {
    if (series.vizInfo.viz.length > 0) {
        const vizConfig = series.vizInfo.viz[0];
        const thresholds = vizConfig.thresholds;
        if (thresholds && thresholds.levels && thresholds.levels.length > 0) {
            return thresholds;
        }
    }

    return null;
}

export function applySingleSensorConfiguration(spec, allSeries) {
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
    applySingleSensorConfiguration(spec, series);
}
