import { shallowMount } from "@vue/test-utils";
import DataView from "@/views/viz/DataView.vue";

jest.mock("@/api/api.js");

const $route = {
    fullPath: {},
};

describe("DataView.vue", () => {
    it("Renders a project name", async () => {
        const dataView = shallowMount(DataView, {
            mocks: {
                $route,
            },
            stubs: ["router-link"],
        });
        const name = "Default FieldKit Project";
        dataView.setData({ projects: [{ name: name }] });
        await dataView.vm.$nextTick();
        expect(dataView.find("#project-name").text()).toBe(name);
        dataView.destroy(); // prevent warnings
    });
});
