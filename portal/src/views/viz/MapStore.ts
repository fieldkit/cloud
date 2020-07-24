export interface Map {
    resize();
    addSource(name: string, source: any);
    addLayer(layer: any);
    fitBounds(bounds: any, options: any);
    getBounds(): any;
    getZoom(): number;
    setZoom(zoom: number);
}

export class MapStore {
    constructor(public readonly maps: { [index: string]: Map } = {}) {}

    public set(key: string, map: Map): Map {
        this.maps[key] = map;
        return map;
    }

    public remove(key: string) {
        delete this.maps[key];
    }

    public get(key: string): Map | null {
        return this.maps[key];
    }
}
