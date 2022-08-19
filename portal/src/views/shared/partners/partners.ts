import { DisplayStation, Project } from "@/store";
import moment from "moment";

import Vue, { Component, PropType } from "vue";

const FieldKitProjectDescription = Vue.extend({
    name: "FieldKitProjectDescription",
    props: {
        showLinksOnMobile: {
            type: Boolean,
            required: true,
        },
        project: {
            type: Object as PropType<Project>,
            required: true,
        },
    },
    data(): {
        isMobileView: boolean;
    } {
        return {
            isMobileView: window.screen.availWidth < 768,
        };
    },
    template: `
        <div>
            <div class="flex flex-al-center">
                <h1 class="detail-title">{{ project.name }}</h1>
                <div class="detail-links" :class="{ 'mobile-visible': showLinksOnMobile }">
                  <router-link :to="{ name: 'viewProject', params: { id: project.id } }" class="link">
                      {{$t('project.dashboard')}} >
                  </router-link>
                </div>
            </div>
            <div class="detail-description">{{ project.description }}</div>
        </div>
    `,
});

const FloodNetProjectDescription = Vue.extend({
    name: "FloodNetProjectDescription",
    props: {
        showLinksOnMobile: {
            type: Boolean,
            required: true,
        },
        project: {
            type: Object as PropType<Project>,
            required: true,
        },
    },
    data(): {
        isMobileView: boolean;
        links: {
            text: string;
            mobileText: string;
            url: string;
        }[];
    } {
        return {
            isMobileView: window.screen.availWidth < 768,
            links: [
                {
                    text: "linkToFloodnet.desktop",
                    mobileText: "linkToFloodnet.mobile",
                    url: "https://www.floodnet.nyc/methodology",
                },
            ],
        };
    },
    template: `
        <div>
            <div class="flex flex-al-center">
                <h1 class="detail-title">{{ project.name }}</h1>
                <div class="detail-links" :class="{ 'mobile-visible': showLinksOnMobile }">
                    <a v-for="link in links" v-bind:key="link.url" :href="link.url" target="_blank" class="link">
                        <template v-if="isMobileView">{{ $t(link.mobileText) }} ></template>
                        <template v-else>{{ $t(link.text) }} ></template>
                    </a>
                    <a v-if="isMobileView" href="https://survey123.arcgis.com/share/b9b1d621d16543378b6d3a6b3e02b424" target="_blank" class="link">
                        {{ $t('floodnet.reportFlood') }} >
                    </a>
                </div>
            </div>
            <div class="detail-description">{{ project.description }}</div>
            <div class="detail-description">
                {{ $t('floodnet.reportFlood') }}:
                <a href="https://survey123.arcgis.com/share/b9b1d621d16543378b6d3a6b3e02b424" target="_blank" class="link">
                    {{ $t('floodnet.survey') }}
                </a>
            </div>
        </div>
    `,
});

export interface PartnerCustomization {
    title: string; // TODO i18n
    class: string;
    icon: string;
    sharing: {
        viz: string;
    };
    nav: {
        viz: {
            back: {
                // project: string;
                map: { label: string };
            };
        };
    };
    projectId: number | null;
    interpolate: (base: string) => string;
    email: {
        subject: string;
    };
    stationLocationName: (station: DisplayStation) => string;
    getStationDeploymentDate: (station: DisplayStation) => string | Date | null;
    viz: {
        groupStation: (station: unknown) => string | null;
    };
    routeAfterLogin: string;
    sidebarNarrow: boolean;
    components: {
        project: Component | null;
    };
}

function getAttribute(station: DisplayStation, name: string): string | null {
    if (station.attributes) {
        const maybeAttribute = station.attributes.find((attr) => attr.name === name);
        if (maybeAttribute) return maybeAttribute.stringValue;
    }
    return null;
}

function getNeighborhood(station: DisplayStation): string | null {
    return getAttribute(station, "Neighborhood");
}

function getBorough(station: DisplayStation): string | null {
    return getAttribute(station, "Borough");
}

function getDeploymentDate(station: DisplayStation): string | null {
    return getAttribute(station, "Deployment Date");
}

export function getPartnerCustomization(): PartnerCustomization | null {
    // dataviz.floodnet.nyc, floodnet.fieldkit.org
    if (window.location.hostname.indexOf("floodnet.") !== 0) {
        return {
            title: "Data Dashboard - FloodNet",
            class: "floodnet",
            icon: "/favicon-floodnet.ico",
            sharing: {
                viz: `Check out this data on FloodNet!`, // TODO i18n
            },
            nav: {
                viz: {
                    back: {
                        // project: "layout.backProjectDashboard",
                        map: { label: "layout.backToSensors" },
                    },
                },
            },
            projectId: 174,
            interpolate: (baseString: string) => {
                return baseString + "floodnet";
            },
            email: {
                subject: "sharePanel.emailSubject.floodnet",
            },
            stationLocationName: (station: DisplayStation) => {
                return getNeighborhood(station) || station.locationName;
            },
            getStationDeploymentDate: (station: DisplayStation) => {
                return getDeploymentDate(station) || (station.deployedAt ? moment(station.deployedAt).format("M/D/YYYY") : "N/A");
            },
            viz: {
                groupStation: (station: DisplayStation): string | null => {
                    return getBorough(station) || null;
                },
            },
            routeAfterLogin: "root",
            sidebarNarrow: true,
            components: {
                project: FloodNetProjectDescription,
            },
        };
    }
    return null;
}

export function getPartnerCustomizationWithDefault(): PartnerCustomization {
    const maybePartner = getPartnerCustomization();
    if (maybePartner) {
        return maybePartner;
    }

    return {
        title: "Data Dashboard - FieldKit",
        class: "fieldkit",
        icon: "/favicon-fieldkit.ico",
        sharing: {
            viz: `Check out this data on FieldKit!`, // TODO i18n
        },
        nav: {
            viz: {
                back: {
                    // project: "layout.backProjectDashboard",
                    map: { label: "layout.backToStations" },
                },
            },
        },
        projectId: null,
        interpolate: (baseString: string) => {
            return baseString;
        },
        email: {
            subject: "sharePanel.emailSubject.fieldkit",
        },
        stationLocationName: (station: DisplayStation) => {
            return station.locationName;
        },
        getStationDeploymentDate: (station: DisplayStation) => {
            return station.deployedAt ? moment(station.deployedAt).format("M/D/YYYY") : "N/A";
        },
        viz: {
            groupStation: (station: DisplayStation): string | null => {
                return null;
            },
        },
        routeAfterLogin: "projects",
        sidebarNarrow: false,
        components: {
            project: FieldKitProjectDescription,
        },
    };
}

export function isCustomisationEnabled(): boolean {
    return getPartnerCustomization() != null;
}

export function interpolatePartner(baseString): string {
    const partnerCustomization = getPartnerCustomization();
    if (partnerCustomization != null) {
        return partnerCustomization.interpolate(baseString);
    }
    return baseString + "fieldkit";
}
