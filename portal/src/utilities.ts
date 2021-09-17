import _ from "lodash";

export function serializePromiseChain(all, fn) {
    return all.reduce((accum, value, index) => {
        return accum.then((allValues) => {
            return Promise.resolve(fn(value, index)).then((singleValue) => {
                allValues.push(singleValue);
                return allValues;
            });
        });
    }, Promise.resolve([]));
}

export function promiseAfter(t, v) {
    return new Promise(function(resolve) {
        setTimeout(resolve.bind(null, v), t);
    });
}

export function hexStringToByteWiseString(str) {
    return str
        .split("")
        .map((c, i) => {
            return (i + 1) % 2 == 0 ? c + " " : c;
        })
        .join("");
}

export function convertBytesToLabel(bytes) {
    // convert to kilobytes or megabytes
    if (bytes < 1000000.0) {
        return Math.round(bytes / 1024.0) + " KB";
    }
    return Math.round(bytes / 1048576.0) + " MB";
}

export function unixNow() {
    return Math.round(new Date().getTime() / 1000);
}

export function convertOldFirmwareResponse(module) {
    // compensate for old firmware
    if (module.name.indexOf("modules") != 0) {
        let name = "modules." + module.name;
        if (name == "modules.water") {
            // this is dicey, but temporary...
            if (module.sensorObjects) {
                name += "." + module.sensorObjects[0].name;
            }
        }
        return name;
    }
    return module.name;
}

export function getModuleImg(module) {
    switch (module.name) {
        case "modules.distance":
        case "distance":
            return "modules/icon-module-distance.svg";
        case "modules.weather":
        case "weather":
            return "modules/icon-module-weather.svg";
        case "modules.water.ec":
        case "water.ec":
            return "modules/icon-module-water-conductivity.svg";
        case "modules.water.ph":
        case "water.ph":
            return "modules/icon-module-water-ph.svg";
        case "modules.water.do":
        case "water.do":
            return "modules/icon-module-dissolved-oxygen.svg";
        case "modules.water.temp":
        case "water.temp":
            return "modules/icon-module-water-temp.svg";
        case "modules.water.orp":
        case "water.orp":
            return "modules/icon-module-water.svg";
        case "modules.water.unknown":
        case "water.unknown":
            return "modules/icon-module-water.svg";
        case "ttn.floodnet":
            return "modules/icon-module-floodnet.png";
        default:
            return "modules/icon-module-generic.svg";
    }
}

export function getBatteryIcon(percentage: number | undefined): string {
    if (!percentage) {
        return "battery/0.svg";
    }

    if (percentage == 0) {
        return "battery/0.svg";
    } else if (percentage <= 20) {
        return "battery/20.svg";
    } else if (percentage <= 40) {
        return "battery/40.svg";
    } else if (percentage <= 60) {
        return "battery/60.svg";
    } else if (percentage <= 80) {
        return "battery/80.svg";
    } else {
        return "battery/100.svg";
    }
}

export function getRunTime(project) {
    const timeUnits = ["seconds", "minutes", "hours", "days", "weeks", "months", "years"];
    const start = new Date(project.start_time);
    let end, runTense;
    if (project.end_time) {
        end = new Date(project.end_time);
        runTense = "Ran for ";
    } else {
        // assume it's still running?
        end = new Date();
        runTense = "Running for ";
    }
    // get interval and convert to seconds
    const interval = (end.getTime() - start.getTime()) / 1000;
    let displayValue = interval;
    let unit = 0;
    // unit is an index into timeUnits
    if (interval < 60) {
        // already set to seconds
    } else if (interval < 3600) {
        // minutes
        unit = 1;
        displayValue /= 60;
        displayValue = Math.round(displayValue);
    } else if (interval < 86400) {
        // hours
        unit = 2;
        displayValue /= 3600;
        displayValue = Math.round(displayValue);
    } else if (interval < 604800) {
        // days
        unit = 3;
        displayValue /= 86400;
        displayValue = Math.round(displayValue);
    } else if (interval < 2628000) {
        // weeks
        unit = 4;
        displayValue /= 604800;
        displayValue = Math.round(displayValue);
    } else if (interval < 31535965) {
        // months
        unit = 5;
        displayValue /= 2628000;
        displayValue = Math.round(displayValue);
    } else {
        // years
        unit = 6;
        displayValue /= 31535965;
        displayValue = Math.round(displayValue);
    }
    return runTense + displayValue + " " + timeUnits[unit];
}

export function tryParseTags(rawTags: string) {
    const sanitized = rawTags.trim();
    if (sanitized.length == 0) {
        return [];
    }
    if (sanitized[0] == "[" || sanitized[0] == "{") {
        try {
            const hopefullyArray = JSON.parse(sanitized);
            const array = _.isArray(hopefullyArray) ? hopefullyArray : [hopefullyArray];
            return array.map((text) => {
                return {
                    text: text,
                };
            });
        } catch (error) {
            console.log(`invalid tags field: '${sanitized}'`);
        }
    }
    return sanitized.split(" ").map((text) => {
        return {
            text: text,
        };
    });
}

export function toSingleValue(v: null | string | (string | null)[]): string | null {
    if (v) {
        if (_.isArray(v) && v.length > 0 && v[0]) {
            return v[0];
        }
        return v as string;
    }
    return null;
}

export function vueTickHack(callback: () => void): void {
    setTimeout(callback, 0);
}
