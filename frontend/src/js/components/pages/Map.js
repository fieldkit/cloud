// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux'

import MapContainer from '../containers/MapContainer';
import type { PointDecorator } from '../../../../../admin/src/js/types/VizTypes';
import type { GeoJSON } from '../../types/MapTypes';
import { FieldKitLogo } from '../icons/Icons';
import '../../../css/map.css';

function generateConstantColor() {
    return {
        type: 'constant',
        colors: [{
            location: 0,
            color: '#facada'
        }],
        dateKey: null,
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
        dateKey: 'temp',
        bounds: null,
    };
}

function generateConstantSize() {
    return {
        type: 'constant',
        dateKey: null,
        bounds: [10, 10],
    };
}

function generateLinearSize() {
    return {
        type: 'linear',
        dateKey: 'temp',
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

type Props = {
    activeExpedition: PropTypes.object.isRequired
};

class Map extends Component {
    props: Props;
    state: {
        pointDecorator: PointDecorator,
        data: GeoJSON,
    };

    constructor(props: Props) {
        super(props);
        this.state = {
            pointDecorator: {},
            data: {}
        } ;
    }

    componentWillMount() {
        const data = {
            type: "FeatureCollection",
            features: [{
                type: "Feature",
                geometry: {
                    type: "Point",
                    coordinates: [-77.02827,38.91427]
                },
                properties: {
                    temp: 32
                }
            },
            {
                type: "Feature",
                geometry: {
                    type: "Point",
                    coordinates: [-77.02013, 38.91538]
                },
                properties: {
                    temp: 45
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
                        <FieldKitLogo />
                    </div>
                </div>
            </div>
        );
    }
}

const mapStateToProps = state => ({
    activeExpedition: state.activeExpedition
});

export default connect(mapStateToProps, {
})(Map);
