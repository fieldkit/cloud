<template>
    <span class="pagination-page" v-on:click="onClick">
        ●
        <!-- {{ pageNumber }} -->
    </span>
</template>

<script>
export default {
    name: "PaginationPage",
    components: {},
    props: {
        pageNumber: {
            type: Number,
            required: true,
        },
    },
    methods: {
        onClick() {
            this.$emit("loadPage", this.pageNumber);
        },
    },
};
</script>

<style scoped>
.pagination-page {
    font-size: 18px;
    color: #d8d8d8;
    margin: 0 2px;
    cursor: pointer;
}
.pagination-page.current {
    font-size: 22px;
    color: #2c3e50;
}
</style>
