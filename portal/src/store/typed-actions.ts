import * as ActionTypes from "@/store/actions";

import { Bookmark } from "@/views/viz/viz";

export class ExportParams {
    public readonly csv: boolean = false;
    public readonly jsonLines: boolean = false;

    constructor(f: { csv?: boolean; jsonLines?: boolean }) {
        this.csv = f.csv || false;
        this.jsonLines = f.jsonLines || false;
    }
}

export class ExportDataAction {
    type = ActionTypes.BEGIN_EXPORT;

    constructor(public readonly bookmark: Bookmark, public readonly params: ExportParams) {}
}

export class ResumeAction {
    type = ActionTypes.LOGIN_RESUME;

    constructor(public readonly token: string) {}
}
