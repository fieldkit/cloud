// @flow weak

import log from 'loglevel';
import * as d3 from 'd3';

type Props = {
    margin: { top: number, left: number, bottom: number, right: number };
    data: {}[];
    node: HTMLElement;
}

export default class D3Scatterplot {

    props: Props;
    state: {
        scales: {
            x: () => {},
            y: () => {},
            z: () => {},
        },
    }

    constructor(props) {
        this.props = props;
        this.state = {
            scales: {
                x: () => {
                },
                y: () => {
                },
                z: () => {
                },
            },
        };
        this.createElements();
    }

    updateData(data) {
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
        const x = [
            d3.min(data, d => d.value1),
            d3.max(data, d => d.value1),
        ];
        const y = [
            d3.min(data, d => d.value2),
            d3.max(data, d => d.value2),
        ];
        const z = [
            d3.min(data, d => d.value3),
            d3.max(data, d => d.value3),
        ];
        return {
            x,
            y,
            z
        };
    }

    getScales(domain) {
        const { node, margin } = this.props;
        if (!domain || !node) {
            return null;
        }

        const width = node.offsetWidth - (margin.left + margin.right);
        const height = node.offsetHeight - (margin.top + margin.bottom);

        const x = d3.scaleLinear().range([0, width]).domain(domain.x);
        const y = d3.scaleLinear().range([height, 0]).domain(domain.y);
        const z = d3.scaleLinear().range([5, 10]).domain(domain.z);

        return {
            x,
            y,
            z
        };
    }

    createElements() {
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

        chart.append('g')
            .attr('class', 'x axis')
            .attr('transform', `translate(0, ${height})`);

        chart.append('g')
            .attr('class', 'y axis');

        chart.append('g').attr('class', 'd3-points');
        this.updateData(data);
    }

    updateElements() {
        const { node, data } = this.props;
        const { scales } = this.state;
        const xAxis = d3.axisBottom(scales.x);
        const yAxis = d3.axisLeft(scales.y);

        const chart = d3.select(node).select('.chart');

        chart.select('g.y')
            .transition()
            .duration(1000)
            .call(yAxis);

        chart.select('g.x')
            .transition()
            .duration(1000)
            .call(xAxis);

        const points = chart.selectAll('.d3-points');
        const point = points.selectAll('.d3-point').data(data);

        point.enter()
            .append('circle')
            .attr('class', 'd3-point')
            .merge(point)
            .transition()
            .duration(1000)
            .attr('cx', d => scales.x(d.value1))
            .attr('cy', d => scales.y(d.value2))
            .attr('r', d => scales.z(d.value3))
            .attr('fill', 'white');

        point.exit().remove();
    }
}
