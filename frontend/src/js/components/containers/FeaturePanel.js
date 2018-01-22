// @flow weak

import _ from 'lodash';
import moment from 'moment';

import React, { Component } from 'react';

import ChartComponent from '../ChartComponent';

const panelStyle: React.CSSProperties = {
    backgroundColor: '#f9f9f9',
    color: "#000",
    borderTopLeftRadius: 2,
    borderTopRightRadius: 2,
    borderBottomLeftRadius: 2,
    borderBottomRightRadius: 2,
};

export default class FeaturePanel extends Component {
    props: {
        style: React.CSSProperties
    }

    constructor() {
        super();
        this.state = {
            chart: null
        };
    }

    componentWillReceiveProps(nextProps) {
        const { feature: oldFeature } = this.props;
        const { feature: newFeature } = nextProps;

        if (oldFeature.properties.id != newFeature.properties.id) {
            this.onHideChart();
        }
    }

    onShowChart(key) {
        const { feature } = this.props;
        this.setState({
            chart: {
                source: feature.properties.source,
                key: key
            }
        });
    }

    onHideChart() {
        this.setState({
            chart: null
        });
    }

    renderProperties(feature) {
        const properties = _(Object.entries(feature.properties.data));

        return (
            <table style={{ padding: '5px', width: '100%' }}>
                <thead>
                    <tr>
                        <th>Key</th>
                        <th>Value</th>
                    </tr>
                </thead>
                <tbody>
                { properties.value().map(([ k , v ]) => (
                    <tr key={k}>
                        <td style={{ width: '50%' }}> <a href="#" onClick={() => this.onShowChart(k)}>{k}</a> </td>
                        <td> {Number(v).toFixed(2)} </td>
                    </tr>
                )) }
                </tbody>
            </table>
        )
    }

    render() {
        const { style, feature } = this.props;
        const { chart } = this.state;

        if (chart) {
            return (
                <div style={{ ...panelStyle, ...style }}>
                    <div style={{ backgroundColor: '#d0d0d0', padding: '5px', fontWeight: 'bold' }} onClick={() => this.onHideChart()}><a href="#">{chart.key} - Back</a></div>
                    <ChartComponent chart={chart} />
                </div>
            );
        }

        const date = moment(new Date(feature.properties['timestamp'])).format('MMM Do YYYY, h:mm:ss a');
        return (
            <div style={{ ...panelStyle, ...style }}>
                <div style={{ backgroundColor: '#d0d0d0', padding: '5px', fontWeight: 'bold' }}>{date}</div>
                {this.renderProperties(feature)}
            </div>
        );
    }
}
