import Vue from "vue";
import { Services } from "@/api";
import * as ActionTypes from "@/store/actions";
import * as MutationTypes from "@/store/mutations";
import { PortalStationFieldNotes } from "@/views/fieldNotes/model";

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
            { commit, dispatch, state }: { commit: any; dispatch: any; state: FieldNotesState },
            payload: { id: number }
        ) => {
            const notes = await services.api.getStationFieldNotes(payload.id);
            commit(MutationTypes.FIELD_NOTES_UPDATE, notes.notes);
        },
        [ActionTypes.ADD_FIELD_NOTE]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: FieldNotesState },
            payload: { stationId: number; note: PortalStationFieldNotes }
        ) => {
            const response = await services.api.addStationFieldNote(payload.stationId, payload.note);
            const combined = [response].concat(state.fieldNotes);
            console.log("radoi add field notes response", response);
            commit(MutationTypes.FIELD_NOTES_UPDATE, combined);
        },
        [ActionTypes.UPDATE_FIELD_NOTE]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: FieldNotesState },
            payload: { stationId: number; note: PortalStationFieldNotes }
        ) => {
            await services.api.updateStationFieldNote(payload.stationId, payload.note);
        },
        [ActionTypes.DELETE_FIELD_NOTE]: async (
            { commit, dispatch, state }: { commit: any; dispatch: any; state: FieldNotesState },
            payload: { stationId: number; noteId: number }
        ) => {
            const response = await services.api.deleteStationFieldNote(payload.stationId, payload.noteId);
            const newNotes = state.fieldNotes.filter(note => note.id !== payload.noteId);
            console.log("radoi delete field notes response", response);
            commit(MutationTypes.FIELD_NOTES_UPDATE, newNotes);
        },
    };
};

const mutations = {
    [MutationTypes.FIELD_NOTES_UPDATE]: (
        state: FieldNotesState,
        payload: PortalStationFieldNotes,
    ) => {
     //   const combined = state.fieldNotes.push(payload);
        console.log("radoi mutating", payload);
        Vue.set(state, "fieldNotes", payload);
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
