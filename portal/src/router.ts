import Vue from "vue";
import Router from "vue-router";
import VueBodyClass from "vue-body-class";

import LoginView from "./views/auth/LoginView.vue";
import DiscourseView from "./views/auth/DiscourseView.vue";
import ResumeView from "./views/auth/ResumeView.vue";
import CreateAccountView from "./views/auth/CreateAccountView.vue";
import RecoverAccountView from "./views/auth/RecoverAccountView.vue";
import ResetPasswordView from "./views/auth/ResetPasswordView.vue";
import UserView from "./views/auth/UserView.vue";

import ProjectsView from "./views/projects/ProjectsView.vue";
import ProjectEditView from "./views/projects/ProjectEditView.vue";
import ProjectUpdateEditView from "./views/projects/ProjectUpdateEditView.vue";
import ProjectView from "./views/projects/ProjectView.vue";

import StationsView from "./views/StationsView.vue";
import ExploreView from "./views/viz/ExploreView.vue";

import NotesView from "./views/notes/NotesView.vue";

import NotificationsView from "./views/notifications/NotificationsView.vue";

import AdminMain from "./views/admin/AdminMain.vue";
import AdminUsers from "./views/admin/AdminUsers.vue";
import AdminStations from "./views/admin/AdminStations.vue";
import Playground from "./views/admin/Playground.vue";

import { Bookmark, deserializeBookmark } from "./views/viz/viz";
import TermsView from "@/views/auth/TermsView.vue";
import { ActionTypes } from "@/store";

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
        props: true,
    },
    {
        path: "/login/discourse",
        name: "discourseLogin",
        component: DiscourseView,
        meta: {
            bodyClass: "blue-background",
            secured: false,
        },
    },
    {
        path: "/login/keycloak",
        name: "loginKeycloak",
        component: ResumeView,
        meta: {
            secured: false,
        },
    },
    {
        path: "/login/:token",
        name: "loginResume",
        component: ResumeView,
        meta: {
            secured: false,
        },
    },
    {
        path: "/spoof",
        name: "spoof",
        component: LoginView,
        props: (route) => {
            return {
                spoofing: true,
            };
        },
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
        path: "/terms",
        name: "terms",
        component: TermsView,
        meta: {
            secured: false,
        },
    },
    {
        path: "/projects/invitation",
        name: "viewInvites",
        component: ProjectsView,
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
        path: "/dashboard/projects/:id/activity",
        name: "viewProjectActivity",
        component: ProjectView,
        props: (route) => {
            return {
                id: Number(route.params.id),
                forcePublic: false,
                activityVisible: true,
            };
        },
        meta: {
            bodyClass: "disable-scrolling",
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
        props: (route) => {
            return {
                id: Number(route.params.id),
            };
        },
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
        props: (route) => {
            return {
                id: Number(route.params.id),
            };
        },
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/stations",
        name: "mapAllStations",
        component: StationsView,
        props: true,
        meta: {
            bodyClass: "map-view",
            secured: true,
        },
    },
    {
        path: "/dashboard/stations/:id(\\d+)",
        name: "mapStation",
        component: StationsView,
        props: (route) => {
            return {
                id: Number(route.params.id),
            };
        },
        meta: {
            bodyClass: "map-view",
            secured: true,
        },
    },
    {
        path: "/dashboard/stations/all/:bounds",
        name: "mapAllStationsBounds",
        component: StationsView,
        props: (route) => {
            return {
                bounds: JSON.parse(route.params.bounds),
            };
        },
        meta: {
            bodyClass: "map-view",
            secured: true,
        },
    },
    {
        path: "/dashboard/stations/:bounds/:id(\\d+)",
        name: "mapStationBounds",
        component: StationsView,
        props: (route) => {
            return {
                id: Number(route.params.id),
                bounds: JSON.parse(route.params.bounds),
            };
        },
        meta: {
            bodyClass: "map-view",
            secured: true,
        },
    },
    {
        path: "/dashboard/explore",
        name: "exploreBookmark",
        component: ExploreView,
        props: (route) => {
            if (route.query.bookmark) {
                return {
                    bookmark: deserializeBookmark(route.query.bookmark),
                };
            }
            return {};
        },
        meta: {},
    },
    {
        path: "/viz",
        name: "exploreShortBookmark",
        component: ExploreView,
        props: (route) => {
            if (route.query.v) {
                return {
                    token: route.query.v,
                };
            }
            return {};
        },
        meta: {},
    },
    {
        path: "/dashboard/explore/:bookmark",
        name: "exportBookmarkOld",
        component: ExploreView,
        props: (route) => {
            return {
                bookmark: deserializeBookmark(route.params.bookmark),
            };
        },
        meta: {},
    },
    {
        path: "/viz/export",
        name: "exportWorkspace",
        component: ExploreView,
        props: (route) => {
            return {
                token: route.query.v,
                exportsVisible: true,
            };
        },
        meta: {
            bodyClass: "disable-scrolling",
            secured: true,
        },
    },
    {
        path: "/viz/share",
        name: "shareWorkspace",
        component: ExploreView,
        props: (route) => {
            return {
                token: route.query.v,
                shareVisible: true,
            };
        },
        meta: {
            bodyClass: "disable-scrolling",
            secured: true,
        },
    },
    {
        path: "/dashboard/projects/:projectId/notes",
        name: "viewProjectNotes",
        component: NotesView,
        props: (route) => {
            return {
                projectId: Number(route.params.projectId),
            };
        },
        meta: {
            secured: true,
        },
    },
    {
        path: "/dashboard/projects/:projectId/notes/:stationId",
        name: "viewProjectStationNotes",
        component: NotesView,
        props: (route) => {
            return {
                projectId: Number(route.params.projectId),
                stationId: Number(route.params.stationId),
            };
        },
        meta: {
            secured: true,
        },
    },
    {
        path: "/notes",
        name: "viewMyNotes",
        component: NotesView,
        props: (route) => {
            return {};
        },
        meta: {
            secured: true,
        },
    },
    {
        path: "/notes/:stationId",
        name: "viewStationNotes",
        component: NotesView,
        props: (route) => {
            return {
                stationId: Number(route.params.stationId),
            };
        },
        meta: {
            secured: true,
        },
    },
    {
        path: "/",
        name: "root",
        component: ProjectsView,
        meta: {
            secured: true,
        },
    },
    {
        path: "/notifications",
        name: "notifications",
        component: NotificationsView,
        meta: {
            secured: true,
        },
    },
    {
        path: "/admin/playground",
        name: "adminPlayground",
        component: Playground,
        props: (route) => {
            return {};
        },
        meta: {
            admin: true,
            secured: true,
        },
    },
    {
        path: "/admin",
        name: "adminMain",
        component: AdminMain,
        props: (route) => {
            return {};
        },
        meta: {
            admin: true,
            secured: true,
        },
    },
    {
        path: "/admin/users",
        name: "adminUsers",
        component: AdminUsers,
        props: (route) => {
            return {};
        },
        meta: {
            admin: true,
            secured: true,
        },
    },
    {
        path: "/admin/stations",
        name: "adminStations",
        component: AdminStations,
        props: (route) => {
            return {};
        },
        meta: {
            admin: true,
            secured: true,
        },
    },
];

