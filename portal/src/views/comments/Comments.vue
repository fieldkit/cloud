<template>
    <section class="container" v-bind:class="{ 'data-view': viewType === 'data' }">
        <header v-if="viewType === 'project'">Notes & Comments</header>

        <SectionToggle leftLabel="Log an event" rightLabel="Comment" @toggle="onSectionToggle" :default="logMode === 'event' ? 'left' : 'right'" v-if="viewType === 'data'">
            <template #left>
                <div class="event-sensor-selector">
                    <label for="allProjectRadio">
                        <div class="event-sensor-radio">
                            <input type="radio" id="allProjectRadio" name="eventLevel" checked />
                            <span class="radio-label">All Project Sensors</span>
                            <p>People will see this event when viewing data for any stations that belong to these projects</p>
                        </div>
                    </label>
                    <label for="allSensorsRadio">
                        <div class="event-sensor-radio">
                            <input type="radio" id="allSensorsRadio" name="eventLevel" />
                            <span class="radio-label">Just These Sensors</span>
                            <p>People will see this event only when viewing data for these stations</p>
                        </div>
                    </label>
                </div>
                <div class="new-comment" :class="{ 'align-center': !user }">
                    <UserPhoto :user="user"></UserPhoto>
                    <template v-if="user">
                        <div class="new-comment-wrap">
                            <Tiptap
                                v-model="newDataEvent.title"
                                placeholder="Event Title"
                                saveLabel="Post"
                                @save="saveDataEvent(newDataEvent)"
                            />
                            <Tiptap
                                v-model="newDataEvent.description"
                                placeholder="Event Description"
                                saveLabel="Post"
                                @save="saveDataEvent(newDataEvent)"
                            />
                        </div>
                    </template>
                    <template v-else>
                        <p class="need-login-msg" @click="test()">
                            {{ $tc("comments.loginToComment.part1") }}
                            <router-link
                                :to="{ name: 'login', query: { after: $route.path, params: JSON.stringify($route.query) } }"
                                class="link"
                            >
                                {{ $tc("comments.loginToComment.part2") }}
                            </router-link>
                            {{ $tc("comments.loginToComment.part3") }}
                        </p>
                        <router-link
                            :to="{ name: 'login', query: { after: $route.path, params: JSON.stringify($route.query) } }"
                            class="button-submit"
                        >
                            {{ $t("login.loginButton") }}
                        </router-link>
                    </template>
                </div>
            </template>
            <template #right>
                <div class="new-comment" :class="{ 'align-center': !user }">
                    <UserPhoto :user="user"></UserPhoto>
                    <template v-if="user">
                        <div class="new-comment-wrap">
                            <Tiptap
                                v-model="newComment.body"
                                placeholder="Join the discussion!"
                                saveLabel="Post"
                                @save="save(newComment)"
                            />
                        </div>
                    </template>
                    <template v-else>
                        <p class="need-login-msg" @click="test()">
                            {{ $tc("comments.loginToComment.part1") }}
                            <router-link
                                :to="{ name: 'login', query: { after: $route.path, params: JSON.stringify($route.query) } }"
                                class="link"
                            >
                                {{ $tc("comments.loginToComment.part2") }}
                            </router-link>
                            {{ $tc("comments.loginToComment.part3") }}
                        </p>
                        <router-link
                            :to="{ name: 'login', query: { after: $route.path, params: JSON.stringify($route.query) } }"
                            class="button-submit"
                        >
                            {{ $t("login.loginButton") }}
                        </router-link>
                    </template>
                </div>
            </template>
        </SectionToggle>
        <!-- TODO: code repeated for project view; componentize -->
        <div class="new-comment" :class="{ 'align-center': !user }" v-if="viewType === 'project'">
            <UserPhoto :user="user"></UserPhoto>
            <template v-if="user">
                <div class="new-comment-wrap">
                    <Tiptap v-model="newComment.body" placeholder="Join the discussion!" saveLabel="Post" @save="save(newComment)" />
                </div>
            </template>
            <template v-else>
                <p class="need-login-msg" @click="test()">
                    {{ $tc("comments.loginToComment.part1") }}
                    <router-link :to="{ name: 'login', query: { after: $route.path, params: JSON.stringify($route.query) } }" class="link">
                        {{ $tc("comments.loginToComment.part2") }}
                    </router-link>
                    {{ $tc("comments.loginToComment.part3") }}
                </p>
                <router-link
                    :to="{ name: 'login', query: { after: $route.path, params: JSON.stringify($route.query) } }"
                    class="button-submit"
                >
                    {{ $t("login.loginButton") }}
                </router-link>
            </template>
        </div>

        <div v-if="errorMessage" class="error">{{ errorMessage }}</div>

        <div v-if="!isLoading && posts.length === 0" class="no-comments">There are no comments yet.</div>
        <div v-if="isLoading" class="no-comments">Loading comments...</div>

        <div class="list" v-if="postsAndEvents && postsAndEvents.length > 0">
            <div class="subheader">
                <span class="comments-counter" v-if="viewType === 'project'">{{ postsAndEvents.length }} comments</span>
                <header v-if="viewType === 'data'">Events & Comments</header>
            </div>
            <transition-group name="fade">
                <div
                    class="comment comment-first-level"
                    v-for="post in postsAndEvents"
                    v-bind:key="post.id"
                    v-bind:id="'comment-id-' + post.id"
                    :ref="post.id"
                >
                    <div class="comment-main" :style="!user ? { 'padding-bottom': '15px' } : {}">
                        <UserPhoto :user="post.author"></UserPhoto>
                        <div class="column-post">
                            <div class="post-header">
                                <span class="author">
                                    {{ post.author.name }}
                                </span>
                                <ListItemOptions
                                    v-if="user && (user.id === post.author.id || user.admin)"
                                    @listItemOptionClick="onListItemOptionClick($event, post)"
                                    :options="getCommentOptions(post)"
                                />
                                <span class="timestamp">{{ formatTimestamp(post.createdAt) }}</span>
                            </div>
                            <Tiptap
                                v-if="post.body"
                                v-model="post.body"
                                :readonly="post.readonly"
                                saveLabel="Save"
                                @save="saveEdit(post.id, post.body)"
                            />
                            <div v-else>
                                <h3>{{ post.title }}</h3>
                                {{ post.description }}
                            </div>
                        </div>
                    </div>
                    <div class="column">
                        <transition-group name="fade" class="comment-replies">
                            <div
                                class="comment"
                                v-for="reply in post.replies"
                                v-bind:key="reply.id"
                                v-bind:id="'comment-id-' + reply.id"
                                :ref="reply.id"
                            >
                                <div class="comment-main">
                                    <UserPhoto :user="reply.author"></UserPhoto>
                                    <div class="column-reply">
                                        <div class="post-header">
                                            <span class="author">
                                                {{ reply.author.name }}
                                            </span>
                                            <ListItemOptions
                                                v-if="user && (user.id === reply.author.id || user.admin)"
                                                @listItemOptionClick="onListItemOptionClick($event, reply)"
                                                :options="getCommentOptions(reply)"
                                            />
                                        </div>
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
                                <div class="new-comment-wrap">
                                    <UserPhoto :user="user"></UserPhoto>
                                    <Tiptap
                                        v-model="newReply.body"
                                        placeholder="Reply to comment"
                                        @save="save(newReply)"
                                        saveLabel="Post"
                                    />
                                </div>
                            </div>
                        </transition>

                        <div v-if="user" class="actions">
                            <button v-if="post.body" @click="addReply(post)">
                                <i class="icon icon-reply"></i>
                                Reply
                            </button>
                            <button v-if="viewType === 'data'" @click="viewDataClick(post)">
                                <i class="icon icon-view-data"></i>
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
import { NewComment, NewDataEvent } from "@/views/comments/model";
import { Comment, DataEvent, DiscussionBase } from "@/views/comments/model";
import { CurrentUser } from "@/api";
import { CommentsErrorsEnum } from "@/views/comments/model";
import ListItemOptions from "@/views/shared/ListItemOptions.vue";
import Tiptap from "@/views/shared/Tiptap.vue";
import SectionToggle from "@/views/shared/SectionToggle.vue";
import { Bookmark } from "@/views/viz/viz";
import { TimeRange } from "@/views/viz/viz/common";
import { ActionTypes } from "@/store";

