var path = require('path')
var webpack = require('webpack')

module.exports = {
  entry: {
    vendor: [path.join(__dirname, 'src', 'js', 'vendors.js')]
  },
  output: {
    path: path.join(__dirname, 'src', 'js', 'dll'),
    filename: 'dll.[name].js',
    library: '[name]'
  },
  module: {
    noParse: [
      /node_modules\/sinon\//
    ],
    loaders: [
      {
        test: /\.json$/,
        loader: 'json-loader'
      }
    ]
  },
  plugins: [
    new webpack.DllPlugin({
      path: path.join(__dirname, 'dll', '[name]-manifest.json'),
      name: '[name]',
      context: path.resolve(__dirname, 'static', 'js')
    }),
    new webpack.optimize.OccurenceOrderPlugin(),
    new webpack.optimize.UglifyJsPlugin()
  ],
  resolve: {
    root: path.resolve(__dirname, 'static', 'js'),
    modulesDirectories: ['node_modules']
  }
}

