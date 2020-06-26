import Vue from "vue";
import Router from "vue-router";
import DataView from "./views/DataView.vue";
import InvitesView from "./views/InvitesView.vue";
import LoginView from "./views/LoginView.vue";
import ProjectsView from "./views/ProjectsView.vue";
import ProjectEditView from "./views/ProjectEditView.vue";
import ProjectUpdateEditView from "./views/ProjectUpdateEditView.vue";
import ProjectView from "./views/ProjectView.vue";
import ResetPasswordView from "./views/ResetPasswordView.vue";
import StationsView from "./views/StationsView.vue";
import UserView from "./views/UserView.vue";

Vue.use(Router);

const routes = [
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
        component: ProjectView,
        props: true,
    },
    {
        path: "/dashboard/project/add",
        name: "addProject",
        component: ProjectEditView,
    },
    {
        path: "/dashboard/projects/:id/edit",
        name: "editProject",
        component: ProjectEditView,
        props: true,
    },
    {
        path: "/dashboard/project-update/add",
        name: "addProjectUpdate",
        component: ProjectUpdateEditView,
        props: true,
    },
    {
        path: "/dashboard/project-updates/:id/edit",
        name: "editProjectUpdate",
        component: ProjectUpdateEditView,
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
];

export default function routerFactory() {
    const router = new Router({
        mode: "history",
        base: process.env.BASE_URL,
        routes: routes,
    });

    router.beforeEach((to, from, chain) => {
        console.log("before", to, from);
        return chain();
    });

    return router;
}
