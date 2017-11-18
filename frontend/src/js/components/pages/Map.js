// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux'

import type { GeoJSON } from '../../types/MapTypes';
import type { ActiveExpedition  } from '../../types';

import MapContainer from '../containers/MapContainer';
import { changePlaybackMode } from '../../actions/index';

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
    activeExpedition: ActiveExpedition,
    visibleFeatures: {
        geojson: GeoJSON,
        focus: mixed
    },
    playbackMode: {},
    changePlaybackMode: () => mixed,
};

class Map extends Component {
    props: Props
    state: {
        pointDecorator: PointDecorator
    };

    constructor(props: Props) {
        super(props);
        this.state = {
            pointDecorator: generatePointDecorator('constant', 'constant')
        };
    }

    render() {
        const { visibleFeatures, playbackMode, changePlaybackMode } = this.props;
        const { pointDecorator } = this.state;

        if (!visibleFeatures.geojson) {
            return (<div></div>);
        }

        return (
            <div className="map page">
                <MapContainer pointDecorator={ pointDecorator }
                    visibleFeatures={ visibleFeatures }
                    playbackMode={ playbackMode }
                    onChangePlaybackMode={ changePlaybackMode.bind(this) } />
            </div>
        );
    }
}

const mapStateToProps = state => ({
    activeExpedition: state.activeExpedition,
    visibleFeatures: state.visibleFeatures,
    playbackMode: state.playbackMode,
});

export default connect(mapStateToProps, {
    changePlaybackMode 
})(Map);
