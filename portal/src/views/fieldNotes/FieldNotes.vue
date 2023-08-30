<template>
    <div class="field-notes-wrap">
        <header class="header">
            <div class="name">{{ $t("fieldNotes.title") }}</div>
            <div class="buttons" v-if="isAuthenticated">
                <button class="button" @click="generatePDF">
                    <i class="icon icon-export"></i>
                    {{ $t("fieldNotes.btnExport") }}
                </button>
            </div>
        </header>

        <div class="new-field-note" v-if="user">
            <UserPhoto :user="user"></UserPhoto>
            <template v-if="user">
                <div class="new-field-note-wrap">
                    <Tiptap
                        @editor-focus="checkEditingFieldNote()"
                        v-model="newNoteText"
                        placeholder="Join the discussion!"
                        saveLabel="Save"
                        @save="save()"
                    />
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

        <div v-if="errorMessage" class="error">{{ errorMessage }}</div>

        <div v-if="!isLoading && !groupedFieldNotes">{{ $tc("fieldNotes.noData") }}</div>
        <div v-if="isLoading">{{ $tc("fieldNotes.loading") }}</div>

        <div class="field-note-list" v-if="groupedFieldNotes" ref="pdfContent">
            <div
                class="field-note-group hidden"
                :ref="'field-note-group-' + index"
                v-for="(monthItems, month, index) in groupedFieldNotes"
                :key="month"
            >
                <div class="month-row" @click="toggleFieldNoteGroup('field-note-group-' + index)">
                    <i class="icon icon-chevron-right"></i>
                    <div class="month-name">{{ getMonthName(month) }} Entries</div>
                    <div class="month-last-updated">Last updated: {{ getMonthLastUpdated(monthItems) }}</div>
                </div>

                <transition-group name="fade">
                    <div v-for="fieldNote in monthItems" :key="fieldNote.id" class="field-note">
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
                        <template>
                            <Tiptap
                                :ref="'note-ref-' + fieldNote.id"
                                :value="fieldNote.body"
                                :readonly="!editingFieldNote || editingFieldNote.id !== fieldNote.id"
                            />
                        </template>
                        <div v-if="!editingFieldNote || (editingFieldNote && editingFieldNote.id !== fieldNote.id)" class="actions">
                            <button v-if="user" @click="editFieldNote(fieldNote)">
                                <i class="icon icon-edit" v-if="canEdit(fieldNote)"></i>
                                {{ $t("fieldNotes.edit") }}
                            </button>
                            <button @click="deleteFieldNote(fieldNote.id)" v-if="canDelete(fieldNote)">
                                <i class="icon icon-trash"></i>
                                {{ $t("fieldNotes.delete") }}
                            </button>
                        </div>

                        <div v-if="editingFieldNote && editingFieldNote.id === fieldNote.id" class="update-actions">
                            <button v-if="user" @click="cancelEdit(fieldNote)">
                                {{ $t("fieldNotes.cancel") }}
                            </button>
                            <button @click="saveEdit(fieldNote)">
                                {{ $t("fieldNotes.update") }}
                            </button>
                        </div>
                    </div>
                </transition-group>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import CommonComponents from "@/views/shared";
import { mapGetters, mapState } from "vuex";
import { ActionTypes, GlobalState } from "@/store";
import Tiptap from "@/views/shared/Tiptap.vue";
import moment from "moment";
import _ from "lodash";
import { PortalStationFieldNotes } from "@/views/fieldNotes/model";
import { jsPDF } from "jspdf";
import { SnackbarStyle } from "@/store/modules/snackbar";

interface GroupedFieldNotes {
    [date: string]: PortalStationFieldNotes[];
}

