// @flow weak

import React, { Component } from 'react';
import PropTypes from "prop-types";
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import MapContainer from '../containers/MapContainer';
import FeaturePanel from '../containers/FeaturePanel';
import ChartComponent from '../ChartComponent';
import SimpleChartContainer from '../containers/ChartContainer';
import CriteriaPanel from '../CriteriaPanel';

import { generatePointDecorator } from '../../common/utilities';

import { loadChartData, changeCriteria } from '../../actions';

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

class ChartsContainer extends Component {
    render() {
        const { chartData } = this.props;

        return <div className="charts">{chartData.charts.map(chart => <div key={chart.id} className="chart"><SimpleChartContainer key={chart.id} chart={chart} /></div>)}</div>;
    }
};

class Source extends Component {
    props: Props;
    state: {};

    constructor(props) {
        super(props)
    }

    onShowChart(chart) {
        const { loadChartData } = this.props;

        loadChartData(chart);

    }

    onCriteriaChanged(newCriteria) {
        const { changeCriteria } = this.props;

        changeCriteria(newCriteria);
    }

    render() {
        const { match, visibleFeatures, chartData } = this.props;
        const { params } = match;
        const { sourceId } = params;

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
                    <div className="project-name"><Link to='/'>FieldKit Project</Link> / Device: {sourceData.source.name}</div>
                    <CriteriaPanel onChangeTimeCiteria={ r => this.onCriteriaChanged(r) }/>
                </div>
                <div className="main-container">
                    <div className="middle-container">
                        <ChartsContainer chartData={chartData} />
                    </div>
                    <div className="side-container">
                        <div className="map">
                            <MapContainer style={{ height: "100%" }} containerStyle={{ width: "100%", height: "100%" }}
                                pointDecorator={ pointDecorator } visibleFeatures={ narrowed } controls={false}
                                playbackMode={ () => false } focusFeature={ () => false }
                                focusSource={ () => false } notifyOfUserMapActivity={ () => false }
                                onChangePlaybackMode={ () => false } />
                        </div>
                        <div className="">
                            <SourceOverview data={sourceData} onShowChart={ this.onShowChart.bind(this) } />
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

const mapStateToProps = state => ({
    visibleFeatures: state.visibleFeatures,
    chartData: state.chartData,
});

export default connect(mapStateToProps, {
    loadChartData,
    changeCriteria
})(Source);
