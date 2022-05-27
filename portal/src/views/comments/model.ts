export interface DiscussionBase {
    id: number;
    author: { id: number; name: string; photo: object };
    bookmark?: string;
    createdAt: number;
    updatedAt: number;
}

export class DiscussionBase {
    id: number;
    author: { id: number; name: string; photo: object };
    bookmark?: string;
    createdAt: number;
    updatedAt: number;

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
}

export class DataEvent extends DiscussionBase {
    title: string;
    description: string;
    start: number;
    end: number;

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
    }
}

export interface NewDataEvent extends NewComment {
    title: string;
    description: string;
}

export enum CommentsErrorsEnum {
    getComments = "Something went wrong loading the comments",
    postComment = "Something went wrong saving your comment.",
    deleteComment = "Something went wrong deleting your comment.",
    editComment = "Something went wrong editing your comment. Your changes will are not saved.",
}