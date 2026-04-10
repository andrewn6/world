import { defineConfig } from "astro/config";
import tailwind from "@astrojs/tailwind";

export default defineConfig({
  site: "https://nijmeh.cloud",
  integrations: [tailwind()],
});
