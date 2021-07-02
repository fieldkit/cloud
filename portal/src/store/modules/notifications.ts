import { ActionContext, Module } from "vuex";
import { Services, SendFunction } from "@/api";
import * as ActionTypes from "../actions";
import { MarkNotificationsSeen } from "../typed-actions";
import { promiseAfter } from "@/utilities";

export interface Notification {
    id: number;
    userId: number;
    postId: number;
    kind: string;
    body: string;
}

const CONNECTED = "CONNECTED";
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
        [ActionTypes.INITIALIZE]: async ({ dispatch, commit }: ActionParameters) => {
            const send = await services.api.listen(async (message: unknown) => {
                commit(NOTIFIED, message as Notification);
                return await Promise.resolve();
            });

            commit(CONNECTED, send);
        },
        [ActionTypes.NOTIFICATIONS_SEEN]: async ({ dispatch, commit, state }: ActionParameters, payload: MarkNotificationsSeen) => {
            if (state.send) {
                await state.send(payload);
            } else {
                // TODO Queue?
            }
        },
    };
};

const mutations = {
    [CONNECTED]: (state: NotificationsState, send: SendFunction) => {
        state.send = send;
    },
    [NOTIFIED]: (state: NotificationsState, payload: Notification) => {
        state.notifications.push(payload);
    },
    [SEEN]: (state: NotificationsState, payload: Notification) => {
        // state.notifications.push(payload);
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
