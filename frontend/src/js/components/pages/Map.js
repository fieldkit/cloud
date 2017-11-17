// @flow weak

import React, { Component } from 'react';
import log from 'loglevel';
import MapContainer from '../containers/MapContainer';
import type { PointDecorator } from '../../../../../admin/src/js/types/VizTypes';
import type { GeoJSON } from '../../types/MapTypes';
import { FieldkitLogo } from '../icons/Icons';
import '../../../css/map.css';

function generateConstantColor() {
    return {
        type: 'constant',
        colors: [{
            location: 0,
            color: '#facada'
        }],
        data_key: null,
        bounds: null,
    };
}

function generateLinearColor() {
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
        data_key: 'temp',
        bounds: null,
    };
}

function generateConstantSize() {
    return {
        type: 'constant',
        data_key: null,
        bounds: [10, 10],
    };
}

function generateLinearSize() {
    return {
        type: 'linear',
        data_key: 'temp',
        bounds: [15, 40],
    };
}

function generatePointDecorator(colorType: string, sizeType: string) {
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

type Props = {};

export default class Map extends Component {
    props: Props;
    state: {
    pointDecorator: PointDecorator,
    data: GeoJSON,
    };

    constructor(props: Props) {
        super(props);
        this.state = {
            pointDecorator: {},
            data: {},
            feature: null,
        } ;
    }

    componentWillMount() {
        /*
        const testSelection = generateSelectionsFromMap({
          id: 'number',
          place: 'string',
          login: 'string',
          location: 'location',
          temp: 'number',
        });

        const testGrouping = {
          operation: 'equal',
          parameter: null,
          source_attribute: 'id',
        };

        const rawData = getMapSampleData();
        const transformedData = transform(rawData, testGrouping, testSelection);
        const data = streamToGeoJSON(transformedData, 'location');
        */

        const data = {
            type: "FeatureCollection",
            features: [{
                type: "Feature",
                geometry: {
                    type: "Point",
                    coordinates: [-77.02827,38.91427]
                },
                properties: {
                    temp: 23
                }
            },
            {
                type: "Feature",
                geometry: {
                    type: "Point",
                    coordinates: [-77.02013, 38.91538]
                },
                properties: {
                    temp: 23
                }
            }]
        }

        const pointDecorator = generatePointDecorator('linear', 'linear');

        this.setState({
            data,
            pointDecorator
        });
    }

    render() {
        const { pointDecorator, data } = this.state;

        return (
            <div className="map page">
                <MapContainer pointDecorator={ pointDecorator } data={ data } />
                <div className="disclaimer-panel">
                    <div className="disclaimer-body">
                        <span className="b">NOTE: </span> Map images have been obtained from a third-party and do not reflect the editorial decisions of National Geographic.
                    </div>
                    <div className="fieldkit-logo">
                        <FieldkitLogo />
                    </div>
                </div>
            </div>
            );
    }
}
