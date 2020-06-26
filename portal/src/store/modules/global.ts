import { StationsState } from "./stations";
import { ClockState } from "./clock";
import { MapState } from "./map";

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
}
