// @flow weak

import _ from 'lodash';
import React, { Component } from 'react';
import { Layer, Feature } from 'react-mapbox-gl';

import type { GeoJSON } from '../../types';

class DataColors {
    colors: Array<string>
    index: number

    constructor() {
        this.colors = [
            '#ECA400',
            '#4CE0D2',
            '#FC60A8',
            '#EF2917',
            '#E07BE0',
            '#ECC8AF',
            '#932F6D',
            '#726DA8',
            '#CE796B',
            '#A0D2DB',
            '#E7AD99',
            '#7D8CC4',
            '#C490D1',
            '#B8336A',
            '#C18C5D',
            '#495867',
            '#008148',
        ];
        this.colors = [
            'rgb(137, 15, 2)',
            'rgb(88, 20, 12)',
            'rgb(153, 68, 10)',
            'rgb(193, 92, 23)',
            'rgb(150, 115, 2)',
            'rgb(204, 163, 0)',
            'rgb(63, 104, 51)',
            'rgb(47, 87, 94)',
            'rgb(100, 176, 200)',
            'rgb(5, 43, 81)',
            'rgb(10, 80, 161)',
            'rgb(88, 68, 119)',
            'rgb(63, 43, 91)',
            'rgb(81, 23, 73)',
            'rgb(226, 77, 66)',
            'rgb(191, 27, 0)',
            'rgb(239, 132, 60)',
            'rgb(244, 213, 152)',
            'rgb(229, 172, 14)',
            'rgb(154, 196, 138)',
            'rgb(80, 134, 66)',
            'rgb(110, 208, 224)',
            'rgb(101, 197, 219)',
            'rgb(10, 67, 124)',
            'rgb(68, 126, 188)',
            'rgb(97, 77, 147)',
            'rgb(214, 131, 206)',
            'rgb(109, 31, 98)',
            'rgb(234, 100, 96)',
            'rgb(224, 117, 45)',
            'rgb(249, 147, 78)',
            'rgb(252, 234, 202)',
            'rgb(234, 184, 57)',
            'rgb(183, 219, 171)',
            'rgb(98, 158, 81)',
            'rgb(112, 219, 237)',
            'rgb(130, 181, 216)',
            'rgb(31, 120, 193)',
            'rgb(174, 162, 224)',
            'rgb(112, 93, 160)',
            'rgb(229, 168, 226)',
            'rgb(150, 45, 130)',
            'rgb(242, 145, 145)',
            'rgb(252, 226, 222)',
            'rgb(249, 186, 143)',
            'rgb(249, 226, 210)',
            'rgb(242, 201, 109)',
            'rgb(224, 249, 215)',
            'rgb(126, 178, 109)',
            'rgb(207, 250, 255)',
            'rgb(186, 223, 244)',
            'rgb(81, 149, 206)',
            'rgb(222, 218, 247)',
            'rgb(128, 110, 183)',
            'rgb(249, 217, 249)',
            'rgb(186, 67, 169)',
        ]
        this.index = 0;
    }

    get() {
        const c = this.colors[this.index];
        this.index = (this.index + 1) % this.colors.length;
        return c;
    }
};

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

        const temporal = _(data).filter(s => s.summary != null && s.summary.spatial != null).map(s => {
            return _(s.summary.temporal).filter(c => s.geometries[c.id]).map(c => {
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
