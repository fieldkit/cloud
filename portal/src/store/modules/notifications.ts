import { ActionContext } from "vuex";
import { Services, SendFunction } from "@/api";
import * as ActionTypes from "../actions";
import { MarkNotificationsSeen } from "../typed-actions";
import { promiseAfter } from "@/utilities";

export interface Notification {
    notificationId: number;
    userId: number;
    postId: number;
    kind: string;
    body: string;
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
    numberOfUnseenNotifications(state: NotificationsState): number {
        return state.notifications.length;
    },
};

const actions = (services: Services) => {
    return {
        [ActionTypes.INITIALIZE]: async ({ commit }: ActionParameters) => {
            const send = await services.api.listenForNotifications(
                async (message: { error?: unknown }) => {
                    if (message.error) {
                        console.log("error", message.error);
                    } else {
                        commit(NOTIFIED, message as Notification);
                    }
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
            if (state.send) {
                if (payload.ids.length == 0) {
                    await state.send(new MarkNotificationsSeen(state.notifications.map((n) => n.notificationId)));
                } else {
                    await state.send(payload);
                }
                commit(SEEN);
            } else {
                // TODO Queue?
            }
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
        state.notifications.push(payload);
    },
    [SEEN]: (state: NotificationsState) => {
        state.notifications.length = 0;
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
