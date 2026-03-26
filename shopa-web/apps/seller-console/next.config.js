/** @type {import('next').NextConfig} */
const config = {
  output: "export",
  reactStrictMode: true,
  transpilePackages: ["@shopa/api-client", "@shopa/types", "@shopa/ui"]
};

module.exports = config;
