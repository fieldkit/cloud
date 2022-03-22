<template>
    <div class="latest-events">
        <div class="event-row">
            <div class="event-container" v-for="annotation, idx in latestEvents" :key="idx">
                <h3>{{ annotation.body }}</h3>
                <div>{{ annotation.start | date }} &ndash; {{ annotation.end | date }}</div>
                <img alt="View Data" class="data-icon" src="@/assets/view-data.svg" width="18" />
                <router-link class="data-link" :to="{name: 'exploreBookmark', query: { bookmark: annotation.bookmark }}">View data</router-link>
                <UserPhoto :user="annotation.user"></UserPhoto>
                <div>{{ annotation.user.name }}</div>
            </div>
        </div>
        <pagination-controls :page="0" :total-pages="4"/>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";
import { Annotation} from "@/store";
// fixme: moment is huge, load module?
import moment from "moment";
import PaginationControls from "./PaginationControls.vue";

export default Vue.extend({
    name: "LatestEvents",
    components: {
        PaginationControls,
        ...CommonComponents,
    },
    props: {

    },
    data(): {
        viewType: string;
    } {
        return {
            viewType: "project"
        };
    },
    computed: {
        latestEvents(): Annotation[] {
            return this.$getters.getEvents.slice(0,4);
        },
    },
    filters: {
        date: (value) => {
            return moment(value).format('M/D/YYYY, h:mm:ss a');
        }
    }
});
</script>

<style lang="scss" scoped>
@import "../../scss/global";

.event-row {
    display: flex;
    flex: 50%;
    flex-wrap: wrap;
    gap: 18px;
    justify-content: space-between;
}
.event-container {
    width: 46%;
    padding: 10px;
    border-radius: 2px;
    border: solid 1px var(--color-border);
}

.data-icon {
    fill: red;
    margin: 10px;
}
.data-link {
    font-size: 14px;
}
</style>
