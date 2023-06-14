<template>
    <div>
        <vue-confirm-dialog />

        <header class="header">
            <div class="name">{{ $t("fieldNotes.title") }}</div>
            <div class="buttons" v-if="isAuthenticated">
                <button class="button">{{ $t("fieldNotes.btnExport") }}</button>
            </div>
        </header>

        <div class="new-field-note" v-if="user">
            <UserPhoto :user="user"></UserPhoto>
            <template v-if="user">
                <div class="new-field-note-wrap">
                    <Tiptap v-model="newNoteText" placeholder="Join the discussion!" saveLabel="Save" @save="save()" />
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

        <div v-if="!isLoading && fieldNotes.length === 0" class="no-comments">{{ $tc("fieldNotes.noData") }}</div>
        <div v-if="isLoading" class="no-comments">{{ $tc("fieldNotes.noNotes") }}</div>

        <div class="list" v-if="fieldNotes && fieldNotes.length > 0">
            <transition-group name="fade">
                <div
                    class="field-note"
                    v-for="fieldNote in fieldNotes"
                    v-bind:key="fieldNote.id"
                    v-bind:id="'field-note-id-' + fieldNote.id"
                    :ref="fieldNote.id"
                >
                    <UserPhoto :user="fieldNote.author"></UserPhoto>
                    <span class="author">
                        {{ fieldNote.author.name }}
                    </span>
                    <div class="timestamps">
                        <div class="timestamp-1">{{ formatTimestamp(fieldNote.createdAt) }}</div>
                        <div class="timestamp-2">
                            {{ $tc("fieldNotes.lastUpdated") }}
                            {{ formatTimestamp(fieldNote.updatedAt) }}
                        </div>
                    </div>
                    <Tiptap
                        v-model="fieldNote.body"
                        :readonly="!editingFieldNote || (editingFieldNote && editingFieldNote.id !== fieldNote.id)"
                        saveLabel="Save"
                        @save="saveEdit(fieldNote.id, fieldNote.body)"
                    />
                    <div v-if="!editingFieldNote || (editingFieldNote && editingFieldNote.id !== fieldNote.id)" class="actions">
                        <button v-if="user" @click="editFieldNote(fieldNote)">
                            <i class="icon icon-edit"></i>
                            {{ $t("fieldNotes.edit") }}
                        </button>
                        <button @click="deleteFieldNote(fieldNote)">
                            <i class="icon icon-trash"></i>
                            {{ $t("fieldNotes.delete") }}
                        </button>
                    </div>

                    <div v-if="editingFieldNote && editingFieldNote.id === fieldNote.id" class="update-actions">
                        <button v-if="user" @click="cancelEdit(fieldNote)">
                            {{ $t("fieldNotes.cancel") }}
                        </button>
                        <button @click="save(fieldNote)">
                            {{ $t("fieldNotes.update") }}
                        </button>
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
import { ActionTypes, AuthenticationRequiredError, GlobalState } from "@/store";
import Tiptap from "@/views/shared/Tiptap.vue";
import moment from "moment";
import _ from "lodash";
import { PortalStationFieldNotes, PortalStationNotes } from "@/views/notes/model";

