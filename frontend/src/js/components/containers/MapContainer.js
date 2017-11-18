// @flow weak

import React, { Component } from 'react';
import ReactMapboxGl, { ScaleControl, ZoomControl, Popup } from 'react-mapbox-gl';

import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../../secrets'

import BubbleMap from '../visualizations/BubbleMap';
import PlaybackControl from '../PlaybackControl';

import type { Coordinates, Bounds, GeoJSONFeature, GeoJSON } from '../../types/MapTypes';
import type { PointDecorator } from '../../../../../admin/src/js/types/VizTypes';

import { FieldKitLogo } from '../icons/Icons';

type Props = {
    pointDecorator: PointDecorator,
    playback: null,
    onChangePlaybackMode: PropTypes.func.isRequired,
    data: GeoJSON,
};

function getFitBounds(data) {
    const lon = data.features.map(f => f.geometry.coordinates[0]);
    const lat = data.features.map(f => f.geometry.coordinates[1]);
    return [[Math.min(...lon), Math.min(...lat)], [Math.max(...lon), Math.max(...lat)]];
}

function getCenter(fitBounds: Bounds) {
    return [(fitBounds[0][0] + fitBounds[1][0]) / 2, (fitBounds[0][1] + fitBounds[1][1]) / 2];
}

export default class MapContainer extends Component {
    props: Props;
    state: {
        fitBounds: Bounds,
        center: Coordinates,
        zoom: [number],
        feature: ?GeoJSONFeature,
    };

    constructor(props: Props) {
        super(props);
        this.state = {
            fitBounds: [[0, 0], [0, 0]],
            center: [0, 0],
            zoom: [12],
            feature: null
        };
    }

    componentWillMount() {
        const { data } = this.props;
        const fitBounds = getFitBounds(data);
        const center = getCenter(fitBounds);
        this.setState({
            fitBounds,
            center
        });
    }

    componentWillReceiveProps(nextProps) {
        const { data: oldData } = this.props;
        const { data: newData } = nextProps;

        if (oldData !== newData) {
            const { focus: oldFocus } = oldData;
            const { focus: newFocus } = newData;

            if (oldFocus !== newFocus) {
                this.setState({
                    center: newFocus.geometry.coordinates
                })
            }
        }
    }

    onMarkerClick(feature) {
        this.setState({
            feature
        });
    }

    onPopupChange() {
        let { feature } = this.state;
        if (feature) { // If we have a popup and we're changing to another one.
            feature = null;
            this.setState({
                feature
            });
        }
    }

    tick(firstFrame) {
        requestAnimationFrame(() => this.tick(false));
    }

    componentDidMount() {
        requestAnimationFrame(() => this.tick(true));
    }

    renderProperties(feature) {
        return (
            <table>
                <thead>
                    <tr>
                        <th>Key</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                { Object.entries(feature.properties).map(([ k , v ]) => (
                    <tr key={k}>
                        <td> {k} </td>
                        <td> {v} </td>
                    </tr>
                )) }
                </tbody>
            </table>
        )
    }

    render() {
        const { pointDecorator, data, onChangePlaybackMode, playbackMode } = this.props;
        const { fitBounds, center, zoom, feature } = this.state;

        if (!data.features) {
            return (<div><h1>Loading</h1></div>);
        }
        else if (data.features.length === 0) {
            return (<div><h1>Empty</h1></div>);
        }

        return (
            <div>
                <ReactMapboxGl accessToken={MAPBOX_ACCESS_TOKEN}
                    style={MAPBOX_STYLE}
                    movingMethod="easeTo" fitBounds={ fitBounds } center={ center }
                    zoom={ zoom } onClick={ this.onPopupChange.bind(this) } containerStyle={ { height: '100vh', width: '100vw', } }>
                    <BubbleMap pointDecorator={ pointDecorator } data={ data } click={ this.onMarkerClick.bind(this) } />
                    <ScaleControl style={ { backgroundColor: 'rgba(0, 0, 0, 0)', left: '12px', bottom: '6px', } } />
                    <ZoomControl className="zoom-control" position={ 'topLeft' } />

                    <PlaybackControl className="playback-control" playback={ playbackMode } onPlaybackChange={ onChangePlaybackMode.bind(this) } />

                    { feature &&
                      <Popup anchor="bottom" offset={ [0, -10] } coordinates={ feature.geometry.coordinates }>
                          <button className="mapboxgl-popup-close-button" type="button" onClick={ this.onPopupChange.bind(this, false) }>
                              ×
                          </button>
                          <div>
                              <span> {this.renderProperties(feature)} </span>
                          </div>
                      </Popup> }
                </ReactMapboxGl>
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
