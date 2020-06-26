import Vue from "vue";
import Vuex from "vuex";
import VCalendar from "v-calendar";
import App from "./App.vue";
import router from "./router";
import i18n from "./i18n";
import storeFactory from "./store";
import * as MutationTypes from "./store/mutations";

Vue.use(VCalendar, {});
Vue.use(Vuex);

Vue.config.productionTip = false;

const store = storeFactory();

store.commit(MutationTypes.INITIALIZE);

new Vue({
    i18n,
    router,
    store,
    render: h => h(App),
}).$mount("#app");
