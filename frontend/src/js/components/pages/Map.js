// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux'

import type { PointDecorator } from '../../../../../admin/src/js/types/VizTypes';
import type { GeoJSON } from '../../types/MapTypes';

import MapContainer from '../containers/MapContainer';
import { changePlaybackMode } from '../../actions/index';

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
    activeExpedition: PropTypes.object.isRequired,
    visibleGeoJson: GeoJSON,
};

class Map extends Component {
    props: Props;
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
        const { visibleGeoJson, playbackMode, changePlaybackMode } = this.props;
        const { pointDecorator } = this.state;

        if (!visibleGeoJson.type) {
            return (<div></div>);
        }

        return (
            <div className="map page">
                <MapContainer pointDecorator={ pointDecorator } data={ visibleGeoJson } playbackMode={ playbackMode } onChangePlaybackMode={ changePlaybackMode.bind(this) } />
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
    activeExpedition: state.activeExpedition,
    visibleGeoJson: state.visibleGeoJson,
    playbackMode: state.playbackMode,
});

export default connect(mapStateToProps, {
    changePlaybackMode 
})(Map);
