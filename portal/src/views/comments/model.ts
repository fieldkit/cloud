export interface Comment {
    id: number;
    author: { id: number; name: string; photo: object };
    bookmark?: string;
    body: string;
    replies: Comment[];
    createdAt: number;
    updatedAt: number;
    readonly: boolean;
}

export class Comment {
    id: number;
    author: { id: number; name: string; photo: object };
    bookmark?: string;
    body: string;
    replies: Comment[];
    createdAt: number;
    updatedAt: number;
    readonly: boolean;

    constructor(
        id: number,
        author: { id: number; name: string; photo: object },
        bookmark: string | undefined,
        body: string,
        createdAt: number,
        updatedAt: number
    ) {
        this.id = id;
        this.author = author;
        this.bookmark = bookmark;
        this.body = body;
        this.createdAt = createdAt;
        this.updatedAt = updatedAt;
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

export enum CommentsErrorsEnum {
    getComments = "Something went wrong loading the comments",
    postComment = "Something went saving your comment.",
    deleteComment = "Something went wrong deleting your comment.",
    editComment = "Something went wrong editing your comment. Your changes will are not saved.",
}
