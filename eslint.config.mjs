import { defineConfig } from "eslint/config";
import core from "ultracite/eslint/core";
import react from "ultracite/eslint/react";
import tanstack from "ultracite/eslint/tanstack";

export default defineConfig([
  {
    ignores: ["**/*.json"],
  },
  {
    extends: [core, react, tanstack],
  },
  {
    files: ["**/vite.config.ts"],
    rules: {
      "sonarjs/file-header": "off",
    },
  },
  {
    files: ["eslint.config.mjs"],
    rules: {
      "import-x/no-rename-default": "off",
      "sonarjs/file-header": "off",
    },
  },
  {
    files: ["apps/server/scripts/**/*.mjs"],
    rules: {
      "import-x/no-nodejs-modules": "off",
      "n/no-process-exit": "off",
      "n/no-sync": "off",
      "no-console": "off",
      "no-magic-numbers": "off",
      "sonarjs/file-header": "off",
      "sonarjs/no-os-command-from-path": "off",
      "sort-keys": "off",
      "unicorn/import-style": "off",
      "unicorn/no-process-exit": "off",
      "unicorn/no-unreadable-array-destructuring": "off",
      "unicorn/prevent-abbreviations": "off",
    },
  },
]);
