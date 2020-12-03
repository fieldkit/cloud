interface SensorRange {
    minimum: number;
    maximum: number;
}

export interface ModuleSensorMeta {
    key: string;
    fullKey: string;
    firmwareKey: string;
    unitOfMeasure: string;
    internal: boolean;
    ranges: SensorRange[];
}

interface Module {
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
    id: number;
    time: number;
    stationId: number;
    sensorId: number;
    location: [number, number] | null;
    value: number;
}

export interface SensorDataResponse {
    summaries: { [index: string]: Summary };
    aggregate: string;
    data: DataRow[];
}

interface StationInfoResponse {
    sensorId: number;
    key: string;
    name: string;
}

export interface SensorInfoResponse {
    stations: { [index: string]: StationInfoResponse[] };
}