export default Vue.extend({
    name: "Comments",
    components: {
        ...CommonComponents,
        ListItemOptions,
        Tiptap,
        SectionToggle,
    },
    props: {
        user: {
            type: Object as PropType<CurrentUser>,
            required: false,
        },
        parentData: {
            type: [Number, Object],
            required: true,
        },
    },
    data(): {
        posts: Comment[];
        dataEvents: DataEvent[];
        isLoading: boolean;
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
        newDataEvent: {
            projectId: number | null;
            bookmark: string | null;
            body: string | null;
            title: string | null;
        };
        errorMessage: string | null;
        logMode: string;
    } {
        return {
            posts: [],
            dataEvents: [],
            isLoading: false,
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
            newDataEvent: {
                projectId: typeof this.parentData === "number" ? this.parentData : null,
                bookmark: null,
                body: "",
                title: "",
            },
            errorMessage: null,
            logMode: "event",
        };
    },
    watch: {
        parentData(): Promise<void> {
            return this.getComments();
        },
    },
    mounted(): Promise<void> {
        this.$store.dispatch(ActionTypes.NEED_DATA_EVENTS, { bookmark: JSON.stringify(this.parentData) }).then(() => {
            this.dataEvents = [...this.$getters.dataEvents, ...this.dataEvents];
        });

        this.placeholder = this.getNewCommentPlaceholder();
        return this.getComments();
    },
    computed: {
        postsAndEvents(): DiscussionBase[] {
            return [...this.posts, ...this.dataEvents].sort(this.sortRecent);
        },
    },
    methods: {
        getNewCommentPlaceholder(): string {
            if (this.viewType === "project") {
                return "Comment on Project";
            } else {
                return "Write a comment about this Data View";
            }
        },
        async saveDataEvent(dataEvent: NewDataEvent): Promise<void> {
            this.errorMessage = null;

            if (this.viewType === "data") {
                const tb: Bookmark = this.parentData;
                dataEvent.bookmark = JSON.stringify(this.parentData);
            }

            const timeRange: TimeRange = this.parentData.allTimeRange;
            dataEvent.start = timeRange.start;
            dataEvent.end = timeRange.end;

            await this.$services.api
                .postDataEvent(dataEvent)
                .then((response: { event: DataEvent }) => {
                    // add data-event to the posts array
                    if (this.dataEvents) {
                        const de = new DataEvent(
                            response.event.id,
                            response.event.author,
                            response.event.bookmark,
                            response.event.createdAt,
                            response.event.updatedAt,
                            response.event.title,
                            response.event.description,
                            response.event.start,
                            response.event.end
                        );
                        this.dataEvents.unshift(de);
                        this.newDataEvent.title = "";
                        this.newDataEvent.description = "";

                        this.$store.dispatch(ActionTypes.NEW_DATA_EVENT, { dataEvent: de });
                    } else {
                        console.log(`posts is null`);
                    }
                })
                .catch((e) => {
                    console.error(e);
                    this.errorMessage = CommentsErrorsEnum.postComment;
                });
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
                .catch((e) => {
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
            this.isLoading = true;
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

                    this.highlightComment();
                })
                .catch(() => {
                    this.errorMessage = CommentsErrorsEnum.getComments;
                })
                .finally(() => {
                    this.isLoading = false;
                });
        },
        viewDataClick(post: Comment) {
            if (post.bookmark) {
                this.$emit("viewDataClicked", JSON.parse(post.bookmark));
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
        highlightComment() {
            this.$nextTick(() => {
                if (location.hash) {
                    const el = document.querySelector(location.hash);

                    if (el) {
                        el.scrollIntoView({ behavior: "smooth" });
                        el.classList.add("highlight");
                        setTimeout(() => {
                            el.classList.remove("highlight");
                        }, 5000);
                    }
                }
            });
        },
        onSectionToggle(evt) {
            if (evt === "right") {
                this.logMode = "comment";
            }
            if (evt === "left") {
                this.logMode = "event";
            }
        },
        sortRecent(a, b) {
            return b.createdAt - a.createdAt;
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

    @include bp-down($xs) {
        margin: 20px -10px 0;
        padding: 0 10px 30px 10px;
    }

    &.data-view {
        margin-top: 0;
        // padding-top: 45px;
        box-shadow: none;
        border: 0;
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

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
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

::v-deep .new-comment {
    @include flex(flex-end);
    padding: 22px 0;
    position: relative;

    @include bp-down($xs) {
        margin: 0 -10px;
        padding: 15px 10px 15px 10px;
    }

    @media screen and (max-width: 320px) {
        flex-wrap: wrap;
    }

    .container.data-view & {
        &:not(.reply) {
            background-color: rgba(#f4f5f7, 0.55);
            padding: 18px 23px 17px 15px;
        }
    }

    &.reply {
        padding: 0 0 0;
        margin: 10px 0 0 0;
        width: 100%;
    }

    &.align-center {
        align-items: center;
    }

    img {
        margin-top: 0 !important;
        width: 30px;
        height: 30px;
    }

    .button-submit {
        margin-left: auto;

        @media screen and (max-width: 320px) {
            width: 100%;
        }
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

        .new-comment-wrap {
            flex: 0 0 calc(100% - 65px);
            flex-direction: column;
            background-color: rgba(#f4f5f7, 0.55);

            .tiptap-container {
                margin-top: 10px;
                background-color: white;
            }
        }
    }

    &-wrap {
        display: flex;
        width: 100%;
        position: relative;
        background-color: #fff;
    }
}
// .ProseMirror p.is-editor-empty:first-child::before {
// 	color: #adb5bd;
// 	content: attr(data-placeholder);
// 	float: left;
// 	height: 0;
// 	pointer-events: none;
// }

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

    &.highlight > div {
        background-color: #a0dbe1;
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
        @include flex(center);
    }

    .icon {
        font-size: 11px;
        margin-right: 5px;
    }

    .icon-view-data {
        font-size: 14px;
        margin-right: 6px;
    }

    .icon-reply {
        margin-top: -2px;
    }
}

.timestamp {
    font-family: $font-family-light;
    flex-shrink: 0;
    margin-left: auto;
    line-height: 1.5;
    padding-left: 10px;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.25s ease-in-out;
}

.error {
    color: $color-danger;
    margin-bottom: 10px;
}

.column-reply,
.column-post {
    position: relative;
    display: flex;
    flex-direction: column;
    width: 100%;
}

.post-header {
    display: flex;
}

.need-login-msg {
    font-size: 16px;
    margin-left: 8px;
    margin-right: 10px;

    @include bp-down($xs) {
        margin-left: 0;
    }

    @media screen and (max-width: 320px) {
        flex: 0 0 calc(100% - 55px);
        margin-right: 0;
    }

    * {
        font-size: 16px;
    }
}

.button-submit {
    width: auto;
    margin-left: auto;
    padding: 0 40px;
}

.no-comments {
    margin-top: 20px;
}
.event-sensor-selector {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    margin-bottom: 15px;
}
.event-sensor-radio {
    max-width: 100vh;
    height: 115px;
    border: solid 1px #d8dce0;
    padding: 15px;
    padding-bottom: 10px;
    margin-left: 10px;
    border-radius: 3px;
    flex: 1;

    p {
        margin-left: 30px;
    }

    .radio-label {
        color: #2c3e50;
        font-size: 18px;
        font-weight: 900;
        margin-left: 10px;
    }
    input:checked {
        background-color: red;
    }
}

.event-sensor-radio > input:checked + div {
    /* (RADIO CHECKED) DIV STYLES */
    background-color: #ffd6bb;
    border: 1px solid #ff6600;
}
</style>
