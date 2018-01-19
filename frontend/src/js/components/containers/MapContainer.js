// @flow weak

import React, { Component } from 'react';
import ReactMapboxGl, { ScaleControl, ZoomControl, Popup } from 'react-mapbox-gl';
import _ from 'lodash';

import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../../secrets'

import BubbleMap from '../visualizations/BubbleMap';
import PlaybackControl from '../PlaybackControl';

import type { Coordinates, Bounds, GeoJSONFeature, GeoJSON } from '../../types/MapTypes';
import type { Focus } from '../../types';

import { FieldKitLogo } from '../icons/Icons';

import FeaturePanel from './FeaturePanel';
import NotificationsPanel from './NotificationsPanel';
import FiltersPanel from './FiltersPanel';

function getFitBounds(geojson: GeoJSON) {
    const lon = geojson.features.map(f => f.geometry.coordinates[0]);
    const lat = geojson.features.map(f => f.geometry.coordinates[1]);
    return [[Math.min(...lon), Math.min(...lat)], [Math.max(...lon), Math.max(...lat)]];
}

function getCenter(fitBounds: Bounds) {
    return [(fitBounds[0][0] + fitBounds[1][0]) / 2, (fitBounds[0][1] + fitBounds[1][1]) / 2];
}

type Props = {
    playbackMode: mixed,
    onChangePlaybackMode: () => mixed,
    notifyOfUserMapActivity: () => mixed,
    pointDecorator: PointDecorator,
    visibleFeatures: {
        geojson: GeoJSON,
        focus: Focus,
    }
};

export default class MapContainer extends Component {
    props: Props
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
        this.onUserActivityThrottled = _.throttle(this.onUserActivity.bind(this), 500, { leading: true });
    }

    componentWillMount() {
        const { visibleFeatures } = this.props;

        const fitBounds = getFitBounds(visibleFeatures.geojson);
        const center = getCenter(fitBounds);
        this.setState({
            fitBounds,
            center
        });
    }

    componentWillReceiveProps(nextProps) {
        const { visibleFeatures: oldData } = this.props;
        const { visibleFeatures: newData } = nextProps;

        if (oldData !== newData) {
            const { focus: oldFocus } = oldData;
            const { focus: newFocus } = newData;

            if (oldFocus !== newFocus) {
                if (newFocus.feature) {
                    this.setState({
                        center: newFocus.feature.geometry.coordinates
                    })
                }
                else if (newFocus.center) {
                    this.setState({
                        center: newFocus.center
                    })
                }
            }
        }
    }

    onMarkerClick(feature) {
        this.setState({
            feature
        });
    }

    onFocusFeature(feature) {
        const { focusFeature } = this.props;
        focusFeature(feature);
        this.setState({
            feature
        });
    }

    onUserActivity(how) {
        const { notifyOfUserMapActivity } = this.props;

        notifyOfUserMapActivity();
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

    /*
    tick(ms, firstFrame) {
        requestAnimationFrame(ms => this.tick(ms, false));
    }

    componentDidMount() {
        requestAnimationFrame(ms => this.tick(ms, true));
    }
    */

    render() {
        const { pointDecorator, visibleFeatures, onChangePlaybackMode, playbackMode } = this.props;
        const { fitBounds, center, zoom, feature } = this.state;

        if (!visibleFeatures.geojson) {
            return (<div><h1>Loading</h1></div>);
        }
        else if (visibleFeatures.geojson.features.length === 0) {
            return (<div><h1>Empty</h1></div>);
        }

        return (
            <div>
                <ReactMapboxGl accessToken={MAPBOX_ACCESS_TOKEN}
                    style={MAPBOX_STYLE}
                    movingMethod="easeTo" fitBounds={ fitBounds } center={ center }
                    onDrag={ this.onUserActivityThrottled } 
                    zoom={ zoom } onClick={ this.onPopupChange.bind(this) } containerStyle={ { height: '100vh', width: '100vw', } }>
                    <BubbleMap pointDecorator={ pointDecorator } data={ visibleFeatures.geojson } click={ this.onMarkerClick.bind(this) } />
                    <ScaleControl style={ { backgroundColor: 'rgba(0, 0, 0, 0)', left: '12px', bottom: '6px', } } />
                    <ZoomControl className="zoom-control" position={ 'topLeft' } />

                    <PlaybackControl className="playback-control" playback={ playbackMode } onPlaybackChange={ onChangePlaybackMode.bind(this) } />

                    <FiltersPanel features visibleFeatures={visibleFeatures} onShowFeature={ this.onFocusFeature.bind(this) } />

                    <NotificationsPanel features={visibleFeatures.focus.features} />

                    { feature &&
                      <Popup anchor="bottom" offset={ [0, -10] } coordinates={ feature.geometry.coordinates } style={{ padding: '0px' }}>
                        <FeaturePanel feature={feature} />
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
