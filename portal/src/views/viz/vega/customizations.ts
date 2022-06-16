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

export function getSeriesThresholds(series: SeriesData): VizThresholds | null {
    if (series.vizInfo.viz.length > 0) {
        const vizConfig = series.vizInfo.viz[0];
        const thresholds = vizConfig.thresholds;
        if (thresholds) {
            thresholds.levels = thresholds.levels.filter((d: VizThresholdLevel) => !d.hidden);

            if (thresholds && thresholds.levels && thresholds.levels.length > 0) {
                return thresholds;
            }
        }
    }

    return null;
}
