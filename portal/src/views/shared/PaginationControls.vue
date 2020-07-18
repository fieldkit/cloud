<template>
    <div class="pagination">
        <div class="button" v-on:click="onPrevious" v-bind:class="{ enabled: canPagePrevious }">Previous</div>
        <div class="pages">
            <div
                v-for="page in pages"
                v-bind:key="page.number"
                class="page"
                v-bind:class="{ selected: page.selected }"
                v-on:click="onPage(page.number)"
            >
                {{ page.number + 1 }}
            </div>
        </div>
        <div class="button" v-on:click="onNext" v-bind:class="{ enabled: canPageNext }">Next</div>
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
    font-size: 14px;
    color: #000000;
    text-align: center;
    padding: 10px;
    background-color: #efefef;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
}
.pagination .button.enabled {
    font-size: 14px;
    font-weight: bold;
    color: #000000;
    text-align: center;
    padding: 10px;
    background-color: #efefef;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}
.footer .button {
    font-size: 18px;
    font-weight: bold;
    text-align: center;
    padding: 10px;
    margin: 20px 0 0 0px;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    color: #000000;
    background-color: #efefef;
}
.pagination {
    display: flex;
    justify-content: space-between;
}
.pages {
    display: flex;
    justify-content: space-evenly;
}
.pages .page {
    font-size: 14px;
    margin-left: 10px;
    margin-right: 10px;
    cursor: pointer;
}
.pages .page.selected {
    font-weight: 900;
}
</style>
