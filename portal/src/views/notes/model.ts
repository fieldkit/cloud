import _ from "lodash";
import moment from "moment";

export class ExistingFieldNote {
    constructor(
        public readonly id: number,
        public readonly key: string,
        public readonly body: string,
        public readonly mediaIds: number[]
    ) {}
}

export class NewFieldNote {
    constructor(public readonly key: string, public readonly body: string, public readonly mediaIds: number[]) {}
}

export class Ids {
    constructor(public readonly mobile: number, public readonly portal: number) {}
}

export interface PortalStationNotes {
    id: number;
    createdAt: number;
    author: { id: number; name: number };
    key: string;
    body: string;
    media: { id: number; key: string; url: string; contentType: string }[];
}

export interface PortalNoteMedia {
    id: number;
    contentType: string;
    url: string;
    key: string;
}

export interface PortalStationNotesReply {
    media: PortalNoteMedia[];
    notes: PortalStationNotes[];
}

export class PatchPortalNote {
    constructor(public readonly creating: NewFieldNote[], public readonly notes: ExistingFieldNote[]) {}
}

export interface PortalPatchNotesPayload {
    notes: PatchPortalNote[];
}

export class NoteHelp {
    constructor(public readonly title: string) {}
}

export class NoteMedia {
    constructor(public readonly key: string) {}

    public static onlyAudio(media: NoteMedia[]): NoteMedia[] {
        return media.filter(NoteMedia.isAudio);
    }

    public static onlyPhotos(media: NoteMedia[]): NoteMedia[] {
        return media.filter(NoteMedia.isPhoto);
    }

    public static isPhoto(nm: NoteMedia): boolean {
        return !NoteMedia.isAudio(nm);
    }

    public static isAudio(nm: NoteMedia): boolean {
        return /(m4a|caf)$/.test(nm.key.toLowerCase());
    }
}

export class NoteForm {
    constructor(
        public readonly body: string,
        public readonly help: NoteHelp,
        public readonly photos: NoteMedia[] = [],
        public readonly audio: NoteMedia[] = []
    ) {}

    public withBody(body: string) {
        return new NoteForm(body, this.help, this.photos, this.audio);
    }
}

export class AddedPhoto {
    public readonly key: string;

    constructor(public readonly type: any, public readonly file: any, public readonly image: any) {
        this.key = moment.utc().toISOString();
    }
}

export class Notes {
    static Keys = ["studyObjective", "sitePurpose", "siteCriteria", "siteDescription"];

    public readonly studyObjective: NoteForm = new NoteForm("", new NoteHelp("Study Objective"));
    public readonly sitePurpose: NoteForm = new NoteForm("", new NoteHelp("Purpose of Site"));
    public readonly siteCriteria: NoteForm = new NoteForm("", new NoteHelp("Site Criteria"));
    public readonly siteDescription: NoteForm = new NoteForm("", new NoteHelp("Site Description"));

    constructor(public readonly addedPhotos: AddedPhoto[] = []) {}

    public get progress(): { total: number; completed: number } {
        const progress = _.map(Notes.Keys, (key) => {
            return this[key].body.length > 0 ? 1 : 0;
        });
        const done = _.sum(progress);
        return {
            total: Notes.Keys.length,
            completed: done,
        };
    }

    public static createFrom(portalNotes: PortalStationNotesReply): Notes {
        return portalNotes.notes.reduce((formNotes, portalNote) => {
            const key = portalNote.key;
            if (!formNotes[key]) {
                throw new Error("unexpected note");
            }
            Object.assign(formNotes[key], {
                photos: NoteMedia.onlyPhotos(portalNote.media),
                audio: NoteMedia.onlyAudio(portalNote.media),
                body: portalNote.body,
            });
            return formNotes;
        }, new Notes());
    }
}

export function mergeNotes(portalNotes: PortalStationNotesReply, notesForm: Notes): PatchPortalNote {
    const portalExisting = _.keyBy(portalNotes.notes, (n) => n.key);
    const localByKey = {
        studyObjective: notesForm.studyObjective,
        sitePurpose: notesForm.sitePurpose,
        siteCriteria: notesForm.siteCriteria,
        siteDescription: notesForm.siteDescription,
    };

    const media = {};

    const modifications = _(localByKey)
        .mapValues((value, key) => {
            const photoIds = value.photos.map((m) => media[m.path]).filter((v) => v);
            const audioIds = value.audio.map((m) => media[m.path]).filter((v) => v);
            const mediaIds = [...photoIds, ...audioIds];

            if (portalExisting[key]) {
                return {
                    creating: null,
                    updating: new ExistingFieldNote(portalExisting[key].id, key, value.body, mediaIds),
                };
            }
            return {
                creating: new NewFieldNote(key, value.body, mediaIds),
                updating: null,
            };
        })
        .values()
        .value();

    const creating = modifications.map((v) => v.creating).filter((v) => v !== null && v.body.length > 0) as NewFieldNote[];
    const updating = modifications.map((v) => v.updating).filter((v) => v !== null) as ExistingFieldNote[];

    return new PatchPortalNote(creating, updating);
}