export default function routerFactory(store) {
    const router = new Router({
        mode: "history",
        base: process.env.BASE_URL,
        routes: routes,
        scrollBehavior(to, from, savedPosition) {
            if (to.name == from.name) {
                return null;
            }
            console.log("scrolling-to-top");
            return { x: 0, y: 0 };
        },
    });

    const vueBodyClass = new VueBodyClass(routes);
    router.beforeEach((to, from, next) => {
        vueBodyClass.guard(to, next);
    });

    router.beforeEach(async (to, from, next) => {
        console.log("nav", from.name, "->", to.name);
        if (from.name === null && (to.name === null || to.name == "login")) {
            console.log("nav", "authenticated", store.getters.isAuthenticated);
            if (store.getters.isAuthenticated) {
                if (!store.getters.isTncValid && to.name != "login") {
                    await store.dispatch(ActionTypes.REFRESH_CURRENT_USER);

                    if (!store.getters.isTncValid) {
                        next("/terms");
                    } else {
                        next("/dashboard");
                    }
                } else {
                    next("/dashboard");
                }
            } else {
                /*
                if (to.name === "login") {
                    next();
                } else {
                    next("/login");
                }
                */
                next();
            }
        } else if (to.matched.some((record) => record.meta.secured)) {
            /*
            if (store.getters.isAuthenticated) {
                if (!store.getters.isTncValid && to.name != "login") {
                    await store.dispatch(ActionTypes.REFRESH_CURRENT_USER);

                    if (!store.getters.isTncValid) {
                        next("/terms");
                    } else {
                        next();
                    }
                } else {
                    next();
                }
            } else {
                const queryParams = new URLSearchParams();
                queryParams.append("after", to.fullPath);
                next("/login?" + queryParams.toString());
            }
            */
            next();
        } else {
            if (to.name === null) {
                if (store.getters.isAuthenticated) {
                    next("/dashboard");
                } else {
                    next("/login");
                }
            } else {
                next();
            }
        }
    });

    return router;
}
