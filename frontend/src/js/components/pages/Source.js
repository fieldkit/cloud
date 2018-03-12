// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import MapContainer from '../containers/MapContainer';
import SimpleChartContainer from '../containers/ChartContainer';
import CriteriaPanel from '../CriteriaPanel';

import { generatePointDecorator } from '../../common/utilities';

import { loadSourceCharts, loadChartData, changeCriteria } from '../../actions';

import '../../../css/source.css';

type Props = {};

class SourceOverview extends Component {
    props: {
        style: React.CSSProperties
    }

    onShowChart(key) {
        const { data, onShowChart } = this.props;
        const { source } = data;

        onShowChart({
            id: ["chart", source.id, key].join('-'),
            sourceId: source.id,
            keys: [key]
        });
    }

    renderReadings(data) {
        const { summary } = data;

        return (
            <table style={{ padding: '5px', width: '100%' }} className="feature-data">
                <thead>
                    <tr>
                        <th>Reading</th>
                    </tr>
                </thead>
                <tbody>
                { summary.readings.map(reading => (
                    <tr key={reading.name}>
                        <td style={{ width: '50%' }}> <div onClick={() => this.onShowChart(reading.name)} style={{ cursor: 'pointer' }}>{reading.name}</div> </td>
                    </tr>
                )) }
                </tbody>
            </table>
        )
    }

    render() {
        const { style, data } = this.props;

        return (
            <div style={{ ...style }} className="feature-panel">
                {this.renderReadings(data)}
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

    getSourceId() {
        const { match } = this.props;
        const { sourceId } = match.params;

        return sourceId;
    }

    componentWillMount() {
        this.props.loadSourceCharts(this.getSourceId());
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
    loadSourceCharts,
    loadChartData,
    changeCriteria
})(showWhenReady(Source, props => {
    const { match, visibleFeatures } = props;
    const { sourceId } = match.params;
    const data = visibleFeatures.sources[sourceId];
    return data && data.source && data.summary;
}));
