import { DisplayStation, StationsState, DisplayProject, MappedStations } from "./stations";
import { ClockState } from "./clock";
import { MapState } from "./map";
import { UserState } from "./user";
import { LayoutState } from "./layout";
import { ExportingState } from "./exporting";
import { NotificationsState } from "./notifications";
import { NotesState } from "@/store";

export interface RouteState {
    name: string | null;
    path: string;
    fullPath: string;
    meta: object;
    params: object;
    query: object;
    hash: string;
    from: RouteState;
}

export interface GlobalState {
    readonly stations: StationsState;
    readonly clock: ClockState;
    readonly map: MapState;
    readonly route: RouteState;
    readonly user: UserState;
    readonly exporting: ExportingState;
    readonly layout: LayoutState;
    readonly notifications: NotificationsState;
    readonly notes: NotesState;
}

export interface GlobalGetters {
    projectsById: { [index: number]: DisplayProject };
    stationsById: { [index: number]: DisplayStation };
    isAuthenticated: boolean;
    isBusy: boolean;
    mapped: MappedStations;
}
