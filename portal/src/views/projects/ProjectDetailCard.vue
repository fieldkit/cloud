<template>
    <div class="project-detail-card" :class="{ 'mobile-expanded': showLinksOnMobile }" @click="showLinksOnMobile = !showLinksOnMobile">
        <div class="photo-container">
            <ProjectPhoto :project="project" :image-size="150" />
        </div>
        <div class="detail-container">
            <component
                v-bind:is="partnerCustomization.components.project"
                :project="project"
                :showLinksOnMobile="showLinksOnMobile"
            ></component>
        </div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { getPartnerCustomizationWithDefault, isCustomisationEnabled, PartnerCustomization } from "@/views/shared/partners";
import { Project } from "@/api/api";
import ProjectPhoto from "@/views/shared/ProjectPhoto.vue";

export default Vue.extend({
    name: "ProjectDetailCard",
    components: {
        ProjectPhoto,
    },
    props: {
        project: {
            type: Object as PropType<Project>,
            required: true,
        },
    },
    data(): {
        showLinksOnMobile: boolean;
        isMobileView: boolean;
    } {
        return {
            showLinksOnMobile: false,
            isMobileView: window.screen.availWidth < 768,
        };
    },
    computed: {
        partnerCustomization(): PartnerCustomization {
            return getPartnerCustomizationWithDefault();
        },
        isPartnerCustomisationEnabled(): boolean {
            return isCustomisationEnabled();
        },
    },
});
</script>

<style scoped lang="scss">
@import "../../scss/project";
@import "../../scss/global";

.project-detail-card {
    display: flex;
    border: 1px solid var(--color-border);
    width: 100%;
    position: absolute;
    top: 66px;
    left: 0;
    box-sizing: border-box;
    background-color: #ffffff;
    text-align: left;
    padding: 27px 30px;
    z-index: $z-index-top;

    body.floodnet & {
        background-color: #f6f9f8;

        @include bp-down($sm) {
            background-color: #ffffff;
        }
    }

    @include bp-down($sm) {
        width: 100%;
        padding: 10px;
        top: 52px;
        align-items: center;
        border-left-width: 0;
        border-right-width: 0;
        z-index: 10;

        &:after {
            content: "";
            width: 0;
            height: 0;
            border-width: 8px 8px 0 8px;
            border-color: #d8d8d8 transparent transparent transparent;
            border-style: solid;
            margin-left: auto;
        }

        &.mobile-expanded:after {
            border-width: 0 8px 8px 8px;
            border-color: transparent transparent #d8d8d8 transparent;
        }
    }

    ::v-deep .link {
        color: $color-primary;
        font-size: 12px;
        letter-spacing: 0.07px;
        text-decoration: initial;
        display: block;

        @include bp-down($sm) {
            font-family: $font-family-medium;
            font-size: 14px;
        }

        body.floodnet & {
            @include bp-up($sm) {
                color: $color-dark;
            }
        }
    }
}

.detail-container {
    overflow: hidden;
}

::v-deep .detail-title {
    font-family: $font-family-bold;
    font-size: 18px;
    margin-top: 0;
    margin-bottom: 2px;
    margin-right: 10px;
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;

    @include bp-down($sm) {
        margin-bottom: 0;
    }
}

::v-deep .detail-description {
    font-family: var(--font-family-light);
    font-size: 14px;
    max-height: 35px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    margin-right: 10px;

    .link {
        display: inline-block;
        font-size: 14px;
        text-decoration: underline;
    }

    @include bp-down($sm) {
        display: none;
    }
}

::v-deep .detail-links {
    @include bp-down($sm) {
        position: fixed;
        top: 104px;
        left: 0;
        background-color: #fff;
        width: 100%;
        padding: 10px 10px 5px 10px;
        opacity: 0;
        visibility: hidden;
        box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.12);

        &.mobile-visible {
            opacity: 1;
            visibility: visible;
        }

        .link {
            color: #1a8ad1;
            border: 1px solid #1a8ad1;
            border-radius: 25px;
            padding: 6px 12px;
            display: inline-block;
            margin-bottom: 5px;
        }
    }
}

.photo-container {
    flex: 0 0 38px;
    height: 38px;
    margin: 0 12px 0 0;

    img {
        border-radius: 2px;
    }

    @include bp-down($sm) {
        flex-basis: 30px;
        height: 30px;
    }
}
</style>
