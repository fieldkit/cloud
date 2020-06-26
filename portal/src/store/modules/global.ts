import { StationsState } from "./stations";
import { ClockState } from "./clock";
import { MapState } from "./map";
import { NavigationState } from "./nav";

export interface GlobalState {
    readonly stations: StationsState;
    readonly clock: ClockState;
    readonly map: MapState;
    readonly nav: NavigationState;
}
