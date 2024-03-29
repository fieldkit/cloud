import _ from "lodash";
import Vue from "vue";
import Vuex from "vuex";
import Vuelidate from "vuelidate";
import VueMeta from "vue-meta";
import VCalendar from "v-calendar";
import VueConfirmDialog from "vue-confirm-dialog";
import Multiselect from "vue-multiselect";
import { BadgePlugin } from "bootstrap-vue";
import VueSilentbox from 'vue-silentbox';

import prettyBytes from "pretty-bytes";
import moment from "moment";
import { sync } from "vuex-router-sync";
import i18n from "./i18n";
import * as MutationTypes from "./store/mutations";
import { Services } from "@/api";
import storeFactory from "./store";
import routerFactory from "./router";
import ConfigurationPlugin from "./config";
import Config from "./secrets";
import App from "./App.vue";
import { format as d3format } from "d3-format";

const services = new Services();

const AssetsPlugin = {
    install(Vue) {
        const loader = require.context("@/assets/", true, /\.(?:svg|png)$/);
        Vue.prototype.$loadAsset = (path: string) => {
            return loader("./" + path);
        };
    },
};

Object.defineProperty(Vue.prototype, "$getters", {
    get: function(this: Vue) {
        return this.$store.getters;
    },
});

Object.defineProperty(Vue.prototype, "$services", {
    get: function(this: Vue) {
        return services;
    },
});

Object.defineProperty(Vue.prototype, "$state", {
    get: function(this: Vue) {
        return this.$store.state;
    },
});

Vue.use(ConfigurationPlugin, Config);
Vue.use(AssetsPlugin, Config);
Vue.use(VueMeta, {
    refreshOnceOnNavigation: true,
});
Vue.use(VCalendar, {});
Vue.use(Vuex);
Vue.use(Vuelidate);
Vue.use(VueConfirmDialog);
Vue.component("vue-confirm-dialog", VueConfirmDialog.default);
Vue.component("multiselect", Multiselect);
Vue.use(BadgePlugin);
Vue.use(VueSilentbox);

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

Vue.filter("prettyDateTime", (value) => {
    if (!value) {
        return "N/A";
    }
    return moment(value).format("M/D/YYYY, HH:mm");
});

Vue.filter("prettyPercentage", (value: number | null) => {
    if (value === null || value === undefined) {
        return "--";
    }
    return value.toFixed(1) + "%";
});

Vue.filter("prettyReadingNarrowSpace", (value: number) => {
    if (!isNaN(value)) {
        if (value % 1 === 0) {
            return value;
        } else {
            if (value >= 10) {
                return d3format(".0f")(value);
            } else if (value >= 1) {
                return d3format(".1f")(value);
            } else {
                // This was using .2s, until we noticed for things like "0.510" we
                // were displaying 510m, which is very confusing.
                return d3format(".2f")(value);
            }
        }
    }
    return value;
});

Vue.filter("prettyNum", (value: number) => {
    if (!isNaN(value)) {
        if (value % 1 === 0) {
            return value;
        } else {
            // This was using .2s, until we noticed for things like "0.510" we
            // were displaying 510m, which is very confusing.
            return d3format(".2f")(value);
        }
    }
    return value;
});

export interface ReadingLike {
    reading: number;
}

Vue.filter("prettyReading", (sensor: ReadingLike) => {
    if (!_.isNumber(sensor.reading)) {
        return "--";
    }
    return sensor.reading.toFixed(1);
});

Vue.filter("prettyCoordinate", (value) => {
    if (!value) {
        return "--";
    }
    return value.toFixed(3) + "°";
});

Vue.filter("prettyDuration", (value) => {
    return _.capitalize(moment.duration(value / 1000, "seconds").humanize());
});

Vue.filter("prettyBytes", (value) => {
    if (!value) {
        return "N/A";
    }
    return prettyBytes(value);
});

const store = storeFactory(services);
store.commit(MutationTypes.INITIALIZE);

const router = routerFactory(store);
sync(store, router);

new Vue({
    i18n,
    router,
    store,
    render: (h) => h(App),
}).$mount("#app");
