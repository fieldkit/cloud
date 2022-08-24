<template>
    <div class="pagination" :class="{ textual: textual }">
        <div class="button prev" v-on:click="onPrevious" v-bind:class="{ enabled: canPagePrevious }">
            <span class="arrow"></span>
        </div>
        <div class="pages" v-if="!textual">
            <div
                v-for="page in pages"
                v-bind:key="page.number"
                class="page"
                v-bind:class="{ selected: page.selected }"
                v-on:click="onPage(page.number)"
            >
                ●
            </div>
        </div>
        <div v-if="textual">{{ page + 1 }} of {{ totalPages }}</div>
        <div class="button next" v-on:click="onNext" v-bind:class="{ enabled: canPageNext }">
            <span class="arrow"></span>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue from "vue";

export default Vue.extend({
    name: "PaginationControls",
    props: {
        page: {
            type: Number,
            required: true,
        },
        totalPages: {
            type: Number,
            required: true,
        },
        maximumPages: {
            type: Number,
            default: 7,
        },
        textual: {
            type: Boolean,
            default: false,
        },
        wrap: {
            type: Boolean,
            default: false,
        },
    },
    computed: {
        pages(this: any) {
            return _.range(0, this.totalPages).map((p) => {
                return {
                    selected: this.page == p,
                    number: p,
                };
            });
        },
        canPageNext(this: any) {
            if (this.wrap) {
                return true;
            } else {
                return this.page < this.totalPages - 1;
            }
        },
        canPagePrevious(this: any) {
            if (this.wrap) {
                return true;
            } else {
                return this.page > 0;
            }
        },
    },
    methods: {
        onPrevious() {
            if (this.wrap && this.page <= 0) {
                this.$emit("new-page", this.totalPages - 1);
            }
            if (this.page > 0) {
                this.$emit("new-page", this.page - 1);
            }
        },
        onNext() {
            if (this.wrap && this.page >= this.totalPages - 1) {
                this.$emit("new-page", 0);
            }
            if (this.page < this.totalPages - 1) {
                this.$emit("new-page", this.page + 1);
            }
        },
        onPage(this: any, page: number) {
            this.$emit("new-page", page);
        },
    },
});
</script>

<style lang="scss" scoped>
@import "src/scss/variables";
@import "src/scss/mixins";

.pagination .button {
    color: #d8d8d8;
    border: 0;
    margin-top: -2px;

    .textual & {
        font-size: 12px;
    }

    .arrow {
        border-style: solid;
        display: flex;
    }

    &.next {
        .arrow {
            border-width: 8px 0 8px 15px;
            border-color: transparent transparent transparent $color-dark;

            @include bp-down($sm) {
                border-width: 6px 0 6px 12px;
            }
        }
    }

    &.prev {
        .arrow {
            border-width: 8px 15px 8px 0;
            border-color: transparent $color-dark transparent transparent;

            @include bp-down($sm) {
                border-width: 6px 12px 6px 0;
            }
        }
    }
}
.button:first-child {
    margin-right: 10px;

    .textual & {
        margin-right: 7px;
        line-height: 10px;
    }
}
.button:last-child {
    margin-left: 10px;

    .textual & {
        margin-left: 7px;
    }
}
.pagination .textual {
    font-size: 12px;
    font-family: var(--font-family-bold);
    align-items: center;
}
.pagination .textual .button {
    padding: 3px;
    margin-bottom: 0;
    user-select: none;
}
.pagination .button.enabled {
    color: #2c3e50;
    cursor: pointer;
}
.pagination {
    display: flex;
    justify-content: center;
}
.pages {
    display: flex;
    justify-content: center;
}
.pages .page {
    font-size: 22px;
    margin-left: 6px;
    margin-right: 6px;
    cursor: pointer;
    color: #d8d8d8;
}
.pages .page.selected {
    color: #2c3e50;
}
</style>
