import { ActionContext } from "vuex";
import { Services, SendFunction, SimpleUser } from "@/api";
import * as ActionTypes from "../actions";
import { promiseAfter } from "@/utilities";
import Vue from "vue";

export interface Annotation {
    type: string;
    start: Date;
    end: Date;
    body: string;
    bookmark?: string;
    user?: { id: number; name: string; photo?: object };
}

export class AnnotationsState {
    constructor(public readonly annotations: Annotation[]) {
        return {
            // fixme: placeholder data
            annotations: [...Array(10).keys()].map( ann => { 
                return {
                    type: "event",
                    start: new Date(),
                    end: new Date(),
                    body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
                    bookmark: "%7B%22v%22:1,%22g%22:%5B%5B%5B%5B%5B129%5D,%5B%5B%22c2Vuc29yXzI=%22,48%5D%5D,%5B1629591476728,1629612598264%5D,%5B%5D,0,-1%5D%5D%5D,%5B%5B%5B%5B129%5D,%5B%5B%22c2Vuc29yXzI=%22,44%5D%5D,%5B1629591476728,1629612598264%5D,%5B%5D,0,-1%5D,%5B%5B127%5D,%5B%5B%22c21pdGgtYW5kLTl0aA==%22,44%5D%5D,%5B1629591476728,1629612598264%5D,%5B%5D,0,-1%5D,%5B%5B128%5D,%5B%5B%22Y2Fycm9sbC1hbmQtNHRo%22,44%5D%5D,%5B1629591476728,1629612598264%5D,%5B%5D,0,-1%5D%5D%5D%5D,%22s%22:%5B%5D%7D",
                    user: { id: 1234, name: "Test Person" },
                }
            })
        };
    }
}

type ActionParameters = ActionContext<AnnotationsState, never>;

const getters = {
    annotations(state: AnnotationsState): Annotation[] {
        return state.annotations;
    },
    getEvents(state: AnnotationsState): Annotation[] {
        return state.annotations.filter( d => d.type === "event");
    }
};

const actions = (services: Services) => {
    return {};
};

const mutations = {};

export const annotations = (services: Services) => {
    const state = () => new AnnotationsState([]);

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
