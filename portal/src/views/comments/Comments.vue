<template>
    <section class="container" v-bind:class="{ 'data-view': viewType === 'data' }">
        <header v-if="viewType === 'project'">{{ $tc("comments.projectHeader") }}</header>

        <SectionToggle
            class="comment-toggle"
            :leftLabel="$tc('comments.sectionToggle.leftLabel')"
            :rightLabel="$tc('comments.sectionToggle.rightLabel')"
            @toggle="onSectionToggle"
            :default="logMode === 'comment' ? 'left' : 'right'"
            v-if="viewType === 'data'"
            :showToggle="(user && user.admin) || (projectUser && projectUser.user && projectUser.role === 'Administrator')"
        >
            <template #left>
                <div class="new-comment" :class="{ 'align-center': !user }">
                    <UserPhoto :user="user"></UserPhoto>
                    <template v-if="user">
                        <div class="new-comment-wrap">
                            <Tiptap
                                v-model="newComment.body"
                                :placeholder="$tc('comments.commentForm.placeholder')"
                                :saveLabel="$tc('comments.commentForm.saveLabel')"
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
            <template #right>
                <div class="event-level-selector">
                    <label for="allProjectRadio">
                        <div class="event-level-radio">
                            <input
                                type="radio"
                                id="allProjectRadio"
                                name="eventLevel"
                                v-model="newDataEvent.allProjectSensors"
                                :value="true"
                                :checked="newDataEvent.allProjectSensors"
                            />
                            <span class="radio-label">
                                {{ $tc(interpolatePartner("comments.eventTypeSelector.allProjectSensors.radioLabel.")) }}
                            </span>

                            <InfoTooltip
                                :message="$tc(interpolatePartner('comments.eventTypeSelector.allProjectSensors.description.'))"
                            ></InfoTooltip>

                            <p>
                                {{ $tc(interpolatePartner("comments.eventTypeSelector.allProjectSensors.description.")) }}
                            </p>
                        </div>
                    </label>
                    <label for="allSensorsRadio">
                        <div class="event-level-radio">
                            <input
                                type="radio"
                                id="allSensorsRadio"
                                name="eventLevel"
                                v-model="newDataEvent.allProjectSensors"
                                :value="false"
                                :checked="!newDataEvent.allProjectSensors"
                            />
                            <span class="radio-label">
                                {{ $tc(interpolatePartner("comments.eventTypeSelector.justTheseSensors.radioLabel.")) }}
                            </span>

                            <InfoTooltip
                                :message="$tc(interpolatePartner('comments.eventTypeSelector.justTheseSensors.description.'))"
                            ></InfoTooltip>

                            <p>
                                {{ $tc(interpolatePartner("comments.eventTypeSelector.justTheseSensors.description.")) }}
                            </p>
                        </div>
                    </label>
                </div>
                <div class="new-comment" :class="{ 'align-center': !user }">
                    <UserPhoto :user="user"></UserPhoto>
                    <template v-if="user">
                        <div class="new-comment-wrap">
                            <Tiptap
                                v-model="newDataEvent.title"
                                :placeholder="$tc('comments.eventForm.title.placeholder')"
                                :saveLabel="$tc('comments.eventForm.title.saveLabel')"
                                @save="saveDataEvent(newDataEvent)"
                            />
                            <Tiptap
                                v-model="newDataEvent.description"
                                :placeholder="$tc('comments.eventForm.description.placeholder')"
                                :saveLabel="$tc('comments.eventForm.description.saveLabel')"
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
        </SectionToggle>
        <!-- TODO: code repeated for project view; componentize -->
        <div class="new-comment" :class="{ 'align-center': !user }" v-if="viewType === 'project'">
            <UserPhoto :user="user"></UserPhoto>
            <template v-if="user">
                <div class="new-comment-wrap">
                    <Tiptap
                        v-model="newComment.body"
                        :placeholder="$tc('comments.commentForm.placeholder')"
                        :saveLabel="$tc('comments.commentForm.saveLabel')"
                        @save="save(newComment)"
                    />
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

        <div v-if="!isLoading && posts.length === 0" class="no-comments">
            {{ viewType === "data" ? $tc("comments.noEventsComments") : $tc("comments.noComments") }}
        </div>
        <div v-if="isLoading" class="no-comments">
            {{ viewType === "data" ? $tc("comments.loadingEventsComments") : $tc("comments.loadingComments") }}
        </div>
        <div class="list" v-if="postsAndEvents && postsAndEvents.length > 0">
            <div class="subheader">
                <span class="comments-counter" v-if="viewType === 'project'">
                    {{ postsAndEvents.length }} {{ $tc("comments.comments") }}
                </span>
                <header v-if="viewType === 'data'">{{ $tc("comments.dataHeader") }}</header>
            </div>
            <transition-group name="fade">
                <div
                    class="comment comment-first-level"
                    v-for="item in postsAndEvents"
                    v-bind:key="(item.type === 'comment' ? 'c' : 'e') + item.id"
                    v-bind:id="(item.body ? 'comment-id-' : 'event-id-') + item.id"
                    :ref="item.id"
                >
                    <div class="comment-main" :style="!user ? { 'padding-bottom': '15px' } : {}">
                        <UserPhoto :user="item.author"></UserPhoto>
                        <div class="column-post">
                            <div class="post-header">
                                <span class="author">
                                    {{ item.author.name }}
                                </span>
                                <span v-if="item.body" class="icon icon-comment"></span>
                                <span v-else class="icon icon-flag"></span>
                                <ListItemOptions
                                    v-if="user && (user.id === item.author.id || user.admin)"
                                    @listItemOptionClick="onListItemOptionClick($event, item)"
                                    :options="getCommentOptions(item)"
                                />
                                <span class="timestamp">{{ formatTimestamp(item.createdAt) }}</span>
                            </div>
                            <Tiptap
                                v-if="item.body"
                                v-model="item.body"
                                :readonly="item.readonly"
                                saveLabel="Save"
                                @save="saveEdit(item.id, item.body)"
                            />
                            <div v-else class="edit-event">
                                <Tiptap
                                    v-model="item.title"
                                    :readonly="item.readonly"
                                    :placeholder="$tc('comments.eventForm.title.placeholder')"
                                    :saveLabel="$tc('comments.eventForm.title.saveLabel')"
                                    @save="saveEditDataEvent(item)"
                                />
                                <div class="event-range">{{ item.start | prettyDateTime }} - {{ item.end | prettyDateTime }}</div>
                                <Tiptap
                                    v-model="item.description"
                                    :readonly="item.readonly"
                                    :placeholder="$tc('comments.eventForm.description.placeholder')"
                                    :saveLabel="$tc('comments.eventForm.description.saveLabel')"
                                    @save="saveEditDataEvent(item)"
                                />
                            </div>
                        </div>
                    </div>
                    <div class="column">
                        <transition-group name="fade" class="comment-replies">
                            <div
                                class="comment"
                                v-for="reply in item.replies"
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
                                            :saveLabel="$tc('comments.reply.saveLabel')"
                                            @save="saveEdit(reply.id, reply.body)"
                                        />
                                    </div>
                                </div>
                            </div>
                        </transition-group>

                        <transition name="fade">
                            <div class="new-comment reply" v-if="newReply && newReply.threadId === item.id">
                                <div class="new-comment-wrap">
                                    <UserPhoto :user="user"></UserPhoto>
                                    <Tiptap
                                        v-model="newReply.body"
                                        :placeholder="$tc('comments.reply.placeholder')"
                                        :saveLabel="$tc('comments.reply.saveLabel')"
                                        @save="save(newReply)"
                                    />
                                </div>
                            </div>
                        </transition>

                        <div class="actions">
                            <button v-if="user && item.body" @click="addReply(item)">
                                <i class="icon icon-reply"></i>
                                {{ $t("comments.actions.reply") }}
                            </button>
                            <button v-if="viewType === 'data'" @click="viewDataClick(item)">
                                <i class="icon icon-view-data"></i>
                                {{ $t("comments.actions.viewData") }}
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
import { DataEventsErrorsEnum, NewComment, NewDataEvent } from "@/views/comments/model";
import { Comment, DataEvent, DiscussionBase } from "@/views/comments/model";
import { CurrentUser, ProjectUser } from "@/api";
import { CommentsErrorsEnum } from "@/views/comments/model";
import ListItemOptions from "@/views/shared/ListItemOptions.vue";
import Tiptap from "@/views/shared/Tiptap.vue";
import { deserializeBookmark } from "../viz/viz";
import SectionToggle from "@/views/shared/SectionToggle.vue";
import { Bookmark } from "@/views/viz/viz";
import { TimeRange } from "@/views/viz/viz/common";
import { ActionTypes, DisplayProject } from "@/store";
import { interpolatePartner, isCustomisationEnabled } from "@/views/shared/partners";
import InfoTooltip from "@/views/shared/InfoTooltip.vue";