export default Vue.extend({
    name: "FieldNotes",
    components: {
        ...CommonComponents,
        Tiptap,
    },
    props: {
        stationName: {
            type: String,
            required: true,
        },
    },
    computed: {
        ...mapGetters({ isAuthenticated: "isAuthenticated" }),
        ...mapState({
            user: (s: GlobalState) => s.user.user,
        }),
        fieldNotes(): PortalStationFieldNotes[] {
            return this.$state.fieldNotes.fieldNotes;
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
        editingFieldNote: PortalStationFieldNotes | null;
    } {
        return {
            groupedFieldNotes: null,
            isLoading: true,
            placeholder: null,
            newNoteText: null,
            errorMessage: null,
            editingFieldNote: null,
        };
    },
    beforeMount(): void {
        this.$store.dispatch(ActionTypes.NEED_FIELD_NOTES, { id: this.stationId });
    },
    watch: {
        fieldNotes() {
            this.groupByMonth();
        },
    },
    methods: {
        checkEditingFieldNote() {
            if (this.editingFieldNote) {
                this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("fieldNotes.finishEditingFirst"),
                    type: SnackbarStyle.fail,
                });
            }
        },
        async save(): Promise<void> {
            if (this.editingFieldNote) {
                await this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("fieldNotes.finishEditingFirst"),
                    type: SnackbarStyle.fail,
                });
                return;
            }
            this.errorMessage = null;
            const note = {
                body: this.newNoteText,
                userId: this.user?.id,
                stationId: this.stationId,
            };

            try {
                await this.$store.dispatch(ActionTypes.ADD_FIELD_NOTE, { stationId: this.stationId, note });
                await this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("fieldNotes.addSuccess"),
                    type: SnackbarStyle.success,
                });
                this.newNoteText = null;
                this.$nextTick(() => {
                    return this.$store.dispatch(ActionTypes.NEED_FIELD_NOTES, { id: this.stationId });
                });
            } catch (e) {
                return this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("somethingWentWrong"),
                    type: SnackbarStyle.fail,
                });
            }
        },
        async saveEdit(fieldNote: PortalStationFieldNotes): Promise<void> {
            const editorRef = this.$refs["note-ref-" + fieldNote.id];
            if (!editorRef || !editorRef[0] || !editorRef[0].editor) {
                console.error("Tip tap ref not found");
                return;
            }

            const editedText = editorRef[0].editor.getHTML();
            const payload = {
                id: this.editingFieldNote?.id,
                body: JSON.stringify(editedText),
                userId: this.user?.id,
                stationId: this.stationId,
            };

            // not the nicest way, but the tiptap editor does not revert itself to the initial body text, so need to reload everything
            if (editorRef[0].editor.isEmpty) {
                await this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("fieldNotes.emptyNoteCannotSave"),
                    type: SnackbarStyle.fail,
                });
                return;
            }

            try {
                this.editingFieldNote = null;
                await this.$store.dispatch(ActionTypes.UPDATE_FIELD_NOTE, { stationId: this.stationId, note: payload });
                await this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("fieldNotes.editSuccess"),
                    type: SnackbarStyle.success,
                });
            } catch (e) {
                return this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("somethingWentWrong"),
                    type: SnackbarStyle.fail,
                });
            }
        },
        editFieldNote(fieldNote: PortalStationFieldNotes) {
            if (this.editingFieldNote) {
                this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                    message: this.$tc("fieldNotes.finishEditingFirstBeforeNewOne"),
                    type: SnackbarStyle.fail,
                });
                return;
            }
            this.editingFieldNote = {
                ...JSON.parse(JSON.stringify(fieldNote)),
                ref: "note-ref-" + fieldNote.id,
            };
        },
        deleteFieldNote(noteId: number) {
            this.$confirm({
                message: this.$tc("fieldNotes.sureDelete"),
                button: {
                    no: this.$tc("no"),
                    yes: this.$tc("yes"),
                },
                callback: async (confirm) => {
                    if (confirm) {
                        try {
                            await this.$store.dispatch(ActionTypes.DELETE_FIELD_NOTE, { stationId: this.stationId, noteId });
                            await this.$store.dispatch(ActionTypes.NEED_FIELD_NOTES, { id: this.stationId });
                            await this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                                message: this.$tc("fieldNotes.deleteSuccess"),
                                type: SnackbarStyle.success,
                            });
                        } catch (e) {
                            await this.$store.dispatch(ActionTypes.SHOW_SNACKBAR, {
                                message: this.$tc("somethingWentWrong"),
                                type: SnackbarStyle.fail,
                            });
                        }
                    }
                },
            });
        },
        formatTimestamp(timestamp: number): string {
            return moment(timestamp).fromNow();
        },
        canEdit(fieldNote: PortalStationFieldNotes) {
            return fieldNote.author.id === this.user?.id;
        },
        // Duplicate functions that currently do the same thing since I expect the condition to change here
        canDelete(fieldNote: PortalStationFieldNotes) {
            return fieldNote.author.id === this.user?.id;
        },
        cancelEdit(fieldNote: PortalStationFieldNotes) {
            const editorRef = this.$refs["note-ref-" + fieldNote.id];
            if (!editorRef || !editorRef[0] || !editorRef[0].editor) {
                console.error("Tip tap ref not found");
                return;
            }
            const currentContent = JSON.stringify(editorRef[0].editor.getJSON());

            if (fieldNote.body !== currentContent) {
                this.$confirm({
                    message: this.$tc("fieldNotes.sureCancelEdit"),
                    button: {
                        no: this.$tc("no"),
                        yes: this.$tc("yes"),
                    },
                    callback: async (confirm) => {
                        if (confirm && this.editingFieldNote) {
                            editorRef[0].editor.commands.setContent(JSON.parse(fieldNote.body));
                            this.editingFieldNote = null;
                        }
                    },
                });
                return;
            }

            this.editingFieldNote = null;
        },
        groupByMonth() {
            if (this.fieldNotes.length > 0) {
                const groupedFieldNotes = _.groupBy(this.fieldNotes, (b) => moment(b.createdAt).startOf("month").format("YYYY/MM"));
                this.groupedFieldNotes = JSON.parse(JSON.stringify(groupedFieldNotes));
            } else {
                this.groupedFieldNotes = null;
            }
            this.$nextTick(() => {
                this.editingFieldNote = null;
                this.isLoading = false;
            });
        },

        getMonthName(month): string {
            return moment(month, "YYYY/MM").format("MMMM");
        },
        getMonthLastUpdated(items): string {
            const mostRecentUpdatedAt = items.reduce((maxUpdatedAt, obj) => {
                return Math.max(maxUpdatedAt, obj.updatedAt);
            }, 0);

            return this.formatTimestamp(mostRecentUpdatedAt);
        },
        toggleFieldNoteGroup(ref) {
            const elements = this.$refs[ref] as (HTMLElement | undefined)[];
            if (elements && elements.length > 0) {
                const element = elements[0];
                if (element instanceof HTMLElement) {
                    element.classList.toggle("hidden");
                }
            }
        },
        async generatePDF() {
            const doc = new jsPDF();
            const elementHTML = this.$refs.pdfContent as HTMLElement;

            const actionsEls = elementHTML.querySelectorAll(".actions");
            const groupEls = elementHTML.querySelectorAll(".field-note-group");

            // prepare HTML
            groupEls.forEach((el) => {
                el.classList.remove("hidden");
            });
            actionsEls.forEach((el) => {
                el.classList.add("hidden");
            });

            doc.html(elementHTML, {
                callback: (doc) => {
                    const fileName = this.stationName + " " + this.$tc("fieldNotes.title") + ".pdf";
                    doc.save(fileName);
                    // display back the actions els
                    actionsEls.forEach((el) => {
                        el.classList.remove("hidden");
                    });
                },
                margin: [10, 10, 10, 10],
                autoPaging: "text",
                x: 0,
                y: 0,
                width: 190, //target width in the PDF document
                windowWidth: 675, //window width in CSS pixels
            });
        },
    },
});
</script>

