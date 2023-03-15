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
                <div class="new-comment-wrap">
                    <Tiptap v-model="newNote.body" placeholder="Join the discussion!" saveLabel="Post" @save="save(newNote)" />
                </div>
            </template>
            <template v-else>
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
            </template>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { mapGetters, mapState } from "vuex";
import { GlobalState } from "@/store";
import { Comment } from "@/views/comments/model";
import Tiptap from "@/views/shared/Tiptap.vue";

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
        posts: Comment[];
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
            posts: [],
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
        save(): void {
            console.log("saved");
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
</style>
