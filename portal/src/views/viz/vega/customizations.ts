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
            if (thresholds && thresholds.levels && thresholds.levels.length > 0) {
                const levels = thresholds.levels;
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
        if (constrained.length > 0) {
            const range = constrained[0];
            if (s.ds.shouldConstrainBy([range.minimum, range.maximum])) {
                console.log("viz:constrained", range, s.ds.graphing.dataRange);
                spec.layer[i].encoding.y.scale.domain = [range.minimum, range.maximum];
            }
        }
    }
}
