var path = require('path')
var webpack = require('webpack')

module.exports = {
  entry: {
    vendor: [path.join(__dirname, 'src', 'js', 'vendors.js')]
  },
  resolve: {
    alias: {
      'webworkify': 'webworkify-webpack',
      'sinon': 'sinon/pkg/sinon',
      'gl-matrix': path.resolve('./node_modules/gl-matrix/dist/gl-matrix.js'),
      'mapbox-gl/js/geo/transform': path.join(__dirname, "/node_modules/mapbox-gl/js/geo/transform"),
      'mapbox-gl': path.join(__dirname, "/node_modules/mapbox-gl/dist/mapbox-gl.js")
    }
  },
  output: {
    path: path.join(__dirname, 'dist', 'dll'),
    filename: 'dll.[name].js',
    library: '[name]'
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
        test: /\.styl$/,
        loader: 'style-loader!css-loader!stylus-loader'
      },
      {
        test: /\.(eot|svg|ttf|otf|woff|woff2)$/,
        loader: 'file?name=fonts/[name].[ext]'
      }
    ]
  },
  plugins: [
    new webpack.DllPlugin({
      path: path.join(__dirname, 'dll', '[name]-manifest.json'),
      name: '[name]',
      context: path.resolve(__dirname, 'src', 'js')
    }),
    new webpack.optimize.OccurenceOrderPlugin(),
    new webpack.optimize.UglifyJsPlugin()
  ]
}

