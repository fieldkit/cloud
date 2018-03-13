// @flow weak

import React, { Component } from 'react';

export default class SourceOverview extends Component {
    props: {
        style?: any,
        data: any,
        onShowChart?: any,
    }

    onShowChart(key) {
        const { data, onShowChart } = this.props;
        const { source } = data;

        if (onShowChart) {
            onShowChart({
                id: ["chart", source.id, key].join('-'),
                sourceId: source.id,
                keys: [key]
            });
        }
    }

    renderReadings(data) {
        const { summary } = data;

        return (
            <table style={{ padding: '5px', width: '100%' }} className="feature-data">
                <thead>
                    <tr>
                        <th>Reading</th>
                    </tr>
                </thead>
                <tbody>
                { summary.readings.map(reading => (
                    <tr key={reading.name}>
                        <td style={{ width: '50%' }}> <div onClick={() => this.onShowChart(reading.name)} style={{ cursor: 'pointer' }}>{reading.name}</div> </td>
                    </tr>
                )) }
                </tbody>
            </table>
        )
    }

    render() {
        const { style, data } = this.props;

        return (
            <div style={{ ...style }} className="feature-panel">
                {this.renderReadings(data)}
            </div>
        );
    }
};

