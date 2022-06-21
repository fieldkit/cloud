<template>
    <div class="section-toggle">
        <hr class="toggle-hr">
        <div class="toggle-wrap">
            <div class="toggle-bg">
                <div class="toggle-controls">
                    <a @click="toggleClickHandler($event, 'left')" :class="{ selected: matchSection('left') }">{{ leftLabel }}</a>
                    <a @click="toggleClickHandler($event, 'right')" :class="{ selected: matchSection('right') }">{{ rightLabel }}</a>
                </div>
            </div>
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
        },
        default: {
            type: String,
            default: "left"
        }
    },
    data(): {
        selectedSection: any;
    } {
        return {
            selectedSection: this.default,
        };
    },
    methods: {
        toggleClickHandler(evt, section) {
            this.selectedSection = section;
            this.$emit("toggle", section);
        },
        matchSection(section) {
            return this.selectedSection === section;
        }
    }
});
</script>

<style lang="scss" scoped>
.toggle-wrap{
    display: flex;
    justify-content: center;
}
.toggle-bg{
    background-color: none;
    border: 1px solid #d8dce0;
    border-radius: 25px;
    align-items: center;
    justify-content: center;
    display: flex;
    margin-bottom: 15px;
    width: 18em;
    background-color: #fff;
    z-index: 100;
}
.toggle-hr {
    border-style: none;
    border-top: 1px solid #d8dce0;
    position: relative;
    top: 2em;
}
.toggle-controls {
    align-items: baseline;
    justify-content: space-between;
    display: flex;
    width: 100%;
    z-index: 10;

    .selected {
        background-color: #2c3e50;
        color: white;
    }
    a {
        padding: 10px 25px 10px 25px;
        user-select: none;
        cursor: pointer;
        border-radius: 25px;
    }
}
</style>
