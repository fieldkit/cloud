// @flow weak

import React, { Component } from 'react';
import * as d3 from 'd3';
import { Layer, Feature } from 'react-mapbox-gl';

import type { GeoJSON } from '../../types';

export default class BubbleMap extends Component {
    props: {
        pointDecorator: PointDecorator,
        data: GeoJSON,
        onClick: () => mixed,
    }
    state: {
        circleColor: {},
        circleRadius: {},
    }

    constructor(props) {
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
            const { dataKey, bounds } = color;
            let min;
            let max;
            if (bounds) {
                min = bounds[0];
                max = bounds[1];
            } else {
                const values = data.features.map(f => f.properties[dataKey]);
                min = d3.min(values);
                max = d3.max(values);
            }
            const interpolate = d3.interpolateNumber(min, max);
            const stops = color.colors.map(c => [
                interpolate(c.location),
                c.color,
            ]);
            circleColor = {
                property: dataKey,
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
            const { dataKey } = size;
            const values = data.features.map(f => f.properties[dataKey]);
            circleRadius = {
                property: dataKey,
                type: 'exponential',
                stops: [
                    [d3.min(values), size.bounds[0]],
                    [d3.max(values), size.bounds[1]],
                ],
            };
        }
        return circleRadius;
    }

    onToggleHover(cursor: string, { map }: { map: any }) {
        map.getCanvas().style.cursor = cursor;
    }

    render() {
        const { data, onClick } = this.props;
        // const { circleColor, circleRadius } = this.state;

        if (!data || data.features.length === 0) {
            return <div></div>;
        }

        const circleColor ={
            type: 'identity',
            property: 'color'
        };

        const circleRadius = {
            type: 'identity',
            property: "size"
        };


        return (
            <div>
                <Layer type="circle" id="circle-markers" paint={ { 'circle-color': circleColor, 'circle-radius': circleRadius } }>
                { data.features.map((f, i) => (
                    <Feature key={i} onClick={ onClick.bind(this, f) }
                            onMouseEnter={ this.onToggleHover.bind(this, 'pointer') }
                            onMouseLeave={ this.onToggleHover.bind(this, '') }
                            properties={ f.properties }
                            coordinates={ f.geometry.coordinates } />
                      )) }
                </Layer>
            </div>
        );
    }
}
