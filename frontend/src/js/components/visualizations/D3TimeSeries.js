// @flow weak

import log from 'loglevel';
import * as d3 from 'd3';
import type { Margin } from '../../types/D3Types';

type Props = {
    margin: Margin,
    data: {}[],
    node: HTMLElement,
};

function parseDates(data) {
    const parseTime = d3.timeParse('%L');
    log.info(data);
}

export default class D3TimeSeries {
    props: Props;
    state: {
        scales: {
            x: () => {},
            y: () => {},
        },
    };

    constructor(props) {
        log.info('D3Chart');
        this.props = props;
        this.state = {
            scales: {
                x: () => {
                },
                y: () => {
                },
            },
        };
        this.createElements();
    }

    updateData(data) {
        log.info('updateData');
        // parseDates(data);
        this.props.data = data;
        const domain = this.getDomain();
        const scales = this.getScales(domain);
        this.setState({
            scales
        });
    }

    setState(state) {
        Object.assign(this.state, state);
        this.updateElements();
    }

    getDomain() {
        const { data } = this.props;
        const x = d3.extent(data, d => d.date);
        const y = [d3.min(data, d => d.value), d3.max(data, d => d.value)];
        log.info(x);
        log.info(y);
        return {
            x,
            y
        };
    }

    getScales(domain) {
        const { node, margin } = this.props;
        if (!domain || !node) {
            return null;
        }

        const width = node.offsetWidth - (margin.left + margin.right);
        const height = node.offsetHeight - (margin.top + margin.bottom);

        const x = d3.scaleTime().range([0, width]).domain(domain.x);
        const y = d3.scaleLinear().range([height, 0]).domain(domain.y);

        return {
            x,
            y
        };
    }

    line() {
        const { scales } = this.state;
        const d3Line = d3.line()
            .x(d => scales.x(d.date))
            .y(d => scales.y(d.value));
        return d3Line;
    }

    createElements() {
        log.info('createElements');
        const { node, margin, data } = this.props;

        const height = node.offsetHeight - (margin.top + margin.bottom);

        const svg = d3
            .select(node)
            .append('svg')
            .attr('class', 'd3')
            .attr('width', node.offsetWidth)
            .attr('height', node.offsetHeight);

        const chart = svg
            .append('g')
            .attr('transform', `translate(${margin.left}, ${margin.top})`)
            .attr('class', 'chart');

        chart.append('g').attr('class', 'x axis').attr('transform', `translate(0, ${height})`);

        chart.append('g').attr('class', 'y axis');

        chart.append('path').attr('class', 'd3-line');
        this.updateData(data);
    }

    updateElements() {
        log.info('updateElements');
        const self = this;
        const { node, data } = this.props;
        const { scales } = this.state;
        log.info(scales);
        const xAxis = d3.axisBottom(scales.x);
        const yAxis = d3.axisLeft(scales.y);

        const chart = d3.select(node).select('.chart');

        chart.select('g.y').transition().duration(1000).call(yAxis);

        chart.select('g.x').transition().duration(1000).call(xAxis);

        const d3Line = d3.line()
            .x(d => {
                log.info(scales.x('Tue Jun 30 2015 00:00:00 GMT-0700 (PDT)'));
                return scales.x(d.date);
            })
            .y(d => scales.y(d.value));

        chart.selectAll('.d3-line')
            .datum(data)
            .transition()
            .duration(1000)
            .attr('d', (d) => {
                return d3Line(d);
            });
    }
}
