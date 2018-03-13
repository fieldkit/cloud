// @flow weak

import React, { Component } from 'react';

import SimpleChartContainer from '../components/ChartsContainer.js';

export default class ChartsContainer extends Component {
    props: {
        chartData: any,
    }

    render() {
        const { chartData } = this.props;

        return <div className="charts">
            {chartData.charts.map(chart => <div key={chart.id} className="chart"><SimpleChartContainer chart={chart} /></div>)}
        </div>;
    }
};
