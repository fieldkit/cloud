module.exports = {
    root: true,
    env: {
        node: true,
    },
    extends: ["plugin:vue/essential", "@vue/prettier"],
    rules: {
        "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
        "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
        "prettier/prettier": [
            "warn",
            {
                "#": "prettier config in here :)",
                printWidth: 140,
                tabWidth: 4,
                trailingComma: "es5",
                semi: true,
                singleQuote: false,
                htmlWhitespaceSensitivity: "ignore",
                endOfLine: "lf",
            },
        ],
    },
    parserOptions: {
        parser: "babel-eslint",
    },
};
