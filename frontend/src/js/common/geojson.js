import _ from 'lodash';

class FkGeoJSONFeature {
    constructor(feature) {
        this.feature = feature;
    }

    time() {
        return new Date(this.feature.properties["timestamp"]);
    }

    coords() {
        return this.feature.geometry.coordinates;
    }

    unwrap() {
        return this.feature;
    }
}

// Love this name. I'm open to another. This basically means give each pair of
// [1, 2, 3, 4] as [1, 2] [2, 3] [3, 4] as you would imagine though also include
// [null, 1] and [4, null]
// This also means you'll get called with null, null for an empty array.
function filterPairwiseIncludingSentinels(arr, func, tx) {
    const sentinel = {};
    const heh =  [ sentinel ].concat(arr).concat([ sentinel ]);
    for (let i = 0; i < heh.length - 2; i++){
        let a = heh[i] === sentinel ? null : heh[i];
        let b = heh[i + 1] === sentinel ? null : heh[i + 1];
        if (func(a, b)) {
            return [ a, b ];
        }
    }
    return null;
}

// Simple linear interpolation.
export function map(val, domain1, domain2, range1, range2) {
    if (domain1 === domain2) {
        return range1;
    }
    return (val - domain1) / (domain2 - domain1) * (range2 - range1) + range1;
}

/**
 * Class for handling GeoJSON data specifically from the fk site. This basically
 * means timestamped GeoJSON data.
 */
export class FkGeoJSON {
    constructor(geojson) {
        // Usually sorted server-side, though this should be fine.
        this.features = _(geojson.features).map(f => new FkGeoJSONFeature(f)).sortBy(f => f.time()).value();
    }

    getFirst() {
        return this.features[0];
    }

    getLast() {
        return this.features[this.features.length - 1];
    }

    getFeaturesWithinTime(time, window) {
        const windowBegin = new Date(time.getTime() - window);
        const windowEnd = new Date(time.getTime() + window);

        return _(this.features).filter(f => {
            return f.time() >= windowBegin && f.time() < windowEnd;
        }).sortBy(f => f.time()).reverse().map(f => f.unwrap()).value();
    }

    getFeaturesSurroundingTime(time) {
        return filterPairwiseIncludingSentinels(this.features, (a, b) => {
            const beforeTime = a != null ? a.time() : 0;
            const afterTime = b != null ? b.time() : Number.MAX_SAFE_INTEGER;
            return time >= beforeTime && time < afterTime;
        });
    }

    getUniqueSourceIds() {
        return _(this.features).map(f => f.unwrap()).map('properties').map('source').map('inputId').uniq().value();
    }

    getSourcesSummary() {
        return _(this.features).map(f => f.unwrap()).groupBy(f => f.properties.source.inputId).mapValues((value) => {
            const sorted = _(value).sortBy(f => f.properties.timestamp).value();
            const lastFeature = sorted[sorted.length - 1];
            return {
                numberOfFeatures: value.length,
                lastFeature: lastFeature
            };
        }).value();
    }

    getCoordinatesAtTime(time) {
        const pair = _(this.getFeaturesSurroundingTime(time)).value();

        const beforeCoordinates = pair[0].coords();
        const afterCoordinates = pair[1].coords();
        const beforeTime = pair[0].time().getTime();
        const afterTime = pair[1].time().getTime();

        const latitude = map(time.getTime(), beforeTime, afterTime, beforeCoordinates[1], afterCoordinates[1])
        const longitude = map(time.getTime(), beforeTime, afterTime, beforeCoordinates[0], afterCoordinates[0])

        return [ longitude, latitude ];
    }
}
