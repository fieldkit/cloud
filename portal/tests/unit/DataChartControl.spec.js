import { createLocalVue, mount } from "@vue/test-utils";
import DataChartControl from "@/components/DataChartControl.vue";
import VueRouter from "vue-router";
import stationSummaryFixture from "./fixtures/stationSummary";
import sensorsFixture from "./fixtures/sensors";
import labelsFixture from "./fixtures/labels";

// some tests use Vue Router with localVue
const localVue = createLocalVue();
localVue.use(VueRouter);
const routes = [
    {
        path: "/dashboard/data/154",
        name: "viewData",
        component: DataChartControl,
    },
];
// and some tests use custom mocks for $route and $router
const $route = {
    query: {},
};
const $router = {
    push: jest.fn(),
};

const range = [new Date("1/2/20"), new Date("2/13/20")];
const DAY = 1000 * 60 * 60 * 24;

beforeEach(() => {
    jest.resetAllMocks();
});

describe("DataChartControl.vue", () => {
    it("Renders chart container if a station exists", async () => {
        const wrapper = mount(DataChartControl, {
            mocks: {
                $route,
            },
            propsData: {
                station: { name: "FieldKit 1" },
            },
        });
        await wrapper.vm.$nextTick();
        expect(wrapper.find("#data-chart-container").isVisible()).toBe(true);
        wrapper.destroy();
    });

    it("Requests one day of data when user selects day", async () => {
        // new start is one day before the end date:
        const newStart = new Date(range[1].getTime() - DAY);
        const wrapper = mount(DataChartControl, {
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
        wrapper.destroy();
    });

    it("Updates the URL when user selects day", async () => {
        // new start is one day before the end date:
        const newStart = new Date(range[1].getTime() - DAY);
        const router = new VueRouter({ routes, mode: "abstract" });
        const wrapper = mount(DataChartControl, {
            localVue,
            router,
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
            },
        });
        await wrapper.vm.$nextTick();
        wrapper.find("[data-time='1']").trigger("click");
        await wrapper.vm.$nextTick();
        expect(wrapper.vm.$route.query.start).toBe(newStart.getTime());
        expect(wrapper.vm.$route.query.end).toBe(range[1].getTime());
        wrapper.destroy();
    });

    it("Updates the URL when user selects sensor", async () => {
        const router = new VueRouter({ routes, mode: "abstract" });
        const wrapper = mount(DataChartControl, {
            localVue,
            router,
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
                labels: labelsFixture,
            },
        });
        expect(wrapper.vm.$route.query["chart-1sensor"]).toBeUndefined();
        wrapper.setProps({
            combinedStationInfo: { stationData: stationSummaryFixture, sensors: sensorsFixture },
        });
        await wrapper.vm.$nextTick();
        wrapper.find(".sensor-selection-dropdown > select > option").element.selected = true;
        wrapper.find(".sensor-selection-dropdown > select").trigger("change");
        await wrapper.vm.$nextTick();
        expect(wrapper.vm.$route.query["chart-1sensor"]).toBe(sensorsFixture[0].key);
        wrapper.destroy();
    });

    it("Adds a chart when the user clicks Compare", async () => {
        const wrapper = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
                labels: labelsFixture,
            },
        });
        wrapper.setProps({
            combinedStationInfo: { stationData: stationSummaryFixture, sensors: sensorsFixture },
        });
        await wrapper.vm.$nextTick();
        expect(wrapper.findAll(".d3Chart").length).toEqual(1);
        wrapper.find(".compare-btn").trigger("click");
        await wrapper.vm.$nextTick();
        expect(wrapper.findAll(".d3Chart").length).toEqual(2);
        wrapper.destroy();
    });

    it("Updates the URL when user adds a chart", async () => {
        const router = new VueRouter({ routes, mode: "abstract" });
        const wrapper = mount(DataChartControl, {
            localVue,
            router,
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
                labels: labelsFixture,
            },
        });
        expect(wrapper.vm.$route.query.numCharts).toBeUndefined();
        wrapper.setProps({
            combinedStationInfo: { stationData: stationSummaryFixture, sensors: sensorsFixture },
        });
        await wrapper.vm.$nextTick();
        wrapper.find(".compare-btn").trigger("click");
        await wrapper.vm.$nextTick();
        expect(wrapper.vm.$route.query.numCharts).toEqual(2);
        wrapper.destroy();
    });

    it("Removes a chart when the user clicks remove button", async () => {
        const wrapper = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
                labels: labelsFixture,
            },
        });
        wrapper.setProps({
            combinedStationInfo: { stationData: stationSummaryFixture, sensors: sensorsFixture },
        });
        await wrapper.vm.$nextTick();
        wrapper.find(".compare-btn").trigger("click");
        await wrapper.vm.$nextTick();
        expect(wrapper.findAll(".d3Chart").length).toEqual(2);
        wrapper.find(".remove-chart > img").trigger("click");
        await wrapper.vm.$nextTick();
        expect(wrapper.findAll(".d3Chart").length).toEqual(1);
        wrapper.destroy();
    });

    it("Updates the URL when user removes a chart", async () => {
        const router = new VueRouter({ routes, mode: "abstract" });
        const wrapper = mount(DataChartControl, {
            localVue,
            router,
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
                labels: labelsFixture,
            },
        });
        wrapper.setProps({
            combinedStationInfo: { stationData: stationSummaryFixture, sensors: sensorsFixture },
        });
        await wrapper.vm.$nextTick();
        expect(wrapper.vm.$route.query.numCharts).toBeUndefined();
        wrapper.find(".compare-btn").trigger("click");
        await wrapper.vm.$nextTick();
        expect(wrapper.vm.$route.query.numCharts).toEqual(2);
        wrapper.find(".remove-chart > img").trigger("click");
        await wrapper.vm.$nextTick();
        expect(wrapper.vm.$route.query.numCharts).toEqual(1);
        wrapper.destroy();
    });

    it("Starts with the correct number of charts based on the URL", async () => {
        const $route = {
            query: { numCharts: 2 },
        };
        const wrapper = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: {
                station: { name: "FieldKit 1" },
                totalTime: range,
                labels: labelsFixture,
            },
        });
        wrapper.setProps({
            combinedStationInfo: { stationData: stationSummaryFixture, sensors: sensorsFixture },
        });
        await wrapper.vm.$nextTick();
        expect(wrapper.findAll(".d3Chart").length).toEqual($route.query.numCharts);
    });
});
