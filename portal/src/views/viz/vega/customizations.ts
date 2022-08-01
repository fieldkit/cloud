import _ from "lodash";
import { expressionFunction } from "vega";
import { SeriesData, VizThresholds, VizThresholdLevel } from "../common";

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

export function getString(d: VizStrings | null): string | null {
    if (d) {
        return d["enUS"] || d["en-US"]; // HACK Portal compatibility.
    }
    return null;
}

export function getAxisLabel(level: VizThresholdLevel): string | null {
    return level.keyLabel ? getString(level.keyLabel) : getString(level.label);
}

export function getSeriesThresholds(series: SeriesData): VizThresholds | null {
    if (series.vizInfo.viz.length > 0) {
        const vizConfig = series.vizInfo.viz[0];
        const thresholds = vizConfig.thresholds;
        if (thresholds) {
            const visible = thresholds.levels.filter((d: VizThresholdLevel) => !d.hidden);
            if (thresholds && thresholds.levels && visible.length > 0) {
                return _.merge({}, thresholds, {
                    levels: visible,
                });
            }
        }
    }

    return null;
}
