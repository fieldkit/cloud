// @flow weak

import React, { Component } from 'react';
import log from 'loglevel';
import ChartContainer from '../containers/ChartContainer';

type Props = {};

function generateScatterplotData() {
    const scatterplotData = [];
    for (let i = 0; i < 10; i++) {
        scatterplotData.push({
            id: i,
            value1: Math.round(Math.random() * 100),
            value2: Math.round(Math.random() * 100),
            value3: Math.round(Math.random() * 20),
        });
    }
    return scatterplotData;
}

function generateTimeSeriesData() {
    const dayInMillis = 86400000;
    const aWeekAgo = Date.now() - (dayInMillis * 7);
    const timeSeriesData = [];
    for (let i = 0; i < 6; i++) {
        timeSeriesData.push({
            id: i,
            value: Math.round(Math.random() * 100),
            date: new Date(aWeekAgo + (i * dayInMillis)),
        });
    }
    log.info(timeSeriesData);
    return timeSeriesData;
}

export default class Sensors extends Component {
    props: Props;
    state: {
    scatterplotData: {}[],
    timeSeriesData: {}[],
    };

    constructor(props: Props) {
        super(props);
        this.state = {
            scatterplotData: generateScatterplotData(),
            timeSeriesData: generateTimeSeriesData(),
        };
    }

    randomize() {
        let { scatterplotData, timeSeriesData } = this.state;
        scatterplotData = generateScatterplotData();
        timeSeriesData = generateTimeSeriesData();
        this.setState({
            scatterplotData,
            timeSeriesData
        });
    }

    render() {
        const { scatterplotData, timeSeriesData } = this.state;
        return (
            <div className="sensors page">
                <ChartContainer type="scatterplot" margin={ { top: 16, right: 16, bottom: 32, left: 32, } } data={ scatterplotData } />
                <ChartContainer type="time-series" margin={ { top: 16, right: 16, bottom: 32, left: 32, } } data={ timeSeriesData } />
                <button onClick={ this.randomize.bind(this) }>Random</button>
            </div>
        );
    }
}
