// From https://matthiashager.com/converting-snake-case-to-camel-case-object-keys-with-javascript
// eslint-disable-next-line
const isObject = function(o: any): boolean {
    return o === Object(o) && !isArray(o) && typeof o !== "function";
};
// eslint-disable-next-line
const isArray = function(a: any): boolean {
    return Array.isArray(a);
};
// eslint-disable-next-line
const toCamel = (s: string): string => {
    return s.replace(/([-_][a-z])/gi, ($1) => {
        return $1
            .toUpperCase()
            .replace("-", "")
            .replace("_", "");
    });
};

// eslint-disable-next-line
export function keysToCamel(o: any): any {
    if (isObject(o)) {
        const n = {};

        Object.keys(o).forEach((k) => {
            // eslint-disable-next-line
            n[toCamel(k)] = keysToCamel(o[k]);
        });

        return n;
    } else if (isArray(o)) {
        // eslint-disable-next-line
        return o.map((i: any) => {
            return keysToCamel(i); // eslint-disable-line
        });
    }

    return o; // eslint-disable-line
}
