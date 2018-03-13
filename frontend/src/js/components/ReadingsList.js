import React, { Component } from 'react';

import LazyLoad from 'react-lazyload';

import SimpleChart from './SimpleChart';

class ReadingOverview extends Component {
    props: {
        reading: {
            name: string
        },
        sourceData: any,
        chartData: {
        },
        loadChartData: () => void,
    }

    componentWillMount() {
        const { reading, loadChartData } = this.props;

        loadChartData(reading.chartDef);
    }

    render() {
        const { reading, chartData } = this.props;
        const data = chartData.queries[reading.chartDef.id];

        return <div className="reading">
            <div className="title">{reading.name}</div>
            <div className="chart">
                <SimpleChart chartDef={reading.chartDef} data={data} />
            </div>
        </div>;
    }
};

type Props = {
    loadChartData: () => void,
    sourceData: any,
    chartData: {
        charts: Array<any>
    },
}

export default class ReadingsList extends Component {
    props: Props

    render() {
        const { sourceData, chartData, loadChartData } = this.props;

        return <div className="readings">
            {sourceData.summary.readings.map((reading, i) => <div key={i} className="reading-container">
                <LazyLoad height={320} overflow={true} offset={100} once>
                    <ReadingOverview reading={reading} loadChartData={loadChartData} chartData={chartData} />
                </LazyLoad>
            </div>)}
        </div>;
    }
};
