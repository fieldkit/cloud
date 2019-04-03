// @flow weak

import _ from 'lodash';
import React, { Component } from 'react';
import { Layer, Feature } from 'react-mapbox-gl';

import type { GeoJSON } from '../../types';

import { DataColors } from '../../common/colors';

export default class ClusterMap extends Component {
    props: {
        data: GeoJSON,
        onClick: () => mixed,
    }

    state: {
    }

    onToggleHover(cursor: string, { map }: { map: any }) {
        map.getCanvas().style.cursor = cursor;
    }

    render() {
        const { data, onClick } = this.props;

        const colors = new DataColors();

        const temporalColor ={
            type: 'identity',
            property: 'color'
        };

        const temporalWidth = 8;

        const temporal = _(data)
            .filter(s => s.summary != null && s.summary.spatial != null)
            .sortBy(s => s.source.id)
            .map(s => {
                return _(s.summary.temporal)
                    .filter(c => s.geometries && s.geometries[c.id])
                    .sortBy(c => c.id)
                    .map(c => {
                        return {
                            geometry: {
                                coordinates: s.geometries[c.id].geometry,
                            },
                            properties: {
                                name: s.summary.name,
                                source: s.summary.id,
                                cluster: c.id,
                                type: 'temporal',
                                color: colors.get(),
                                numberOfFeatures: c.numberOfFeatures
                            }
                        };
                    })
                    .value();
            })
            .flatten()
            .value();

        const spatialColor ={
            type: 'identity',
            property: 'color'
        };

        const spatialRadius = {
            property: "numberOfFeatures",
            stops: [
                [500, 10],
                [1000, 15],
                [5000, 20],
                [10000, 25],
            ]
        };

        const spatial = _(data).filter(s => s.summary != null && s.summary.spatial != null).map(s => {
            return _(s.summary.spatial).map(c => {
                return {
                    geometry: {
                        coordinates: c.centroid,
                    },
                    properties: {
                        name: s.summary.name,
                        source: s.summary.id,
                        cluster: c.id,
                        type: 'spatial',
                        color: colors.get(),
                        numberOfFeatures: c.numberOfFeatures
                    }
                };
            }).value();
        }).flatten().value();

        return (
            <div>
                <Layer type="circle" id="spatial-clusters" paint={ { 'circle-color': spatialColor, 'circle-radius': spatialRadius } }>
                    { spatial.map((f, i) => (
                        <Feature key={i} onClick={ () => onClick.bind(this, f) }
                            onMouseEnter={ this.onToggleHover.bind(this, 'pointer') }
                            onMouseLeave={ this.onToggleHover.bind(this, '') }
                            properties={ f.properties }
                            coordinates={ f.geometry.coordinates } />
                    )) }
                </Layer>

                <Layer type="line" id="temporal-clusters" paint={ { 'line-color': temporalColor, 'line-width': temporalWidth } }>
                    { temporal.map((f, i) => (
                        <Feature key={i} onClick={ () => onClick.bind(this, f) }
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
