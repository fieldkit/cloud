// @flow weak

import React, { Component } from 'react';
import D3Scatterplot from '../visualizations/D3Scatterplot';
import D3TimeSeries from '../visualizations/D3TimeSeries';
import type { Margin } from '../../types/D3Types';

type Props = {
    type: string;
    margin: Margin;
    data: {}[],
};

function chartStringToClass(type, props) {
    switch (type) {
    case 'scatterplot':
        return (new D3Scatterplot(props));
    case 'time-series':
        return (new D3TimeSeries(props));
    default:
        return null;
    }
}

export default class ChartContainer extends Component {
    props: Props;
    node: HTMLElement;
    chart: D3Scatterplot;

    componentDidMount() {
        const { type, data, margin } = this.props;
        const { node } = this;
        this.chart = chartStringToClass(type, {
            node,
            margin,
            data
        });
    }

    componentWillReceiveProps(nextProps: Props) {
        const { chart } = this;
        const { data } = nextProps;
        chart.updateData(data);
    }

    render() {
        return <div className="chart" ref={ node => (this.node = node) } />;
    }
}
