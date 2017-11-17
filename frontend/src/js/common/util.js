// @flow weak

// import {Location} from '../../../../../caboodle/src/proto/flow';
import { Location } from '../common/flow';
import D3Scatterplot from '../components/visualizations/D3Scatterplot';
import D3TimeSeries from '../components/visualizations/D3TimeSeries';

export function getProjectSlug(): string {
    return window.location.hostname.split('.')[0];
}

export function chartStringToClass(type, props) {
    switch (type) {
        case 'scatterplot':
            return (new D3Scatterplot(props));
        case 'time-series':
            return (new D3TimeSeries(props));
        default:
            return null;
    }
}

export function getMapSampleData() {
    return [[
        {
            id: 0,
            place: 'The coffee bar',
            login: 'espresso',
            temp: 0,
            location: new Location({
                latitude: 38.91427,
                longitude: -77.02827,
            }),
        },
        {
            id: 1,
            place: 'Bistro Bohem',
            login: '2027355895',
            location: new Location({
                latitude: 38.91538,
                longitude: -77.02013,
            }),
            temp: 20,
        },
        {
            id: 2,
            place: 'Black Cat',
            login: 'luckycat',
            location: new Location({
                latitude: 38.91458,
                longitude: -77.03155,
            }),
            temp: 100,
        },
        {
            id: 3,
            place: 'Snap',
            login: 'nutella1',
            location: new Location({
                latitude: 38.92239,
                longitude: -77.04227,
            }),
            temp: 60,
        },
        {
            id: 4,
            place: 'Columbia Heights Coffee',
            login: 'FAIRTRADE1',
            location: new Location({
                latitude: 38.93222,
                longitude: -77.02854,
            }),
            temp: 50,
        },
        {
            id: 5,
            place: 'Azi\'s Cafe',
            login: 'sunny',
            location: new Location({
                latitude: 38.90842,
                longitude: -77.02419,
            }),
            temp: 30,
        },
        {
            id: 1,
            place: 'Blind Dog Cafe',
            login: 'baxtercantsee',
            location: new Location({
                latitude: 38.91931,
                longitude: -77.02518,
            }),
            temp: 80,
        },
        {
            id: 6,
            place: 'Le Caprice',
            login: 'baguette',
            location: new Location({
                latitude: 38.93260,
                longitude: -77.03304,
            }),
            temp: 20,
        },
        {
            id: 7,
            place: 'Filter',
            login: '',
            location: new Location({
                latitude: 38.91368,
                longitude: -77.04509,
            }),
            temp: 10,
        },
        {
            id: 8,
            place: 'Peregrine',
            login: 'espresso',
            location: new Location({
                latitude: 38.88516,
                longitude: -76.99656,
            }),
            temp: 100,
        },
        {
            id: 9,
            place: 'Tryst',
            login: 'coupetnt',
            location: new Location({
                latitude: 38.921894,
                longitude: -77.042438,
            }),
            temp: 50,
        },
        {
            id: 10,
            place: 'The Coupe',
            login: 'voteforus',
            location: new Location({
                latitude: 38.93206,
                longitude: -77.02821,
            }),
            temp: 20,
        },
        {
            id: 11,
            place: 'Big Bear Cafe',
            login: '',
            location: new Location({
                latitude: 38.91275,
                longitude: -77.01239,
            }),
            temp: 10,
        },
    ]];
}
