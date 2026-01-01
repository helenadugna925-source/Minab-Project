import { defineNuxtConfig } from "nuxt/config";

// Default to local Hasura from docker-compose if GQL_HOST is not provided
const defaultGraphqlEndpoint =
  process.env.GQL_HOST || "http://localhost:8080/v1/graphql";

export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      // Used on the client and in plugins/composables
      GQL_HOST: defaultGraphqlEndpoint,
      API_BASE: process.env.API_BASE || "http://localhost:8082",
      PLACE_API_URL: process.env.PLACE_API_URL,
      PLACE_API_KEY: process.env.PLACE_API_KEY,
    },
  },
  devtools: { enabled: true },
  css: ["~/assets/css/main.css"],
  modules: ["@nuxtjs/apollo", "@nuxtjs/leaflet"],
  apollo: {
    autoImports: true,
    authHeader: "Authorization",
    clients: {
      default: {
        // Point Apollo to your own Hasura instance
        httpEndpoint: defaultGraphqlEndpoint,
      },
    },
  },
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
});