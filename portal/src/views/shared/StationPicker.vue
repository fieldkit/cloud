<template>
    <div class="station-picker">
        <div class="dialog">
            <div class="close-button" v-on:click="onClose">
                <img alt="Close" src="../../assets/close.png" />
            </div>
        </div>
        <div class="header">
            <div class="title">Add Station</div>
            <div class="search">
                <input v-model="search" label="Search" @input="onSearch(search)" placeholder="Search Stations" />
            </div>
        </div>
        <div class="main">
            <div class="child" v-for="station in visible" v-bind:key="station.id">
                <StationPickerStation :station="station" :selected="selected === station.id" @selected="(ev) => onSelected(station)" />
            </div>
        </div>
        <div class="pagination">
            <div class="button" v-on:click="onPrevious" v-bind:class="{ enabled: canPagePrevious }">Previous</div>
            <div class="pages">
                <div
                    v-for="page in pages"
                    v-bind:key="page.number"
                    class="page"
                    v-bind:class="{ selected: page.selected }"
                    v-on:click="onPage(page.number)"
                >
                    {{ page.number + 1 }}
                </div>
            </div>
            <div class="button" v-on:click="onNext" v-bind:class="{ enabled: canPageNext }">Next</div>
        </div>
        <div class="footer">
            <div class="button" v-on:click="onAdd" v-bind:class="{ enabled: selected }">Add Station</div>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import StationPickerStation from "./StationPickerStation.vue";
import CommonComponents from "@/views/shared";
import { DisplayStation } from "@/store/modules/stations";

export default Vue.extend({
    name: "StationPicker",
    components: {
        ...CommonComponents,
        StationPickerStation,
    },
    props: {
        stations: {
            type: Array,
            required: true,
        },
    },
    data() {
        return {
            selected: null,
            paging: {
                number: 0,
                size: 6,
            },
            search: "",
        };
    },
    computed: {
        pages(this: any) {
            return _.range(0, this.totalPages()).map((p) => {
                return {
                    selected: this.paging.number == p,
                    number: p,
                };
            });
        },
        canPageNext(this: any) {
            return this.paging.number < this.totalPages() - 1;
        },
        canPagePrevious(this: any) {
            return this.paging.number > 0;
        },
        visible(this: any) {
            const start = this.paging.number * this.paging.size;
            const end = start + this.paging.size;
            return this.filteredStations().slice(start, end);
        },
    },
    methods: {
        filteredStations() {
            if (this.search.length === 0) {
                return this.stations;
            }
            return _.filter(this.stations, (station) => {
                return station.name.toLowerCase().indexOf(this.search.toLowerCase()) >= 0;
            });
        },
        totalPages() {
            return Math.ceil(this.filteredStations().length / this.paging.size);
        },
        onSelected(this: any, station: DisplayStation) {
            this.selected = station.id;
        },
        onClose() {
            this.$emit("close");
        },
        onAdd() {
            if (this.selected) {
                this.$emit("add", this.selected);
            }
        },
        onPrevious() {
            if (this.paging.number > 0) {
                this.paging.number--;
            }
        },
        onNext() {
            if (this.paging.number < this.totalPages() - 1) {
                this.paging.number++;
            }
        },
        onPage(this: any, page: number) {
            this.paging.number = page;
        },
        onSearch(this: any, search: string) {
            // console.log("search", search);
        },
    },
});
</script>

<style scoped>
.station-picker {
    display: flex;
    flex-direction: column;
}
.station-picker .dialog {
    display: flex;
}
.dialog .close-button {
    margin-left: auto;
    cursor: pointer;
}
.station-picker .header {
    display: flex;
}
.header .title {
    font-weight: 500;
    font-size: 20px;
    color: #2c3e50;
}
.header .search {
    margin-left: auto;
}
.station-picker .main {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    min-height: 164px; /* Eh */
}
.main .child {
    margin-bottom: 10px;
}
.station-picker .pagination {
    display: flex;
}
.station-picker .footer {
    display: flex;
}
.pagination .button {
    font-size: 14px;
    color: #000000;
    text-align: center;
    padding: 10px;
    background-color: #efefef;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
}
.pagination .button.enabled {
    font-size: 14px;
    font-weight: bold;
    color: #000000;
    text-align: center;
    padding: 10px;
    background-color: #efefef;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    cursor: pointer;
}
.footer .button {
    font-size: 18px;
    font-weight: bold;
    text-align: center;
    padding: 10px;
    margin: 20px 0 0 0px;
    border: 1px solid rgb(215, 220, 225);
    border-radius: 4px;
    color: #000000;
    background-color: #efefef;
}
.footer .button.enabled {
    color: #ffffff;
    background-color: #ce596b;
    cursor: pointer;
}
.pagination {
    display: flex;
    justify-content: space-between;
}
.pages {
    display: flex;
    justify-content: space-evenly;
}
.pages .page {
    font-size: 14px;
    margin-left: 10px;
    margin-right: 10px;
    cursor: pointer;
}
.pages .page.selected {
    font-weight: 900;
}
input {
    font-size: inherit;
    padding: 0.3em;
    border: 2px solid #efefef;
    border-radius: 3px;
}
.header {
    margin-bottom: 1em;
}
</style>
