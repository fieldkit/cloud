import * as ActionTypes from "@/store/actions";

import { Bookmark } from "@/views/viz/viz";

export class ExportParams {
    constructor(public readonly csv: boolean) {}
}

export class ExportDataAction {
    type = ActionTypes.BEGIN_EXPORT;

    constructor(public readonly bookmark: Bookmark, public readonly params: ExportParams = { csv: true }) {}
}
