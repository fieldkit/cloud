<template>
    <div class="view-type-container" :class="{ 'list-toggled': viewType === 'list' }">
        <!--
            <label class="toggle-btn">
                <input type="checkbox" v-model="recentMapMode" />
                <span :class="{ active: !recentMapMode }">{{ $t("map.toggle.current") }}</span>
                <i></i>
                <span :class="{ active: recentMapMode }">{{ $t("map.toggle.recent") }}</span>
            </label>
        -->
        <div class="view-type">
            <router-link
                v-for="route in routes"
                v-bind:key="route.viewType"
                v-bind:class="{ active: route.viewType === viewType }"
                :to="{ name: route.name, params: route.params }"
            >
                {{ $t(route.label) }}
            </router-link>
        </div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";

export enum MapViewType {
    map = "map",
    list = "list",
}

export default Vue.extend({
    name: "MapViewTypeToggle",
    props: {
        routes: {
            required: true,
            type: Array as PropType<
                Array<{
                    name: string;
                    label: string;
                    params: object;
                    viewType: string;
                }>
            >,
        },
    },
    computed: {
        viewType(): MapViewType {
            if (this.$route.meta?.viewType) {
                return this.$route.meta.viewType;
            }
            return MapViewType.map;
        },
    },
});
</script>

<style lang="scss" scoped>
@import "src/scss/mixins";
.view-type {
    width: 115px;
    height: 38px;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);
    border: solid 1px #f4f5f7;
    background-color: #ffffff;
    cursor: pointer;
    margin-left: 30px;
    @include flex(center, center);

    @include bp-down($sm) {
        width: 98px;
        margin-left: 5px;
    }

    @media screen and (max-width: 350px) {
        width: 88px;
    }

    &-container {
        z-index: $z-index-top;
        margin: 0;
        box-sizing: border-box;
        @include flex(center, center);
        @include position(absolute, 90px 25px null null);

        @include bp-down($sm) {
            @include position(absolute, 115px 10px null null);
        }

        &.list-toggled {
            @include bp-down($sm) {
                left: 0;
                width: 100%;
                justify-content: flex-end;
                padding: 0 10px;
            }
        }
    }

    > a {
        flex-basis: 50%;
        height: 100%;
        font-family: $font-family-light;
        @include flex(center, center);

        &:nth-of-type(1) {
            border-right: solid 1px $color-border;
        }

        @include bp-down($sm) {
            font-size: 14px;
        }

        &.active {
            font-family: $font-family-bold;
        }
    }

    .icon {
        font-size: 18px;
    }
}

.toggle-btn {
    cursor: pointer;
    z-index: $z-index-top;
    position: relative;
    font-size: 14px;
    -webkit-tap-highlight-color: transparent;
    font-family: $font-family-medium !important;
    height: 38px;
    display: flex;
    align-items: center;

    @include bp-down($sm) {
        .view-type-container:not(.list-toggled) & {
            background-color: #fff;
            padding: 0 10px;
            box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.13);
            border: solid 1px #f4f5f7;
        }
    }

    @media screen and (max-width: 350px) {
        padding: 0 8px;
    }

    * {
        font-family: $font-family-medium !important;
    }

    span {
        opacity: 0.5;

        &.active {
            opacity: 1;
            color: $color-fieldkit-primary;

            body.floodnet & {
                color: $color-dark;
            }
        }
    }
}
.toggle-btn i {
    position: relative;
    display: inline-block;
    margin: 0 10px;
    width: 27px;
    height: 16px;
    background-color: #e6e6e6;
    border-radius: 20px;
    vertical-align: text-bottom;
    transition: all 0.3s linear;
    user-select: none;
}
.toggle-btn i::before {
    content: "";
    position: absolute;
    left: 0;
    width: 27px;
    background-color: #fff;
    border-radius: 11px;
    transform: translate3d(2px, 2px, 0) scale3d(1, 1, 1);
    transition: all 0.25s linear;
}
.toggle-btn i::after {
    content: "";
    position: absolute;
    left: 0;
    width: 12px;
    height: 12px;
    background-color: #fff;
    border-radius: 50%;
    transform: translate3d(2px, 2px, 0);
    transition: all 0.2s ease-in-out;
}
.toggle-btn:active i::after {
    width: 28px;
    transform: translate3d(2px, 2px, 0);
}
.toggle-btn:active input:checked ~ i::after {
    transform: translate3d(16px, 2px, 0);
}
.toggle-btn input {
    display: none;
}
.toggle-btn input ~ i {
    background-color: $color-primary;
}
.toggle-btn input:checked ~ i::before {
    transform: translate3d(18px, 2px, 0) scale3d(0, 0, 0);
}
.toggle-btn input:checked ~ i::after {
    transform: translate3d(13px, 2px, 0);
}
</style>
