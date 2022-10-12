export interface DiscussionBase {
    id: number;
    author: { id: number; name: string; photo: object };
    bookmark?: string;
    createdAt: number;
    updatedAt: number;
    readonly?: boolean;
}

export class DiscussionBase {
    id: number;
    author: { id: number; name: string; photo: object };
    bookmark?: string;
    createdAt: number;
    updatedAt: number;
    readonly?: boolean;

    constructor(
        id: number,
        author: { id: number; name: string; photo: object },
        bookmark: string | undefined,
        createdAt: number,
        updatedAt: number
    ) {
        this.id = id;
        this.author = author;
        this.bookmark = bookmark;
        this.createdAt = createdAt;
        this.updatedAt = updatedAt;
    }
}

export interface Comment extends DiscussionBase {
    body: string;
    replies: Comment[];
    readonly: boolean;
    type?: string;
}

export class Comment extends DiscussionBase {
    body: string;
    replies: Comment[];
    readonly: boolean;

    constructor(
        id: number,
        author: { id: number; name: string; photo: object },
        bookmark: string | undefined,
        body: string,
        createdAt: number,
        updatedAt: number
    ) {
        super(id, author, bookmark, createdAt, updatedAt);
        this.body = body;
        this.replies = [];
        this.readonly = true;
        this.type = 'comment';
    }
}

export interface NewComment {
    projectId: number;
    bookmark: string;
    body: string;
    threadId?: number;
}

export interface DataEvent extends DiscussionBase {
    title: string;
    description: string;
    start: number;
    end: number;
    readonly?: boolean;
    type?: string;
}

export class DataEvent extends DiscussionBase {
    title: string;
    description: string;
    start: number;
    end: number;
    readonly?: boolean;
    type?: string;

    constructor(
        id: number,
        author: { id: number; name: string; photo: object },
        bookmark: string | undefined,
        createdAt: number,
        updatedAt: number,
        title: string,
        description: string,
        start: number,
        end: number
    ) {
        super(id, author, bookmark, createdAt, updatedAt);
        this.title = title;
        this.description = description;
        this.start = start;
        this.end = end;
        this.readonly = true;
        this.type = 'event';
    }
}

export interface NewDataEvent extends NewComment {
    title: string;
    description: string;
    allProjectSensors?: boolean;
}

export enum CommentsErrorsEnum {
    getComments = "Something went wrong loading the comments",
    postComment = "Something went wrong saving your comment.",
    deleteComment = "Something went wrong deleting your comment.",
    editComment = "Something went wrong editing your comment. Your changes will are not saved.",
}

export enum DataEventsErrorsEnum {
    getDataEvents = "Something went wrong loading the events",
    postDataEvent = "Something went wrong saving your event.",
    deleteDataEvent = "Something went wrong deleting your event.",
    editDataEvent = "Something went wrong editing your event. Your changes will are not saved.",
}