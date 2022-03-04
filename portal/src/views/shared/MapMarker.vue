<template>
    <div>&nbsp;
        <svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
            <circle :fill="color" cx="100" cy="100" r="100"/>
            <text>{{value}}</text>
        </svg>
    </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
    name: "MapMarker",
    data(): {

    } {
        return {

        };
    },
    props: {
        value: {
            type: Number,
        },
        color: {
            type: String
        },
     }
});
</script>

<style lang="scss" scoped>

.map-marker {
    background-color: red;
    height: '25px';
    width: '25px';
    cursor: pointer;
}

</style>
