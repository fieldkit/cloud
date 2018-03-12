// @flow weak

import React, { Component } from 'react';
import PropTypes from "prop-types";

class Range {
    constructor(title, range) {
        this.title = title;
        this.range = range;
    }
}

export default class CriteriaPanel extends Component {
    constructor(props) {
        super(props)
        this.state = {
            open: false
        }
    }

    onChangeTimeCiteria(newRange) {
        const { onChangeTimeCiteria } = this.props;

        const now = new Date();
        const start = new Date();
        start.setTime(now.getTime() + (newRange * 1000));

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
        const ranges = [
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
                new Range('Last 3 hours', -60 * 60 * 3),
                new Range('Last 6 hours', -60 * 60 * 6),
                new Range('Last 12 hours', -60 * 60 * 12),
                new Range('Last 24 hours', -60 * 60 * 24),
            ]
        ];

        return <div className="drawer">
            {ranges.map((column, j) => <ul key={j}>
                    {column.map((r, i) => {
                        return <li key={i} onClick={ this.onChangeTimeCiteria.bind(this, r.range) }>{r.title}</li>;
                    })}
                </ul>
            )}
        </div>;
    }

    render() {
        return <div className="criteria-panel">
            <div onClick={this.onOpenToggle.bind(this)} className="button">Critera</div>
            {this.state.open && this.renderOpen()}
        </div>;
    }
}
