import { mount } from "@vue/test-utils";
import DataView from "@/views/DataView.vue";

const $route = {
    query: {},
};

describe("DataView.vue", () => {
    it("Renders a station name", async () => {
        const wrapper = mount(DataView, {
            stubs: ["router-link"],
            mocks: {
                $route,
            },
        });
        const name = "FieldKit 1";
        wrapper.setData({ station: { name: name } });
        await wrapper.vm.$nextTick();
        expect(wrapper.find("#station-name").text()).toBe(name);
    });
});
