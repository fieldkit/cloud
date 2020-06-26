import Vue from "vue";
import Vuex from "vuex";
import VCalendar from "v-calendar";
import { sync } from "vuex-router-sync";
import i18n from "./i18n";
import * as MutationTypes from "./store/mutations";
import storeFactory from "./store";
import routerFactory from "./router";
import App from "./App.vue";

Vue.use(VCalendar, {});
Vue.use(Vuex);

Vue.config.productionTip = false;

const router = routerFactory();
const store = storeFactory();
const unsync = sync(store, router);

store.commit(MutationTypes.INITIALIZE);

new Vue({
    i18n,
    router,
    store,
    render: h => h(App),
}).$mount("#app");
