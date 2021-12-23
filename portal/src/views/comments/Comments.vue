<template>
    <section class="container" v-bind:class="{ 'data-view': viewType === 'data' }">
        <header v-if="viewType === 'project'">Notes & Comments</header>

        <form @submit.prevent="save(newComment)" class="new-comment">
            <UserPhoto v-if="user" :user="user"></UserPhoto>
            <TextAreaField class="new-comment-input" :placeholder="placeholder" v-model="newComment.body" />
            <button type="submit" class="new-comment-submit" v-if="newComment.body.length > 0">Post</button>
        </form>

        <div v-if="errorMessage" class="error">{{ errorMessage }}</div>

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
                        <div class="column">
                            <span class="timestamp">{{ formatTimestamp(post.createdAt) }}</span>
                            <span class="author">
                                {{ post.author.name }}
                                <ListItemOptions
                                    v-if="user.id === post.author.id || user.admin"
                                    @listItemOptionClick="onListItemOptionClick($event, post)"
                                    :options="getCommentOptions(post)"
                                ></ListItemOptions>
                            </span>
                            <div class="break-word" v-if="post.readonly">{{ post.body }}</div>
                            <TextAreaField v-else class="body" v-model="post.body" />
                            <button type="submit" class="new-comment-submit" v-if="!post.readonly" @click="saveEdit(post.id, post.body)">
                                Save
                            </button>
                        </div>
                    </div>
                    <div class="column">
                        <transition-group name="fade" class="comment-replies">
                            <div class="comment" v-for="reply in post.replies" v-bind:key="reply.id">
                                <div class="comment-main">
                                    <UserPhoto :user="reply.author"></UserPhoto>
                                    <div class="column">
                                        <span class="author">
                                            {{ reply.author.name }}
                                            <ListItemOptions
                                                v-if="user.id === post.author.id || user.admin"
                                                @listItemOptionClick="onListItemOptionClick($event, reply)"
                                                :options="getCommentOptions(post)"
                                            ></ListItemOptions>
                                        </span>
                                        <div v-if="reply.readonly">{{ reply.body }}</div>
                                        <TextAreaField v-else class="body" v-model="reply.body" />
                                        <button
                                            type="submit"
                                            class="new-comment-submit"
                                            v-if="!reply.readonly"
                                            @click="saveEdit(reply.id, reply.body)"
                                        >
                                            Save
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </transition-group>

                        <transition name="fade">
                            <form
                                @submit.prevent="save(newReply)"
                                class="new-comment reply"
                                v-if="newReply && newReply.threadId === post.id"
                            >
                                <UserPhoto :user="user"></UserPhoto>
                                <TextAreaField class="body" placeholder="Reply to comment" v-model="newReply.body" />
                                <button type="submit" class="new-comment-submit" v-if="newReply.body">Post</button>
                            </form>
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
import ListItemOptions from "@/views/shared/ListItemOptions.vue";

export default Vue.extend({
    name: "Comments",
    components: {
        ...CommonComponents,
        ListItemOptions,
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
                    setTimeout(function() {
                        document.addEventListener(
                            "click",
                            function() {
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
        saveEdit(commentID: number, body: string) {
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
        onListItemOptionClick(event: string, post: Comment): void {
            if (event === "edit-comment") {
                this.startEditing(post);
            }
            if (event === "delete-comment") {
                this.deleteComment(post.id);
            }
        },
        getCommentOptions(post: Comment): { label: string; event: string }[] {
            if (this.user.id === post.author.id) {
                return [
                    {
                        label: "Edit post",
                        event: "edit-comment",
                    },
                    {
                        label: "Delete post",
                        event: "delete-comment",
                    },
                ];
            }

            return [
                {
                    label: "Delete post",
                    event: "delete-comment",
                },
            ];
        },
        test(event: any) {
            console.log("COMPSOING", event.target.composing);
        },
    },
});
</script>

<style lang="scss" scoped>
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
    padding: 0 20px 30px 20px;
    background: #fff;
    border-radius: 1px;
    border: 1px solid $color-border;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);

    &.data-view {
        margin-top: 0;
        padding-top: 45px;
        box-shadow: none;
        border: 0;
    }

    @include bp-down($xs) {
        margin: 20px -10px 0;
        padding: 0 10px 30px 10px;
    }
}

header {
    @include flex(center, space-between);
    padding: 13px 0;
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
        padding: 0;
    }
}

.list {
    overflow: hidden;

    .data-view & {
        margin-top: 30px;
        @include bp-down($xs) {
            margin-top: 10px;
        }
    }
}

.new-comment {
    @include flex(center);
    padding: 22px 0;
    position: relative;

    @include bp-down($xs) {
        margin: 0 -10px;
        padding: 15px 10px 15px 10px;
    }

    .container.data-view & {
        &:not(.reply) {
            background-color: rgba(#f4f5f7, 0.55);
            padding: 18px 23px 17px 15px;

            .new-comment-submit {
                right: 25px;
            }
        }
    }

    &.reply {
        padding: 0 0 0;
        margin: 10px 0 0 0;
        width: 100%;
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

            @include bp-down($xs) {
                width: 42px;
                height: 42px;
            }
        }
    }

    &-input {
        width: 100%;
        display: flex;

        ::v-deep textarea {
            margin: 0;
            box-sizing: border-box;
            padding: 14px 72px 14px 13px;
            border-radius: 2px;
            border: solid 1px $color-border;
            outline: none;
            width: 100%;
            font-weight: 500;
            resize: none;
            overflow: hidden;

            &::placeholder {
                color: #cccdcf;
            }

            @include bp-down($xs) {
                padding: 7px 40px 7px 7px;
            }
        }
    }

    &-submit {
        @include position(absolute, 50% 10px null null);
        @include flex(center);
        height: 45px;
        padding: 0 10px;
        transform: translateY(-50%);
        font-weight: 900;

        .data-view &,
        .reply & {
            right: 0;
        }
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

::v-deep .options-trigger {
    opacity: 0;
    visibility: hidden;
}

.body {
    max-width: unset;
    font-family: $font-family-light;
    outline: none;
    border: solid 1px $color-border;
    width: calc(100% - 40px);
    overflow-wrap: break-word;

    &[readonly] {
        border: none;
        max-width: 550px;
        min-height: unset;
    }

    ::v-deep textarea {
        border: 0;
        margin: 0;
        padding: 7px 45px 0 7px;
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
        &:nth-of-type(2) {
            padding-left: 36px;
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
    overflow-wrap: break-word;

    @include attention() {
        ::v-deep .options-trigger {
            opacity: 1;
            visibility: visible;
        }
    }
}

.column {
    @include flex(flex-start);
    width: 100%;
    flex-direction: column;
    position: relative;

    > * {
        overflow-wrap: anywhere;
    }
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

.error {
    color: $color-danger;
    margin-bottom: 10px;
}

::v-deep textarea {
    overflow-y: hidden;
    resize: none;
}
</style>
