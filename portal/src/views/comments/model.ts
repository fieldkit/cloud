export interface Comment {
    id: number;
    author: object;
    bookmark?: string;
    body: string;
    replies: Comment[];
    createdAt: number;
    updatedAt: number;
}

export interface NewComment {
    projectId: number;
    bookmark: string;
    body: string;
    threadId?: number;
}
