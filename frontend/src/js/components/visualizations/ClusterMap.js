// @flow weak

import _ from 'lodash';
import React, { Component } from 'react';
import { Layer, Feature } from 'react-mapbox-gl';

import type { GeoJSON } from '../../types/MapTypes';

type Props = {
    data: GeoJSON,
};

class DataColors {
    constructor() {
        this.colors = [
            '#ECA400',
            '#4CE0D2',
            '#FC60A8',
            '#EF2917',
            // '#E07BE0',
            // '#ECC8AF',
            // '#932F6D',
            // '#726DA8',
            // '#CE796B',
            // '#A0D2DB',
            // '#E7AD99',
            // '#7D8CC4',
            // '#C490D1',
            // '#B8336A',
            // '#C18C5D',
            // '#495867'
            '#008148',
        ];
        this.index = 0;
    }

    get() {
        const c = this.colors[this.index];
        this.index = (this.index + 1) % this.colors.length;
        return c;
    }
};

export default class ClusterMap extends Component {
    props: Props;
    state: {
    };

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

        const temporal = _(data).map(s => {
            return _(s.summary.temporal).map(c => {
                return {
                    geometry: {
                        coordinates: c.geometry,
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
            }).value();
        }).flatten().value();

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

        const spatial = _(data).map(s => {
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
