export function getUpdatedDate(station) {
    let date = null;
    // try more accurate first: last_uploads
    if (station.last_uploads) {
        const dataUploads = station.last_uploads.filter(u => {
            return u.type == "data";
        });
        date = dataUploads[0].time;
    } else {
        if (!station.status_json) {
            return "N/A";
        }
        date = station.status_json.updated;
    }
    if (!date) {
        return "N/A";
    }
    const d = new Date(date);
    return d.toLocaleDateString("en-US");
}

export function serializePromiseChain(all, fn) {
    return all.reduce((accum, value, index) => {
        return accum.then(allValues => {
            return fn(value, index).then(singleValue => {
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
        module.name = "modules." + module.name;
        if (module.name == "modules.water") {
            // this is dicey, but temporary...
            module.name += "." + module.sensorObjects[0].name;
        }
    }
    return module.name;
}

export function getModuleImg(module) {
    switch (module.name) {
        case "modules.distance":
        case "distance":
            return "Icon_Distance_Module.png";
        case "modules.weather":
        case "weather":
            return "Icon_Weather_Module.png";
        case "modules.water.ec":
        case "water.ec":
            return "Icon_WaterConductivity_Module.png";
        case "modules.water.ph":
        case "water.ph":
            return "Icon_WaterpH_Module.png";
        case "modules.water.do":
        case "water.do":
            return "Icon_DissolvedOxygen_Module.png";
        case "modules.water.temp":
        case "water.temp":
            return "Icon_WaterTemp_Module.png";
        case "modules.water.orp":
        case "water.orp":
            return "Icon_Water_Module.png";
        case "modules.water.unknown":
        case "water.unknown":
            return "Icon_Water_Module.png";
        default:
            return "Icon_Generic_Module.png";
    }
}

export function getRunTime(project) {
    const timeUnits = ["seconds", "minutes", "hours", "days", "weeks", "months", "years"];
    let start = new Date(project.start_time);
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
