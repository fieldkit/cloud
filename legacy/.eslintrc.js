module.exports = {
    "rules": {
        "indent": ["off"],
        "no-extra-parens": 0,
        "react/jsx-uses-vars": 1,
        "react/jsx-uses-react": 1,
    },
    "parserOptions": {
        "ecmaVersion": 6,
        "sourceType": "module",
        "ecmaFeatures": {
            "modules": true,
            "jsx": true,
        }
    },
    "env": {
        "node": false,
        "browser": true,
        "es6": true,
    },
    "parser": "babel-eslint",
    "plugins": [
        "eslint-plugin-react"
    ]
}
