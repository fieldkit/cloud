import Vue from "vue";
import Router from "vue-router";
import Login from "./views/LoginView.vue";
import DataView from "./views/DataView.vue";
import ProjectsView from "./views/ProjectsView.vue";
import StationsView from "./views/StationsView.vue";

Vue.use(Router);

export default new Router({
    mode: "history",
    base: process.env.BASE_URL,
    routes: [
        {
            path: "/",
            alias: "/login",
            name: "login",
            component: Login
        },
        {
            path: "/dashboard",
            name: "projects",
            component: ProjectsView
        },
        {
            path: "/dashboard/projects/:id",
            name: "projectById",
            component: ProjectsView,
            props: true
        },
        {
            path: "/dashboard/stations",
            name: "stations",
            component: StationsView,
            props: true
        },
        {
            path: "/dashboard/data",
            name: "data",
            component: DataView,
            props: true
        },
        {
            path: "/dashboard/data/:id",
            name: "dataById",
            component: DataView,
            props: true
        }
    ]
});
