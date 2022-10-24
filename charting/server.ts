import _ from "lodash";
import express from "express";
import * as vega from "vega";
import * as vegaLite from "vega-lite";
import axios from "axios";

const app = express();

import { keysToCamel } from "./json-tools";

import { TimeRange, QueriedData, VizInfo, SeriesData, VizSensor, DataSetSeries } from "./common";

import { SensorInfoResponse, ModuleSensorMeta } from "./api";

import { ChartSettings } from "./vega/SpecFactory";
import { TimeSeriesSpecFactory } from "./vega/TimeSeriesSpecFactory";
import { HistogramSpecFactory } from "./vega/HistogramSpecFactory";
import { RangeSpecFactory } from "./vega/RangeSpecFactory";

interface ResponseRow {
    vizSensor: VizSensor;
    sensor: ModuleSensorMeta;
    station: { name: string; location: [number, number] | null };
}

const port = Number(process.env.FIELDKIT_PORT || 8081);
const baseUrl = process.env.FIELDKIT_BASE_URL || `http://127.0.0.1:8080`;

class Chart {
    constructor(public readonly settings: ChartSettings) {}

    prepare(metaResponses, dataResponses) {
        const allSeries = metaResponses.map((row: ResponseRow, index: number) => {
            const { vizSensor, sensor, station } = row;

            console.log("chart", vizSensor, sensor, station);
            const strings = sensor.strings["enUs"];
            const name = strings["label"] || "Unknown";
            const axisLabel = strings["axisLabel"] || "";
            const data = dataResponses[index][0];
            const scale = [];
            const chartLabel = "";
            const vizInfo = new VizInfo(
                sensor.key,
                scale,
                station,
                sensor.unitOfMeasure,
                sensor.key,
                name,
                sensor.viz || [],
                sensor.ranges,
                chartLabel,
                axisLabel
            );

            return new SeriesData(data.key, data.timeRange, new DataSetSeries(vizSensor, data), data, vizInfo);
        });

        return this.finalize(allSeries);
    }

    finalize(allSeries): unknown[] {
        return [];
    }
}

class TimeSeriesChart extends Chart {
    finalize(allSeries): any[] {
        const factory = new TimeSeriesSpecFactory(allSeries, this.settings);
        const spec = factory.create();
        return [spec];
    }
}

class RangeChart extends Chart {
    finalize(allSeries): any[] {
        const factory = new RangeSpecFactory(allSeries, this.settings);
        const spec = factory.create();
        return [vegaLite.compile(spec).spec];
    }
}

class HistogramChart extends Chart {
    finalize(allSeries): any[] {
        const factory = new HistogramSpecFactory(allSeries, this.settings);
        const spec = factory.create();
        return [vegaLite.compile(spec as any).spec];
    }
}

const chartCtors = [TimeSeriesChart, HistogramChart, RangeChart, TimeSeriesChart];

const statusHandler = (req, res) => {
    res.send({
        server_name: process.env.FIELDKIT_SERVER_NAME,
        version: process.env.FIELDKIT_VERSION,
        name: process.env.HOSTNAME,
        tag: process.env.TAG,
        git: { hash: process.env.GIT_HASH },
    });
};

const queryArgument = (args): string => {
    return args;
};

type Handler = (key: string, qd: unknown) => void;

interface KeyedHandler {
    key: string;
    index: number;
    url: string;
    handle: Handler;
}

app.get("/", statusHandler);

app.get("/charting", statusHandler);

