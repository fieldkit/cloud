<template>
    <section class="container" v-bind:class="{ 'data-view': viewType === 'data' }">
        <header v-if="viewType === 'project'">Notes & Comments</header>

        <div class="new-comment">
            <UserPhoto v-if="user" :user="user"></UserPhoto>
            <Tiptap v-model="newComment.body" placeholder="Join the discussion!" saveLabel="Post" @save="save(newComment)" />
        </div>

        <div v-if="!errorMessage" class="error">{{ errorMessage }}</div>

        <div v-if="posts.length === 0">There are no comments yet.</div>

        <div class="list" v-if="posts && posts.length > 0">
            <div class="subheader">
                <span class="comments-counter" v-if="viewType === 'project'">{{ posts.length }} comments</span>
                <header v-if="viewType === 'data'">Notes & Comments</header>
            </div>
            <transition-group name="fade">
                <div class="comment comment-first-level" v-for="post in posts" v-bind:key="post.id">
                    <div class="comment-main">
                        <UserPhoto :user="post.author"></UserPhoto>
                        <div class="column-post">
                            <span class="timestamp">{{ formatTimestamp(post.createdAt) }}</span>
                            <span class="author">
                                {{ post.author.name }}
                                <i
                                    class="icon-ellipsis options-trigger"
                                    v-if="user.id === post.author.id || user.admin"
                                    @click="showCommentOptions($event)"
                                ></i>
                                <div class="options-btns">
                                    <button @click="startEditing(post)" v-if="user.id == post.author.id">Edit Post</button>
                                    <button @click="deleteComment(post.id)">Delete Post</button>
                                </div>
                            </span>
                            <Tiptap v-model="post.body" :readonly="post.readonly" saveLabel="Save" @save="saveEdit(post.id, post.body)" />
                        </div>
                    </div>
                    <div class="column">
                        <transition-group name="fade" class="comment-replies">
                            <div class="comment" v-for="reply in post.replies" v-bind:key="reply.id">
                                <div class="comment-main">
                                    <UserPhoto :user="reply.author"></UserPhoto>
                                    <div class="column-reply">
                                        <span class="author">
                                            {{ reply.author.name }}
                                            <i
                                                class="icon-ellipsis options-trigger"
                                                v-if="user.id === reply.author.id || user.admin"
                                                @click="showCommentOptions($event)"
                                            ></i>
                                            <div class="options-btns">
                                                <button @click="startEditing(reply)" v-if="user.id == reply.author.id">Edit Post</button>
                                                <button @click="deleteComment(reply.id)">Delete Post</button>
                                            </div>
                                        </span>
                                        <Tiptap
                                            v-model="reply.body"
                                            :readonly="reply.readonly"
                                            saveLabel="Save"
                                            @save="saveEdit(reply.id, reply.body)"
                                        />
                                    </div>
                                </div>
                            </div>
                        </transition-group>

                        <transition name="fade">
                            <div class="new-comment reply" v-if="newReply && newReply.threadId === post.id">
                                <UserPhoto :user="user"></UserPhoto>
                                <Tiptap v-model="newReply.body" placeholder="Reply to comment" @save="save(newReply)" saveLabel="Post" />
                            </div>
                        </transition>

                        <div class="actions">
                            <button @click="addReply(post)">
                                <img src="@/assets/icon-reply.svg" />
                                Reply
                            </button>
                            <button v-if="viewType === 'data'" @click="viewDataClick(post)">
                                <img src="@/assets/icon-view-data.svg" />
                                View Data
                            </button>
                        </div>
                    </div>
                </div>
            </transition-group>
        </div>
    </section>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";
import moment from "moment";
import { NewComment } from "@/views/comments/model";
import { Comment } from "@/views/comments/model";
import { CurrentUser } from "@/api";
import { CommentsErrorsEnum } from "@/views/comments/model";
import Tiptap from "@/views/shared/Tiptap.vue";

