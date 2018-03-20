import React, { Component } from 'react';

import * as d3 from "d3";
import { NodeGroup } from 'resonance';

function polarToCartesian(centerX, centerY, radius, degrees) {
    const radians = (degrees - 90) * Math.PI / 180.0;
    return {
        x: centerX + (radius * Math.cos(radians)),
        y: centerY + (radius * Math.sin(radians))
    };
}

export class FancyMenu extends Component {
    constructor(props) {
        super(props);

        const radius = 120;
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

    getOptionsData() {
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

    getCenterData() {
        const { center } = this.props;
        if (center && center.text) {
            return [{
                key: 'key-center',
                text: center.text,
            }];
        }
        return [];
    }

    // TODO Better scaling for the text width. Possible truncate if they're too long?
    applyFontScaling(text) {
        if (text.length >= 12) {
            return {
                radius: 40,
                fontSize: 10,
            };
        }
        else if (text.length >= 8) {
            return {
                radius: 40,
                fontSize: 14,
            };
        }
        return {
            radius: 40,
            fontSize: 20,
        };
    }

    renderOption(data, i) {
        const { radius } = this.state;
        const center = polarToCartesian(0, 0, radius, data.angle);
        const fontScaling = this.applyFontScaling(data.text);

        return {
            g: {
                key: data.key,
                id: data.key,
            },
            background: {
                id: data.key + "-bg",
                cx: [0, center.x],
                cy: [0, center.y],
                r: fontScaling.radius,
                opacity: [0, 0.7],
                style: {
                    fill: 'black',
                    strokeWidth: 2,
                    stroke: 'black',
                },
            },
            text: {
                x: [0, center.x],
                y: [0, center.y],
                textAnchor: 'middle',
                alignmentBaseline: 'middle',
                fontFamily: 'sans-serif',
                fontWeight: 100,
                fontSize: fontScaling.fontSize + 'px',
                fill: 'white',
                opacity: [0, 1.0],
            }
        };
    }

    renderCenter(data, i) {
        const fontScaling = this.applyFontScaling(data.text);
        return {
            g: {
                key: data.key,
                id: data.key,
            },
            background: {
                id: data.key + "-bg",
                cx: 0,
                cy: 0,
                r: [0, fontScaling.radius],
                opacity: [0, 0.7],
                style: {
                    fill: 'black',
                    strokeWidth: 2,
                    stroke: 'black',
                },
            },
            text: {
                x: 0,
                y: 0,
                textAnchor: 'middle',
                alignmentBaseline: 'middle',
                fontFamily: 'sans-serif',
                fontWeight: 100,
                fontSize: fontScaling.fontSize + 'px',
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
        const optionsData = this.getOptionsData();

        optionsData.forEach((d, i) => {
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
                    const { options } = this.props;
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

        const optionsData = visible ? this.getOptionsData() : [];
        const centerData = visible ? this.getCenterData() : [];

        return <div style={{ width: "100%", height: "100%", zIndex: visible ? 20000 : -10000, position: 'fixed', transform: `translate(${position.x - (width / 2)}px, ${position.y - (height / 2)}px)` }} onClick={ this.onClosed.bind(this) }>
            <svg width={width} height={height} preserveAspectRatio="none">
                <g transform={`translate(${width / 2}, ${height / 2})`}>
                    <NodeGroup data={centerData} keyAccessor={(d) => d.key}
                        enter={(d, i) => {
                            return {
                                ...this.renderCenter(d, i),
                                timing: {
                                    duration: 200,
                                    delay: i * 50,
                                    ease: d3.easeCircleInOut
                                },
                            };
                        }}
                        start={(d, i) => {
                            return {
                                ...this.renderCenter(d, i),
                            };
                        }}
                        update={(d, i) => {
                            return {
                                ...this.renderCenter(d, i),
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
                    <NodeGroup data={optionsData} keyAccessor={(d) => d.key}
                        enter={(d, i) => {
                            return {
                                ...this.renderOption(d, i),
                                timing: {
                                    duration: 200,
                                    delay: i * 50,
                                    ease: d3.easeCircleInOut
                                },
                            };
                        }}
                        start={(d, i) => {
                            return {
                                ...this.renderOption(d, i),
                            };
                        }}
                        update={(d, i) => {
                            return {
                                ...this.renderOption(d, i),
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

