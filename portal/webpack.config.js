// webpackage.config.js

"use strict";

const { DefinePlugin } = require("webpack");
const { VueLoaderPlugin } = require("vue-loader");

const PUBLIC_PATH = process.env.PUBLIC_PATH || "/";

module.exports = {
    mode: "development",
    entry: ["./src/app.js"],
    output: {
        publicPath: PUBLIC_PATH
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                use: "vue-loader"
            }
        ]
    },
    plugins: [
        new VueLoaderPlugin(),
        new DefinePlugin({
            "process.env.PUBLIC_PATH": JSON.stringify(PUBLIC_PATH)
        }),
        new webpack.ProvidePlugin({
            mapboxgl: 'mapbox-gl',
        })
    ]
};
