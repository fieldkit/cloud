// From https://matthiashager.com/converting-snake-case-to-camel-case-object-keys-with-javascript
const isArray = function(a) {
    return Array.isArray(a);
};
const isObject = function(o) {
    return o === Object(o) && !isArray(o) && typeof o !== "function";
};
const toCamel = s => {
    return s.replace(/([-_][a-z])/gi, $1 => {
        return $1
            .toUpperCase()
            .replace("-", "")
            .replace("_", "");
    });
};

export function keysToCamel(o) {
    if (isObject(o)) {
        const n = new Proxy(o, {
            get(target, name, receiver) {
                if (typeof name === "string" && !/^_/.test(name)) {
                    const camelName = toCamel(name);
                    if (camelName !== name) {
                        if (false) {
                            const err = new Error();
                            console.warn("style violation", name, err.stack);
                        } else {
                            console.warn("style violation", name);
                        }
                    }
                }
                return target[name];
            },
        });

        return n;
    } else if (isArray(o)) {
        return o.map(i => keysToCamel(i));
    }

    return o;
}
