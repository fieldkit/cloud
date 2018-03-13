// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux'

import type { ActiveExpedition, GeoJSON, GeoJSONFeature } from '../../types';

import MapContainer from '../containers/MapContainer';

import { notifyOfUserMapActivity, changePlaybackMode, focusFeature, focusSource } from '../../actions';
import { generatePointDecorator } from '../../common/utilities';

import '../../../css/map.css';

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
    notifyOfUserMapActivity: () => mixed,
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

    // Not sure how I feel about this. I'm open to suggestions.
    componentWillMount() {
        // $FlowFixMe
        document.body.style.overflow = "hidden";
    }

    componentWillUnmount() {
        // $FlowFixMe
        document.body.style.overflow = null;
    }

    render() {
        const { visibleFeatures, playbackMode, notifyOfUserMapActivity, focusFeature, focusSource, changePlaybackMode } = this.props;
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
                        notifyOfUserMapActivity={ notifyOfUserMapActivity }
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
    focusSource
})(Map);
