<template>
    <div class="section-toggle">
        <div class="toggle-controls">
            <a @click="toggleClickHandler($event, 'left')" :class="{ selected: matchSection('left') }">{{ leftLabel }}</a>
            <a @click="toggleClickHandler($event, 'right')" :class="{ selected: matchSection('right') }">{{ rightLabel }}</a>
        </div>
        <slot name="left" v-if="matchSection('left')"></slot>
        <slot name="right" v-if="matchSection('right')"></slot>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "SectionToggle",
    props: {
        leftLabel: {
            type: String,
            required: true
        },
        rightLabel: {
            type: String,
            required: true
        }

    },
    data(): {
        selectedSection: any;
    } {
        return {
            selectedSection: "left",
        };
    },
    methods: {
        toggleClickHandler(evt, section) {
            this.selectedSection = section
        },
        matchSection(section) {
            return this.selectedSection === section
        }
    }
});
</script>

<style lang="scss" scoped>
.toggle-controls {
    align-items: center;
    justify-content: center;
    display: flex;
    margin-bottom: 15px;
    
    .selected {
        background-color: #2c3e50;
        color: white;
    }
    a {
        padding: 10px;
        user-select: none;
        cursor: pointer;
        border-radius: 25px;
    }
}
</style>
