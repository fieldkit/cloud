<template>
    <StandardLayout>
        <div class="features-view">
            <div>
                TsDB
                <button v-if="features.tsdb" @click="toggleBackend()">Disable</button>
                <button v-else @click="toggleBackend()">Enable</button>
            </div>
            <div>
                Tiny Charts
                <button v-if="features.tinyCharts" @click="toggleTinyCharts()">Disable</button>
                <button v-else @click="toggleTinyCharts()">Enable</button>
            </div>
            {{ features }}
        </div>
    </StandardLayout>
</template>

<script>
import StandardLayout from "./StandardLayout";
import { getFeaturesEnabled } from "@/utilities";

export default {
    name: "FeaturesView",
    components: {
        StandardLayout,
    },
    data: () => {
        return {
            features: getFeaturesEnabled(),
        };
    },
    computed: {
        /*
        user() {
            return this.$store.state.user.user;
        },
        */
    },
    methods: {
        toggleBackend() {
            if (this.features.tsdb) {
                delete window.localStorage["fk:backend"];
            } else {
                window.localStorage["fk:backend"] = "tsdb";
            }
            this.features = getFeaturesEnabled();
        },
        toggleTinyCharts() {
            if (this.features.tinyCharts) {
                delete window.localStorage["fk:tiny-charts"];
            } else {
                window.localStorage["fk:tiny-charts"] = "true";
            }
            this.features = getFeaturesEnabled();
        },
    },
};
</script>

<style scoped>
.features-view div {
    padding: 1em;
}
</style>
