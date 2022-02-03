const _ = require("lodash");
const express = require("express");
const vega = require("vega");
const vegaLite = require("vega-lite");
const axios = require("axios");
const app = express();

const port = 8081;
const baseUrl = `http://127.0.0.1:8080`;

const chartConfig = require("./vega/chartConfig.json");
const lineSpec = require("./vega/line.v1.json");
const histogramSpec = require("./vega/histogram.v1.json");
const rangeSpec = require("./vega/range.v1.json");
const doubleLineSpec = require("./vega/doubleLine.v1.json");

const locale = require("./en.json");
const localizedSensors = _(locale.modules)
  .map((m, moduleKey) => {
    return _(m.sensors)
      .map((sensorName, sensorKey) => {
        const normalizedKey = sensorKey
          .split(".")
          .map((p) => _.camelCase(p).replace("10M", "10m").replace("2M", "2m"))
          .join(".");

        if (moduleKey.indexOf("wh.") == 0) {
          const fullKey = [moduleKey, normalizedKey].join(".");
          return [fullKey, sensorName];
        }

        const fullKey = ["fk", moduleKey, normalizedKey].join(".");
        return [fullKey, sensorName];
      })
      .value();
  })
  .flatten()
  .fromPairs()
  .value();

app.get("/", (req, res) => {
  res.send({ server_name: "", version: "", name: "", git: { hash: "" } });
});

app.get("/charts/rendered", async (req, res) => {
  // TODO Authorization header

  res.setHeader("Content-Type", "image/png");

  try {
    console.log(`charting:query`, req.query);

    if (!req.query.bookmark) throw new Error("charting: Bookmark is required");

    const spec = _.cloneDeep(doubleLineSpec);
    const bookmark = JSON.parse(req.query.bookmark);
    const w = req.query.w || 800;
    const h = req.query.h || 300;

    console.log(`charting:bookmark`, bookmark);

    const allQueries = _.flattenDeep(
      bookmark.g.map((g1) => {
        return g1.map((g2) => {
          return g2.map((viz) => {
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

                    if (!sensorKeysById[String(sensorId)])
                      throw new Error("charting: Missing sensor");

                    const sensorKey = sensorKeysById[String(sensorId)][0].key;
                    const sensors = _(data.modules)
                      .map((m) => m.sensors)
                      .flatten()
                      .groupBy((s) => s.full_key)
                      .value();

                    const byKey = sensors[sensorKey];
                    if (byKey.length == 0)
                      throw new Error("charting: Missing sensor meta");
                    const sensor = byKey[0];

                    const name = localizedSensors[sensorKey];

                    console.log(
                      `charting:handle-meta(${key}) sensor-id=${sensorId} sensor-key=${sensorKey} uom='${sensor.unit_of_measure}'`
                    );

                    const make = () => {
                      if (sensor.unit_of_measure) {
                        return `${name} (${sensor.unit_of_measure})`;
                      }
                      return `${name}`;
                    };

                    spec.layer[index].encoding.y.title = make();
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
                    spec.layer[index].data = {
                      name: `table_${index}`,
                      values: data.data,
                    };
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
    console.log(`charting:data-queries`, uniqueQueries);

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

    // spec.layer[0].data = { name: "table0", values: dataLeft.data };
    // spec.layer[0].encoding.y.title = "Left";
    // spec.layer[1].data = { name: "table1", values: dataRight.data };
    // spec.layer[1].encoding.y.title = "Right";
    spec.config = chartConfig;
    spec.width = w;
    spec.height = h;

    const vegaSpec = vegaLite.compile(spec);
    const parsedSpec = vega.parse(vegaSpec.spec);
    const view = new vega.View(parsedSpec, {
      logger: vega.logger(vega.Debug, "error"),
      renderer: "none",
    }).finalize();

    const canvas = await view.toCanvas();
    const stream = canvas.createPNGStream();

    stream.pipe(res);
  } catch (error) {
    console.log(`charting: ${error.message}`);
  }
});

app.listen(port, () => {
  console.log(`charting: listening on port ${port}`);
});
