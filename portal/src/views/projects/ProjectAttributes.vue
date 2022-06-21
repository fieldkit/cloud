<template>
    <div class="notes-form">
        <div class="header">
            <div class="name">{{ $t("attributes.title") }}</div>
        </div>
        <div class="site-notes">
            <form id="form">
                <div class="note-editor" v-for="attribute in attributes" v-bind:key="attribute.attributeId">
                    <div class="title">{{ attribute.name }}</div>
                    <div class="field">
                        <template v-if="attribute.stringValue">
                            {{ attribute.stringValue }}
                        </template>
                        <template v-else>
                            <span class="no-data-yet">
                                {{ $t("attributes.noValueYet") }}
                            </span>
                        </template>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>

<script lang="ts">
import { ProjectAttribute } from "@/api/api";
import Vue, { PropType } from "vue";

export default Vue.extend({
    name: "ProjectAttributes",
    props: {
        attributes: {
            type: Array as PropType<ProjectAttribute[]>,
            required: true,
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/notes";

.title {
    font-size: 16px;
    font-weight: 500;

    @include bp-down($xs) {
        font-size: 14px;
    }
}
</style>
