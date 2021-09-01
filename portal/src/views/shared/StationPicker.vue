<template>
    <div class="station-picker">
        <div class="dialog">
            <div class="close-button" v-on:click="onClose">
                <img alt="Close" src="@/assets/icon-close.svg" />
            </div>
        </div>
        <div class="header">
            <div class="title">{{ title }}</div>
            <div>
                <span class="selected-counter">{{ $t("selected") }} ({{ selected.length }})</span>
                <input
                    class="search"
                    v-model="search"
                    :label="$t('project.stations.search')"
                    @input="onSearch(search)"
                    :placeholder="$t('project.stations.search')"
                />
            </div>
        </div>
        <div class="main">
            <StationPickerStation
                v-for="station in visible"
                v-bind:key="station.id"
                :station="station"
                :selected="selected.includes(station.id)"
                @selected="(ev) => onSelected(station)"
            />
        </div>
        <PaginationControls :page="paging.number" :totalPages="totalPages()" @new-page="onNewPage" />
        <div class="footer">
            <button class="button-solid" v-on:click="onCtaClick" v-bind:class="{ enabled: selected }">{{ actionText }}</button>
        </div>
    </div>
</template>

<script lang="ts">
import _ from "lodash";
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";
import StationPickerStation from "./StationPickerStation.vue";
import PaginationControls from "./PaginationControls.vue";
import { DisplayStation } from "@/store";

export enum StationPickerActionType {
    add = "add",
    remove = "remove",
}

export default Vue.extend({
    name: "StationPicker",
    components: {
        ...CommonComponents,
        StationPickerStation,
        PaginationControls,
    },
    props: {
        title: {
            type: String,
            required: true,
        },
        actionType: {
            type: String,
            required: true,
        },
        actionText: {
            type: String,
            required: true,
        },
        stations: {
            type: Array as PropType<DisplayStation[]>,
            required: true,
        },
        filter: {
            type: Function as PropType<(station: DisplayStation) => boolean>,
            default: (station) => true,
        },
    },
    data(): {
        selected: number[];
        paging: {
            number: number;
            size: number;
        };
        search: string;
    } {
        return {
            selected: [],
            paging: {
                number: 0,
                size: window.screen.availWidth > 500 ? 6 : 3,
            },
            search: "",
        };
    },
    computed: {
        visible(): DisplayStation[] {
            const start = this.paging.number * this.paging.size;
            const end = start + this.paging.size;
            return this.filteredStations().slice(start, end);
        },
    },
    methods: {
        filteredStations(): DisplayStation[] {
            if (this.search.length === 0) {
                return this.stations.filter((station) => this.filter(station));
            }
            return _.filter(
                this.stations.filter((station) => this.filter(station)),
                (station) => {
                    console.log("fff", station);
                    return station.name.toLowerCase().indexOf(this.search.toLowerCase()) >= 0;
                }
            );
        },
        totalPages(): number {
            return Math.ceil(this.filteredStations().length / this.paging.size);
        },
        onSelected(station: DisplayStation): void {
            const idAtIndex = this.selected.findIndex((value) => value === station.id);
            if (idAtIndex > -1) {
                this.selected.splice(idAtIndex, 1);
            } else {
                this.selected.push(station.id);
            }
        },
        onClose(): void {
            this.$emit("close");
        },
        onCtaClick(): void {
            if (!this.selected) {
                return;
            }
            this.$emit("ctaClick", this.selected);
        },
        onNewPage(page: number): void {
            this.paging.number = page;
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/mixins";
@import "../../scss/global";

.station-picker {
    display: flex;
    flex-direction: column;
    position: relative;
}
.station-picker .dialog {
    display: flex;
    @include position(absolute, -48px -1px null null);

    @include bp-down($xs) {
        @include position(absolute, -25px -1px null null);
    }
}
.dialog .close-button {
    margin-left: auto;
    cursor: pointer;

    img {
        height: 19px;
    }
}
.station-picker .header {
    margin-bottom: 1em;
    flex-wrap: wrap;
    @include flex(center, space-between);
}
.header .title {
    font-weight: 500;
    font-size: 20px;
    color: #2c3e50;

    @include bp-down($xs) {
        flex-basis: 100%;
        margin-bottom: 10px;
    }

    ~ div {
        @include bp-down($xs) {
            width: 100%;
            @include flex(center);
        }
    }
}
.header .search {
    font-size: 14px;
    width: 247px;
    padding: 7px 10px;
    border: solid 1px #cccdcf;

    @include bp-down($xs) {
        width: unset;
        flex: 1;
    }
}
.station-picker .main {
    display: flex;
    flex-wrap: wrap;
    padding-bottom: 5px;
    margin: 0 -8px;
    width: calc(100% + 19px);

    @include bp-down($xs) {
        width: 100%;
        margin: 0;
    }
}
.station-picker .pagination {
    display: flex;
    margin-top: 40px;

    @include bp-down($xs) {
        margin-top: 0;
    }
}
.station-picker .footer {
    @include flex(center, center);
    margin-top: 35px;
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
input {
    font-size: inherit;
    padding: 0.3em;
    border: 2px solid #efefef;
    border-radius: 3px;
}

.selected-counter {
    font-size: 14px;
    margin-right: 17px;
}
</style>
