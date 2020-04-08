import { createLocalVue, mount } from "@vue/test-utils";
import DataChartControl from "@/components/DataChartControl.vue";
import VueRouter from "vue-router";

// some tests use Vue Router with localVue
const localVue = createLocalVue();
localVue.use(VueRouter);
// and some tests use custom mocks for $route and $router
const $route = {
    query: {},
};
const $router = {
    push: jest.fn(),
};

const DAY = 1000 * 60 * 60 * 24;

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

    it("Updates the url start param when user selects day", async () => {
        const range = [new Date("1/2/20"), new Date("2/13/20")];
        // new start is one day before the end date:
        const newStart = new Date(range[1].getTime() - DAY);
        const routes = [
            {
                path: "/dashboard/data/154",
                name: "viewData",
                component: DataChartControl,
                props: true,
            },
        ];
        const router = new VueRouter({
            routes,
        });
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
    });
});
