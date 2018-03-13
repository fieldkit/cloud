// @flow weak

import React, { Component } from 'react';

import SimpleChart from './ChartsContainer';

type Props = {
    chartData: {
        charts: Array<any>
    },
}

export default class ChartsContainer extends Component {
    props: Props

    render() {
        const { chartData } = this.props;

        return <div className="charts">
            {chartData.charts.map(chart => <div key={chart.id} className="chart"><SimpleChart chartDef={chart} /></div>)}
        </div>;
    }
};
