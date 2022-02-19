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

function getString(d) {
    return d["enUS"] || d["en-US"]; // HACK Portal compatibility.
}

export function applySensorMetaConfiguration(spec, series) {
    for (let i = 0; i < series.length; ++i) {
        const s = series[i];

        if (s.vizInfo.viz.length > 0) {
            const vizConfig = s.vizInfo.viz[0];
            const thresholds = vizConfig.thresholds;
            const levels = thresholds.levels;
            if (levels && levels.length > 0) {
                const thresholdLayers = levels
                    .map((d, i) => {
                        return {
                            transform: [
                                {
                                    calculate: "datum.value <= " + d.value + " ? datum.value : null",
                                    as: "layerValue" + i,
                                },
                                {
                                    calculate: "datum.layerValue" + i + " <= " + d.value + " ? '" + getString(d.label) + "' : null",
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
                                        domain: levels.map((d) => getString(d.label)),
                                        range: levels.map((d) => d.color),
                                    },
                                },
                            },
                            mark: {
                                type: "line",
                                interpolate: "monotone",
                                tension: 1,
                            },
                        };
                    })
                    .reverse();

                spec.layer[i].layer = thresholdLayers;
            }
        }

        if (series.length == 1) {
            const suffix = s.vizInfo.unitOfMeasure;
            spec.layer[1].encoding.tooltip[0].format = suffix;
        }

        const constrained = s.vizInfo.constrainedRanges;
        if (s.ds.constrainDataAxis && constrained.length > 0) {
            const range = constrained[0];
            spec.layer[i].encoding.y.scale.domain = [range.minimum, range.maximum];
        }
    }
}
