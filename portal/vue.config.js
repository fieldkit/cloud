// vue.config.js

const PUBLIC_PATH = process.env.PUBLIC_PATH || "/";

module.exports = {
    publicPath: PUBLIC_PATH,
    outputDir: "build",
    pluginOptions: {
        i18n: {
            locale: "en",
            fallbackLocale: "en",
            localeDir: "locales",
            enableInSFC: true,
        },
    },
    devServer: {
        hot: false,
    },
};
