import Vue from "vue";
import {Services} from "@/api";
import {PortalNoteMedia, PortalStationFieldNotes, PortalStationNotes} from "@/views/notes/model";
import * as ActionTypes from "@/store/actions";
import * as MutationTypes from "@/store/mutations";

export class FieldNotesState {
    fieldNotes: PortalStationFieldNotes[];
}

const getters = {
    fieldNotes(state: FieldNotesState): PortalStationFieldNotes[] {
        return state.fieldNotes;
    },
};

const actions = (services: Services) => {
    return {
        [ActionTypes.NEED_FIELD_NOTES]: async (
            {commit, dispatch, state}: { commit: any; dispatch: any; state: FieldNotesState },
            payload: { id: number }
        ) => {
            const notes = await services.api.getStationFieldNotes(payload.id);
            console.log("radoi field notes response", notes);
            commit(MutationTypes.FIELD_NOTES_UPDATE, notes);
        },
        [ActionTypes.ADD_FIELD_NOTE]: async (
            {commit, dispatch, state}: { commit: any; dispatch: any; state: FieldNotesState },
            payload: { stationId: number, note: PortalStationFieldNotes }
        ) => {
            const response = await services.api.addStationFieldNote(payload.stationId, payload.note);
            console.log("radoi add field notes response", response);
            commit(MutationTypes.FIELD_NOTES_UPDATE, response);
        },
    };
};

const mutations = {
    [MutationTypes.FIELD_NOTES_UPDATE]: (
        state: FieldNotesState,
        payload: {
            fieldNotes: PortalStationFieldNotes[];
        }
    ) => {
        Vue.set(state, "fieldNotes", payload.fieldNotes);
    },
};

export const fieldNotes = (services: Services) => {
    const state = () => new FieldNotesState();

    return {
        namespaced: false,
        state,
        getters,
        actions: actions(services),
        mutations,
    };
};
