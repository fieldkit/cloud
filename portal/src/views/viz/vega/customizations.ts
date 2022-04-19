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
        thresholds.levels = thresholds.levels.filter(d => !d.hidden );

        if (thresholds && thresholds.levels && thresholds.levels.length > 0) {
            return thresholds;
        }
    }

    return null;
}
