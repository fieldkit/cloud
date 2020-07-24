<template>
    <div>
        <div class="project-tags" v-if="parsed.length > 0">
            <span v-for="tag in parsed" class="project-tag" v-bind:key="tag.text">{{ tag.text }}</span>
        </div>
        <div v-else></div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { tryParseTags } from "@/utilities";

export default Vue.extend({
    name: "DisplayProjectTags",
    props: {
        tags: {
            type: String,
            required: true,
        },
    },
    computed: {
        parsed() {
            return tryParseTags(this.tags);
        },
    },
});
</script>

<style scoped>
.project-tags {
    display: flex;
    align-items: baseline;
}
.project-tag {
    background-color: #f4f5f7;
    color: #2c3e50;
    font-size: 13px;
    border-radius: 10px;
    display: flex;
    padding: 6px 10px;
    margin: 2px;
}
</style>
