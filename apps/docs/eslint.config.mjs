import baseConfig from "@formbricks/eslint-config/base.js";
import reactConfig from "@formbricks/eslint-config/react-internal.js";

/** @type {import('eslint').Linter.Config[]} */
export default [
  ...baseConfig,
  ...reactConfig,
  {
    ignores: [
      "build/**",
      ".docusaurus/**",
      "docs/api/**", // Generated OpenAPI docs
      "node_modules/**",
    ],
  },
];
