// @flow weak

import _ from 'lodash';
import React, { Component } from 'react';
import { Source, Layer, Feature } from 'react-mapbox-gl';

import type { GeoJSON } from '../../types/MapTypes';

type Props = {
    data: GeoJSON,
};

function onToggleHover(cursor) {
    const body = document.body;
    if (body) {
        body.style.cursor = cursor;
    }
}

export default class ClusterMap extends Component {
    props: Props;
    state: {
    };

    render() {
        const { data, click } = this.props;
        const { visibleFeatures } = this.props;

        const temporalColor = '#F5B041';

        const temporalRadius = 20;

        const spatialColor = '#884EA0';

        const spatialRadius = {
            property: "numberOfFeatures",
            stops: [
                [500, 10],
                [1000, 15],
                [5000, 20],
                [10000, 25],
            ]
        };

        const temporal = _(data).map(r => r.summary.temporal).flatten().map(c => {
            return {
                geometry: {
                    coordinates: c.geometry,
                },
                properties: {
                }
            };
        }).value();

        const spatial = _(data).map(r => r.summary.spatial).flatten().map(c => {
            return {
                geometry: {
                    coordinates: c.centroid,
                },
                properties: {
                    numberOfFeatures: c.numberOfFeatures
                }
            };
        }).value();

        return (
            <div>
                <Layer type="circle" id="spatial-clusters" paint={ { 'circle-color': spatialColor, 'circle-radius': spatialRadius } }>
                    { spatial.map((f, i) => (
                        <Feature key={i} onClick={ () => console.log(f) }
                            onMouseEnter={ onToggleHover.bind(this, 'pointer') }
                            onMouseLeave={ onToggleHover.bind(this, '') }
                            properties={ f.properties }
                            coordinates={ f.geometry.coordinates } />
                    )) }
                </Layer>

                <Layer type="line" id="temporal-clusters" paint={ { 'line-color': temporalColor, 'line-width': 8 } }>
                    { temporal.map((f, i) => (
                        <Feature key={i} onClick={ () => console.log(f) }
                            onMouseEnter={ onToggleHover.bind(this, 'pointer') }
                            onMouseLeave={ onToggleHover.bind(this, '') }
                            properties={ f.properties }
                            coordinates={ f.geometry.coordinates } />
                    )) }
                </Layer>
            </div>
        );
    }
}
