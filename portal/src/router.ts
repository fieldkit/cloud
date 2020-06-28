import Vue from "vue";
import Router from "vue-router";
import LoginView from "./views/auth/LoginView.vue";
import DataView from "./views/viz/DataView.vue";
import InvitesView from "./views/InvitesView.vue";
import ProjectsView from "./views/ProjectsView.vue";
import ProjectEditView from "./views/ProjectEditView.vue";
import ProjectUpdateEditView from "./views/ProjectUpdateEditView.vue";
import ProjectView from "./views/ProjectView.vue";
import ResetPasswordView from "./views/auth/ResetPasswordView.vue";
import StationsView from "./views/StationsView.vue";
import UserView from "./views/UserView.vue";

Vue.use(Router);

const routes = [
    {
        path: "/login",
        name: "login",
        component: LoginView,
        meta: {
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
        path: "/dashboard/user/reset",
        name: "reset",
        component: ResetPasswordView,
        meta: {
            secured: false,
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
        props: route => {
            return {
                id: Number(route.params.id),
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
            secured: true,
        },
    },
    {
        path: "/dashboard/stations/:id",
        name: "viewStation",
        component: StationsView,
        props: route => {
            return {
                id: Number(route.params.id),
            };
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
];

export default function routerFactory(store) {
    const router = new Router({
        mode: "history",
        base: process.env.BASE_URL,
        routes: routes,
    });

    router.beforeEach((to, from, chain) => {
        console.log("nav", from.name, "->", to.name);
        if (from.name === null && (to.name === null || to.name == "login")) {
            console.log("nav", "authenticated", store.getters.isAuthenticated);
            if (store.getters.isAuthenticated) {
                chain("/dashboard");
            } else {
                if (to.name === "login") {
                    chain();
                } else {
                    chain("/login");
                }
            }
        } else if (to.matched.some(record => record.meta.secured)) {
            if (store.getters.isAuthenticated) {
                chain();
                return;
            }
            chain("/login");
        } else {
            chain();
        }
    });

    return router;
}
