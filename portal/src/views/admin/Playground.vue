<template>
    <StandardLayout>
        <button v-on:click="onToggle">Modal</button>
        <StationPickerModal :stations="stations" @close="onToggle" v-if="modalOpen" />
    </StandardLayout>
</template>

<script lang="ts">
import Vue from "vue";
import StandardLayout from "../StandardLayout.vue";
import CommonComponents from "@/views/shared";

import StationPicker from "@/views/shared/StationPicker.vue";
import StationPickerModal from "@/views/shared/StationPickerModal.vue";

import { mapState, mapGetters } from "vuex";
import * as ActionTypes from "@/store/actions";
import { GlobalState } from "@/store/modules/global";

export default Vue.extend({
    name: "Playground",
    components: {
        StandardLayout,
        ...CommonComponents,
        StationPickerModal,
    },
    props: {},
    data: () => {
        return {
            modalOpen: false,
        };
    },
    computed: {
        ...mapState({
            stations: (s: GlobalState) => s.stations.user.stations,
        }),
    },
    methods: {
        onToggle() {
            this.modalOpen = !this.modalOpen;
        },
    },
});
</script>

<style scoped>
.container {
    display: flex;
    flex-direction: column;
    padding: 20px;
    max-width: 900px;
}
</style>
