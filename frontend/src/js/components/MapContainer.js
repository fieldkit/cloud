// @flow weak

import _ from 'lodash';

import React, { Component } from 'react';
import ReactMapboxGl, { ScaleControl, ZoomControl } from 'react-mapbox-gl';

import type { Focus, StyleSheet, Coordinates, GeoJSON } from '../types';

import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../secrets';

import BubbleMap from './visualizations/BubbleMap';
import ClusterMap from './visualizations/ClusterMap';
import PlaybackControl from './PlaybackControl';
import FiltersPanel from './FiltersPanel';

import { FancyMenu } from './FancyMenu';

import { FieldKitLogo } from './icons/Icons';

const Map = ReactMapboxGl({
    accessToken: MAPBOX_ACCESS_TOKEN
});

type Props = {
    containerStyle?: StyleSheet,
    style?: StyleSheet,
    playbackMode: mixed,
    focusSource: () => mixed,
    onChangePlaybackMode: () => mixed,
    onUserActivity: () => mixed,
    loadMapFeatures: () => mixed,
    pointDecorator: PointDecorator,
    controls: ?bool,
    visibleFeatures: {
        geojson: GeoJSON,
        focus: Focus,
        sources: any,
    }
};

export default class MapContainer extends Component {
    props: Props
    state: {
        center: ?Coordinates,
        zoom: [number],
        menu: any,
    }

    onUserActivityThrottled: () => mixed

    constructor(props: Props) {
        super(props);
        this.state = {
            center: null,
            zoom: [14],
            menu: null,
        };
        this.onUserActivityThrottled = _.throttle(this.onUserActivity.bind(this), 500, { leading: true });
    }

    _createState(props) {
        const { visibleFeatures } = props;
        const { focus } = visibleFeatures;

        if (focus.center) {
            return {
                center: focus.center,
                zoom: [focus.center.length === 3 ? focus.center[2] : 14],
            };
        }
        return {};
    }

    componentWillMount() {
        this.setState(this._createState(this.props));
    }

    componentWillReceiveProps(nextProps) {
        const { visibleFeatures: oldData } = this.props;
        const { visibleFeatures: newData } = nextProps;

        if (oldData !== newData) {
            const { focus: oldFocus } = oldData;
            const { focus: newFocus } = newData;

            if (oldFocus !== newFocus) {
                this.setState(this._createState(nextProps));
            }
        }
    }

    onMarkerClick(thing) {
        console.log(thing);
    }

    onFocusSource(source) {
        const { focusSource } = this.props;

        focusSource(source);
    }

    onUserActivity(map) {
        const { onUserActivity, loadMapFeatures } = this.props;

        const bounds = map.getBounds();
        loadMapFeatures({
            ne: bounds.getNorthEast().toArray(),
            sw: bounds.getSouthWest().toArray(),
        });
        onUserActivity(map);
    }

    onClick(target, ev) {
        const selected = _.uniqBy(target.queryRenderedFeatures(ev.point), f => f.properties.name);
        const coordinates = ev.lngLat;

        if (selected.length === 0) {
            this.setState({
                menu: null
            });
            return;
        }

        this.setState({
            menu: {
                ...{
                    mouse: ev.point,
                    coordinates: coordinates,
                    selected: selected,
                },
                ...this.getMenu(selected)
            }
        });
    }

    getMenu(selected) {
        return {
            options: this.getOptions(selected),
            center: this.getCenter(selected),
        };
    }

    getCenter(selected) {
        if (selected.length === 1) {
            return {
                text: selected[0].properties.name,
            }
        }
        return null;
    }

    getOptions(selected) {
        if (selected.length === 1) {
            return [
                {
                    title: "Graph",
                    onClick: () => {
                    }
                },
                {
                    title: "Playback",
                    onClick: () => {
                    }
                },
                {
                    title: "Hide",
                    onClick: () => {
                    }
                }
            ];
        }

        return selected.map(feature => {
            return {
                title: feature.properties.name,
                onClick: () => {
                    const { menu } = this.state;
                    this.setState({
                        menu: { ...menu, ...this.getMenu([feature]) },
                    });
                    return true;
                }
            };
        });
    }

    onMenuClose() {
        this.setState({
            menu: null
        });
    }

    renderRadialMenu() {
        const { menu } = this.state;

        const visible = menu != null && menu.mouse != null;
        const position = visible ? menu.mouse : { x: 0, y: 0};

        return <FancyMenu visible={visible} options={menu ? menu.options : []} center={menu ? menu.center : null} position={position} onClosed={ this.onMenuClose.bind(this) } />
    }

    render() {
        const { pointDecorator, visibleFeatures, onChangePlaybackMode, playbackMode, containerStyle, style, controls } = this.props;
        const { center, zoom } = this.state;

        if (!center) {
            return <div></div>;
        }

        return (
            <div style={style}>
                <Map style={MAPBOX_STYLE} containerStyle={containerStyle}
                    movingMethod="easeTo" center={ center } zoom={ zoom }
                    onClick={ this.onClick.bind(this) }
                    onZoomEnd={ this.onUserActivityThrottled.bind(this) }
                    onDrag={ this.onUserActivityThrottled.bind(this) }
                    onDragEnd={ this.onUserActivityThrottled.bind(this) }>

                    <ClusterMap data={ visibleFeatures.sources } onClick={ this.onMarkerClick.bind(this) } />

                    <BubbleMap pointDecorator={ pointDecorator } data={ visibleFeatures.geojson } onClick={ this.onMarkerClick.bind(this) } />

                    { controls && <ScaleControl style={ { backgroundColor: 'rgba(0, 0, 0, 0)', left: '12px', bottom: '6px', } } />}

                    { controls && false && <ZoomControl className="zoom-control" position={ 'topLeft' } />}

                    { controls && <PlaybackControl className="playback-control" playback={ playbackMode } onPlaybackChange={ onChangePlaybackMode.bind(this) } />}

                    { controls && <FiltersPanel visibleFeatures={ visibleFeatures } onShowSource={ this.onFocusSource.bind(this) } onShowFeature={ () => console.log(arguments) } />}

                    {this.renderRadialMenu()}
                </Map>
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
