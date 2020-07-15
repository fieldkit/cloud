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
    /*
    color: white;
    padding: 5px;
    margin-right: 5px;
	margin-bottom: 5px;
	*/
    background-color: #0a67aa;
    color: #fff;
    border-radius: 2px;
    display: flex;
    padding: 3px 5px;
    margin: 2px;
    font-size: 0.85em;
}
</style>
