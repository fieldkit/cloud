<template>
    <div class="pagination-container">
        <Button class="pagination-btn" :disabled="previousDisabled" v-on:click="previousPage">
            ◀︎
        </Button>
        <PaginationPage
            v-for="page in pages"
            :class="{ current: page === currentPage }"
            :key="page"
            :pageNumber="page"
            @loadPage="onLoadPage"
        />
        <Button class="pagination-btn" :disabled="nextDisabled" v-on:click="nextPage">
            ▶︎
        </Button>
    </div>
</template>

<script>
import PaginationPage from "../components/PaginationPage";

export default {
    name: "Pagination",
    components: {
        PaginationPage,
    },
    props: {
        currentPage: {
            type: Number,
            required: true,
        },
        pageCount: {
            type: Number,
            required: true,
        },
        numVisiblePages: {
            type: Number,
            default: 5,
        },
    },
    computed: {
        previousDisabled() {
            return this.currentPage === 1;
        },
        nextDisabled() {
            return this.currentPage === this.pageCount;
        },
        pages() {
            let pages = [];
            for (var i = 1; i < this.pageCount + 1; i++) {
                pages.push(i);
            }
            return pages;
        },
    },
    methods: {
        nextPage() {
            this.$emit("nextPage");
        },
        previousPage() {
            this.$emit("previousPage");
        },
        onLoadPage(value) {
            this.$emit("loadPage", value);
        },
    },
};
</script>

<style scoped>
.pagination-btn {
    color: #2c3e50;
    font-size: 14px;
    border: none;
    background: none;
}
.pagination-btn:disabled {
    color: #d8d8d8;
}
</style>
