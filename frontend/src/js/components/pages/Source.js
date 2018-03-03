// @flow weak

import React, { Component } from 'react';
import PropTypes from "prop-types";
import { connect } from 'react-redux';

import MapContainer from '../containers/MapContainer';
import FeaturePanel from '../containers/FeaturePanel';
import ChartComponent from '../ChartComponent';

import { generatePointDecorator } from '../../common/utilities';

import '../../../css/source.css';

type Props = {};

class SourceOverview extends Component {
    render() {
        const { data, onShowChart } = this.props;
        const { source, summary, lastFeature } = data;

        return (
            <div className="source-overview">
                <FeaturePanel feature={lastFeature} onShowChart={ onShowChart } />
            </div>
        );
    }
};

class Source extends Component {
    props: Props;
    state: {};

    constructor(props) {
        super(props)

        this.state = {
            chart: null
        };
    }

    // Not sure how I feel about this. I'm open to suggestions.
    componentWillMount() {
        document.body.style.overflow = "hidden";
    }

    componentWillUnmount() {
        document.body.style.overflow = null;
    }

    onShowChart(chart) {
        this.setState({
            chart
        })
    }

    render() {
        const { match, visibleFeatures } = this.props;
        const { params } = match;
        const { sourceId } = params;
        const { chart } = this.state;

        const sourceData = visibleFeatures.sources[sourceId];
        if (!sourceData || !sourceData.summary || !sourceData.source || !sourceData.lastFeature) {
            return <div></div>;
        }

        const newSources = {};
        newSources[sourceId] = sourceData;
        const narrowed = {
            geojson: { features: [] },
            sources: newSources,
            focus: {
                center: sourceData.source.centroid,
            }
        };

        const pointDecorator = generatePointDecorator('constant', 'constant');
        return (
            <div className="source page">
                <div className="header">
                    <div className="project-name">FieldKit Project / Device: {sourceData.source.name}</div>
                </div>
                <div>
                    <section className="wrapper dashboard">
                        <div className="row top">
                            <div className="column-6 map">
                                <MapContainer style={{ height: "100%" }} containerStyle={{ width: "100%", height: "100%" }}
                                    pointDecorator={ pointDecorator } visibleFeatures={ narrowed } controls={false}
                                    playbackMode={ () => false } focusFeature={ () => false }
                                    focusSource={ () => false } notifyOfUserMapActivity={ () => false }
                                    onChangePlaybackMode={ () => false } />

                                <div className="chart">
                                    { chart && <ChartComponent chart={chart} /> }
                                </div>
                            </div>
                            <div className="column-6 source">
                                <SourceOverview data={sourceData} onShowChart={ this.onShowChart.bind(this) } />
                            </div>
                        </div>
                    </section>
                </div>
            </div>
        );
    }
}

const mapStateToProps = state => ({
    visibleFeatures : state.visibleFeatures
});

export default connect(mapStateToProps, {

})(Source);
