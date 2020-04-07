import { shallowMount } from "@vue/test-utils";
import StationsView from "@/views/StationsView.vue";

describe("StationsView.vue", () => {
    it("Renders a map", () => {
        const wrapper = shallowMount(StationsView);
        const map = wrapper.find(".stations-map");
        expect(map.isVisible()).toBe(true);
    });
});
