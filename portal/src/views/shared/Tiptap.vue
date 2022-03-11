<template>
    <div :class="'tiptap-container' + (readonly ? ' tiptap-reading' : ' tiptap-editing')">
        <div class="tiptap-row">
            <div class="tiptap-main"><editor-content :editor="editor" /></div>
            <div class="tiptap-side" v-if="!readonly && !empty">
                <button type="submit" @click="onSave">{{ saveLabel }}</button>
            </div>
        </div>
        <div
            v-if="characters && editor"
            :class="{ 'character-count': true, 'character-count--warning': editor.getCharacterCount() === limit }"
        >
            <svg height="20" width="20" viewBox="0 0 20 20" class="character-count__graph">
                <circle r="10" cx="10" cy="10" fill="#e9ecef" />
                <circle
                    r="5"
                    cx="10"
                    cy="10"
                    fill="transparent"
                    stroke="currentColor"
                    stroke-width="10"
                    :stroke-dasharray="`calc(${percentage} * 31.4 / 100) 31.4`"
                    transform="rotate(-90) translate(-20)"
                />
                <circle r="6" cx="10" cy="10" fill="white" />
            </svg>
            <div class="character-count__text">{{ editor.getCharacterCount() }}/{{ limit }} characters</div>
        </div>
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
import CharacterCount from "@tiptap/extension-character-count";
import Mention from "@tiptap/extension-mention";
import Placeholder from "@tiptap/extension-placeholder";
import MentionList from "../comments/MentionList.vue";
import tippy from "tippy.js";

export default Vue.extend({
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
        limit: number;
    } {
        return {
            editor: null,
            limit: 280,
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
                this.editor.commands.setContent(value);
            }
        },
    },
    computed: {
        percentage(): number {
            if (this.editor) {
                return Math.round((100 / this.limit) * this.editor.getCharacterCount());
            }
            return 0;
        },
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
                Placeholder,
                ModifyEnter,
                CharacterCount.configure({
                    limit: this.limit,
                }),
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

                                    popup = tippy("body", {
                                        getReferenceClientRect: props.clientRect,
                                        appendTo: () => document.body,
                                        content: component.element,
                                        showOnCreate: true,
                                        interactive: true,
                                        trigger: "manual",
                                        placement: "bottom-start",
                                    });
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
            }
        },
    },
});
</script>
<style lang="scss">
.tiptap-container {
    width: 100%;
}

.tiptap-editing {
    border: 1px solid #bfbfbf;
    padding: 0.3em 13px 0.3em 13px;
}

.tiptap-reading {
    border: 1px solid transparent;
    /*padding: 0.3em 1.4em 0.3em 0.8em;*/

    p {
        margin: 0em;
    }
}

/* Basic editor styles */
.ProseMirror-focused {
    outline: none;
}

.ProseMirror {
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
        word-break: break-all;
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
    align-items: center;
    justify-content: space-between;

    .tiptap-main {
        width: 100%;
    }

    .tiptap-side {
        flex-shrink: 0;

        button {
            background-color: transparent;
            border: 0;
            font-weight: 900;
            font-size: 14px;
        }
    }
}
</style>
