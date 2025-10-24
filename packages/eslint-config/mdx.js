import pluginMdx from 'eslint-plugin-mdx';
import baseConfig from './base.js';

/**
 * ESLint configuration for MDX files with embedded JSX/TSX.
 * @type {import("eslint").Linter.Config[]}
 */
export default [
  ...baseConfig,
  {
    files: ['**/*.mdx'],
    plugins: { mdx: pluginMdx },
    processor: pluginMdx.createProcessor(),
    settings: {
      'mdx/code-blocks': true,
    },
    rules: {
      // MDX recommended brings in sensible defaults
      ...pluginMdx.configs.recommended.rules,
    },
  },
];