app.get("/charting/rendered", async (req, res, next) => {
    // TODO Authorization header

    try {
        console.log(`charting: query`);

        if (!req.query.bookmark) {
            res.status(400).send("bad request");
            return;
        }

        const bookmark = JSON.parse(queryArgument(req.query.bookmark));
        const size = {
            w: Number(req.query.w || 800 * 2),
            h: Number(req.query.h || 418 * 2),
        };
        const settings = new ChartSettings(TimeRange.eternity, size, size);

        console.log(`charting: bookmark`, JSON.stringify(bookmark));

        const charts: Chart[] = [];

        const vizStationIds: number[] = _.flattenDeep(
            bookmark.g.map((g1) => {
                return g1.map((g2) => {
                    return g2.map((viz) => {
                        return viz[0].map((vizSensor, index) => {
                            return vizSensor[0];
                        });
                    });
                });
            })
        );

        // Query for all station information mentioned in the Bookmark, we need
        // to do this before building other queries so that we can resolve
        // modules to the station that owns them.
        const stationIds = [...vizStationIds, ...bookmark.s];
        const queryParams = new URLSearchParams();
        queryParams.append("stations", stationIds.join(","));
        const metaResponse = await axios({ url: `${baseUrl}/meta/stations?` + queryParams.toString() });
        const meta: SensorInfoResponse = keysToCamel(metaResponse.data) as SensorInfoResponse;

        // Build the map from modules to stations.
        const modulesToStations = _(Object.values(meta.stations))
            .flatten()
            .map((metaRow) => [metaRow.moduleId, metaRow.stationId])
            .fromPairs()
            .value();

        const allQueries: KeyedHandler[] = _.flattenDeep(
            bookmark.g.map((g1) => {
                return g1.map((g2) => {
                    return g2.map((viz) => {
                        console.log(`charting: viz`, JSON.stringify(viz));

                        const chartTypeBookmark = viz[3];

                        const chartCtor = chartCtors[chartTypeBookmark];
                        if (!chartCtor) throw new Error("charting: Unknown chart type");

                        const chart = new chartCtor(settings);
                        const chartIndex = charts.length;

                        charts.push(chart);

                        return viz[0].map((vizSensor, index) => {
                            const when = viz[1];

                            // This may not be the station id we should be
                            // querying against, and is instead the one that the
                            // user found the sensor via/through.
                            const moduleId = vizSensor[1][0];
                            const sensorId = vizSensor[1][1];
                            const stationId = modulesToStations[moduleId] || vizSensor[0];

                            const dataParams = new URLSearchParams();
                            dataParams.append("start", when[0].toString());
                            dataParams.append("end", when[1].toString());
                            dataParams.append("stations", stationId.toString());
                            dataParams.append("sensors", vizSensor[1].join(","));
                            dataParams.append("resolution", "1000");
                            dataParams.append("complete", "true");

                            return [
                                {
                                    index: chartIndex,
                                    key: "all-meta",
                                    url: `${baseUrl}/sensors`,
                                    handle: (key, data) => {
                                        const sensorKeysById = _(data.sensors)
                                            .groupBy((r) => r.id)
                                            .value();

                                        if (!sensorKeysById[String(sensorId)]) {
                                            console.log(`charting: sensors: ${JSON.stringify(_.keys(sensorKeysById))}`);
                                            throw new Error(`charting: Missing sensor: ${sensorId}`);
                                        }

                                        const sensorKey = sensorKeysById[String(sensorId)][0].key;
                                        const sensors = _(data.modules)
                                            .map((m) => m.sensors)
                                            .flatten()
                                            .groupBy((s) => s.fullKey)
                                            .value();

                                        const byKey = sensors[sensorKey];
                                        if (byKey.length == 0) {
                                            throw new Error(`charting: Missing sensor meta: ${sensorKey}`);
                                        }

                                        const sensor = byKey[0];

                                        console.log(
                                            `charting: handle-meta(${key}) sensor-id=${sensorId} sensor-key=${sensorKey} uom='${sensor.unitOfMeasure}'`
                                        );

                                        const station = {
                                            name: "STATION",
                                            location: "LOCATION",
                                        };

                                        return {
                                            vizSensor: vizSensor,
                                            sensor: sensor,
                                            station: station,
                                        };
                                    },
                                },
                                {
                                    index: chartIndex,
                                    key: dataParams.toString(),
                                    url: `${baseUrl}/sensors/data?${dataParams.toString()}`,
                                    handle: (key, data) => {
                                        console.log(`charting: handle-data(${key})`);
                                        return new QueriedData(key, new TimeRange(when[0], when[1]), data);
                                    },
                                },
                            ];
                        });
                    });
                });
            })
        );

        const handlersByKey = _(allQueries)
            .groupBy((q) => q.key)
            .value();

        const uniqueQueries = _.uniqBy(allQueries, (q) => q.key);
        console.log(`charting: data-queries`, uniqueQueries.length);

        const responses = await Promise.all(
            uniqueQueries.map((axiosQuery) =>
                axios({ url: axiosQuery.url }).then((response) => {
                    const handlers = handlersByKey[axiosQuery.key].map((q) => q.handle);
                    const params = handlers.map((h) => h(axiosQuery.key, keysToCamel(response.data)));

                    return {
                        key: axiosQuery.key,
                        index: axiosQuery.index,
                        data: response.data,
                        handlers: handlers,
                        handled: params,
                    };
                })
            )
        );

        const byChartIndex = _(responses)
            .groupBy((r) => r.index)
            .mapValues((r) => r.map((c) => c.handled))
            .mapValues((r, k) => {
                const chart = charts[k];
                return chart.prepare(r[0], r.slice(1));
            })
            .value();

        const specs = _.flatten(_.values(byChartIndex));
        if (specs.length == 0) throw new Error(`viz: No charts`);

        const parsedSpec = vega.parse(specs[0]);
        const view = new vega.View(parsedSpec, {
            logger: vega.logger(vega.Debug, "error"),
            renderer: "none",
        }).finalize();

        const canvas = await view.toCanvas();

        (canvas as any).toBuffer((err, buffer) => {
            if (err) {
                console.log("charting: error", err);
                return;
            }

            console.log("charting: buffer", buffer.length);

            res.setHeader("Content-Type", "image/png");

            res.end(buffer);

            console.log("charting: done");
        });
    } catch (error) {
        console.log(`charting: error`, error.message);
        next(error);
    }
});

app.listen(port, () => {
    console.log(`charting: listening port=${port} base=${baseUrl}`);
});
