<template>
    <div :class="'tiptap-container' + (readonly ? ' tiptap-reading' : ' tiptap-editing')">
        <p v-if="placeholder && !readonly && editor && editor.isEmpty" class="placeholder">{{ placeholder }}</p>
        <div class="tiptap-row">
            <div ref="contentContainer" class="tiptap-main" :class="{ truncated: readonly }">
                <editor-content :editor="editor" />
            </div>
            <div class="tiptap-side" v-if="!readonly && !empty && showSaveButton">
                <button type="submit" @click="onSave">{{ saveLabel }}</button>
            </div>
        </div>
        <a v-if="seeMore" class="see-more" @click="toggleSeeMore(true)">{{ $t("seeMore") }}</a>
        <a v-if="seeLess" class="see-more" @click="toggleSeeMore(false)">{{ $t("seeLess") }}</a>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";
import { ResizeAuto } from "./ResizeAuto";
import { Editor, JSONContent, EditorContent, VueRenderer, Extension } from "@tiptap/vue-2";
import Document from "@tiptap/extension-document";
import Paragraph from "@tiptap/extension-paragraph";
import Text from "@tiptap/extension-text";
import Mention from "@tiptap/extension-mention";
import Placeholder from "@tiptap/extension-placeholder";
import MentionList from "../comments/MentionList.vue";
import tippy, { Props } from "tippy.js";

export default Vue.extend({
    name: "TipTap",
    components: {
        // ResizeAuto,
        EditorContent,
    },
    props: {
        value: {
            // type: Object as PropType<JSONContent | string>,
            required: true,
        },
        characters: {
            type: Boolean,
            default: false,
        },
        readonly: {
            type: Boolean,
            default: false,
        },
        showSaveButton: {
            type: Boolean,
            default: true,
        },
        placeholder: {
            type: String,
            default: "",
        },
        saveLabel: {
            type: String,
            default: "Post",
        },
    },
    data(): {
        editor: Editor | null;
        seeMore: boolean;
        seeLess: boolean;
    } {
        return {
            editor: null,
            seeMore: false,
            seeLess: false,
        };
    },
    watch: {
        readonly(value: boolean): void {
            if (this.editor) {
                this.editor.setOptions({ editable: !value });
            }
        },
        value(value: string): void {
            if (this.editor) {
                if (JSON.stringify(this.editor.getJSON()) === JSON.stringify(value)) return;
                this.editor.commands.setContent(value);
            }
        },
    },
    computed: {
        empty(): boolean {
            return this.editor == null || this.editor.getCharacterCount() == 0;
        },
    },
    mounted() {
        const services = this.$services;

        const changed = (value) => {
            this.$emit("input", value);
        };
        const saved = (editor, ...args) => {
            if (!editor.isEmpty) {
                this.$emit("save", editor.getJSON());
            }
        };

        const CustomNewLine = Extension.create({
            name: "newline",
            addCommands() {
                return {
                    addNewline: () => ({ state, dispatch }) => {
                        const { schema, tr } = state;
                        const paragraph = schema.nodes.paragraph;

                        const transaction = tr
                            .deleteSelection()
                            .replaceSelectionWith(paragraph.create(), true)
                            .scrollIntoView();
                        if (dispatch) dispatch(transaction);
                        return true;
                    },
                } as never;
            },
            addKeyboardShortcuts() {
                return {
                    "Shift-Enter": () => (this.editor.commands as any).addNewline(),
                };
            },
        });

        const ModifyEnter = Extension.create({
            addKeyboardShortcuts() {
                return {
                    Enter: (...args) => {
                        saved(this.editor);
                        // return true prevents default behaviour
                        return true;
                    },
                };
            },
        });

        function asContent(v: unknown): JSONContent | null {
            if (_.isString(v)) {
                if (v.length == 0) {
                    return null;
                }
                return JSON.parse(v);
            }
            return v as JSONContent;
        }

        this.editor = new Editor({
            editable: !this.readonly,
            content: asContent(this.value),
            extensions: [
                Document,
                Paragraph,
                Text,
                ModifyEnter,
                CustomNewLine,
                Mention.configure({
                    HTMLAttributes: {
                        class: "mention",
                    },
                    suggestion: {
                        items: (props: { query: string; editor: Editor }): any[] => {
                            if (props.query.length > 0) {
                                return (services.api.mentionables(props.query).then((mentionables) => {
                                    console.log("mentionables", mentionables);
                                    return mentionables.users;
                                }) as unknown) as any[];
                            } else {
                                return (Promise.resolve([]) as unknown) as any[];
                            }
                        },
                        render: () => {
                            let component;
                            let popup;

                            return {
                                onStart: (props) => {
                                    console.log("mentions-start", props);

                                    component = new VueRenderer(MentionList, {
                                        parent: this,
                                        propsData: props,
                                    });

                                    const newProps: Partial<Props> = {
                                        getReferenceClientRect: null,
                                        appendTo: () => document.body,
                                        content: component.element,
                                        showOnCreate: true,
                                        interactive: true,
                                        trigger: "manual",
                                        placement: "bottom-start",
                                    };

                                    const clientRect = props.clientRect;
                                    if (clientRect) {
                                        newProps.getReferenceClientRect = () => {
                                            const rect = clientRect();
                                            if (!rect) {
                                                throw new Error();
                                            }
                                            return rect;
                                        };
                                    }

                                    popup = tippy("body", newProps);
                                },
                                onUpdate(props) {
                                    console.log("mentions-update", props);

                                    component.updateProps(props);

                                    popup[0].setProps({
                                        getReferenceClientRect: props.clientRect,
                                    });
                                },
                                onKeyDown(props) {
                                    return component.ref?.onKeyDown(props);
                                },
                                onExit() {
                                    popup[0].destroy();
                                    component.destroy();
                                },
                            };
                        },
                    },
                }),
            ],
            onUpdate({ editor }) {
                changed(editor.getJSON());
            },
            onBlur({ editor }) {
                console.log("editor-blur");
            },
            onFocus({ editor }) {
                console.log("editor-focus");
            },
        });

        setTimeout(() => {
            this.truncate();
        });
    },
    beforeDestroy() {
        if (this.editor) {
            this.editor.destroy();
        }
    },
    methods: {
        onChange(...args) {
            console.log("on-change", args);
        },
        onSave() {
            if (this.editor && !this.editor.isEmpty) {
                this.$emit("save");
                this.editor.commands.clearContent();
            }
        },
        truncate() {
            if (!this.readonly) {
                return false;
            }

            const contentContainerEl = this.$refs.contentContainer as HTMLElement;

            if (contentContainerEl.offsetHeight < contentContainerEl.scrollHeight) {
                this.seeMore = true;
                contentContainerEl.classList.add("truncated");
            } else {
                this.seeMore = false;
                contentContainerEl.classList.remove("truncated");
            }
        },
        toggleSeeMore(show: boolean) {
            const contentContainerEl = this.$refs.contentContainer as HTMLElement;
            contentContainerEl.classList.toggle("truncated");
            this.seeMore = !show;
            this.seeLess = show;
        },
    },
});
</script>
<style lang="scss">
@import "../../scss/global";

