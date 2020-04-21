import { createLocalVue, mount } from "@vue/test-utils";
import DataChartControl from "@/components/DataChartControl.vue";
import VueRouter from "vue-router";
import stationSummaryFixture from "./fixtures/stationSummary";
import sensorsFixture from "./fixtures/sensors";
import stationFixture from "./fixtures/station";

jest.mock("@/api/api.js");

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

const treeSelectOptions = [
    {
        id: "FieldKit1ph",
        label: "pH",
        isDefaultExpanded: true,
    },
];
const range = [new Date("1/2/20"), new Date("2/13/20 23:59:59")];
const props = {
    treeSelectOptions: treeSelectOptions,
    allSensors: sensorsFixture,
};
const initData = {
    "123": {
        data: stationSummaryFixture,
        sensors: sensorsFixture,
        timeRange: range,
        station: stationFixture,
    },
};
const DAY = 1000 * 60 * 60 * 24;

beforeEach(() => {
    jest.resetAllMocks();
});

function getNumberOfCharts(params) {
    let numCharts = 0;
    params.forEach(p => {
        // each chart has one station
        if (p.indexOf("station") > -1) {
            numCharts += 1;
        }
    });
    return numCharts;
}

describe("DataChartControl.vue", () => {
    it("Renders chart controls if a station exists", async () => {
        const dataChartControl = mount(DataChartControl, {
            mocks: {
                $route,
            },
            propsData: props,
        });
        await dataChartControl.vm.$nextTick();
        expect(dataChartControl.find("#chart-controls").isVisible()).toBe(true);
        dataChartControl.destroy();
    });

    it("Requests one day of data when user selects day", async () => {
        // new start is one day before the end date:
        const newStart = new Date(range[1].getTime() - DAY);
        const dataChartControl = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: props,
        });
        dataChartControl.vm.initializeChart(initData, "123", "chart-1");
        await dataChartControl.vm.$nextTick();
        dataChartControl.find("[data-time='1']").trigger("click");
        await dataChartControl.vm.$nextTick();
        // this emitted event triggers an api query for the specified range
        expect(dataChartControl.emitted().timeChanged[0][0].start).toEqual(newStart);
        dataChartControl.destroy();
    });

    it("Updates the URL when user selects day", async () => {
        // new start is one day before the end date:
        const newStart = new Date(range[1].getTime() - DAY);
        const router = new VueRouter({ routes, mode: "abstract" });
        const dataChartControl = mount(DataChartControl, {
            localVue,
            router,
            propsData: props,
        });
        dataChartControl.vm.initializeChart(initData, "123", "chart-1");
        await dataChartControl.vm.$nextTick();
        dataChartControl.find("[data-time='1']").trigger("click");
        await dataChartControl.vm.$nextTick();
        expect(dataChartControl.vm.$route.query["chart-1start"]).toBe(newStart.getTime());
        expect(dataChartControl.vm.$route.query["chart-1end"]).toBe(range[1].getTime());
        dataChartControl.destroy();
    });

    it("Updates the URL when user selects sensor", async () => {
        const router = new VueRouter({ routes, mode: "abstract" });
        const dataChartControl = mount(DataChartControl, {
            localVue,
            router,
            propsData: props,
        });
        expect(dataChartControl.vm.$route.query["chart-1sensor"]).toBeUndefined();
        dataChartControl.vm.initializeChart(initData, "123", "chart-1");
        await dataChartControl.vm.$nextTick();
        // const arrow = dataChartControl.find(".vue-treeselect__control-arrow-container");
        // arrow.trigger("click");
        // await dataChartControl.vm.$nextTick();
        // console.log(dataChartControl.find(".vue-treeselect__option"));
        // // can't figure out how to programmatically select an option, so just calling method
        dataChartControl.vm.onTreeSelect(sensorsFixture[0], "chart-1");
        expect(dataChartControl.vm.$route.query["chart-1sensor"]).toBe(sensorsFixture[0].key);
        dataChartControl.destroy();
    });

    it("Adds a chart when the user clicks Compare", async () => {
        const dataChartControl = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: props,
        });
        dataChartControl.vm.initializeChart(initData, "123", "chart-1");
        await dataChartControl.vm.$nextTick();
        expect(dataChartControl.findAll(".d3Chart").length).toEqual(1);
        dataChartControl.find(".compare-btn").trigger("click");
        await dataChartControl.vm.$nextTick();
        expect(dataChartControl.findAll(".d3Chart").length).toEqual(2);
        dataChartControl.destroy();
    });

    it("Updates the URL when user adds a chart", async () => {
        const router = new VueRouter({ routes, mode: "abstract" });
        const dataChartControl = mount(DataChartControl, {
            localVue,
            router,
            propsData: props,
        });
        let numCharts = getNumberOfCharts(Object.keys(dataChartControl.vm.$route.query));
        expect(numCharts).toBe(0);
        dataChartControl.vm.initializeChart(initData, "123", "chart-1");
        await dataChartControl.vm.$nextTick();
        dataChartControl.find(".compare-btn").trigger("click");
        await dataChartControl.vm.$nextTick();
        numCharts = getNumberOfCharts(Object.keys(dataChartControl.vm.$route.query));
        expect(numCharts).toBe(2);
        // expect(dataChartControl.vm.$route.query.numCharts).toEqual(2);
        dataChartControl.destroy();
    });

    it("Removes a chart when the user clicks remove button", async () => {
        const dataChartControl = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: props,
        });
        dataChartControl.vm.initializeChart(initData, "123", "chart-1");
        await dataChartControl.vm.$nextTick();
        dataChartControl.find(".compare-btn").trigger("click");
        await dataChartControl.vm.$nextTick();
        expect(dataChartControl.findAll(".d3Chart").length).toEqual(2);
        dataChartControl.find(".remove-chart > img").trigger("click");
        await dataChartControl.vm.$nextTick();
        expect(dataChartControl.findAll(".d3Chart").length).toEqual(1);
        dataChartControl.destroy();
    });

    it("Updates the URL when user removes a chart", async () => {
        const router = new VueRouter({ routes, mode: "abstract" });
        const dataChartControl = mount(DataChartControl, {
            localVue,
            router,
            propsData: props,
        });
        let numCharts = getNumberOfCharts(Object.keys(dataChartControl.vm.$route.query));
        expect(numCharts).toBe(0);
        dataChartControl.vm.initializeChart(initData, "123", "chart-1");
        await dataChartControl.vm.$nextTick();
        dataChartControl.find(".compare-btn").trigger("click");
        await dataChartControl.vm.$nextTick();
        numCharts = getNumberOfCharts(Object.keys(dataChartControl.vm.$route.query));
        expect(numCharts).toBe(2);
        dataChartControl.find(".remove-chart > img").trigger("click");
        await dataChartControl.vm.$nextTick();
        numCharts = getNumberOfCharts(Object.keys(dataChartControl.vm.$route.query));
        expect(numCharts).toBe(1);
        dataChartControl.destroy();
    });

    it("Starts with the correct number of charts based on the URL", async () => {
        const $route = {
            query: { "chart-1station": 154, "chart-2station": 177 },
        };
        const dataChartControl = mount(DataChartControl, {
            mocks: {
                $route,
                $router,
            },
            propsData: props,
        });
        dataChartControl.vm.initialize();
        await dataChartControl.vm.$nextTick();
        // just checking that it requests data for the correct number of charts
        expect(dataChartControl.emitted()["stationAdded"].length).toBe(2);
    });
});
