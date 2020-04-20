import { shallowMount } from "@vue/test-utils";
import StationsView from "@/views/StationsView.vue";

jest.mock("@/api/api");

describe("StationsView.vue", () => {
    it("Renders a map", () => {
        const stationsView = shallowMount(StationsView);
        const map = stationsView.find(".stations-map");
        expect(map.isVisible()).toBe(true);
        stationsView.destroy(); // prevent warnings
    });
});
