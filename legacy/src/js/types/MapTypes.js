// @flow weak

export type Coordinates = [number, number];

export type Bounds = [Coordinates, Coordinates];

export type GeoJSONFeature = {
    type: string,
    properties: {},
    geometry: {
        type: string,
        coordinates: Coordinates,
    },
};

export type GeoJSON = {
    type: string,
    features: GeoJSONFeature[],
};
