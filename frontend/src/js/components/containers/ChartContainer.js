// @flow weak

import _ from 'lodash';

import React, { Component } from 'react';

import { Loading } from '../Loading';

import { TimeSeries } from "pondjs";
import { Charts, ChartContainer, ChartRow, YAxis, LineChart, styler } from "react-timeseries-charts";

type Props = {
    chart: {},
    geojson: {},
};

export default class SimpleChartContainer extends Component {
    props: Props;

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
        const myStyler = styler([
            {key: "value1", color: "#2ca02c", width: 1},
            {key: "value", color: "#9467bd", width: 2}
        ]);
        return (
            <div style={{ padding: '10px' }}>
                <ChartContainer timeRange={timeseries.timerange()} width={480} showGrid={true}>
                    <ChartRow height="200">
                        <YAxis id="y" label={chart.key} min={timeseries.min()} max={timeseries.max()} width="70" type="linear" format=",.2f"/>
                        <Charts>
                            <LineChart axis="y" series={timeseries} style={myStyler} />
                        </Charts>
                    </ChartRow>
                </ChartContainer>
            </div>
        );
    }
}
