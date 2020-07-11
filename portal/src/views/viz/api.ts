interface SensorRange {
    minimum: number;
    maximum: number;
}

export interface ModuleSensor {
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
    sensors: ModuleSensor[];
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
    location: number[] | null;
    value: number;
}

export interface SensorDataResponse {
    summaries: { [index: string]: Summary };
    data: DataRow[];
}

export interface SensorInfoResponse {
    stations: { [index: string]: { key: string; sensorId: number }[] };
}
