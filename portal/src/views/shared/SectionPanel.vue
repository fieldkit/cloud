<template>
    <section class="container">
        <header :class="{ divider: divider }">{{ title }}</header>
        <slot></slot>
    </section>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import CommonComponents from "@/views/shared";

export default Vue.extend({
    name: "SectionPanel",
    components: {
        ...CommonComponents,
    },
    props: {
        title: {
            type: String,
            required: true,
        },
        divider: {
            type: Boolean,
            default: true
        }
    },
});
</script>

<style lang="scss" scoped>
@import "../../scss/global";

* {
    box-sizing: border-box;
    font-size: 14px;
}

.hide {
    display: none;
}

.container {
    margin-top: 20px;
    padding: 0 20px 30px 20px;
    background: #fff;
    border-radius: 1px;
    border: 1px solid $color-border;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.05);

    @include bp-down($xs) {
        margin: 20px -10px 0;
        padding: 0 10px 30px 10px;
    }

    &.data-view {
        margin-top: 0;
        // padding-top: 45px;
        box-shadow: none;
        border: 0;
    }
}

header {
    @include flex(center, space-between);
    padding: 13px 0;
    font-size: 20px;
    font-weight: 500;
    
    .divider {
        border-bottom: 1px solid $color-border;
    }

    .data-view & {
        font-size: 18px;
        height: auto;
        border: none;
    }

    body.floodnet & {
        font-family: $font-family-floodnet-bold;
    }
}

.subheader {
    @include flex(center, space-between);
    border-top: 1px solid $color-border;
    border-bottom: 1px solid $color-border;
    padding: 15px 0;

    .data-view & {
        border-top: none;
        padding: 0;
    }
}


.timestamp {
    font-family: $font-family-light;
    flex-shrink: 0;
    margin-left: auto;
    line-height: 1.5;
    padding-left: 10px;
}



</style>