.tiptap-container {
    width: 100%;
    text-align: justify;
    box-sizing: border-box;
    position: relative;
}

.tiptap-editing {
    border-radius: 2px;
    border: solid 1px #d8dce0;
    max-height: 70vh;
    overflow-y: auto;
    padding-left: 10px;
    padding-right: 80px;

    @include bp-down($sm) {
        max-height: 60vh;
    }

    @include bp-down($xs) {
        padding-right: 60px;
    }
}

/* Basic editor styles */
.ProseMirror-focused {
    outline: none;
}

.ProseMirror {
    @supports (-webkit-touch-callout: none) {
        // iOS only - prevent zoom in behaviour
        font-size: 16px;
    }

    > * + * {
        margin-top: 0.75em;
    }

    h1,
    h2,
    h3,
    h4,
    h5,
    h6 {
        line-height: 1.1;
    }

    p {
        word-break: normal;
        margin: 14px 0;
    }
}

/* Placeholder (at the top) */
.ProseMirror p.is-editor-empty:first-child::before {
    content: attr(data-placeholder);
    float: left;
    color: #ced4da;
    pointer-events: none;
    height: 0;
}

.tiptap-reading {
    border: 1px solid transparent;

    p {
        &:first-of-type {
            margin-top: 0;
        }

        &:last-of-type {
            margin-bottom: 0;
        }
    }
}

.mention {
    color: #a975ff;
    background-color: rgba(#a975ff, 0.1);
    border-radius: 0.3rem;
    padding: 0.1rem 0.3rem;
}

.character-count {
    margin-top: 1rem;
    display: flex;
    align-items: center;
    color: #68cef8;

    &--warning {
        color: #fb5151;
    }

    &__graph {
        margin-right: 0.5rem;
    }

    &__text {
        color: #868e96;
    }
}

.tiptap-row {
    display: flex;
    flex-direction: row;
    align-items: flex-end;
    justify-content: space-between;

    .tiptap-main {
        width: 100%;

        &.truncated {
            display: -webkit-box;
            -webkit-line-clamp: 8;
            -webkit-box-orient: vertical;
            overflow: hidden;
        }
    }

    .tiptap-side {
        flex-shrink: 0;
        padding: 12px 0;
        position: absolute;
        bottom: 2px;
        right: 25px;

        body.floodnet & {
            padding: 10px 0;
        }

        @include bp-down($sm) {
            right: 10px;
        }

        button {
            background-color: transparent;
            border: 0;
            font-weight: 900;
            font-size: 14px;
        }
    }
}

.see-more {
    font-size: 14px;
    line-height: 1;
    cursor: pointer;
    color: var(--color-primary);
    background-color: transparent;
    border: 0;
    font-family: var(--font-family-bold);
    white-space: nowrap;
    display: inline-block;
    padding: 10px 0;

    @at-root body.floodnet & {
        color: var(--color-dark);
    }
}

.placeholder {
    position: absolute;
    opacity: 0.25;
    top: 0;
    left: 10px;
}
</style>
