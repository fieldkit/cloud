<template>
    <i class="icon-ellipsis options-trigger" @click="show($event)">
        <div class="options-btns">
            <button v-for="option in options" v-bind:key="option.event" @click="$emit('listItemOptionClick', option.event)">
                <span v-if="option.icon" class="options-icon">
                    <i class="icon" :class="option.icon"></i>
                </span>
                {{ option.label }}
            </button>
        </div>
    </i>
</template>

<script lang="ts">
export interface ListItemOption {
    label: string;
    event: string;
    icon?: string;
}

import Vue, { PropType } from "vue";

export default Vue.extend({
    name: "ListItemOptions",
    props: {
        options: {
            type: Array as PropType<ListItemOption[]>,
            required: true,
        },
    },
    methods: {
        show(event: MouseEvent): void {
            if (event.target) {
                const optionsMenu = (event.target as HTMLElement).firstElementChild;

                if (!optionsMenu) {
                    return;
                }

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
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/global";

.icon-ellipsis {
    display: block;
    cursor: pointer;

    &:after {
        @include flex(flex-end);
        content: "...";
        height: 7px;
        font-size: 32px;
        font-family: "icomoon";
    }
}

button {
    padding: 0;
    border: 0;
    outline: 0;
    box-shadow: none;
    cursor: pointer;
    background: transparent;
}

.options {
    &-trigger {
        // @include position(absolute, 0 -40px null null);
        padding-left: 10px;
        transition: opacity 0.33s;
    }

    &-btns {
        @include position(absolute, 0 null null null);
        opacity: 0;
        visibility: hidden;
        padding: 8px 12px;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);
        background: #fff;
        z-index: $z-index-top;
        transition: opacity 0.33s;
        min-width: 100px;
        cursor: initial;

        &.visible {
            opacity: 1;
            visibility: visible;
        }

        > * {
            display: flex;
            white-space: nowrap;
            font-family: var(--font-family-medium);
            cursor: pointer;
            text-align: left;
            width: 100%;
            padding: 7px 0;
        }
    }

    &-icon {
        width: 16px;
        margin-right: 8px;
        margin-top: -1px;
        @include flex(center, center);

        .icon {
            font-size: 16px;
        }
    }
}
</style>
