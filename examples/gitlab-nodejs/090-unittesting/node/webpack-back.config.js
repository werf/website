/* eslint no-unused-vars: "warn" */

const path = require('path');
const nodeExternals = require('webpack-node-externals');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = {
  target: "node",
  entry: {
    app: ["./src/server/server.js"]
  },
  node: {
    global: false,
    __filename: false,
    __dirname: false,
  },
  output: {
    path: path.join(__dirname, "dist"),
    publicPath: "/",
    filename: "[name].js",
  },
  externals: [nodeExternals()],
  plugins: [
    new CopyPlugin({
      patterns: [
        { from: 'src/views', to: 'views' },
      ],
    }),
  ],
};