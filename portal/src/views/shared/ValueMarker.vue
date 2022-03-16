<template>
    <div class="marker-container" @click="onClick">
        <span class="value-label">{{ (value !== undefined) ? value : "?" | prettyNum }}</span>
        <svg viewBox="0 0 36 36">
            <circle class="marker-circle" cx="18" cy="18" r="16" :fill="color"/>
        </svg>
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import * as d3 from "d3";

export default Vue.extend({
    name: "ValueMarker",
    props: {
        color: {
            type: String,
            default: "#ccc",
        },
        value: {
            type: Number
        },
        id: {
            type: Number
        }
    },
    mounted () {
        if(this.value === undefined){
            this.color = "#ccc";
        }
    },
    methods: {
        onClick () {
            this.$emit("marker-click", { id: this.id });
        }
    },
    filters: {
        prettyNum: function (value) {
            if(!isNaN(value)){
                if(value % 1 === 0){
                    return value;
                }
                else return d3.format(".2s")(value);
            }
            else return value;
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
    left: 0px;
    top: 27px;
    text-align: center;
    color: white;
    font-weight: bold;

}
.marker-container {
    width: 32px;
    height: 32px;
    cursor: pointer;
    text-align: center;
    top: -18px;
}
</style>
