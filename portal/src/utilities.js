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
