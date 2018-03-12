import type { Coordinates, Bounds, GeoJSONFeature, GeoJSON } from '../../types/MapTypes';

export function getFitBounds(geojson: GeoJSON) {
    const lon = geojson.features.map(f => f.geometry.coordinates[0]);
    const lat = geojson.features.map(f => f.geometry.coordinates[1]);
    return [[Math.min(...lon), Math.min(...lat)], [Math.max(...lon), Math.max(...lat)]];
}

export function getCenter(fitBounds: Bounds) {
    return [(fitBounds[0][0] + fitBounds[1][0]) / 2, (fitBounds[0][1] + fitBounds[1][1]) / 2];
}

export function generateConstantColor() {
    return {
        type: 'constant',
        colors: [{
            location: 0,
            color: '#C0392B'
        }],
        dateKey: null,
        bounds: null,
    };
}

export function generateLinearColor() {
    return {
        type: 'linear',
        colors: [
            {
                location: 0.0,
                color: 'rgb(255, 255, 178)'
            },
            {
                location: 0.25,
                color: 'rgb(254, 204, 92)'
            },
            {
                location: 0.50,
                color: 'rgb(253, 141, 60)'
            },
            {
                location: 0.75,
                color: 'rgb(240, 59, 32)'
            },
            {
                location: 1.0,
                color: 'rgb(189, 0, 38)'
            } ,
        ],
        dateKey: 'temp',
        bounds: null,
    };
}

export function generateConstantSize() {
    return {
        type: 'constant',
        dateKey: null,
        bounds: [10, 10],
    };
}

export function generateLinearSize() {
    return {
        type: 'linear',
        dateKey: 'temp',
        bounds: [15, 40],
    };
}

export function generatePointDecorator(colorType: string, sizeType: string) {
    return {
        points: {
            color: colorType === 'constant' ? generateConstantColor() : generateLinearColor(),
            size: sizeType === 'constant' ? generateConstantSize() : generateLinearSize(),
            sprite: 'circle.png',
        },
        title: '',
        type: 'point',
    };
}
