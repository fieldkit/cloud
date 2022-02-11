const _ = require("lodash");
const express = require("express");
const vega = require("vega");
const vegaLite = require("vega-lite");
const os = require("os");
const axios = require("axios");
const app = express();

const chartConfig = require("./vega/chartConfig.json");
const lineSpec = require("./vega/line.v1.json");
const histogramSpec = require("./vega/histogram.v1.json");
const rangeSpec = require("./vega/range.v1.json");
const doubleLineSpec = require("./vega/doubleLine.v1.json");

const port = Number(process.env.FIELDKIT_PORT || 8081);
const baseUrl = process.env.FIELDKIT_BASE_URL || `http://127.0.0.1:8080`;

vega.expressionFunction("fkHumanReadable", (datum) => {
  if (_.isUndefined(datum)) {
    return "N/A";
  }
  if (this.valueSuffix) {
    return `${datum.toFixed(3)} ${this.valueSuffix}`;
  }
  return `${datum.toFixed(3)}`;
});

class Chart {
  sensor(index, label, units) {
    throw new Error("charting: NOT IMPLEMENTED");
  }

  thresholds(index, data) {
    throw new Error("charting: NOT IMPLEMENTED");
  }

  data(index, data) {
    throw new Error("charting: NOT IMPLEMENTED");
  }
}

class TimeSeriesChart extends Chart {
  constructor(viz) {
    super();
    this.double = viz[0].length > 1;
    if (this.double) {
      this.spec = _.cloneDeep(doubleLineSpec);
    } else {
      this.spec = _.cloneDeep(lineSpec);
    }
  }

  sensor(index, label, units) {
    if (this.double) {
      this.spec.layer[index].encoding.y.title = label;
    } else {
      if (index == 0) {
        this.spec.layer[0].encoding.y.axis.title = label;
      }
    }
  }

  thresholds(index, vizConfig) {
    const EnglishLocale = "en-US"; // TODO This has a dash, portal strips this.
    const thresholdLayers = vizConfig.levels
      .map((d, i) => {
        return {
          transform: [
            {
              calculate: "datum.value <= " + d.value + " ? datum.value : null",
              as: "layerValue" + i,
            },
            {
              calculate:
                "datum.layerValue" +
                i +
                " <= " +
                d.value +
                " ? '" +
                d.label[EnglishLocale] +
                "' : null",
              as: vizConfig.label[EnglishLocale],
            },
          ],
          encoding: {
            y: { field: "layerValue" + i },
            stroke: {
              field: vizConfig.label[EnglishLocale],
              legend: {
                orient: "top",
              },
              scale: {
                domain: vizConfig.levels.map((d) => d.label[EnglishLocale]),
                range: vizConfig.levels.map((d) => d.color),
              },
            },
          },
          mark: {
            type: "line",
            interpolate: "monotone",
            tension: 1,
          },
        };
      })
      .reverse();

    this.spec.layer[0].layer = thresholdLayers;
  }

  data(index, data) {
    if (this.double) {
      this.spec.layer[index].data = {
        name: `table_${index}`,
        values: data.data,
      };
    } else {
      this.spec.data = { name: `table`, values: data.data };
    }
  }
}

class RangeChart extends Chart {
  constructor(viz) {
    super();
    this.spec = _.cloneDeep(rangeSpec);
  }

  sensor(index, label, units) {
    this.spec.encoding.y.axis.title = label;
  }

  thresholds(index, data) {
    // IGNORED
  }

  data(index, data) {
    this.spec.data = { name: "table", values: data.data };
  }
}

class HistogramChart extends Chart {
  constructor(viz) {
    super();
    this.spec = _.cloneDeep(histogramSpec);
  }

  sensor(index, label, units) {
    this.spec.encoding.x.axis.title = label;
  }

  thresholds(index, data) {
    // IGNORED
  }

  data(index, data) {
    this.spec.data = { name: "table", values: data.data };
  }
}

const charts = [TimeSeriesChart, HistogramChart, RangeChart, TimeSeriesChart];

const statusHandler = (req, res) => {
  res.send({
    server_name: process.env.FIELDKIT_SERVER_NAME,
    version: process.env.FIELDKIT_VERSION,
    name: process.env.HOSTNAME,
    tag: process.env.TAG,
    git: { hash: process.env.GIT_HASH },
  });
};

app.get("/", statusHandler);

app.get("/charting", statusHandler);

