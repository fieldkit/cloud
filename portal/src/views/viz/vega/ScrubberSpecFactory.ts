import _ from "lodash";
import { ChartSettings, SeriesData } from "./SpecFactory";
import { TimeRange } from "../common";

export { ChartSettings };

export class ScrubberSpecFactory {
    constructor(private readonly allSeries, private readonly settings: ChartSettings = ChartSettings.Container, private readonly dataEvents = null) {}

    create() {
        const first = this.allSeries[0]; // TODO
        const xDomainsAll = this.allSeries.map((series: SeriesData) => series.queried.timeRange);
        // We ignore extreme ranges here because of this.settings.timeRange
        const allRanges = [...xDomainsAll, this.settings.timeRange.toArray()];
        const timeRangeAll = TimeRange.mergeArraysIgnoreExtreme(allRanges).toArray();

        console.log("DATA EVENTS SCRUBBER", this.dataEvents)

        return {
            $schema: "https://vega.github.io/schema/vega/v5.json",
            description: "FK Scrubber Spec",
            autosize:
            {
                type: "fit-x",
            },
            background: "white",
            padding: 5,
            height: 50,
            style: "cell",
            config:
            {
                axis:
                {
                    labelFont: "Avenir Light",
                    labelFontSize: 12,
                    labelColor: "#2c3e50",
                    titleColor: "#2c3e50",
                    titleFont: "Avenir Light",
                    titleFontSize: 14,
                    titlePadding: 10,
                    tickSize: 10,
                    tickOpacity: 0,
                    domain: false
                },
                axisX:
                {
                    title: null,
                    tickSize: 20
                },
                style:
                {
                    cell:
                    {
                        fill: "#f4f5f7",
                        stroke: "transparent"
                    }
                }
            },
            data: [
                {
                    name: "brush_store"
                },
                {
                    name: "table",
                    values: first.queried.data,
                },
                {
                    name: "data_0",
                    source: "table",
                    transform:
                    [
                        {
                            type: "formula",
                            expr: "toDate(datum[\"time\"])",
                            as: "time"
                        }
                    ]
                },
                {
                    name: "data_1",
                    source: "data_0",
                    transform:
                    [
                        {
                            type: "impute",
                            field: "value",
                            groupby:
                            [],
                            key: "time",
                            method: "value",
                            value: 0
                        },
                        {
                            type: "stack",
                            groupby:
                            [
                                "time"
                            ],
                            field: "value",
                            sort:
                            {
                                field:
                                [],
                                order:
                                []
                            },
                            as:
                            [
                                "value_start",
                                "value_end"
                            ],
                            offset: "zero"
                        }
                    ]
                },
                {
                    name: "data_2",
                    source: "data_0",
                    transform:
                    [
                        {
                            type: "filter",
                            expr: "!length(data(\"brush_store\")) || vlSelectionTest(\"brush_store\", datum)"
                        },
                        {
                            type: "impute",
                            field: "value",
                            groupby:
                            [],
                            key: "time",
                            method: "value",
                            value: 0
                        },
                        {
                            type: "stack",
                            groupby:
                            [
                                "time"
                            ],
                            field: "value",
                            sort:
                            {
                                field:
                                [],
                                order:
                                []
                            },
                            as:
                            [
                                "value_start",
                                "value_end"
                            ],
                            offset: "zero"
                        }
                    ]
                },
                {
                    name: "data_events",
                    values: this.dataEvents,
                },
            ],
            signals:
            [
                {
                    name: "width",
                    init: "isFinite(containerSize()[0]) ? containerSize()[0] : 200",
                    on:
                    [
                        {
                            update: "isFinite(containerSize()[0]) ? containerSize()[0] : 200",
                            events: "window:resize"
                        }
                    ]
                },
                {
                    name: "unit",
                    value:
                    {},
                    on:
                    [
                        {
                            "events": "mousemove",
                            "update": "isTuple(group()) ? group() : unit"
                        }
                    ]
                },
                {
                    name: "brush",
                    update: "vlSelectionResolve(\"brush_store\", \"union\")"
                },
                // TODO block all other events when scrubber handle clicked
                {
                    name: "scrub_handle_left",
                    value:
                    {},
                    on:
                    [
                        {
                            events:
                            {
                                type: "mouseup",
                                marktype: "symbol",
                                markname: "left_scrub"
                            },
                            update: "scrub_handle_left + 1"
                        }
                    ]
                },
                {
                    name: "scrub_handle_right",
                    value:
                    {},
                    on:
                    [
                        {
                            events:
                            {
                                type: "mouseup",
                                marktype: "symbol",
                                markname: "right_scrub"
                            },
                            update: "scrub_handle_right + 1"
                        }
                    ]
                },
                {
                    name: "brush_x",
                    value:
                    [],
                    on:
                    [
                        {
                            // Update brush xy on area mousedown
                            events:
                            {
                                source: "scope",
                                type: "mousedown",
                                filter:
                                [
                                    //"!event.item || event.item.mark.name !== \"brush_brush\" || event.item.name !== \"left_scrub\" || event.item.name !== \"right_scrub\""
                                    "!event.item || event.item.mark.name !== \"brush_brush\""
                                ]
                            },
                            update: "[x(unit), x(unit)]"
                        },
                        {
                            //Update right extent of brush on mouse down
                            // TODO: filter does not seem to prevent mousedown on scrubber handle from happening
                            events:
                            {
                                source: "window",
                                type: "mousemove",
                                consume: true,
                                between:
                                [
                                    {
                                        source: "scope",
                                        type: "mousedown",
                                        filter:
                                        [
                                            //"!event.item || event.item.mark.name !== \"brush_brush\" || event.item.name !== \"left_scrub\" || event.item.name !== \"right_scrub\""
                                            "!event.item || event.item.mark.name !== \"brush_brush\""
                                        ]
                                    },
                                    {
                                        source: "window",
                                        type: "mouseup"
                                    }
                                ]
                            },
                            update: "[brush_x[0], clamp(x(unit), 0, width)]"
                        },
                        {
                            events:
                            {
                                signal: "brush_scale_trigger"
                            },
                            update: "[scale(\"x\", brush_time[0]), scale(\"x\", brush_time[1])]"
                        },
                        {
                            events:
                            [
                                {
                                    source: "view",
                                    type: "dblclick"
                                }
                            ],
                            update: "[0, 0]"
                        },
                        {
                            events:
                            {
                                signal: "brush_translate_delta"
                            },
                            update: "clampRange(panLinear(brush_translate_anchor.extent_x, brush_translate_delta.x / span(brush_translate_anchor.extent_x)), 0, width)"
                        },
                        {
                            events:
                            {
                                signal: "brush_zoom_delta"
                            },
                            update: "clampRange(zoomLinear(brush_x, brush_zoom_anchor.x, brush_zoom_delta), 0, width)"
                        }
                    ]
                },
                {
                    name: "brush_time",
                    on:
                    [
                        {
                            events:
                            {
                                signal: "brush_x"
                            },
                            update: "brush_x[0] === brush_x[1] ? null : invert(\"x\", brush_x)"
                        }
                    ]
                },
                {
                    name: "brush_scale_trigger",
                    value:
                    {},
                    on:
                    [
                        {
                            events:
                            [
                                {
                                    scale: "x"
                                }
                            ],
                            update: "(!isArray(brush_time) || (+invert(\"x\", brush_x)[0] === +brush_time[0] && +invert(\"x\", brush_x)[1] === +brush_time[1])) ? brush_scale_trigger : {}"
                        }
                    ]
                },
                {
                    name: "brush_tuple",
                    on:
                    [
                        {
                            events:
                            [
                                {
                                    signal: "brush_time"
                                }
                            ],
                            update: "brush_time ? {unit: \"layer_0\", fields: brush_tuple_fields, values: [brush_time]} : null"
                        }
                    ]
                },
                {
                    name: "brush_tuple_fields",
                    value:
                    [
                        {
                            field: "time",
                            channel: "x",
                            type: "R"
                        }
                    ]
                },
                {
                    name: "brush_translate_anchor",
                    value:
                    {},
                    on:
                    [
                        {
                            // start brush area translation
                            events:
                            [
                                {
                                    source: "scope",
                                    type: "mousedown",
                                    markname: "brush_brush",
                                }
                            ],
                            update: "{x: x(unit), y: y(unit), extent_x: slice(brush_x)}"
                        }
                    ]
                },
                {
                    name: "brush_translate_delta",
                    value:
                    {},
                    on:
                    [
                        {
                            // translate brush area after mousdown
                            events:
                            [
                                {
                                    source: "window",
                                    type: "mousemove",
                                    consume: true,
                                    between:
                                    [
                                        {
                                            source: "scope",
                                            type: "mousedown",
                                            markname: "brush_brush",
                                            filter: "event.item.name !== \"right_scrub\""
                                        },
                                        {
                                            source: "window",
                                            type: "mouseup"
                                        }
                                    ]
                                }
                            ],
                            update: "{x: brush_translate_anchor.x - x(unit), y: brush_translate_anchor.y - y(unit)}"
                        }
                    ]
                },
                {
                    name: "brush_zoom_anchor",
                    on:
                    [
                        {
                            events:
                            [
                                {
                                    source: "scope",
                                    type: "wheel",
                                    consume: true,
                                    markname: "brush_brush"
                                }
                            ],
                            update: "{x: x(unit), y: y(unit)}"
                        }
                    ]
                },
                {
                    name: "brush_zoom_delta",
                    on:
                    [
                        {
                            events:
                            [
                                {
                                    source: "scope",
                                    type: "wheel",
                                    consume: true,
                                    markname: "brush_brush"
                                }
                            ],
                            force: true,
                            update: "pow(1.001, event.deltaY * pow(16, event.deltaMode))"
                        }
                    ]
                },
                {
                    name: "brush_modify",
                    on:
                    [
                        {
                            events:
                            {
                                signal: "brush_tuple"
                            },
                            update: "modify(\"brush_store\", brush_tuple, true)"
                        }
                    ]
                }
            ],
            marks:
            [
                {
                    name: "brush_brush_bg",
                    type: "rect",
                    clip: true,
                    encode:
                    {
                        enter:
                        {
                            fill:
                            {
                                value: "#333"
                            },
                            fillOpacity:
                            {
                                value: 0.125
                            }
                        },
                        update:
                        {
                            x:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    signal: "brush_x[0]"
                                },
                                {
                                    value: 0
                                }
                            ],
                            y:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    value: 0
                                },
                                {
                                    value: 0
                                }
                            ],
                            x2:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    signal: "brush_x[1]"
                                },
                                {
                                    value: 0
                                }
                            ],
                            y2:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    field:
                                    {
                                        group: "height"
                                    }
                                },
                                {
                                    value: 0
                                }
                            ]
                        }
                    }
                },
                {
                    name: "layer_0_marks",
                    type: "area",
                    style:
                    [
                        "area"
                    ],
                    sort:
                    {
                        field: "datum[\"time\"]"
                    },
                    interactive: true,
                    from:
                    {
                        data: "data_1"
                    },
                    encode:
                    {
                        update:
                        {
                            orient:
                            {
                                value: "vertical"
                            },
                            fill:
                            {
                                value: "#DCDEDF"
                            },
                            description:
                            {
                                signal: "\"Time: \" + (timeFormat(datum[\"time\"], '%b %d, %Y')) + \"; value: \" + (format(datum[\"value\"], \"\"))"
                            },
                            x:
                            {
                                scale: "x",
                                field: "time"
                            },
                            y:
                            {
                                scale: "y",
                                field: "value_end"
                            },
                            y2:
                            {
                                scale: "y",
                                field: "value_start"
                            },
                            defined:
                            {
                                signal: "isValid(datum[\"time\"]) && isFinite(+datum[\"time\"]) && isValid(datum[\"value\"]) && isFinite(+datum[\"value\"])"
                            }
                        }
                    }
                },
                {
                    name: "layer_1_marks",
                    type: "area",
                    style:
                    [
                        "area"
                    ],
                    sort:
                    {
                        "field": "datum[\"time\"]"
                    },
                    interactive: false,
                    from:
                    {
                        data: "data_2"
                    },
                    encode:
                    {
                        update:
                        {
                            orient:
                            {
                                value: "vertical"
                            },
                            fill:
                            {
                                value: "#52b5e4"
                            },
                            description:
                            {
                                signal: "\"Time: \" + (timeFormat(datum[\"time\"], '%b %d, %Y')) + \"; value: \" + (format(datum[\"value\"], \"\"))"
                            },
                            x:
                            {
                                scale: "x",
                                field: "time"
                            },
                            y:
                            {
                                scale: "y",
                                field: "value_end"
                            },
                            y2:
                            {
                                scale: "y",
                                field: "value_start"
                            },
                            defined:
                            {
                                signal: "isValid(datum[\"time\"]) && isFinite(+datum[\"time\"]) && isValid(datum[\"value\"]) && isFinite(+datum[\"value\"])"
                            }
                        }
                    }
                },
                {
                    name: "brush_brush",
                    type: "rect",
                    clip: true,
                    encode:
                    {
                        enter:
                        {
                            fill:
                            {
                                value: "transparent"
                            }
                        },
                        update:
                        {
                            x:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    signal: "brush_x[0]"
                                },
                                {
                                    value: 0
                                }
                            ],
                            y:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    value: 0
                                },
                                {
                                    value: 0
                                }
                            ],
                            x2:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    signal: "brush_x[1]"
                                },
                                {
                                    value: 0
                                }
                            ],
                            y2:
                            [
                                {
                                    test: "data(\"brush_store\").length && data(\"brush_store\")[0].unit === \"layer_0\"",
                                    field:
                                    {
                                        group: "height"
                                    }
                                },
                                {
                                    value: 0
                                }
                            ],
                            stroke:
                            [
                                {
                                    test: "brush_x[0] !== brush_x[1]",
                                    value: "white"
                                },
                                {
                                    value: null
                                }
                            ]
                        }
                    }
                },
                {
                    type: "rect",
                    interactive: false,
                    encode: {
                      enter: {
                        y: {value: 0},
                        height: {value: 50},
                        width: {value: 2},
                        fill: "transparent"
                      },
                      update: {
                        x: {signal: "brush_x[0]"},
                        fill: {value: "white"}
                      }
                    }
                },
                {
                    type: "rect",
                    interactive: false,
                    encode: {
                      enter: {
                        y: {value: 0},
                        height: {value: 50},
                        width: {value: 2},
                        fill: {value: "white"}
                      },
                      update: {
                        x: {signal: "brush_x[1]"}
                      }
                    }
                },
                {
                    type: "symbol",
                    interactive: true,
                    clip: true,
                    name: "right_scrub",
                    encode: {
                      enter: {
                        yc: {value: 25},
                        fill: "transparent",
                        size: {value: 100}
                      },
                      update: {
                        xc: {signal: "brush_x[1] + 1"},
                        fill: {value: "#b6b6b6"},
                        stroke: {value: "#999"}
                      }
                    }
                },
                {
                    type: "symbol",
                    interactive: true,
                    clip: true,
                    name: "left_scrub",
                    encode: {
                      enter: {
                        yc: {value: 25},
                        fill: "transparent",
                        size: {value: 100}
                      },
                      update: {
                        xc: {signal: "brush_x[0] + 1"},
                        fill: {value: "#b6b6b6"},
                        stroke: {value: "#999"}
                      }
                    }
                },
                {
                    type: "symbol",
                    interactive: false,
                    clip: true,
                    name: "de_circle",
                    from:
                    {
                        data: "data_events"
                    },
                    encode: {
                      enter: {
                        yc: {value: 0},
                        fill: {value: "white"},
                        stroke: {value: "#999"},
                        size: {value: 200},
                      },
                      update: {
                        "x": {"scale": "x", "field": "start"}
                      }
                    }
                },
                {
                    type: "symbol",
                    interactive: false,
                    clip: true,
                    name: "de_flag",
                    from:
                    {
                        data: "data_events"
                    },
                    encode: {
                      enter: {
                        yc: {value: 0},
                        fill: {value: "#52b5e4"},
                        size: {value: 100},
                        shape: {value: "triangle-right"},
                      },
                      update: {
                        "x": {"scale": "x", "field": "start"}
                      }
                    }
                },
            ],
            scales:
            [
                {
                    name: "x",
                    type: "time",
                    domain: timeRangeAll,
                    range:
                    [
                        0,
                        {
                            signal: "width"
                        }
                    ]
                },
                {
                    name: "y",
                    type: "linear",
                    domain:
                    {
                        fields:
                        [
                            {
                                "data": "data_1",
                                "field": "value_start"
                            },
                            {
                                "data": "data_1",
                                "field": "value_end"
                            },
                            {
                                "data": "data_2",
                                "field": "value_start"
                            },
                            {
                                "data": "data_2",
                                "field": "value_end"
                            }
                        ]
                    },
                    range:
                    [
                        {
                            signal: "height"
                        },
                        0
                    ],
                    nice: true,
                    zero: true
                }
            ],
            axes:
            [
                {
                    scale: "x",
                    orient: "bottom",
                    grid: false,
                    title: "Time",
                    formatType: "time",
                    labelPadding: -14,
                    tickCount: 8,
                    titlePadding: 5,
                    labelFlush: true,
                    labelOverlap: true,
                    zindex: 0
                },
                {
                    scale: "y",
                    orient: "left",
                    grid: false,
                    title: "value",
                    labelOpacity: 0,
                    titleOpacity: 0,
                    labelOverlap: true,
                    tickCount:
                    {
                        signal: "ceil(height/40)"
                    },
                    zindex: 0
                }
            ],
        };
    }
}
