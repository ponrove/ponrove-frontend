import prettierPluginSvelte from "prettier-plugin-svelte";
import prettierPluginTailwindcss from "prettier-plugin-tailwindcss";

export default {
  useTabs: true,
  singleQuote: true,
  trailingComma: "none",
  printWidth: 100,
  pluginSearchDirs: false,
  plugins: [prettierPluginSvelte, prettierPluginTailwindcss],
  overrides: [
    {
      files: "*.svelte",
      options: {
        parser: "svelte",
      },
    },
  ],
};
