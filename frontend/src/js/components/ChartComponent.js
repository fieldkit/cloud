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
        const { loadChartData, chart } = this.props;

        loadChartData(chart);
    }

    componentWillReceiveProps(nextProps) {
        const { chart: oldChart, loadChartData } = this.props;
        const { chart: newChart } = nextProps;

        if (newChart != oldChart) {
            loadChartData(newChart);
        }
    }

    render() {
        const { chart, data } = this.props;

        return <SimpleChartContainer chart={chart} data={data} />
    }
}

const mapStateToProps = state => ({
    data: state.chartData
});

export default connect(mapStateToProps, {
    loadChartData
})(ChartComponent);
