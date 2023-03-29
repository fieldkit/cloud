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
                    <Tiptap v-model="newNote.body" placeholder="Join the discussion!" saveLabel="Post" @save="save(newNote)" />
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
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { mapGetters, mapState } from "vuex";
import { GlobalState } from "@/store";
import { Comment, CommentsErrorsEnum, NewComment } from "@/views/comments/model";
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

            await this.$services.api
                .postComment(note)
                .then((response: { post: Comment }) => {
                    this.newNote.body = "";
                    // add the comment to the replies array
                })
                .catch((e) => {
                    this.errorMessage = CommentsErrorsEnum.postComment;
                });
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

::v-deep .new-note {
    @include flex(flex-end);
    padding: 22px 20px;
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
</style>
