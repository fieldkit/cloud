<template>
    <i class="icon-ellipsis options-trigger" @click="show($event)">
        <div class="options-btns">
            <button v-for="option in options" v-bind:key="option.event" @click="$emit('listItemOptionClick', option.event)">
                {{ option.label }}
            </button>
        </div>
    </i>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "ListItemOptions",
    props: {
        options: {
            type: Array,
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
        height: 17px;
        font-size: 32px;
        font-family: var(--font-family-bold);
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
        padding: 10px;
        box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.5);
        background: #fff;
        z-index: $z-index-top;
        transition: opacity 0.33s;
        min-width: 100px;

        &.visible {
            opacity: 1;
            visibility: visible;
        }

        > * {
            display: block;
            white-space: nowrap;
            font-family: var(--font-family-bold);
            cursor: pointer;
            text-align: left;
            width: 100%;

            &:not(:last-of-type) {
                margin-bottom: 10px;
            }
        }
    }
}
</style>
