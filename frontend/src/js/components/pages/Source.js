// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import MapContainer from '../containers/MapContainer';
import CriteriaPanel from '../CriteriaPanel';
import SourceOverview from './SourceOverview';
import ChartsContainer  from '../ChartsContainer';

import { generatePointDecorator } from '../../common/utilities';

import { loadChartData, changeCriteria } from '../../actions';

import '../../../css/source.css';

class Source extends Component {
    props: {
        loadChartData: any => void,
        changeCriteria: any => void,
        visibleFeatures: any,
        chartData: any,
        match: any,
        style: any,
    };

    state: {};

    getSourceId() {
        const { match } = this.props;
        const { sourceId } = match.params;

        return sourceId;
    }

    onShowChart(chart) {
        this.props.loadChartData(chart);
    }

    onCriteriaChanged(newCriteria) {
        this.props.changeCriteria(newCriteria);
    }

    render() {
        const { visibleFeatures, chartData } = this.props;

        const sourceId = this.getSourceId();

        const sourceData = visibleFeatures.sources[sourceId];
        if (!sourceData || !sourceData.summary || !sourceData.source) {
            return <div></div>;
        }

        const newSources = {};
        newSources[sourceId] = sourceData;

        console.log(sourceData);

        const narrowed = {
            geojson: { type: '', features: [] },
            sources: newSources,
            focus: {
                center: sourceData.source.centroid,
                feature: null,
                time: 0,
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

function showWhenReady(WrappedComponent, isReady) {
    return class extends React.Component {
        render() {
            if (!isReady(this.props)) {
                return <div>Loading</div>;
            }

            return <WrappedComponent {...this.props} />
        }
    }
}

export default connect(mapStateToProps, {
    loadChartData,
    changeCriteria
})(showWhenReady(Source, props => {
    const { match, visibleFeatures } = props;
    const { sourceId } = match.params;
    const data = visibleFeatures.sources[sourceId];
    return data && data.source && data.summary;
}));
