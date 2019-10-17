<template>
    <div id="notes-list-container" v-if="this.selectedSensor">
        <div id="notes-title">Notes & Comments</div>
        <p>No notes yet</p>
    </div>
</template>

<script>
export default {
    name: "NotesList",
    props: ["station", "selectedSensor"],
    data: () => {
        return {};
    },
    methods: {}
};
</script>

<style scoped>
#notes-list-container {
    float: left;
    width: 54%;
    margin: 10px 10px 10px 5px;
    padding: 10px;
    border: 1px solid rgb(230, 230, 230);
}
#notes-title {
    margin: 10px 0;
    padding-bottom: 10px;
    border-bottom: 1px solid rgb(200, 200, 200);
}
</style>