interface GroupedFieldNotes {
    [date: string]: PortalStationFieldNotes[];
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
        fieldNotes(): PortalStationFieldNotes[] {
            return [
                {
                    id: 0,
                    author: { id: 1, name: "Christine", photo: {} },
                    body:
                        '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms."}]}]}',
                    createdAt: 1668775672000,
                    updatedAt: 1668775672000,
                },
                {
                    id: 1,
                    author: { id: 2, name: "Christine", photo: {} },
                    body:
                        '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms."}]}]}',
                    createdAt: 1768775672000,
                    updatedAt: 1868775672000,
                },
                {
                    id: 3,
                    author: { id: 1, name: "Christine", photo: {} },
                    body:
                        '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms."}]}]}',
                    createdAt: 1668775672050,
                    updatedAt: 1668775672050,
                },
                {
                    id: 3,
                    author: { id: 1, name: "Christine", photo: {} },
                    body:
                        '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms."}]}]}',
                    createdAt: 1665745349000,
                    updatedAt: 1665745349000,
                },
            ];
            //   return this.$state.notes.fieldNotes;
        },
        stationId(): number {
            return parseInt(this.$route.params.stationId, 10);
        },
    },
    data(): {
        groupedFieldNotes: GroupedFieldNotes[] | null;
        isLoading: boolean;
        placeholder: string | null;
        newNoteText: string | null;
        errorMessage: string | null;
        //    editingFieldNote: FieldNote | null;
    } {
        return {
            groupedFieldNotes: null,
            /*fieldNotes: [
                {
                    id: 0,
                    author: { id: 1, name: "Christine", photo: {} },
                    body:
                        '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms."}]}]}',
                    createdAt: 1668775672000,
                    updatedAt: 1668775672000,
                },
                {
                    id: 1,
                    author: { id: 2, name: "Christine", photo: {} },
                    body:
                        '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"In the last 24 hours there have been some connectivity issues, the station data was interrupted by severe storms."}]}]}',
                    createdAt: 1768775672000,
                    updatedAt: 1868775672000,
                },
            ],*/
            isLoading: false,
            placeholder: null,
            newNoteText: null,
            errorMessage: null,
            //    editingFieldNote: null,
        };
    },
    beforeMount(): void {
        this.$store.dispatch(ActionTypes.NEED_FIELD_NOTES, { id: this.stationId });
    },
    mounted() {
        this.groupByMonth();
    },
    methods: {
        async save(): Promise<void> {
            this.errorMessage = null;
            console.log("new note radoi", this.newNoteText);
            const note = {
                body: this.newNoteText,
                userId: this.user?.id,
                stationId: this.stationId,
            };
            await this.$store.dispatch(ActionTypes.ADD_FIELD_NOTE, { stationId: this.stationId, note });
        },
        saveEdit(): void {
            console.log("rr");
        },
        /*getFieldNotesOptions(post: FieldNote): { label: string; event: string }[] {
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
        editFieldNote(fieldNote: FieldNote) {
            this.editingFieldNote = JSON.parse(JSON.stringify(fieldNote));
        },*/
        deleteFieldNote() {
            console.log("r");
        },
        formatTimestamp(timestamp: number): string {
            return moment(timestamp).fromNow();
        },
        /*cancelEdit(fieldNote: FieldNote) {
            if (JSON.stringify(fieldNote.body) !== JSON.stringify(this.editingFieldNote?.body)) {
                console.log("radoi diff");
                this.$confirm({
                    message: this.$tc("fieldNotes.sureCancelEdit"),
                    button: {
                        no: this.$tc("no"),
                        yes: this.$tc("yes"),
                    },
                    callback: async (confirm) => {
                        if (confirm) {
                            this.editingFieldNote = null;
                        }
                    },
                });
                return;
            }

            this.editingFieldNote = null;
        },*/
        groupByMonth() {
            const groupedItems = _.groupBy(this.fieldNotes, (b) =>
                moment(b.createdAt)
                    .startOf("month")
                    .format("YYYY/MM")
            );

            console.log("radoi groupe", groupedItems);
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/global";
@import "../../scss/notes";

.new-field-note {
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

        ::v-deep .tiptap-container {
            flex: 0 0 100%;
            margin-left: 0;
        }
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
    flex: 0 0 calc(100% - 52px);
    margin-left: 52px;
}

.field-note {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    width: 100%;
    padding: 30px 0;
    border-top: solid 1px #d8dce0;
}
.author {
    font-size: 18px;
    color: #6a6d71;
}

.timestamps {
    margin-left: auto;
    margin-top: -10px;
    color: #818388;
    text-align: right;
}

.timestamp-1 {
    font-size: 14px;
}

.timestamp-2 {
    font-size: 10px;
    margin-top: 2px;
}

button {
    background-color: transparent;
    border: 0;
    font-size: 14px;
    padding: 0;
    margin-right: 12px;
}

.actions {
    padding-left: 52px;
    margin-top: 15px;
}

.update-actions {
    color: #818388;
    margin-left: auto;
    margin-right: -12px;
    margin-top: 12px;

    button:nth-of-type(2) {
        color: $color-fieldkit-primary;
        font-weight: 900;
    }
}
</style>
