// @flow weak

import _ from 'lodash';
import moment from 'moment';

import React, { Component } from 'react';

import { Loading } from '../Loading';

import { TimeSeries, TimeRange } from "pondjs";
import { Charts, ChartContainer, ChartRow, YAxis, LineChart, styler, Resizable } from "react-timeseries-charts";

type Props = {
    chart: {},
    geojson: {},
};

export default class SimpleChartContainer extends Component {
    props: Props;

    handleChartResize(width) {
        this.setState({ width });
    }

    render() {
        const { chart, geojson } = this.props;

        if (geojson.loading) {
            return (<Loading />);
        }

        const samples = _(geojson.geo).map(f => {
            return [
                f.properties.timestamp,
                Number(f.properties.data[chart.key])
            ];
        }).value();

        const data = {
            name: 'Data',
            columns: [ "time", "value" ],
            points: samples
        };
        const timeseries = new TimeSeries(data);
        const lineStyle = styler([
            { key: "value", color: "#9467bd", width: 1 }
        ]);
        /*
        const begin = moment();
        const end =  moment().subtract(3, 'days');
        const tr = new TimeRange(begin, end);
        */
        const tr = timeseries.timerange();
        return (
            <div style={{ padding: '10px' }}>
                <Resizable>
                    <ChartContainer timeRange={tr} showGrid={false} onChartResize={width => this.handleChartResize(width)} showGrid={true}>
                        <ChartRow height="180">
                            <YAxis id="y" label={chart.key} min={timeseries.min()} max={timeseries.max()} width="100" type="linear" format=",.2f" labelOffset={-40} />
                            <Charts>
                                <LineChart axis="y" series={timeseries} style={lineStyle} />
                            </Charts>
                        </ChartRow>
                    </ChartContainer>
                </Resizable>
            </div>
        );
    }
}
