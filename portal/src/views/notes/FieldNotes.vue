<template>
    <div>
        <header class="header">
            <div class="name">{{ $t("fieldNotes.title") }}</div>
            <div class="buttons" v-if="isAuthenticated">
                <button class="button">{{ $t("fieldNotes.btnExport") }}</button>
            </div>
        </header>

        <div class="new-note" :class="{ 'align-center': !user }">
            <UserPhoto :user="user"></UserPhoto>
            <template v-if="user">
                <div class="new-note-wrap">
                    <Tiptap v-model="newNote.body" placeholder="Join the discussion!" saveLabel="fieldNote" @save="save(newNote)" />
                </div>
            </template>
            <!--            <template v-else>
                <p class="need-login-msg">
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
            </template>-->
        </div>

        <div v-if="errorMessage" class="error">{{ errorMessage }}</div>

        <div v-if="!isLoading && fieldNotes.length === 0" class="no-comments">There are no comments yet.</div>
        <div v-if="isLoading" class="no-comments">Loading comments...</div>

        <div class="list" v-if="fieldNotes && fieldNotes.length > 0">
            <transition-group name="fade">
                <div
                    class="comment comment-first-level"
                    v-for="fieldNote in fieldNotes"
                    v-bind:key="fieldNote.id"
                    v-bind:id="'comment-id-' + fieldNote.id"
                    :ref="fieldNote.id"
                >
                    <div class="comment-main" :style="!user ? { 'padding-bottom': '15px' } : {}">
                        <UserPhoto :user="fieldNote.author"></UserPhoto>
                        <div class="column-fieldNote">
                            <div class="fieldNote-header">
                                <span class="author">
                                    {{ fieldNote.author.name }}
                                </span>
                                <ListItemOptions
                                    v-if="user && (user.id === fieldNote.author.id || user.admin)"
                                    @listItemOptionClick="onListItemOptionClick($event, fieldNote)"
                                    :options="getFieldNotesOptions(fieldNote)"
                                />
                                <span class="timestamp">{{ formatTimestamp(fieldNote.createdAt) }}</span>
                            </div>
                            <Tiptap
                                v-model="fieldNote.body"
                                :readonly="fieldNote.readonly"
                                saveLabel="Save"
                                @save="saveEdit(fieldNote.id, fieldNote.body)"
                            />
                        </div>
                    </div>
                    <div class="column">
                        <transition-group name="fade" class="comment-replies">
                            <div
                                class="comment"
                                v-for="reply in fieldNote.replies"
                                v-bind:key="reply.id"
                                v-bind:id="'comment-id-' + reply.id"
                                :ref="reply.id"
                            >
                                <div class="comment-main">
                                    <UserPhoto :user="reply.author"></UserPhoto>
                                    <div class="column-reply">
                                        <div class="fieldNote-header">
                                            <span class="author">
                                                {{ reply.author.name }}
                                            </span>
                                            <ListItemOptions
                                                v-if="user && (user.id === reply.author.id || user.admin)"
                                                @listItemOptionClick="onListItemOptionClick($event, reply)"
                                                :options="getFieldNotesOptions(reply)"
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
                            <div class="new-comment reply" v-if="newReply && newReply.threadId === fieldNote.id">
                                <div class="new-comment-wrap">
                                    <UserPhoto :user="user"></UserPhoto>
                                    <Tiptap
                                        v-model="newReply.body"
                                        placeholder="Reply to comment"
                                        @save="save(newReply)"
                                        saveLabel="fieldNote"
                                    />
                                </div>
                            </div>
                        </transition>

                        <div class="actions">
                            <button v-if="user" @click="addReply(fieldNote)">
                                <i class="icon icon-reply"></i>
                                {{ $t("comments.actions.reply") }}
                            </button>
                            <button v-if="viewType === 'data'" @click="viewDataClick(fieldNote)">
                                <i class="icon icon-view-data"></i>
                                {{ $t("comments.actions.viewData") }}
                            </button>
                        </div>
                    </div>
                </div>
            </transition-group>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { mapGetters, mapState } from "vuex";
import { GlobalState } from "@/store";
import {Comment, CommentsErrorsEnum, FiledNotesErrorsEnum, NewComment} from "@/views/comments/model";
import Tiptap from "@/views/shared/Tiptap.vue";

interface FieldNote {
    id: number;
    author: { id: number; name: string; photo: object };
    body: string;
    createdAt: number;
    updatedAt: number;
}

export default Vue.extend({
    name: "FieldNotes",
    components: {
        ...CommonComponents,
        Tiptap,
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
        }),
    },
    data(): {
        fieldNotes: FieldNote[];
        isLoading: boolean;
        placeholder: string | null;
        newNote: {
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
            fieldNotes: [
                {
                    author: "Christine",
                    body:
                        "In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms.",
                    date: "9/17/22 11:02",
                },
                {
                    author: "Christine",
                    body:
                        "In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms.",
                    date: "9/17/22 11:02",
                },
            ],
            isLoading: false,
            placeholder: null,
            newNote: {
                projectId: null,
                bookmark: null,
                body: "",
            },
            newReply: {
                projectId: null,
                bookmark: null,
                body: "",
                threadId: null,
            },
            errorMessage: null,
        };
    },
    methods: {
        async save(note: any): Promise<void> {
            this.errorMessage = null;
            console.log("new note radoi", note);
            // TODO: Achil api

           /* await this.$services.api
                .postComment(note)
                .then((response: { fieldNote: FiledN }) => {
                    this.newNote.body = "";
                    // add the comment to the replies array
                })
                .catch((e) => {
                    this.errorMessage = FiledNotesErrorsEnum.getFiledNotes;
                });*/
        },
        getFieldNotesOptions(post: FieldNote): { label: string; event: string }[] {
            if (this.user && this.user.id === post.author.id) {
                return [
                    {
                        label: "Edit",
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
                    label: "Delete",
                    event: "delete-comment",
                },
            ];
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/global";
@import "../../scss/notes";

.new-note {
    display: flex;
    align-items: center;
    padding: 25px 0;
    position: relative;

    @include bp-down($xs) {
        padding: 15px 10px;
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

    .button-submit {
        margin-left: auto;

        @media screen and (max-width: 320px) {
            width: 100%;
        }
    }

    &-wrap {
        display: flex;
        width: 100%;
        position: relative;
        background-color: #fff;
    }
}

::v-deep .default-user-icon {
    width: 36px;
    height: 36px;
    margin: 0 15px 0 0;
}

.new-comment-wrap {
    width: 100%;
}

::v-deep .tiptap-container {
    box-sizing: border-box;
}
</style>