app.get("/charting/rendered", async (req, res, next) => {
  // TODO Authorization header

  res.setHeader("Content-Type", "image/png");

  try {
    console.log(`charting:query`);

    if (!req.query.bookmark) {
      res.status(400).send("bad request");
      return;
    }

    const bookmark = JSON.parse(req.query.bookmark);
    const w = req.query.w || 800;
    const h = req.query.h || 418;

    console.log(`charting:bookmark`, JSON.stringify(bookmark));

    const specs = [];

    const allQueries = _.flattenDeep(
      bookmark.g.map((g1) => {
        return g1.map((g2) => {
          return g2.map((viz) => {
            console.log(`charting:viz`, JSON.stringify(viz));

            const chartTypeBookmark = viz[3];

            const chartCtor = charts[chartTypeBookmark];
            if (!chartCtor) throw new Error("charting: Unknown chart type");

            const chart = new chartCtor(viz);

            const spec = chart.spec;
            spec.config = chartConfig;
            spec.width = w;
            spec.height = h;
            specs.push(spec);

            return viz[0].map((vizSensor, index) => {
              const when = viz[1];

              const stationId = vizSensor[0];
              const sensorId = vizSensor[1][1];

              const metaParams = new URLSearchParams();
              metaParams.append("stations", stationId.toString());

              const dataParams = new URLSearchParams();
              dataParams.append("start", when[0].toString());
              dataParams.append("end", when[1].toString());
              dataParams.append("stations", stationId.toString());
              dataParams.append("sensors", vizSensor[1].join(","));
              dataParams.append("resolution", "1000");
              dataParams.append("complete", "true");

              return [
                {
                  key: metaParams.toString(),
                  url: `${baseUrl}/sensors`,
                  handle: (key, data) => {
                    const sensorKeysById = _(data.sensors)
                      .groupBy((r) => r.id)
                      .value();

                    if (!sensorKeysById[String(sensorId)]) {
                      console.log(
                        `charting: sensors: ${JSON.stringify(
                          _.keys(sensorKeysById)
                        )}`
                      );
                      throw new Error(`charting: Missing sensor: ${sensorId}`);
                    }

                    const sensorKey = sensorKeysById[String(sensorId)][0].key;
                    const sensors = _(data.modules)
                      .map((m) => m.sensors)
                      .flatten()
                      .groupBy((s) => s.full_key)
                      .value();

                    const byKey = sensors[sensorKey];
                    if (byKey.length == 0) {
                      throw new Error(
                        `charting: Missing sensor meta: ${sensorKey}`
                      );
                    }
                    const sensor = byKey[0];

                    const name = sensor.strings["en-us"]["label"] || "Unknown";

                    if (sensor.viz.length > 0) {
                      // TODO Check viz Name, no need now until other chart types support this.
                      chart.thresholds(index, sensor.viz[0].thresholds);
                    }

                    console.log(
                      `charting:handle-meta(${key}) sensor-id=${sensorId} sensor-key=${sensorKey} uom='${sensor.unit_of_measure}'`
                    );

                    const make = () => {
                      if (sensor.unit_of_measure) {
                        return `${name} (${sensor.unit_of_measure})`;
                      }
                      return `${name}`;
                    };

                    chart.sensor(index, make(), sensor.unit_of_measure);
                  },
                },
                /*
                {
                  key: metaParams.toString(),
                  url: `${baseUrl}/sensors/data?${metaParams.toString()}`,
                  handle: (key, data) => {
                    const moduleId = vizSensor[1][0];
                    const meta = _.first(
                      data.stations[stationId].filter(
                        (row) => row.moduleId == moduleId
                      )
                    );

                    console.log(
                      `charting:handle-meta(${key}) ${moduleId}`,
                      meta
                    );
                  },
                },
                */
                {
                  key: dataParams.toString(),
                  url: `${baseUrl}/sensors/data?${dataParams.toString()}`,
                  handle: (key, data) => {
                    console.log(`charting:handle-data(${key})`);
                    chart.data(index, data);
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
    console.log(`charting:data-queries`, uniqueQueries.length);

    const responses = await Promise.all(
      uniqueQueries.map((axiosQuery) =>
        axios(axiosQuery).then((response) => {
          const handlers = handlersByKey[axiosQuery.key].map((q) => q.handle);

          return {
            key: axiosQuery.key,
            data: response.data,
            handlers: handlers,
            handled: handlers.map((h) => h(axiosQuery.key, response.data)),
          };
        })
      )
    );

    console.log(`charting:data-queries-done`, responses);

    const vegaSpec = vegaLite.compile(specs[0]);
    const parsedSpec = vega.parse(vegaSpec.spec);
    const view = new vega.View(parsedSpec, {
      logger: vega.logger(vega.Debug, "error"),
      renderer: "none",
    }).finalize();

    const canvas = await view.toCanvas();
    const stream = canvas.createPNGStream();

    stream.pipe(res);
  } catch (error) {
    console.log(`charting:error`, error.message);
    next(error);
  }
});

app.listen(port, () => {
  console.log(`charting: listening port=${port} base=${baseUrl}`);
});