export default Vue.extend({
    name: "Comments",
    components: {
        ...CommonComponents,
        Tiptap,
    },
    props: {
        user: {
            type: Object as PropType<CurrentUser>,
            required: true,
        },
        parentData: {
            type: [Number, Object],
            required: true,
        },
    },
    data(): {
        posts: Comment[];
        placeholder: string | null;
        viewType: string;
        newComment: {
            projectId: number | null;
            bookmark: string | null;
            body: string | null;
        };
        newReply: {
            projectId: number | null;
            bookmark: string | null;
            body: string | null;
            threadId: number | null;
        };
        errorMessage: string | null;
    } {
        return {
            posts: [],
            placeholder: null,
            viewType: typeof this.$props.parentData === "number" ? "project" : "data",
            newComment: {
                projectId: typeof this.parentData === "number" ? this.parentData : null,
                bookmark: null,
                body: "",
            },
            newReply: {
                projectId: typeof this.parentData === "number" ? this.parentData : null,
                bookmark: null,
                body: "",
                threadId: null,
            },
            errorMessage: null,
        };
    },
    watch: {
        parentData(): Promise<void> {
            return this.getComments();
        },
    },
    mounted(): Promise<void> {
        this.placeholder = this.getNewCommentPlaceholder();
        return this.getComments();
    },
    methods: {
        getNewCommentPlaceholder(): string {
            if (this.viewType === "project") {
                return "Comment on Project";
            } else {
                return "Write a comment about this Data View";
            }
        },
        async save(comment: NewComment): Promise<void> {
            this.errorMessage = null;
            if (this.viewType === "data") {
                comment.bookmark = JSON.stringify(this.parentData);
            }
            await this.$services.api
                .postComment(comment)
                .then((response: { post: Comment }) => {
                    this.newComment.body = "";
                    // add the comment to the replies array
                    if (comment.threadId) {
                        if (this.posts) {
                            this.posts
                                .filter((post) => post.id === comment.threadId)[0]
                                .replies.push(
                                    new Comment(
                                        response.post.id,
                                        response.post.author,
                                        response.post.bookmark,
                                        response.post.body,
                                        response.post.createdAt,
                                        response.post.updatedAt
                                    )
                                );
                            this.newReply.body = "";
                        } else {
                            console.warn(`posts is null`);
                        }
                    } else {
                        // add it to the posts array
                        if (this.posts) {
                            this.posts.unshift(
                                new Comment(
                                    response.post.id,
                                    response.post.author,
                                    response.post.bookmark,
                                    response.post.body,
                                    response.post.createdAt,
                                    response.post.updatedAt
                                )
                            );
                            this.newComment.body = "";
                        } else {
                            console.log(`posts is null`);
                        }
                    }
                })
                .catch(() => {
                    this.errorMessage = CommentsErrorsEnum.postComment;
                });
        },
        formatTimestamp(timestamp: number): string {
            return moment(timestamp).fromNow();
        },
        addReply(post: Comment): void {
            if (this.newReply.body && post.id === this.newReply.threadId) {
                return;
            }
            this.errorMessage = null;
            this.newReply.threadId = post.id;
            this.newReply.body = "";
        },
        async getComments(): Promise<void> {
            await this.$services.api
                .getComments(this.parentData)
                .then((data) => {
                    this.posts = [];
                    data.posts.forEach((post) => {
                        this.posts.push(new Comment(post.id, post.author, post.bookmark, post.body, post.createdAt, post.updatedAt));

                        post.replies.forEach((reply) => {
                            this.posts[this.posts.length - 1].replies.push(
                                new Comment(reply.id, reply.author, reply.bookmark, reply.body, reply.createdAt, reply.updatedAt)
                            );
                        });
                    });
                })
                .catch(() => {
                    this.errorMessage = CommentsErrorsEnum.getComments;
                });
        },
        viewDataClick(post: Comment) {
            if (post.bookmark) {
                this.$emit("viewDataClicked", JSON.parse(post.bookmark));
            }
        },
        showCommentOptions(event: MouseEvent) {
            if (event.target) {
                const optionsMenu = (event.target as HTMLElement).nextElementSibling;

                if (!(optionsMenu as HTMLElement).classList.contains("visible")) {
                    (optionsMenu as HTMLElement).classList.add("visible");
                    setTimeout(function () {
                        document.addEventListener(
                            "click",
                            function () {
                                (optionsMenu as HTMLElement).classList.remove("visible");
                            },
                            {
                                once: true,
                            }
                        );
                    }, 1);
                }
            }
        },
        deleteComment(commentID: number) {
            this.$services.api
                .deleteComment(commentID)
                .then((response) => {
                    if (response) {
                        this.getComments();
                    } else {
                        this.errorMessage = CommentsErrorsEnum.deleteComment;
                    }
                })
                .catch(() => {
                    this.errorMessage = CommentsErrorsEnum.deleteComment;
                });
        },
        startEditing(post: Comment) {
            post.readonly = false;
        },
        saveEdit(commentID: number, body: Record<string, unknown>) {
            this.$services.api
                .editComment(commentID, body)
                .then((response) => {
                    if (response) {
                        this.getComments();
                    } else {
                        this.errorMessage = CommentsErrorsEnum.deleteComment;
                    }
                })
                .catch(() => {
                    this.errorMessage = CommentsErrorsEnum.editComment;
                });
        },
    },
});
</script>

<style lang="scss">
@import "../../scss/global";

button {
    padding: 0;
    border: 0;
    outline: 0;
    box-shadow: none;
    cursor: pointer;
    background: transparent;
}

* {
    box-sizing: border-box;
    font-size: 14px;
}
.hide {
    display: none;
}

