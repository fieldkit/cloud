// webpackage.config.js

"use strict";

const { DefinePlugin } = require("webpack");
const { VueLoaderPlugin } = require("vue-loader");

const PUBLIC_PATH = process.env.PUBLIC_PATH || "/";

module.exports = {
    mode: "development",
    entry: ["./src/app.js"],
    output: {
        publicPath: PUBLIC_PATH,
    },
    module: {
        loaders: [{ test: /\.css$/, loader: "css-loader" }],
        rules: [
            {
                test: /\.vue$/,
                use: "vue-loader",
            },
            // this will apply to both plain `.scss` files
            // AND `<style lang="scss">` blocks in `.vue` files
            {
                test: /\.scss$/,
                use: ["vue-style-loader", "css-loader", "sass-loader"],
            },
        ],
    },
    plugins: [
        new VueLoaderPlugin(),
        new DefinePlugin({
            "process.env.PUBLIC_PATH": JSON.stringify(PUBLIC_PATH),
        }),
        new webpack.ProvidePlugin({
            mapboxgl: "mapbox-gl",
        }),
    ],
};
