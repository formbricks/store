import js from '@eslint/js';
import eslintConfigPrettier from 'eslint-config-prettier';
import turboPlugin from 'eslint-plugin-turbo';
import tseslint from 'typescript-eslint';
import importPlugin from 'eslint-plugin-import';
import promisePlugin from 'eslint-plugin-promise';
import checkFilePlugin from 'eslint-plugin-check-file';
import prettierPlugin from 'eslint-plugin-prettier';

/**
 * A shared ESLint configuration for the repository.
 *
 * @type {import("eslint").Linter.Config[]}
 * */
export default [
  js.configs.recommended,
  eslintConfigPrettier,
  ...tseslint.configs.recommended,
  // Avoid using non-flat configs; enable promise rules via flat preset if available
  // and manually merge other plugin rule sets
  // Some plugins don't yet ship flat presets; we add their rules explicitly below
  promisePlugin.configs['flat/recommended'],
  {
    plugins: {
      turbo: turboPlugin,
      import: importPlugin,
      promise: promisePlugin,
      // sonarjs: sonarjsPlugin, // Temporarily disabled until version alignment
      'check-file': checkFilePlugin,
      prettier: prettierPlugin,
    },
    settings: {
      // Enable TypeScript-aware import resolution across the monorepo
      'import/resolver': {
        typescript: {},
        node: {},
      },
    },
    rules: {
      'turbo/no-undeclared-env-vars': 'warn',
      // Import hygiene
      'import/no-duplicates': 'error',
      'import/newline-after-import': 'warn',
      // SonarJS rules can be enabled after aligning versions across ESLint and TS-ESLint
      // File naming convention: enforce kebab-case for TS/TSX
      'check-file/filename-naming-convention': [
        'error',
        { '**/*.{ts,tsx}': 'KEBAB_CASE' },
        { ignoreMiddleExtensions: true },
      ],
      // Promise best practices beyond the recommended set can be enabled here if desired
      // "promise/always-return": "warn",
      // Prettier as an ESLint rule (surfaces formatting issues per .prettierrc)
      'prettier/prettier': 'error',
    },
  },
  {
    ignores: ['dist/**'],
  },
];
