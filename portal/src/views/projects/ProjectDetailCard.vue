<template>
    <div class="project-detail-card">
        <div class="photo-container">
            <ProjectPhoto :project="project" :image-size="150" />
        </div>
        <div class="detail-container">
            <div class="flex flex-al-center">
                <h1 class="detail-title">{{ project.name }}</h1>
                <router-link v-if="!isPartnerCustomisationEnabled()" :to="{ name: 'viewProject', params: { id: id } }" class="link">
                    Project Dashboard >
                </router-link>
                <a v-for="link in partnerCustomization.links" v-bind:key="link.url" :href="link.url" target="_blank" class="link">
                    {{ $t(link.text) }} >
                </a>
            </div>
            <div class="detail-description">{{ project.description }}</div>
            <div v-if="isPartnerCustomisationEnabled()" class="detail-description">
                {{ $t("floodnetDescription.text") }}
                <a href="https://survey123.arcgis.com/share/b9b1d621d16543378b6d3a6b3e02b424" class="link">
                    {{ $t("floodnetDescription.link") }}
                </a>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import Vue, { PropType } from "vue";
import { getPartnerCustomizationWithDefault, isCustomisationEnabled } from "@/views/shared/partners";
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
    methods: {
        isPartnerCustomisationEnabled(): boolean {
            return isCustomisationEnabled();
        },
        partnerCustomization() {
            return getPartnerCustomizationWithDefault();
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
    border-radius: 3px;
    z-index: $z-index-top;
    width: 100%;
    position: absolute;
    top: 66px;
    left: 0;
    box-sizing: border-box;
    background-color: #ffffff;
    text-align: left;
    padding: 27px 30px;

    @include bp-down($sm) {
        width: 100%;
        position: fixed;
        border-top-right-radius: 10px;
        border-top-left-radius: 10px;
        bottom: 0px;
        text-align: center;
        padding-bottom: 10px;
        top: auto;
        right: auto;
        align-items: center;
        justify-content: center;
    }

    body.floodnet & {
        background-color: #f6f9f8;
    }

    .link {
        color: $color-primary;
        font-size: 12px;
        letter-spacing: 0.07px;
        text-decoration: initial;

        body.floodnet & {
            color: $color-dark;
        }
    }
}
.detail-title {
    font-family: var(--font-family-bold);
    font-size: 18px;
    margin-top: 0;
    margin-bottom: 2px;
    margin-right: 10px;
}
.detail-container {
    width: 75%;
}
.detail-description {
    font-family: var(--font-family-light);
    font-size: 14px;
    max-height: 35px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    max-height: 35px;
    overflow: hidden;
    margin-right: 10px;
}

.photo-container {
    width: 38px;
    height: 38px;
    margin: 0 12px 0 0;

    img {
        border-radius: 2px;
    }

    @include bp-down($sm) {
        display: none;
    }
}
</style>
