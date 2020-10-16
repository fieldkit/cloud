import _ from "lodash";

export interface HasLocation {
    readonly latitude: number | null;
    readonly longitude: number | null;
}

export type LngLat = [number, number];

function lngLatMinimum(c: LngLat[]): LngLat {
    return [_.min(_.map(c, (v) => v[0])), _.min(_.map(c, (v) => v[1]))];
}

function lngLatMaximum(c: LngLat[]): LngLat {
    return [_.max(_.map(c, (v) => v[0])), _.max(_.map(c, (v) => v[1]))];
}

export class Location implements HasLocation {
    constructor(public readonly latitude: number, public readonly longitude) {}

    maximum(other: Location): Location {
        return new Location(
            other.latitude > this.latitude ? other.latitude : this.latitude,
            other.longitude > this.longitude ? other.longitude : this.longitude
        );
    }

    minimum(other: Location): Location {
        return new Location(
            other.latitude < this.latitude ? other.latitude : this.latitude,
            other.longitude < this.longitude ? other.longitude : this.longitude
        );
    }

    clone(): Location {
        return new Location(this.latitude, this.longitude);
    }

    lngLat(): LngLat {
        return [this.longitude, this.latitude];
    }
}

export class BoundingRectangle {
    constructor(public min: LngLat | null = null, public max: LngLat | null = null) {}

    public empty(): boolean {
        return this.min == null || this.max == null;
    }

    public include(l: LngLat): BoundingRectangle {
        this.min = this.min == null ? l : lngLatMinimum([this.min, l]);
        this.max = this.max == null ? l : lngLatMaximum([this.max, l]);
        return this;
    }

    public includeAll(l: LngLat[]): BoundingRectangle {
        return l.reduce((b, c) => b.include(c), this);
    }

    public contains(l: LngLat): boolean {
        if (this.min == null || this.max == null) {
            return false;
        }
        return l[1] >= this.min[1] && l[0] >= this.min[0] && l[1] <= this.max[1] && l[0] <= this.max[0];
    }

    public isEmpty(): boolean {
        return this.min == null || this.max == null;
    }

    private calculateMargin(margin: number | undefined): number {
        if (this.isSingleCoordinate()) {
            return 1000;
        }
        const maximum = _.max([this.max[0] - this.min[0], this.max[1] - this.min[1]]);
        return maximum * 50000;
    }

    public zoomOutOrAround(defaultCenter: LngLat, margin: number | undefined = undefined): BoundingRectangle {
        if (this.isEmpty()) {
            return BoundingRectangle.around(defaultCenter, margin);
        }
        if (this.isSingleCoordinate()) {
            return BoundingRectangle.around(this.min, this.calculateMargin(margin));
        }
        return this.zoomOut(this.calculateMargin(margin));
    }

    public zoomOut(margin: number): BoundingRectangle {
        /*
		At 38 degrees North latitude:
		One degree of latitude equals approximately 364,000 feet (69
		miles), one minute equals 6,068 feet (1.15 miles), and
		one-second equals 101 feet.

		One-degree of longitude equals 288,200 feet (54.6 miles), one
		minute equals 4,800 feet (0.91 mile), and one second equals 80
		feet.
		*/
        const FeetPerLatitude = 364000; /* ft per degree */
        const FeetPerLongitude = 288200; /* ft per degree */
        const latitudeMargin = margin / FeetPerLatitude;
        const longitudeMargin = margin / FeetPerLongitude;
        const min: LngLat = [this.min[0] - longitudeMargin, this.min[1] - latitudeMargin];
        const max: LngLat = [this.max[0] + longitudeMargin, this.max[1] + latitudeMargin];
        return new BoundingRectangle(min, max);
    }

    public expandIfSingleCoordinate(defaultCenter: LngLat, margin: number): BoundingRectangle {
        if (this.isEmpty()) {
            return BoundingRectangle.around(defaultCenter, margin);
        }
        if (this.isSingleCoordinate()) {
            return BoundingRectangle.around(this.min, margin);
        }
        return this;
    }

    public isSingleCoordinate(): boolean {
        return this.min[0] === this.max[0] || this.min[1] === this.max[1];
    }

    public static around(center: LngLat, margin: number) {
        return new BoundingRectangle(center, center).zoomOut(margin);
    }
}

export class MapCenter {
    constructor(public readonly location: Location, public readonly bounds: BoundingRectangle, public readonly zoom: number) {}
}
