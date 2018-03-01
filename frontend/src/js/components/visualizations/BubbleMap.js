// @flow weak

import React, { Component } from 'react';
import * as d3 from 'd3';
import { Source, Layer, Feature } from 'react-mapbox-gl';

import type { GeoJSON } from '../../types/MapTypes';

type Props = {
    pointDecorator: PointDecorator,
    data: GeoJSON,
    click: () => mixed,
};

function onToggleHover(cursor) {
    const body = document.body;
    if (body) {
        body.style.cursor = cursor;
    }
}

export default class BubbleMap extends Component {
    props: Props;
    state: {
        circleColor: {},
        circleRadius: {},
    };

    constructor(props: Props) {
        super(props);
        this.state = {
            circleColor: {},
            circleRadius: {},
        };
    }

    componentWillMount() {
        let { circleColor, circleRadius } = this.state;
        circleColor = this.fkColorToMapboxCircleColor();
        circleRadius = this.fkSizeToMapboxCircleRadius();
        this.setState({
            circleColor,
            circleRadius
        });
    }

    fkColorToMapboxCircleColor() {
        const { pointDecorator, data } = this.props;
        const { color } = pointDecorator.points;
        let circleColor = {};
        if (color.type === 'constant') {
            circleColor = color.colors[0].color;
        } else {
            const { dateKey, bounds } = color;
            let min;
            let max;
            if (bounds) {
                min = bounds[0];
                max = bounds[1];
            } else {
                const values = data.features.map(f => f.properties[dateKey]);
                min = d3.min(values);
                max = d3.max(values);
            }
            const interpolate = d3.interpolateNumber(min, max);
            const stops = color.colors.map(c => [
                interpolate(c.location),
                c.color,
            ]);
            circleColor = {
                property: dateKey,
                type: 'exponential',
                stops,
            };
        }
        return circleColor;
    }

    fkSizeToMapboxCircleRadius() {
        const { pointDecorator, data } = this.props;
        const { size } = pointDecorator.points;
        let circleRadius = {};
        if (size.type === 'constant') {
            circleRadius = size.bounds[0];
        } else {
            const { dateKey } = size;
            const values = data.features.map(f => f.properties[dateKey]);
            circleRadius = {
                property: dateKey,
                type: 'exponential',
                stops: [
                    [d3.min(values), size.bounds[0]],
                    [d3.max(values), size.bounds[1]],
                ],
            };
        }
        return circleRadius;
    }

    render() {
        const { data, click } = this.props;
        const { circleColor, circleRadius } = this.state;

        if (!data || data.features.length === 0) {
            return <div></div>;
        }

        return (
            <div>
                <Layer type="circle" id="circle-markers" paint={ { 'circle-color': circleColor, 'circle-radius': circleRadius } }>
                { data.features.map((f, i) => (
                        <Feature key={i} onClick={ click.bind(this, f) }
                            onMouseEnter={ onToggleHover.bind(this, 'pointer') }
                            onMouseLeave={ onToggleHover.bind(this, '') }
                            coordinates={ f.geometry.coordinates } />
                      )) }
                </Layer>
            </div>
        );
    }
}
