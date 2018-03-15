import React, { Component } from 'react';

import * as d3 from "d3";
import { NodeGroup } from 'resonance';

function polarToCartesian(centerX, centerY, radius, angleInDegrees) {
    var angleInRadians = (angleInDegrees - 90) * Math.PI / 180.0;
    return {
        x: centerX + (radius * Math.cos(angleInRadians)),
        y: centerY + (radius * Math.sin(angleInRadians))
    };
}

export class FancyMenu extends Component {
    constructor(props) {
        super(props);

        const radius = 100;
        const width = radius * 4;
        const height = radius * 4;

        this.state = {
            radius: radius,
            width: width,
            height: height,
        }
    }

    componentWillReceiveProps(nextProps) {

    }

    componentDidMount() {
        /*
        <div ref="menu"></div>

        const { radius, width, height } = this.state;

        const nodes = [];
        for (let i = 0; i < 5; i++) {
            const angle = (i / (5 / 1)) * Math.PI * 2;
            const x = (radius * Math.cos(angle));
            const y = (radius * Math.sin(angle));
            nodes.push({'id': i, 'x': x, 'y': y});
        }

        const context = d3.select(this.refs.menu)
            .append('svg')
            .attr("width", width)
            .attr("height", height)
            .append("g")
            .attr('transform', `translate(${width / 2}, ${height / 2})`);

        const items = context.selectAll('circle')
            .data(nodes)
            .enter();

        items.append("circle")
            .attr("fill", "black")
            .attr('r', 40)
            .style("fill-opacity", 0.8)
            .attr('cx', function (d, i) {
                return d.x;
            })
            .attr('cy', function (d, i) {
                return d.y;
            });


        items.append('text')
            .attr("fill", "white")
            .attr('x', function (d, i) {
                return d.x;
            })
            .attr('y', function (d, i) {
                return d.y + 6;
            })
            .text(function(d) {
                return "Item";
            })
            .style('font-family', "sans-serif")
            .style("text-anchor", "middle")
            .style("font-size", "20px")
            .style("font-weight", "100");
        */
    }

    getData() {
        const { options } = this.props;
        const angle = i => (i / options.length) * 360.0;

        return options.map((option, i) => {
            return {
                key: 'key-' + i,
                text: option.title,
                angle: angle(i),
            };
        });
    }

    renderObject(data, i) {
        const { radius } = this.state;
        const center = polarToCartesian(0, 0, radius, data.angle);

        // TODO Better scaling for the text width. Possible truncate if they're too long?
        let fontSize = 20;
        let backgroundRadius = 40;
        if (data.text.length >= 12) {
            fontSize = 10;
        }
        else if (data.text.length >= 8) {
            fontSize = 16;
        }

        return {
            g: {
                key: data.key,
                id: data.key,
            },
            background: {
                id: data.key + "-bg",
                cx: [0, center.x],
                cy: [0, center.y],
                r: backgroundRadius,
                opacity: [0, 0.7],
                style: {
                    fill: 'black',
                    strokeWidth: 2,
                    stroke: 'black',
                },
            },
            text: {
                x: [0, center.x],
                y: [0, center.y + 6],
                textAnchor: 'middle',
                fontFamily: 'sans-serif',
                fontWeight: 100,
                fontSize: fontSize + 'px',
                fill: 'white',
                opacity: [0, 1.0],
            }
        };
    }

    onClosed() {
        const { onClosed } = this.props;

        onClosed();
    }

    componentDidUpdate() {
        const { options } = this.props;
        const data = this.getData();

        data.forEach((d, i) => {
            d3.select('#' + d.key)
                .on('mouseover', function () {
                    d3.select(`#${d.key}-bg`).style('opacity', 1);
                    document.body.style.cursor = "pointer";
                })
                .on('mouseout', function () {
                    d3.select(`#${d.key}-bg`).style('opacity', 0.7);
                    document.body.style.cursor = "default";
                })
                .on('click', ev => {
                    d3.event.stopPropagation();
                    if (options[i].onClick) {
                        if (!options[i].onClick()) {
                            this.onClosed();
                        }
                    }
                });
        })
    }

    render() {
        const { visible, position } = this.props;
        const { width, height } = this.state;

        const data = visible ? this.getData() : [];

        return <div style={{ width: "100%", height: "100%", zIndex: visible ? 20000 : -10000, position: 'fixed', transform: `translate(${position.x - (width / 2)}px, ${position.y - (height / 2)}px)` }} onClick={ this.onClosed.bind(this) }>
            <svg width={width} height={height} preserveAspectRatio="none">
                <g transform={`translate(${width / 2}, ${height / 2})`}>
                    <NodeGroup data={data} keyAccessor={(d) => d.key}
                        enter={(d, i) => {
                            return {
                                ...this.renderObject(d, i),
                                timing: {
                                    duration: 200,
                                    delay: i * 50,
                                    ease: d3.easeCircleInOut
                                },
                            };
                        }}
                        start={(d, i) => {
                            return {
                                ...this.renderObject(d, i)
                            };
                        }}
                        update={(d, i) => {
                            return {
                                ...this.renderObject(d, i),
                                timing: {
                                    duration: 400,
                                    delay: i * 100,
                                    ease: d3.easeCircleInOut
                                },
                                events: {
                                    start() {
                                    },
                                    interrupt() {
                                    },
                                    end() {
                                    },
                                }
                            };
                        }}>
                        {(nodes) => {
                            return (
                                <g>
                                    {nodes.map( ({ key, data, state} ) => {
                                        return (
                                            <g {...state.g}>
                                                <circle {...{ ...state.background, style: {...state.background.style } } } />
                                                <text {...state.text}>{data.text}</text>
                                            </g>
                                        );
                                    })}
                                </g>
                            );
                        }}
                    </NodeGroup>
                </g>
            </svg>
        </div>;
    }
};

