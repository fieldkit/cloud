import { Services } from "@/api";
import * as ActionTypes from "../actions";

export interface Notification {
    id: number;
    userId: number;
    postId: number;
    kind: string;
    body: string;
}

export class NotificationsState {
    constructor(public readonly notifications: Notification[]) {}
}

const getters = {
    numberOfUnseenNotifications(state: NotificationsState): number {
        return state.notifications.length;
    },
};

const NOTIFIED = "NOTIFIED";

const actions = (services: Services) => {
    return {
        [ActionTypes.INITIALIZE]: async ({ dispatch, commit }: { dispatch: any; commit: any }) => {
            await services.api.listen(async (message: unknown) => {
                commit(NOTIFIED, message as Notification);
                return await Promise.resolve();
            });
        },
    };
};

const mutations = {
    [NOTIFIED]: (state: NotificationsState, payload: Notification) => {
        state.notifications.push(payload);
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
