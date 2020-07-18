<template>
    <div class="pagination">
        <div class="button" v-on:click="onPrevious" v-bind:class="{ enabled: canPagePrevious }">◀</div>
        <div class="pages">
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
        <div class="button" v-on:click="onNext" v-bind:class="{ enabled: canPageNext }">▶</div>
    </div>
</template>

<script lang="ts">
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

<style scoped>
.pagination .button {
    font-size: 22px;
    color: #d8d8d8;
    justify-content: center;
    align-items: center;
}
.pagination .button:first-child {
    margin-right: 10px;
}
.pagination .button:last-child {
    margin-left: 10px;
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
    font-weight: 900;
    color: #2c3e50;
}
</style>
