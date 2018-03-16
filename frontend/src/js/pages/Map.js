// @flow weak

import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux'

import type { ActiveExpedition, GeoJSON, GeoJSONFeature } from '../types';

import { notifyOfUserMapActivity, changePlaybackMode, focusFeature, focusSource, loadMapFeatures } from '../actions';
import { generatePointDecorator } from '../common/utilities';

import MapContainer from '../components/MapContainer';

import '../../css/map.css';

type Props = {
    activeExpedition: ActiveExpedition,
    visibleFeatures: {
        geojson: GeoJSON,
        focus: {
            feature: ?GeoJSONFeature,
            time: number,
        },
        sources: any,
    },
    playbackMode: {},
    focusFeature: () => mixed,
    focusSource: () => mixed,
    changePlaybackMode: () => mixed,
    loadMapFeatures: () => mixed,
    notifyOfUserMapActivity: () => mixed,
};

class Map extends Component {
    props: Props
    state: {
        pointDecorator: PointDecorator
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
            pointDecorator: generatePointDecorator('constant', 'constant')
        };
    }

    // Not sure how I feel about this. I'm open to suggestions.
    componentWillMount() {
        // $FlowFixMe
        document.body.style.overflow = "hidden";
    }

    componentWillUnmount() {
        // $FlowFixMe
        document.body.style.overflow = null;
    }

    onUserActivity(map) {
        this.props.notifyOfUserMapActivity(map);
        this.updateCoordinates(map);
    }

    updateCoordinates(map) {
        const newZoom = map.getZoom();
        const mapCenter = map.getCenter();
        const newCenter = [ mapCenter.lng, mapCenter.lat, newZoom ];
        this.context.router.history.push({
            pathname: '/',
            search: '?center=' + newCenter.join(","),
        })
    }

    render() {
        const { visibleFeatures, playbackMode, focusFeature, focusSource, changePlaybackMode, loadMapFeatures } = this.props;
        const { pointDecorator } = this.state;

        return (
            <div className="map page">
                <div className="header">
                    <div className="project-name">FieldKit Project</div>
                </div>
                <div>
                    <MapContainer style={{ }} containerStyle={{ width: "100vw", height: "100vh" }} controls={true}
                        pointDecorator={ pointDecorator }
                        visibleFeatures={ visibleFeatures }
                        playbackMode={ playbackMode }
                        focusFeature={ focusFeature }
                        focusSource={ focusSource }
                        onUserActivity={ this.onUserActivity.bind(this) }
                        loadMapFeatures={ loadMapFeatures }
                        onChangePlaybackMode={ changePlaybackMode } />
                </div>
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
    notifyOfUserMapActivity,
    changePlaybackMode,
    focusFeature,
    focusSource,
    loadMapFeatures
})(Map);
