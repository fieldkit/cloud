import { shallowMount } from "@vue/test-utils";
import DataView from "@/views/DataView.vue";

describe("DataView.vue", () => {
    it("Renders a station name", async () => {
        const wrapper = shallowMount(DataView, {
            stubs: ["router-link"],
        });
        const name = "FieldKit 1";
        wrapper.setData({ station: { name: name } });
        await wrapper.vm.$nextTick();
        expect(wrapper.find("#station-name").text()).toBe(name);
        wrapper.destroy(); // prevent warnings
    });
});
