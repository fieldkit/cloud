import Vue from "vue";

export class Configuration {
    constructor(public readonly baseUrl: string, public readonly sso = true) {}
}

const plugin = {
    install(Vue, config: any) {
        Vue.prototype.$config = new Configuration(config.API_HOST);
        Vue.prototype.$seriousError = (error) => {
            console.log("serious error", error);
        };
    },
};

export default plugin;
