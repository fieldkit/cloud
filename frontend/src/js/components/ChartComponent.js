// @flow weak

import React, { Component } from 'react';
import { connect } from 'react-redux'

import { loadChartData } from '../actions';

import SimpleChartContainer from './containers/ChartContainer';

type Props = {
    data: {}[],
};

class ChartComponent extends Component {
    props: Props;

    componentWillMount() {
        const { loadChartData , chart } = this.props;
        loadChartData(chart);
    }

    render() {
        const { chart, geojson } = this.props;

        return (
            <SimpleChartContainer chart={chart} geojson={geojson} />
        );
    }
}

const mapStateToProps = state => ({
    geojson: state.chartData
});

export default connect(mapStateToProps, {
    loadChartData
})(ChartComponent);
