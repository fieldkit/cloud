module.exports = {
    preset: "@vue/cli-plugin-unit-jest/presets/typescript-and-babel",
    moduleNameMapper: {
        "\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$": "<rootDir>/__mocks__/fileMock.js",
        "\\.(css|less)$": "<rootDir>/__mocks__/styleMock.js",
        "\\api": "<rootDir>/__mocks__/api.js",
    },
    transform: {
        "^.+\\.ts?$": "ts-jest",
        ".*\\.(vue)$": "@vue/vue2-jest",
    },
    testMatch: ["**/*.spec.(js|jsx|ts|tsx)|**/__tests__/*.(js|jsx|ts|tsx)"],
    moduleFileExtensions: ["ts", "js", "vue", "json"],
};
