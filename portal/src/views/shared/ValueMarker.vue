<template>
    <div class="marker-container" @click="onClick">
        <span class="value-label">{{value}}</span>
        <svg viewBox="0 0 36 36">
            <circle class="marker-circle" cx="18" cy="18" r="16" :fill="color"/>
        </svg>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "ValueMarker",
    props: {
        color: {
            type: String,
            default: "#fff",
        },
        value: {
            type: Number
        },
        id: {
            type: Number
        }
    },
    methods: {
        onClick () {
            this.$emit("marker-click", { id: this.id });
        }
    }
});
</script>

<style scoped>
.marker-circle {
    stroke: white;
    stroke-opacity: 0.5;
    stroke-width: 5;
}
.value-label {
    position: relative;
    left: 8px;
    top: 28px;
    text-align: center;
    color: white;
    font-family: "Avenir"; 
    font-weight: bold;

}
.marker-container {
    width: 32px;
    height: 32px;
    cursor: pointer;
}
</style>
