<template>
    <div class="latest-events">
        <div class="event-row">
            <div class="event-container" v-for="(annotation, idx) in displayEvents" :key="idx">
                <h3 class="event-title">{{ annotation.body }}</h3>
                <div>{{ annotation.start | date }} &ndash; {{ annotation.end | date }}</div>
                <div class="row">
                    <div class="data-link">
                        <img alt="View Data" class="data-icon" src="@/assets/view-data.svg" width="18" />
                        <router-link class="data-link" :to="{ name: 'exploreBookmark', query: { bookmark: annotation.bookmark } }">
                            View data
                        </router-link>
                    </div>
                    <div class="user-details">
                        <UserPhoto :user="annotation.user"></UserPhoto>
                        <span class="user-name">{{ annotation.user.name }}</span>
                    </div>
                </div>
            </div>
        </div>
        <pagination-controls :page="currentPage" :total-pages="numPages" @new-page="onNewPage" />
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";
import { Annotation } from "@/store";
// fixme: moment is huge, load module?
import moment from "moment";
import PaginationControls from "./PaginationControls.vue";

export default Vue.extend({
    name: "LatestEvents",
    components: {
        PaginationControls,
        ...CommonComponents,
    },
    props: {},
    data(): {
        viewType: string;
        currentPage: number;
    } {
        return {
            viewType: "project",
            currentPage: 0,
        };
    },
    computed: {
        latestEvents(): Annotation[] {
            return this.$getters.getEvents;
        },
        numPages(): number {
            return Math.ceil(this.latestEvents.length / 4);
        },
        displayEvents(): Annotation[] {
            return this.latestEvents.slice(this.currentPage, this.currentPage + 4);
        },
    },
    methods: {
        onNewPage(evt) {
            this.currentPage = evt;
        },
    },
    filters: {
        date: (value) => {
            return moment(value).format("M/D/YYYY, h:mm:ss a");
        },
    },
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
.event-title {
    max-width: 100%;
    text-overflow: ellipsis;
    white-space: nowrap;
    height: 1.1em;
    overflow: hidden;
}
.data-icon {
    fill: red;
    margin: 10px;
}
.data-link {
    font-size: 14px;
    display: flex;
    align-items: center;
}
.user-details {
    display: flex;
    align-items: center;
    margin-right: 1em;
}
.row {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
}
</style>
