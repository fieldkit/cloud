declare module "*.vue" {
    import Vue = require("@/store/strong-vue");
    export default Vue;
}

declare module "*.png" {
    const value: any;
    export = value;
}

declare module "*.svg" {
    const value: any;
    export = value;
}
