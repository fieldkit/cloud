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

import { RadialMenu } from './RadialMenu';

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
        feature: any,
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

    onUserActivity(how) {
        const { onUserActivity } = this.props;

        onUserActivity(how);
    }

    onClick(target, ev) {
        const features = target.queryRenderedFeatures(ev.point);
        const coordinates = ev.lngLat;

        let options = this.clusterOptions(features);
        if (features.length > 1) {
            options = this.selectFeatureOptions(features);
        }

        this.setState({
            menu: {
                mouse: ev.point,
                coordinates: coordinates,
                features: features,
                options: options,
            }
        });
    }

    onMenuClose() {
        const { menu } = this.state;


        if (menu.nextOptions) {
            this.setState({
                menu: {
                    mouse: menu.mouse,
                    coordinates: menu.coordinates,
                    features: menu.features,
                    options: menu.nextOptions,
                }
            });
        }
        else {
            this.setState({
                menu: null
            });
        }
    }

    renderRadialMenu() {
        const { menu } = this.state;

        if (!menu || menu.features.length === 0) {
            return <div></div>;
        }

        const { options } = menu;

        return <RadialMenu key={options.key} delay={80} duration={400} strokeWidth={2} innerRadius={20} outerRadius={120}
                stroke={"rgba(255, 255, 255, 1)"} fill={"rgba(0, 0, 0, 0.8)"}
                buttons={options.buttons} buttonFunctions={options.functions}
                position={menu.mouse} onClosed={ this.onMenuClose.bind(this) } />;
    }

    selectFeatureOptions(features) {
        const options = _(features).map(f => f.properties).map(f => {
            return {
                title: f.name,
                function: () => {
                    const { menu } = this.state;
                    const nextMenu = { ...menu, ...{ features: [f], nextOptions: this.clusterOptions([f]) } };
                    this.setState({
                        menu: nextMenu
                    });
                }
            };
        });

        return {
            key: 1,
            buttons: options.map(o => o.title).value(),
            functions: options.map(o => o.function).value()
        };
    }

    clusterOptions(features) {
        return {
            key: 0,
            buttons: [
                "Data",
                "Graph",
                "Summary"
            ],
            functions: [
                () => console.log("Data"),
                () => console.log("Graph"),
                () => console.log("Summary"),
            ]
        };
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
