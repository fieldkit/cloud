import Vue from "vue";
import Vuex from "vuex";
import VCalendar from "v-calendar";
import moment from "moment";
import { sync } from "vuex-router-sync";
import i18n from "./i18n";
// import * as ActionTypes from "./store/actions";
import * as MutationTypes from "./store/mutations";
import storeFactory from "./store";
import routerFactory from "./router";
import ConfigurationPlugin from "./config";
import Config from "./secrets";
import App from "./App.vue";

const AssetsPlugin = {
    install(Vue) {
        const loader = require.context("@/assets/", true, /\.png$/);
        Vue.prototype.$loadAsset = path => {
            return loader("./" + path);
        };
    },
};

Vue.use(ConfigurationPlugin, Config);
Vue.use(AssetsPlugin, Config);
Vue.use(VCalendar, {});
Vue.use(Vuex);

Vue.config.productionTip = false;

Vue.filter("prettyDate", value => {
    if (!value) {
        return "N/A";
    }
    return moment(value).format("MM/DD/YYYY");
});

Vue.filter("prettyReading", sensor => {
    if (!sensor.reading) {
        return "--";
    }
    return sensor.reading.toFixed(1);
});

Vue.filter("prettyCoordinate", value => {
    if (!value) {
        return "--";
    }
    return value.toFixed(3);
});

const store = storeFactory();
store.commit(MutationTypes.INITIALIZE);

const router = routerFactory(store);
sync(store, router);

new Vue({
    i18n,
    router,
    store,
    render: h => h(App),
}).$mount("#app");
