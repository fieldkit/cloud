module.exports = {
    root: true,
    env: {
        node: true
    },
    extends: ["plugin:vue/essential", "@vue/prettier"],
    rules: {
        "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
        "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
        "prettier/prettier": [
            "warn",
            {
                "#": "prettier config in here :)",
                printWidth: 110,
                tabWidth: 4
            }
        ]
    },
    parserOptions: {
        parser: "babel-eslint"
    }
};
