import Vue from "vue";
import Vuex from "vuex";
import VCalendar from "v-calendar";
import { sync } from "vuex-router-sync";
import i18n from "./i18n";
import * as ActionTypes from "./store/actions";
import * as MutationTypes from "./store/mutations";
import storeFactory from "./store";
import routerFactory from "./router";
import ConfigurationPlugin from "./config";
import Config from "./secrets";
import App from "./App.vue";

Vue.use(ConfigurationPlugin, Config);
Vue.use(VCalendar, {});
Vue.use(Vuex);

Vue.config.productionTip = false;

const store = storeFactory();
const router = routerFactory(store);
const unsync = sync(store, router);

new Vue({
    i18n,
    router,
    store,
    render: h => h(App),
}).$mount("#app");
