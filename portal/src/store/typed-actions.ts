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

export class LoginOidcAction {
    type = ActionTypes.LOGIN_OIDC;

    constructor(
        public readonly token: string | null,
        public readonly state: string,
        public readonly sessionState: string,
        public readonly code: string
    ) {}
}

export class DiscourseParams {
    constructor(public readonly sso: string, public readonly sig: string) {}
}

export interface LoginFields {
    readonly email: string;
    readonly password: string;
}

export class LoginDiscourseAction {
    type = ActionTypes.LOGIN_DISCOURSE;

    constructor(
        public readonly token: string | null,
        public readonly login: LoginFields | null,
        public readonly discourse: DiscourseParams
    ) {}
}

export class MarkNotificationsSeen {
    type = ActionTypes.NOTIFICATIONS_SEEN;

    constructor(public readonly ids: number[]) {}
}
