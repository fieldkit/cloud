import _ from "lodash";
import { expressionFunction } from "vega";

expressionFunction("fkHumanReadable", (datum) => {
    if (_.isUndefined(datum)) {
        return "N/A";
    }
    /*
    if (this.valueSuffix) {
        return `${datum.toFixed(3)} ${this.valueSuffix}`;
    }
    */
    return `${datum.toFixed(3)}`;
});

function getString(d) {
    return d["enUS"] || d["en-US"]; // HACK Portal compatibility.
}

export interface SensorDetails {
    key: string;
    internal: boolean;
    fullKey: string;
    firmwareKey: string;
    unitOfMeasure: string;
    ranges: { minimum: number; maximum: number }[];
}

export function applySensorMetaConfiguration(spec, series) {
    for (let i = 0; i < series.length; ++i) {
        const s = series[i];
        if (s.vizInfo.viz.length == 0) {
            continue;
        }

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

        const constrained = s.vizInfo.constrainedRanges;
        if (s.ds.constrainDataAxis && constrained.length > 0) {
            const range = constrained[0];
            spec.layer[i].encoding.y.scale.domain = [range.minimum, range.maximum];
        }
    }
}
