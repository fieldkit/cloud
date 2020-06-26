import Vue from "vue";

export class Configuration {
    constructor(public readonly baseUrl: string) {}
}

const plugin = {
    install(Vue, config: any) {
        Vue.prototype.$config = new Configuration(config.API_HOST);
    },
};

export default plugin;
