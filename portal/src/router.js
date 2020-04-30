import Vue from "vue";
import Router from "vue-router";
import DataView from "./views/DataView.vue";
import InvitesView from "./views/InvitesView.vue";
import LoginView from "./views/LoginView.vue";
import ProjectsView from "./views/ProjectsView.vue";
import ResetPasswordView from "./views/ResetPasswordView.vue";
import StationsView from "./views/StationsView.vue";
import UserView from "./views/UserView.vue";

Vue.use(Router);

export default new Router({
    mode: "history",
    base: process.env.BASE_URL,
    routes: [
        {
            path: "/",
            alias: "/login",
            name: "login",
            component: LoginView,
        },
        {
            path: "/projects/invitation",
            name: "viewInvites",
            component: InvitesView,
            props: true,
        },
        {
            path: "/dashboard/user",
            name: "user",
            component: UserView,
        },
        {
            path: "/dashboard/user/reset",
            name: "reset",
            component: ResetPasswordView,
        },
        {
            path: "/dashboard/",
            name: "projects",
            component: ProjectsView,
        },
        {
            path: "/dashboard/projects/:id",
            name: "viewProject",
            component: ProjectsView,
            props: true,
        },
        {
            path: "/dashboard/projects/:id/edit",
            name: "editProject",
            component: ProjectsView,
            props: true,
        },
        {
            path: "/dashboard/stations",
            name: "stations",
            component: StationsView,
            props: true,
        },
        {
            path: "/dashboard/stations/:id",
            name: "viewStation",
            component: StationsView,
            props: true,
        },
        {
            path: "/dashboard/data",
            name: "viewData",
            component: DataView,
            props: true,
        },
    ],
});
