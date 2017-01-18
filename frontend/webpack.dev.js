var webpack = require('webpack')
var path = require('path')
var HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  cache: true,
  context: path.join(__dirname, ''),
  devtool: 'eval',
  entry: ['./src/js/index.js'],
  resolve: {
    alias: {
      'webworkify': 'webworkify-webpack',
      'sinon': 'sinon/pkg/sinon',
      'gl-matrix': path.resolve('./node_modules/gl-matrix/dist/gl-matrix.js'),
      'mapbox-gl/js/geo/transform': path.join(__dirname, "/node_modules/mapbox-gl/js/geo/transform"),
      'mapbox-gl': path.join(__dirname, "/node_modules/mapbox-gl/dist/mapbox-gl.js")
    }
  },
  module: {
    noParse: [
      /node_modules\/sinon\//
    ],
    loaders: [
      { 
        test: /\.css$/, 
        loader: 'style-loader!css-loader' 
      },
      {
        test: /\.scss$/,
        loaders: ["style-loader", "css-loader", "sass-loader"]
      },
      {
        test: /\.(jpe?g|png|gif|svg)$/i,
        loaders: [
          'file?hash=sha512&digest=hex&name=[hash].[ext]',
          'image-webpack?bypassOnDebug&optimizationLevel=7&interlaced=false'
        ]
      },
      {
        test: /\.hbs$/,
        loader: 'handlebars-loader'
      },
      {
        test: /\.json$/,
        loader: 'json-loader'
      },
      {
        test: /\.jsx?$/,
        exclude: /(node_modules|bower_components)/,
        loader: 'babel-loader',
        include: [
            path.join(__dirname, 'src', 'js')
        ],
        query: {
          cacheDirectory: true,
          presets: ['react', 'es2015', 'stage-0'],
          plugins: ['react-html-attrs', 'transform-class-properties', 'transform-decorators-legacy']
        }
      },
      {
        test: /\.(eot|svg|ttf|otf|woff|woff2)$/,
        loader: 'file?name=fonts/[name].[ext]',
        include: [
            path.join(__dirname, 'src', 'fonts')
        ]
      },
      {
        test: /\.svg(\?v=\d+\.\d+\.\d+)?$/, 
        loader: 'file-loader?mimetype=image/svg+xml',
        include: [
            path.join(__dirname, 'src', 'fonts')
        ]
      },
      {test: /\.woff(\?v=\d+\.\d+\.\d+)?$/, loader: "file-loader?mimetype=application/font-woff"},
      {test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/, loader: "file-loader?mimetype=application/font-woff"},
      {test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/, loader: "file-loader?mimetype=application/octet-stream"},
      {test: /\.eot(\?v=\d+\.\d+\.\d+)?$/, loader: "file-loader"}
    ]
  },
  devServer: {
    host: '0.0.0.0',
    port: 8081,
    proxy: {
      '/api/v1/*': 'http://localhost:8080'
    },
    historyApiFallback: true
  },
  output: {
    path: path.join(__dirname, 'dist'),
    filename: '[name].[hash].js'
  },
  plugins: [
    new webpack.DllReferencePlugin({
      context: path.join(__dirname),
      manifest: require('./src/vendors/dll/vendor-manifest.json')
    }),
    new HtmlWebpackPlugin({
      title: 'FIELDKIT',
      filename: 'index.html',
      template: 'src/templates/index.hbs',
      hash: true,
      inject: false
    })
  ]
}
