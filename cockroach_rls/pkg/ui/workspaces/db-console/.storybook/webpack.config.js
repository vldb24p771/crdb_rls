// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

const custom = require("../webpack.config");
const path = require("path");

const appConfig = custom({dist: "ccl"}, {mode: "development"});

module.exports = async ({ config, mode }) => {
  return {
    ...config,
    resolve: {
      ...config.resolve,
      modules: [
        path.resolve(__dirname, "..", "ccl"),
        path.resolve(__dirname, ".."),
        "node_modules",
      ],
      extensions: appConfig.resolve.extensions,
      alias: appConfig.resolve.alias,
    },
    module: {
      rules: [
        ...appConfig.module.rules,
      ]
    },
    plugins: [
      ...config.plugins,
      // Import 'vendors' library only
      appConfig.plugins[2]
    ]
  };
};
