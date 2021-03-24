const path = require('path')
const webpack = require('webpack') // eslint-disable-line
const HtmlWebpackPlugin = require('html-webpack-plugin')
const VueLoaderPlugin = require('vue-loader/lib/plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin');

module.exports = {
  entry: './src/main.js',
  output: {
    path: path.resolve('./dist'),
    publicPath: '/guides/dist',
    filename: '[name].[chunkhash].js'
  },
  optimization: {
    splitChunks: {
      cacheGroups: {
        vendor: {
          chunks: 'initial',
          test: path.resolve(__dirname, 'node_modules'),
          name: 'vendor',
          enforce: true
        }
      }
    }
  },
  mode: 'development',
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
      {
        test: /\.css$/,
        use: [
          'vue-style-loader',
          'css-loader'
        ]
      },
      {
        test: /\.scss$/,
        use: [
          'vue-style-loader',
          'css-loader',
          'sass-loader'
        ]
      },
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: file => (
          /node_modules/.test(file) &&
          !/\.vue\.js/.test(file)
        )
      }
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      filename: '_includes/scripts.html',
      template: path.resolve(__dirname, './_include_templates/_scripts.html'),
      debug: process.env.DEBUG,
      send_metrics: process.env.SENDMETRICS,
    }),
    new HtmlWebpackPlugin({
      filename: '_includes/breadcrumbs.html',
      template: path.resolve(__dirname, './_include_templates/_breadcrumbs.html'),
    }),
    new HtmlWebpackPlugin({
      filename: '_includes/sidebar.html',
      template: path.resolve(__dirname, './_include_templates/_sidebar.html'),
    }),
    new HtmlWebpackPlugin({
      filename: '_includes/en/landing-tiles.html',
      template: path.resolve(__dirname, './_include_templates/en/landing-tiles.html'),
    }),
    new HtmlWebpackPlugin({
      filename: '_includes/ru/landing-tiles.html',
      template: path.resolve(__dirname, './_include_templates/ru/landing-tiles.html'),
    }),
    new VueLoaderPlugin(),
    new CleanWebpackPlugin(),
    new webpack.EnvironmentPlugin(['LANG'])
  ],
  resolve: {
    alias: {
      vue$: 'vue/dist/vue.esm.js'
    },
    extensions: ['*', '.js', '.vue', '.json']
  }
}

if (process.env.NODE_ENV === 'production') {
  module.exports.mode = 'production'
}
