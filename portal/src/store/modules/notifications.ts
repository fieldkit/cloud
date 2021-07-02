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

const getters = {};

const actions = (services: Services) => {
    return {
        [ActionTypes.INITIALIZE]: async ({ dispatch }: { dispatch: any }) => {
            await services.api.listen(async (message: unknown) => {
                return await Promise.resolve();
            });
        },
    };
};

const mutations = {};

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
