import Vuex from "vuex";
import { shallowMount, createLocalVue } from "@vue/test-utils";
import StationsView from "@/views/StationsView.vue";
import storeFactory from "@/store";

jest.mock("@/api/api.js");

const localVue = createLocalVue();

localVue.use(Vuex);

describe("StationsView.vue", () => {
    let store;

    beforeEach(() => {
        store = storeFactory();
    });

    it("Renders a map", () => {
        const stationsView = shallowMount(StationsView, { store, localVue });
        const map = stationsView.find("#summary-and-map");
        expect(map.isVisible()).toBe(true);
        stationsView.destroy(); // prevent warnings
    });
});
