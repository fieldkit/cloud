export interface HasLocation {
    readonly latitude: number | null;
    readonly longitude: number | null;
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

    lngLat(): number[] {
        return [this.longitude, this.latitude];
    }
}

export class BoundingRectangle {
    constructor(public min: Location | null = null, public max: Location | null = null) {}

    empty(): boolean {
        return this.min == null || this.max == null;
    }

    include(l: Location): BoundingRectangle {
        this.min = this.min == null ? l.clone() : this.min.minimum(l);
        this.max = this.max == null ? l.clone() : this.max.maximum(l);
        return this;
    }

    contains(l: Location): boolean {
        if (this.min == null || this.max == null) {
            return false;
        }
        return (
            l.latitude >= this.min.latitude &&
            l.longitude >= this.min.longitude &&
            l.latitude <= this.max.latitude &&
            l.longitude <= this.max.longitude
        );
    }

    isEmpty(): boolean {
        return this.min == null || this.max == null;
    }

    expandIfSingleCoordinate(defaultCenter: Location, margin: number): BoundingRectangle {
        if (this.isEmpty()) {
            return BoundingRectangle.around(defaultCenter, margin);
        }
        if (this.isSingleCoordinate()) {
            return BoundingRectangle.around(this.min, margin);
        }
        return this;
    }

    isSingleCoordinate(): boolean {
        return this.min.latitude == this.max.latitude || this.max.longitude == this.max.longitude;
    }

    static around(center: Location, margin: number) {
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
        const min = new Location(center.latitude - latitudeMargin, center.longitude - longitudeMargin);
        const max = new Location(center.latitude + latitudeMargin, center.longitude + longitudeMargin);
        return new BoundingRectangle(min, max);
    }
}

export class MapCenter {
    constructor(public readonly location: Location, public readonly bounds: BoundingRectangle, public readonly zoom: number) {}
}
