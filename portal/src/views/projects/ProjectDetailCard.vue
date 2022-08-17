<template>
    <div class="project-detail-card" :class="{ 'mobile-expanded': showLinksOnMobile }" @click="showLinksOnMobile = !showLinksOnMobile">
        <div class="photo-container">
            <ProjectPhoto :project="project" :image-size="150" />
        </div>
        <div class="detail-container">
            <div class="flex flex-al-center">
                <h1 class="detail-title">{{ project.name }}</h1>
                <div class="detail-links" :class="{ 'mobile-visible': showLinksOnMobile }">
                    <router-link
                        v-if="!isPartnerCustomisationEnabled"
                        :to="{ name: 'viewProject', params: { id: project.id } }"
                        class="link"
                    >
                        Project Dashboard >
                    </router-link>
                    <a v-for="link in partnerCustomization.links" v-bind:key="link.url" :href="link.url" target="_blank" class="link">
                        {{ $t(link.text) }} >
                    </a>
                </div>
            </div>
            <div class="detail-description">{{ project.description }}</div>

            <div
                v-if="partnerCustomization.templates.extraProjectDescription"
                v-html="partnerCustomization.templates.extraProjectDescription"
            ></div>

            <component v-if="partnerCustomization.components.project" v-bind:is="partnerCustomization.components.project"></component>
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
    } {
        return {
            showLinksOnMobile: false,
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

    .link {
        color: $color-primary;
        font-size: 12px;
        letter-spacing: 0.07px;
        text-decoration: initial;

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

.detail-title {
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

    @include bp-down($sm) {
        display: none;
    }
}

.detail-links {
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
