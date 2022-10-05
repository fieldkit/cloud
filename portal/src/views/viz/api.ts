export type StationID = number;
export type ModuleID = string;
export type SensorID = number;
export type Stations = StationID[];
export type SensorSpec = [ModuleID, SensorID];
export type Sensors = SensorSpec[];
export type VizSensor = [StationID, SensorSpec];

export interface SensorRange {
    minimum: number;
    maximum: number;
    constrained: boolean | null;
}

export interface VizThresholdLevel {
    label: { [index: string]: string };
    keyLabel: { [index: string]: string };
    mapKeyLabel: { [index: string]: string };
    plainLabel: { [index: string]: string };
    value: number;
    color: string;
    hidden: boolean;
    start: number;
}

export interface VizThresholds {
    label: { [index: string]: string };
    levels: VizThresholdLevel[];
}

export interface VizConfig {
    name: string;
    disabled: boolean;
    thresholds: VizThresholds;
    minimumGap: number | undefined;
}

export type SensorStrings = { [index: string]: Record<string, string> };

export interface ModuleSensorMeta {
    key: string;
    fullKey: string;
    firmwareKey: string;
    unitOfMeasure: string;
    internal: boolean;
    order: number;
    ranges: SensorRange[];
    viz: VizConfig[];
    strings: SensorStrings;
    aggregationFunction: string;
}

export interface Module {
    id: number;
    key: string;
    internal: boolean;
    header: { manufacturer: number; kind: number; version: number; allKinds: number[] };
    sensors: ModuleSensorMeta[];
}

export interface SensorId {
    id: number;
    key: string;
}

export interface SensorsResponse {
    sensors: SensorId[];
    modules: Module[];
}

export interface Summary {
    numberRecords: number;
}

export interface DataRow {
    time: number;
    stationId: number;
    sensorId: number;
    moduleId: string;
    location: [number, number] | null;
    value: number | null; // TODO Remove
    avg: number | undefined;
    min: number | undefined;
    max: number | undefined;
    last: number | undefined;
}

export interface SensorDataResponse {
    data: DataRow[];
    bucketSize: number;
    bucketSamples: number;
    dataEnd: number | null;
}

export interface StationInfoResponse {
    stationId: number;
    stationName: string;
    stationLocation: [number, number];
    moduleId: ModuleID;
    moduleKey: string;
    sensorId: number;
    sensorKey: string;
    sensorReadAt: string;
}

export interface SensorInfoResponse {
    stations: { [index: string]: StationInfoResponse[] };
}
