// From https://matthiashager.com/converting-snake-case-to-camel-case-object-keys-with-javascript
const isArray = function(a) {
    return Array.isArray(a);
};
const isObject = function(o) {
    return o === Object(o) && !isArray(o) && typeof o !== "function";
};
const toCamel = (s) => {
    return s.replace(/([-_][a-z])/gi, ($1) => {
        return $1
            .toUpperCase()
            .replace("-", "")
            .replace("_", "");
    });
};

export function keysToCamelWithWarnings(o) {
    if (isObject(o)) {
        const n = new Proxy(o, {
            get(target, name) {
                if (typeof name === "string" && !/^_/.test(name)) {
                    const camelName = toCamel(name);
                    if (camelName !== name) {
                        console.warn("style violation", name);
                    }
                }
                return keysToCamelWithWarnings(target[name]);
            },
        });

        return n;
    } else if (isArray(o)) {
        return o.map((i) => keysToCamelWithWarnings(i));
    }

    return o;
}

export function keysToCamel(o) {
    if (isObject(o)) {
        const n = {};

        Object.keys(o).forEach((k) => {
            n[toCamel(k)] = keysToCamel(o[k]);
        });

        return n;
    } else if (isArray(o)) {
        return o.map((i) => keysToCamel(i));
    }

    return o;
}
