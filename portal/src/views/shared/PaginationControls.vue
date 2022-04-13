<template>
    <div class="pagination" :class="{textual: textual}">
        <div class="button" v-on:click="onPrevious" v-bind:class="{ enabled: canPagePrevious }">◀</div>
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
        <div v-if="textual">{{ page + 1}} of {{ totalPages }}</div>
        <div class="button" v-on:click="onNext" v-bind:class="{ enabled: canPageNext }">▶</div>
    </div>
</template>

<script lang="ts">
import { findParentNodeClosestToPos } from "@tiptap/core";
import _ from "lodash";
import Vue, { PropType } from "vue";

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
            return this.page < this.totalPages - 1;
        },
        canPagePrevious(this: any) {
            return this.page > 0;
        },
    },
    methods: {
        onPrevious() {
            if (this.page > 0) {
                this.$emit("new-page", this.page - 1);
            }
        },
        onNext() {
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
.pagination .button {
    font-size: 22px;
    color: #d8d8d8;
    justify-content: center;
    align-items: center;

    .textual & {
        font-size: 12px;
    }
}
.button:first-child {
    margin-right: 10px;

    .textual &{
        margin-right: 5px;
        line-height: 10px;
    }
}
.button:last-child {
    margin-left: 10px;

    .textual &{
        margin-left: 5px;
    }
}
.pagination .textual {
    font-size: 12px;
    align-items: center;
}
.pagination .textual .button {
    padding: 0;
    border: 0;
    margin-bottom: 0;
    user-select: none;
    font-size: 12px;
    color: #d8d8d8;
    justify-content: center;
    align-items: center;
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
