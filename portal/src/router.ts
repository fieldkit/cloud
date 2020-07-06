import Vue from "vue";
import Router from "vue-router";
import VueBodyClass from "vue-body-class";

import LoginView from "./views/auth/LoginView.vue";
import CreateAccountView from "./views/auth/CreateAccountView.vue";
import RecoverAccountView from "./views/auth/RecoverAccountView.vue";
import ResetPasswordView from "./views/auth/ResetPasswordView.vue";

import DataView from "./views/viz/DataView.vue";
import InvitesView from "./views/InvitesView.vue";
import ProjectsView from "./views/ProjectsView.vue";
import ProjectEditView from "./views/ProjectEditView.vue";
import ProjectUpdateEditView from "./views/ProjectUpdateEditView.vue";
import ProjectView from "./views/ProjectView.vue";
import StationsView from "./views/StationsView.vue";
import UserView from "./views/UserView.vue";
import ExploreView from "./views/viz/ExploreView.vue";

Vue.use(Router);

const routes = [
    {
        path: "/login",
        name: "login",
        component: LoginView,
        meta: {
            bodyClass: "blue-background",
            secured: false,
        },
    },
    {
        path: "/register",
        name: "register",
        component: CreateAccountView,
        meta: {
            bodyClass: "blue-background",
            secured: false,
        },
    },
    {
        path: "/recover",
        name: "recover",
        component: RecoverAccountView,
        meta: {
            bodyClass: "blue-background",
            secured: false,
        },
    },
    {
        path: "/recover/complete",
        name: "reset",
        component: ResetPasswordView,
        meta: {
            bodyClass: "blue-background",
            secured: false,
        },
    },
    {
        path: "/projects/invitation",
        name: "viewInvites",
        component: InvitesView,
        props: true,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/user",
        name: "user",
        component: UserView,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/user",
        name: "editUser",
        component: UserView,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard",
        name: "projects",
        component: ProjectsView,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/projects/:id",
        name: "viewProject",
        component: ProjectView,
        props: (route) => {
            return {
                id: Number(route.params.id),
                forcePublic: false,
            };
        },
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/projects/:id/public",
        name: "viewProjectPublic",
        component: ProjectView,
        props: (route) => {
            return {
                id: Number(route.params.id),
                forcePublic: true,
            };
        },
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/project/add",
        name: "addProject",
        component: ProjectEditView,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/projects/:id/edit",
        name: "editProject",
        component: ProjectEditView,
        props: true,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/project-update/add",
        name: "addProjectUpdate",
        component: ProjectUpdateEditView,
        props: true,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/project-updates/:id/edit",
        name: "editProjectUpdate",
        component: ProjectUpdateEditView,
        props: true,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/stations",
        name: "stations",
        component: StationsView,
        props: true,
        meta: {
            bodyClass: "map-view",
            secured: true,
        },
    },
    {
        path: "/dashboard/stations/:id",
        name: "viewStation",
        component: StationsView,
        props: (route) => {
            return {
                id: Number(route.params.id),
            };
        },
        meta: {
            bodyClass: "map-view",
        },
    },
    {
        path: "/dashboard/data",
        name: "viewData",
        component: DataView,
        props: true,
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/explore",
        name: "explore",
        component: ExploreView,
        props: true,
        meta: {
            secured: true,
        },
    },
];

export default function routerFactory(store) {
    const router = new Router({
        mode: "history",
        base: process.env.BASE_URL,
        routes: routes,
    });

    const vueBodyClass = new VueBodyClass(routes);
    router.beforeEach((to, from, next) => {
        vueBodyClass.guard(to, next);
    });

    router.beforeEach((to, from, next) => {
        console.log("nav", from.name, "->", to.name);
        if (from.name === null && (to.name === null || to.name == "login")) {
            console.log("nav", "authenticated", store.getters.isAuthenticated);
            if (store.getters.isAuthenticated) {
                next("/dashboard");
            } else {
                if (to.name === "login") {
                    next();
                } else {
                    next("/login");
                }
            }
        } else if (to.matched.some((record) => record.meta.secured)) {
            if (store.getters.isAuthenticated) {
                next();
                return;
            }
            next("/login");
        } else {
            next();
        }
    });

    return router;
}
