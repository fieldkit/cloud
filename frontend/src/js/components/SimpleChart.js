// @flow weak

import _ from 'lodash';
import moment from 'moment';

import React, { Component } from 'react';

import { ResponsiveContainer, XAxis, YAxis, CartesianGrid, Legend, AreaChart, Area } from 'recharts';

import { Loading } from './Loading';

type Props = {
    chartDef: any,
    data: any,
};

export default class SimpleChart extends Component {
    props: Props;

    render() {
        const { chartDef, data } = this.props;

        if (data.loading) {
            return <Loading />;
        }

        const samples = _(data.series[0].rows).map(r => {
            const row = {
                x: r[0],
            };
            _.each(chartDef.keys, (key, i) => {
                row[key] = r[i + 1];
            });
            return row;
        }).value();

        return <div>
            <ResponsiveContainer width="100%" height={300}>
                <AreaChart data={samples} margin={{top: 5, right: 30, left: 20, bottom: 5}}>
                    <XAxis dataKey="x"  tickFormatter={t => moment.unix(t).format('HH:mm')} />
                    <YAxis />
                    <CartesianGrid stroke="#626262" />
                    <Legend verticalAlign="top" height={36} />
                    { chartDef.keys.map(k => <Area key={k} type="monotone" dataKey={k} dot={false} animationDuration={100} stroke="#82ca9d" strokeWidth={2} fill="#82ca9d" fillOpacity={0.3} />)}
                </AreaChart>
            </ResponsiveContainer>
        </div>;
    }
}
