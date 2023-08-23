import Vue from "vue";
import {Services} from "@/api";
import {PortalNoteMedia, PortalStationNotes} from "@/views/notes/model";
import * as ActionTypes from "@/store/actions";
import * as MutationTypes from "@/store/mutations";

export class NotesState {
    media: PortalNoteMedia[];
    notes: PortalStationNotes[];
    readOnly: boolean;
}

const getters = {
    notes(state: NotesState): PortalStationNotes[] {
        return state.notes;
    },
    media(state: NotesState): PortalNoteMedia[] {
        return state.media;
    },
};

const actions = (services: Services) => {
    return {
        [ActionTypes.NEED_NOTES]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: NotesState },
            payload: { id: number }
        ) => {
            commit(MutationTypes.NOTES_CLEAR);
            const notes = await services.api.getStationNotes(payload.id);
            commit(MutationTypes.NOTES_UPDATE, notes);
        },
    };
};

const mutations = {
    [MutationTypes.NOTES_UPDATE]: (
        state: NotesState,
        payload: {
            media: PortalNoteMedia[];
            notes: PortalStationNotes[];
            station: { readOnly: boolean };
        }
    ) => {
        Vue.set(state, "notes", payload.notes);
        Vue.set(state, "media", payload.media);
        Vue.set(state, "readOnly", payload.station.readOnly);
    },

    [MutationTypes.NOTES_CLEAR]: (
        state: NotesState,
    ) => {
        Vue.set(state, "notes", null);
        Vue.set(state, "media", null);
        Vue.set(state, "readOnly", true);
    },
};

export const notes = (services: Services) => {
    const state = () => new NotesState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
