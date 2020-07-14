import Vue from "vue";
import Vuex from "vuex";
import Vuelidate from "vuelidate";
import VCalendar from "v-calendar";
import moment from "moment";
import { sync } from "vuex-router-sync";
import i18n from "./i18n";
import * as MutationTypes from "./store/mutations";
import storeFactory from "./store";
import routerFactory from "./router";
import ConfigurationPlugin from "./config";
import Config from "./secrets";
import App from "./App.vue";

const AssetsPlugin = {
    install(Vue) {
        const loader = require.context("@/assets/", true, /\.png$/);
        Vue.prototype.$loadAsset = (path) => {
            return loader("./" + path);
        };
    },
};

Object.defineProperty(Vue.prototype, "$getters", {
    get: function(this: Vue) {
        return this.$store.getters;
    },
});

Object.defineProperty(Vue.prototype, "$state", {
    get: function(this: Vue) {
        return this.$store.state;
    },
});

Vue.use(ConfigurationPlugin, Config);
Vue.use(AssetsPlugin, Config);
Vue.use(VCalendar, {});
Vue.use(Vuex);
Vue.use(Vuelidate);

Vue.config.productionTip = false;

Vue.filter("prettyTime", (value) => {
    if (!value) {
        return "N/A";
    }
    return moment(value).format("M/D/YYYY HH:mm:ss");
});

Vue.filter("prettyDate", (value) => {
    if (!value) {
        return "N/A";
    }
    return moment(value).format("M/D/YYYY");
});

Vue.filter("prettyReading", (sensor) => {
    if (!sensor.reading) {
        return "--";
    }
    return sensor.reading.toFixed(1);
});

Vue.filter("prettyCoordinate", (value) => {
    if (!value) {
        return "--";
    }
    return value.toFixed(3);
});

Vue.filter("prettyDuration", (value) => {
    return moment.duration(value / 1000, "seconds").humanize();
});

const store = storeFactory();
store.commit(MutationTypes.INITIALIZE);

const router = routerFactory(store);
sync(store, router);

new Vue({
    i18n,
    router,
    store,
    render: (h) => h(App),
}).$mount("#app");
