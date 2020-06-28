import Vuex from "vuex";
import { shallowMount, createLocalVue } from "@vue/test-utils";
import storeFactory from "@/store";
import DataView from "@/views/viz/DataView.vue";
import * as MutationTypes from "@/store/mutations";

jest.mock("@/api/api.js");
jest.mock("@/api/tokens.js");

const localVue = createLocalVue();

localVue.use(Vuex);

const $route = {
    fullPath: {},
};

describe("DataView.vue", () => {
    let store;

    beforeEach(() => {
        store = storeFactory();
    });

    it("Renders a project name", async () => {
        const dataView = shallowMount(DataView, {
            store,
            localVue,
            mocks: {
                $route,
            },
            stubs: ["router-link"],
        });
        const name = "Default FieldKit Project";
        store.commit("HAVE_USER_PROJECTS", [{ name: name }]);
        await dataView.vm.$nextTick();
        expect(dataView.find("#project-name").text()).toBe(name);
        dataView.destroy(); // prevent warnings
    });
});