export default Vue.extend({
    name: "Comments",
    components: {
        ...CommonComponents,
        ListItemOptions,
        Tiptap,
        SectionToggle,
        InfoTooltip,
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
            allProjectSensors: boolean;
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
                allProjectSensors: true,
                bookmark: null,
                body: "",
                title: "",
            },
            errorMessage: null,
            logMode: "comment",
        };
    },
    watch: {
        async parentData(): Promise<void> {
            await this.getDataEvents();
            return this.getComments();
        },
        $route() {
            this.highlightComment();
        },
    },
    async mounted(): Promise<void> {
        const projectId = typeof this.parentData === "number" ? null : this.parentData?.p[0];

        if (projectId) {
            await this.$store.dispatch(ActionTypes.NEED_PROJECT, { id: projectId });
            await this.$getters.projectsById[projectId];
        }
        this.placeholder = this.getNewCommentPlaceholder();

        await this.getDataEvents();
        return this.getComments();
    },
    computed: {
        postsAndEvents(): DiscussionBase[] {
            return [...this.posts, ...this.dataEvents].sort(this.sortRecent);
        },
        projectUser(): ProjectUser | null {
            const projectId = typeof this.parentData === "number" ? null : this.parentData?.p[0];

            if (projectId) {
                const displayProject = this.$getters.projectsById[projectId];
                return displayProject?.users?.filter((user) => user.user.id === this.user?.id)[0];
            }

            return null;
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
                .then((response) => {
                    if (response) {
                        this.newDataEvent.title = "";
                        this.newDataEvent.description = "";

                        this.getDataEvents();
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
                this.$emit("viewDataClicked", deserializeBookmark(post.bookmark));
                window.scrollTo({ top: 0, left: 0, behavior: "smooth" });
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
        startEditing(item: Comment | DataEvent) {
            item.readonly = false;
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
        async getDataEvents(): Promise<void> {
            if (typeof this.parentData === "number") {
                this.dataEvents = [];
                return;
            }
            this.isLoading = true;
            await this.$store
                .dispatch(ActionTypes.NEED_DATA_EVENTS, { bookmark: JSON.stringify(this.parentData) })
                .catch(() => {
                    this.errorMessage = DataEventsErrorsEnum.getDataEvents;
                })
                .finally(() => {
                    this.isLoading = false;
                });
            const dataEvents = this.$getters.dataEvents;
            this.dataEvents = [];
            dataEvents.forEach((event) => {
                this.dataEvents.push(
                    new DataEvent(
                        event.id,
                        event.author,
                        event.bookmark,
                        event.createdAt,
                        event.updatedAt,
                        event.title ? JSON.parse(event.title) : event.title,
                        event.description ? JSON.parse(event.description) : event.description,
                        event.start,
                        event.end
                    )
                );
            });
        },
        saveEditDataEvent(dataEvent: DataEvent) {
            this.$services.api
                .updateDataEvent(dataEvent)
                .then((response) => {
                    if (response) {
                        this.newDataEvent.title = "";
                        this.newDataEvent.description = "";
                        this.getDataEvents();
                    }
                })
                .catch((e) => {
                    console.error(e);
                    this.errorMessage = DataEventsErrorsEnum.postDataEvent;
                });
        },
        deleteDataEvent(dataEventID: number) {
            this.$services.api
                .deleteDataEvent(dataEventID)
                .then((response) => {
                    if (response) {
                        this.getDataEvents();
                    } else {
                        this.errorMessage = DataEventsErrorsEnum.deleteDataEvent;
                    }
                })
                .catch(() => {
                    this.errorMessage = CommentsErrorsEnum.deleteComment;
                });
        },
        onListItemOptionClick(event: string, item: Comment | DataEvent): void {
            if (event === "edit-comment") {
                this.startEditing(item);
            }
            if (event === "delete-comment") {
                if (item.type === "comment") {
                    this.deleteComment(item.id);
                }
                if (item.type === "event") {
                    this.deleteDataEvent(item.id);
                }
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
            if (evt === "left") {
                this.logMode = "comment";
            }
            if (evt === "right") {
                this.logMode = "event";
            }
        },
        sortRecent(a, b) {
            return b.createdAt - a.createdAt;
        },
        interpolatePartner(baseString) {
            return interpolatePartner(baseString);
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
    padding: 0 0 30px 0;
    background: #fff;
    border-radius: 1px;
    border: 1px solid $color-border;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);

    @include bp-down($xs) {
        margin: 20px -10px 0;
        padding: 0 0 30px 0;
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
    padding: 13px 20px;
    border-bottom: 1px solid $color-border;
    font-size: 20px;
    font-weight: 500;

    @include bp-down($xs) {
        padding: 13px 10px;
    }

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
    padding: 15px 20px;

    @include bp-down($xs) {
        padding: 15px 10px;
    }

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
    @include flex(flex-start);
    padding: 22px 20px;
    position: relative;
    margin-left: 20px;
    margin-right: 20px;

    @include bp-down($xs) {
        margin: 0 -10px;
        padding: 15px 10px;
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
                background-color: white;

                &:nth-of-type(2) {
                    margin-top: 10px;
                }
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
    padding: 15px 20px 0 20px;
    position: relative;
    flex-wrap: wrap;

    @include bp-down($xs) {
        padding: 15px 10px 0 10px;
    }

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
        background-color: #f5fbfc;
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
    align-items: center;
    margin-bottom: 5px;

    .icon {
        font-size: 12px;
        margin-left: 5px;

        @include bp-up($md) {
            margin-top: -2px;
        }
    }
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
    margin-left: 20px;

    @include bp-down($xs) {
        margin-left: 10px;
    }
}
.event-level-selector {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    margin-bottom: 15px;

    .info {
        display: none;
    }

    @include bp-down($xs) {
        flex-direction: column;

        label {
            width: 100%;
        }

        .info {
            display: inline-block;
            float: right;
        }

        ::v-deep .info-content {
            right: 0;
        }
    }
}
.event-level-radio {
    width: 340px;
    height: 115px;
    border: solid 1px #d8dce0;
    padding: 15px;
    padding-bottom: 10px;
    margin-left: 10px;
    border-radius: 3px;
    flex: 1;

    @include bp-down($xs) {
        width: calc(100% - 20px);
        height: auto;
        margin-right: 10px;
        margin-bottom: 5px;
    }

    p {
        margin-left: 30px;

        @include bp-down($xs) {
            display: none;
        }
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
.comment-toggle {
    margin-top: 20px;
    //margin-left: 20px;
}

.edit-event {
    > * {
        margin-top: 10px;
    }

    .tiptap-container:first-child {
        font-weight: bold;
    }
}

.event-range {
    margin-top: 0px;
    font-size: 12px;
}

.icon-flag,
.icon-comment,
.icon-view-data {
    &::before {
        body.floodnet & {
            color: $color-floodnet-dark;
        }
    }
}
</style>
