import _ from 'lodash';

export class GeoRect {
    sw: Array<number>
    ne: Array<number>

    constructor(bounds) {
        this.ne = bounds.ne;
        this.sw = bounds.sw;
    }

    contains(r) {
        return (r.ne[0] < this.ne[0] && r.ne[0] > this.sw[0]) &&
               (r.sw[0] < this.ne[0] && r.sw[0] > this.sw[0]) &&
               (r.ne[1] < this.ne[1] && r.ne[1] > this.sw[1]) &&
               (r.sw[1] < this.ne[1] && r.sw[1] > this.sw[1]);
    }

    enlarge(factor) {
        const w = this.ne[0] - this.sw[0];
        const h = this.ne[1] - this.sw[1];
        const newNe = [
            this.ne[0] + (w / 2),
            this.ne[1] + (h / 2),
        ];
        const newSw = [
            this.sw[0] - (w / 2),
            this.sw[1] - (h / 2),
        ];
        return new GeoRect({
            ne: newNe,
            sw: newSw,
        });
    }
};

export class GeoRectSet {
    rects: Array<GeoRect>

    constructor(rects) {
        this.rects = _(rects).map(r => new GeoRect(r)).value();
    }

    add(r) {
        this.rects.push(r);
    }

    contains(q) {
        for (let r of this.rects) {
            if (r.contains(q)) {
                return true;
            }
        }
        return false;
    }
};

