// @flow weak

import _ from 'lodash';

import React, { Component } from 'react';

class Range {
    title: string
    range: number
    active: bool

    constructor(title, range, active) {
        this.title = title;
        this.range = range;
        this.active = active || false;
    }
}

type Props = {
    onChangeTimeCiteria: any,
};

type State = {
    open: bool,
};

export default class CriteriaPanel extends Component {
    props: Props
    state: State
    ranges: Array<Array<Range>>

    constructor(props) {
        super(props)

        this.state = {
            open: false
        }

        this.ranges = [
            [
                new Range('Last 2 days', -86400 * 2),
                new Range('Last 7 days', -86400 * 7),
                new Range('Last 30 days', -86400 * 30),
                new Range('Yesterday', -86400),
            ],
            [
                new Range('Last 5 minutes', -60 * 5),
                new Range('Last 15 minutes', -60 * 15),
                new Range('Last 30 minutes', -60 * 30),
                new Range('Last hour', -60 * 60),
                new Range('Last 3 hours', -60 * 60 * 3, true),
                new Range('Last 6 hours', -60 * 60 * 6),
                new Range('Last 12 hours', -60 * 60 * 12),
                new Range('Last 24 hours', -60 * 60 * 24),
            ]
        ];
    }

    componentWillMount() {
        this.onChangeTimeCiteria(this.getActiveRange());
    }

    getActiveRange() {
        return _(this.ranges).flatten().filter(r => r.active).first();
    }

    onChangeTimeCiteria(newRange) {
        const { onChangeTimeCiteria } = this.props;

        const oldActive = this.getActiveRange()
        oldActive.active = false;
        newRange.active = true;

        const now = new Date();
        const start = new Date();
        start.setTime(now.getTime() + (newRange.range * 1000));

        onChangeTimeCiteria({
            startTime: Math.trunc(start.getTime() / 1000),
            endTime: Math.trunc(now.getTime() / 1000),
        });

        this.setState({
            open: false
        })
    }

    onOpenToggle() {
        this.setState({
            open: !this.state.open
        })
    }

    renderOpen() {
        return <div className="drawer">
            {this.ranges.map((column, j) => <ul key={j}>
                    {column.map((r, i) => {
                        return <li key={i} onClick={ this.onChangeTimeCiteria.bind(this, r) }>{r.title}</li>;
                    })}
                </ul>
            )}
        </div>;
    }

    render() {
        const active = this.getActiveRange();
        return <div className="criteria-panel">
            <div onClick={this.onOpenToggle.bind(this)} className="button">{active.title}</div>
            {this.state.open && this.renderOpen()}
        </div>;
    }
}
