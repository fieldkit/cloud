// @flow weak

import _ from 'lodash';

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactMapboxGl, { ScaleControl, ZoomControl } from 'react-mapbox-gl';

import { MAPBOX_ACCESS_TOKEN, MAPBOX_STYLE } from '../../secrets';

import BubbleMap from '../visualizations/BubbleMap';
import ClusterMap from '../visualizations/ClusterMap';
import PlaybackControl from '../PlaybackControl';

import FeaturePanel from './FeaturePanel';
import NotificationsPanel from './NotificationsPanel';
import ChartComponent from '../ChartComponent';
import FiltersPanel from './FiltersPanel';

import { RadialMenu } from '../RadialMenu';

import { FieldKitLogo } from '../icons/Icons';

import type { Coordinates, Bounds, GeoJSONFeature, GeoJSON } from '../../types/MapTypes';
import type { Focus } from '../../types';

type Props = {
    playbackMode: mixed,
    focusSource: () => mixed,
    onChangePlaybackMode: () => mixed,
    notifyOfUserMapActivity: () => mixed,
    pointDecorator: PointDecorator,
    visibleFeatures: {
        geojson: GeoJSON,
        focus: Focus,
    }
};

const panelContainerStyle = {
    backgroundColor: '#f9f9f9',
    color: "#000",
    position: 'absolute',
    zIndex: 10,
    display: 'flex',
    flexDirection: 'column',
    boxShadow: '0px 1px 4px rgba(0, 0, 0, .3)',
};

const panelHeaderStyle = {
    padding: '5px',
    backgroundColor: 'rgb(196, 196, 196)',
    borderBottom: '1px solid rgb(128, 128, 128)',
    fontWeight: 'bold',
};

const panelBodyStyle = {
};

export class MapBottom extends Component {
    positionStyle() {
        const { sidePanelVisible } = this.props;
        if (!sidePanelVisible) {
            return { top: 'auto', right: 0, bottom: 0, left: 0, height: '250px', width: '100%' };
        }
        return { top: 'auto', right: 300, bottom: 0, left: 0, height: '250px' };
    }

    render() {
        const { onHide, children } = this.props;
        const position = this.positionStyle();

        return (
            <div style={{...panelContainerStyle, ...position}}>
                <div style={{ ...panelHeaderStyle }}><span style={{ cursor: 'pointer' }} onClick={() => onHide()}>Close</span></div>
                <div style={{ ...panelBodyStyle }}>
                    {children}
                </div>
            </div>
        );
    }
}

export class MapRight extends Component {
    render() {
        const { onHide, children } = this.props;
        const position = { top: 0, right: 0, bottom: 'auto', left: 'auto', width: '300px', height: '100%' };

        return (
            <div style={{...panelContainerStyle, ...position}}>
                <div style={{ ...panelHeaderStyle }}><span style={{ cursor: 'pointer' }} onClick={() => onHide()}>Close</span></div>
                <div style={{ ...panelBodyStyle }}>
                    {children}
                </div>
            </div>
        );
    }
}

const Map = ReactMapboxGl({
    accessToken: MAPBOX_ACCESS_TOKEN
});

export default class MapContainer extends Component {
    props: Props
    state: {
        fitBounds: Bounds,
        center: Coordinates,
        zoom: [number],
        feature: ?GeoJSONFeature,
    };

    static contextTypes = {
        router: PropTypes.shape({
            history: PropTypes.shape({
                push: PropTypes.func.isRequired,
                replace: PropTypes.func.isRequired,
            }).isRequired
        }).isRequired
    };

    constructor(props: Props) {
        super(props);
        this.state = {
            fitBounds: null,
            center: null,
            zoom: [14],
            feature: null,
            panels: {
                sidePanelVisible: false,
                bottomPanelVisible: false,
            }
        };
        this.onUserActivityThrottled = _.throttle(this.onUserActivity.bind(this), 500, { leading: true });
    }

    _createState(props) {
        const { visibleFeatures } = props;
        const { focus } = visibleFeatures;

        if (focus.feature) {
            return {
                center: focus.feature.geometry.coordinates
            };
        }
        else if (focus.center) {
            return {
                center: focus.center
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

    onMarkerClick(feature) {
        this.setState({
            feature,
            panels: {
                sidePanelVisible: true,
            },
        });
    }

    onFocusFeature(feature) {
        const { focusFeature } = this.props;

        focusFeature(feature);

        this.setState({
            feature,
            panels: {
                sidePanelVisible: true,
            },
        });
    }

    onShowChart(chart) {
        this.setState({
            chart,
            panels: {
                sidePanelVisible: true,
                bottomPanelVisible: true,
            },
        });
    }

    onFocusSource(source) {
        const { focusSource } = this.props;

        focusSource(source);
    }

    onUserActivity(how) {
        const { notifyOfUserMapActivity } = this.props;

        notifyOfUserMapActivity();
    }

    renderPanels() {
        const { visibleFeatures, playbackMode } = this.props;
        const { panels, chart, feature } = this.state;

        const { mode } = playbackMode;
        const paused = mode === "Pause";

        return (
            <div>
                { !paused && <NotificationsPanel features={visibleFeatures.focus.features} sidePanelVisible={panels.sidePanelVisible} />}
                { panels.sidePanelVisible && <MapRight onHide={() => this.setState({ panels: { ...panels, sidePanelVisible: false, bottomPanelVisible: false }})}>
                  <FeaturePanel feature={feature} onShowChart={chart => this.onShowChart(chart)}/>
                </MapRight> }
                { panels.bottomPanelVisible && <MapBottom onHide={() => this.setState({ panels: { ...panels, bottomPanelVisible: false }})} sidePanelVisible={panels.sidePanelVisible}>
                  <ChartComponent chart={chart} />
                </MapBottom> }
            </div>
        );
    }

    onMouseMove(target, ev) {
    }

    onMouseOut(target, ev) {
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
        const { fitBounds, center, zoom } = this.state;

        if (!center) {
            return <div></div>;
        }

        return (
            <div style={style}>
                <Map style={MAPBOX_STYLE} containerStyle={containerStyle}
                    movingMethod="easeTo" center={ center } zoom={ zoom } fitBounds={ fitBounds }
                    onMouseMove={ this.onMouseMove.bind(this) }
                    onMouseOut={ this.onMouseOut.bind(this) }
                    onClick={ this.onClick.bind(this) }
                    onDrag={ this.onUserActivityThrottled.bind(this) }>

                    <ClusterMap data={ visibleFeatures.sources } onClick={ this.onMarkerClick.bind(this) } />

                    <BubbleMap pointDecorator={ pointDecorator } data={ visibleFeatures.geojson } onClick={ this.onMarkerClick.bind(this) } />

                    { controls && <ScaleControl style={ { backgroundColor: 'rgba(0, 0, 0, 0)', left: '12px', bottom: '6px', } } />}

                    { controls && false && <ZoomControl className="zoom-control" position={ 'topLeft' } />}

                    { controls && <PlaybackControl className="playback-control" playback={ playbackMode } onPlaybackChange={ onChangePlaybackMode.bind(this) } />}

                    { controls && <FiltersPanel visibleFeatures={ visibleFeatures } onShowSource={ this.onFocusSource.bind(this) } onShowFeature={ this.onFocusFeature.bind(this) } />}

                    {this.renderRadialMenu()}
                </Map>
                { controls && this.renderPanels()}
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
