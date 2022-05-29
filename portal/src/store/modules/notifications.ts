import { ActionContext } from "vuex";
import { Services, SendFunction, SimpleUser } from "@/api";
import * as ActionTypes from "../actions";
import { MarkNotificationsSeen } from "../typed-actions";
import { promiseAfter } from "@/utilities";
import Vue from "vue";

export interface Notification {
    notificationId: number;
    userId: number;
    postId: number;
    kind: string;
    body: string;
    projectId?: string;
    bookmark?: string;
    user?: { id: number; name: string; photo: object };
}

const CONNECTED = "CONNECTED";
const DISCONNECTED = "DISCONNECTED";
const NOTIFIED = "NOTIFIED";
const SEEN = "SEEN";

export class NotificationsState {
    constructor(public readonly notifications: Notification[], public send: SendFunction | null = null) {}
}

type ActionParameters = ActionContext<NotificationsState, never>;

const getters = {
    notifications(state: NotificationsState): Notification[] {
        return state.notifications;
    },
    numberOfUnseenNotifications(state: NotificationsState): number {
        return state.notifications.length;
    },
};

const actions = (services: Services) => {
    return {
        [ActionTypes.INITIALIZE]: async ({ commit }: ActionParameters) => {
            const send = await services.api.listenForNotifications(
                async (message: Notification) => {
                    commit(NOTIFIED, message);
                    return await Promise.resolve();
                },
                async (connected) => {
                    if (!connected) {
                        commit(DISCONNECTED);
                    }
                    await promiseAfter(5000);
                    return await Promise.resolve();
                }
            );

            commit(CONNECTED, send);
        },
        [ActionTypes.NOTIFICATIONS_SEEN]: async ({ commit, state }: ActionParameters, payload: MarkNotificationsSeen) => {
            await services.api.seenNotifications({ ids: state.notifications.map((n) => n.notificationId) });

            commit(SEEN);
        },
    };
};

const mutations = {
    [DISCONNECTED]: (state: NotificationsState) => {
        state.send = null;
    },
    [CONNECTED]: (state: NotificationsState, send: SendFunction) => {
        state.send = send;
    },
    [NOTIFIED]: (state: NotificationsState, payload: Notification) => {
        if (
            !state.notifications.find((notification) => {
                return notification.notificationId === payload.notificationId;
            })
        ) {
            state.notifications.push(payload);
        }
    },
    [SEEN]: (state: NotificationsState) => {
        Vue.set(state, "notifications", []);
    },
};

export const notifications = (services: Services) => {
    const state = () => new NotificationsState([]);

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