.container {
    margin-top: 20px;
    padding: 0 25px 30px 20px;
    background: #fff;
    border-radius: 1px;
    border: 1px solid $color-border;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);

    &.data-view {
        margin-top: 0;
        padding-top: 45px;
        padding-bottom: 22px;
        box-shadow: none;
        border: 0;
    }
}

header {
    @include flex(center, space-between);
    height: 52px;
    border-bottom: 1px solid $color-border;
    font-size: 20px;
    font-weight: 500;

    .data-view & {
        font-size: 18px;
        height: auto;
        border: none;
    }
}

.subheader {
    @include flex(center, space-between);
    border-top: 1px solid $color-border;
    border-bottom: 1px solid $color-border;
    padding: 15px 0;

    .data-view & {
        border-top: none;
    }
}

.list {
    overflow-y: hidden;

    .data-view & {
        margin-top: 30px;
        width: 60%;
    }
}

.new-comment {
    @include flex(top);
    padding: 22px 0;
    position: relative;

    .container.data-view & {
        &:not(.reply) {
            background-color: rgba(#f4f5f7, 0.55);
            padding: 18px 23px 17px 15px;

            .new-comment-submit {
                right: 30px;
            }
        }
    }

    &.reply {
        padding: 0;
        margin-top: 10px;
        width: 100%;

        input {
            height: 35px;
        }
    }

    img {
        margin-top: 0;
        width: 30px;
        height: 30px;
    }

    &:not(.reply) {
        img {
            width: 46px;
            height: 46px;
        }
    }

    input {
        height: 45px;
        padding: 14px 72px 12px 13px;
        border-radius: 2px;
        border: solid 1px $color-border;
        outline: none;
        width: 100%;
        font-weight: 500;

        &::placeholder {
            color: #cccdcf;
        }
    }

    &-submit {
        @include position(absolute, 50% 10px null null);
        @include flex(center);
        height: 45px;
        padding: 0 10px;
        transform: translateY(-50%);
        font-weight: 900;
    }
}

.comments-counter {
    font-family: $font-family-light;
}

.author {
    font-size: 16px;
    font-weight: 500;
    margin-bottom: 5px;
    margin-top: 2px;
    position: relative;
}

.options {
    &-trigger {
        @include position(absolute, 0 -40px null null);
        opacity: 0;
        visibility: hidden;
        transition: opacity 0.33s;
    }

    &-btns {
        @include position(absolute, 0 calc(-100px - 50px) null null);
        opacity: 0;
        visibility: hidden;
        padding: 10px;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);
        background: #fff;
        z-index: $z-index-top;
        transition: opacity 0.33s;
        width: 100px;

        &.visible {
            opacity: 1;
            visibility: visible;
        }

        > * {
            display: block;
            white-space: nowrap;
            font-family: $font-family-bold;
            cursor: pointer;

            &:not(:last-of-type) {
                margin-bottom: 10px;
            }
        }
    }
}

.body {
    max-width: unset;
    font-family: $font-family-light;
    outline: none;
    border: solid 1px $color-border;
    min-height: 35px;
    width: calc(100% - 40px);
    overflow-wrap: break-word;

    &[readonly] {
        border: none;
        max-width: 550px;
        min-height: unset;
    }

    + .new-comment-submit {
        transform: translateY(calc(-50% + 12px));
        right: -10px;
    }
}

.comment {
    @include flex(flex-start);
    flex: 100%;
    padding: 15px 0 0;
    position: relative;
    flex-wrap: wrap;

    &-first-level {
        border-bottom: 1px solid $color-border;
    }

    &::v-deep .default-user-icon {
        margin-top: 0;
        width: 30px;
        height: 30px;
    }

    .column {
        /* padding-top: 1em; */

        &:nth-of-type(2) {
            padding-left: 42px;
        }
    }

    &.highlight {
        opacity: 0.1;
        background-color: #a0dbe1;
        @include position(absolute, 0 null null 0);
    }
}

.comment-replies {
    width: 100%;

    .column {
        border-bottom: none;
    }
}

.comment-main {
    display: flex;
    flex: 100%;

    @include attention() {
        .options-trigger {
            opacity: 1;
            visibility: visible;
        }
    }
}

.column,
.column-post,
.column-reply {
    @include flex(flex-start);
    width: 100%;
    flex-direction: column;
    position: relative;
}

.actions {
    margin: 15px 0;
    user-select: none;
    @include flex();

    button {
        font-weight: 500;
        margin-right: 20px;
        @include flex(flex-start);
    }

    img {
        width: 14px;
        margin-right: 6px;
    }
}

.timestamp {
    @include position(absolute, 0 0 null null);
    font-family: $font-family-light;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.25s ease-in-out;
}

.icon-ellipsis {
    display: block;
    cursor: pointer;

    &:after {
        @include flex(flex-end);
        content: "...";
        height: 17px;
        font-size: 32px;
        font-family: $font-family-bold;
    }
}
.error {
    color: $color-danger;
    margin-bottom: 10px;
}
</style>
