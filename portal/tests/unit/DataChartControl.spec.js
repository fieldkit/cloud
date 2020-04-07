import { mount } from "@vue/test-utils";
import DataChartControl from "@/components/DataChartControl.vue";

const $route = {
    query: {},
};

const $router = {
    push: jest.fn(),
};

let wrapper = null;

afterEach(() => {
    wrapper.destroy();
});

describe("DataChartControl.vue", () => {
    it("Renders chart container if a station exists", async () => {
        wrapper = mount(DataChartControl, {
            mocks: {
                $route,
            },
            propsData: {
                station: { name: "FieldKit 1" },
            },
        });
        await wrapper.vm.$nextTick();
        expect(wrapper.find("#data-chart-container").isVisible()).toBe(true);
    });

    it("Requests one day of data when user selects day", async () => {
        const DAY = 1000 * 60 * 60 * 24;
        const range = [new Date("1/2/20"), new Date("2/13/20")];
        // new start is one day before the end date:
        const newStart = new Date(range[1].getTime() - DAY);
        wrapper = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
            },
        });
        await wrapper.vm.$nextTick();
        wrapper.find("[data-time='1']").trigger("click");
        await wrapper.vm.$nextTick();
        // this emitted event triggers an api query for the specified range
        expect(wrapper.emitted().timeChanged[0][0].start).toEqual(newStart);
    });
});