<style scoped lang="scss">
@import "src/scss/global";
@import "src/scss/notes";

.new-field-note {
    display: flex;
    align-items: center;
    padding: 25px 0;
    position: relative;

    @include bp-down($xs) {
        padding: 20px 0;
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

.button {
    @include bp-down($xs) {
        transform: none;
        margin-right: 0;
    }
}

::v-deep .default-user-icon {
    width: 36px;
    height: 36px;
    margin: 0 15px 0 0;
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
    flex-basis: 100%;
    width: 100%;
    padding: 30px 20px;
    margin-left: -20px;
    border-top: solid 1px #d8dce0;

    @include bp-down($xs) {
        padding: 20px 10px;
        margin-left: -10px;
    }

    .field-note-group.hidden & {
        display: none;
    }
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

    &.hidden {
        display: none;
    }
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

.month-row {
    flex: 0 0 100%;
    width: 100%;
    padding: 0 20px;
    margin-left: -20px;
    display: flex;
    align-items: center;
    height: 88px;
    border-top: solid 1px #d8dce0;
    cursor: pointer;

    @include bp-down($xs) {
        margin-left: -10px;
        padding: 0 10px;
        height: 68px;
    }

    .icon-chevron-right {
        font-size: 32px;
        transform: rotate(270deg);
        transition: all 250ms;
        margin-top: -5px;

        .field-note-group.hidden & {
            transform: rotate(90deg);
        }
    }

    .month-name {
        color: #6a6d71;
        font-size: 18px;
    }

    .month-last-updated {
        color: #818388;
        font-size: 14px;
        margin-left: auto;
    }

    &.expanded {
        .icon-chevron-right {
            transform: rotate(270deg);
        }
    }
}

.field-notes-wrap {
    @include bp-down($sm) {
        padding: 20px 10px 20px;
    }
}

.field-note-list {
    margin-bottom: -20px;

    ::v-deep .tiptap-side {
        display: none;
    }
}

.icon-export:before {
    color: var(--color-dark);
    margin-right: 8px;
}
</style>
