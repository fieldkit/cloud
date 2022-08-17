import { DisplayStation } from "@/store";
import moment from "moment";

import Vue, { Component } from "vue";

const FloodNetProjectLinks = Vue.extend({
    name: "FloodNetProjectLinks",
    template: `
    <div class="example" v-if="false">
        Here's an example, this can also be in a separate file and imported.
        The v-if is there to hide this example, probably would be left off in a real use case.
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
    links: {
        text: string;
        url: string;
    }[];
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
    templates: {
        [key: string]: string;
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
    if (window.location.hostname.indexOf("floodnet.") >= 0) {
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
            links: [
                {
                    text: "linkToFloodnet",
                    url: "https://www.floodnet.nyc/methodology",
                },
            ],
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
                project: FloodNetProjectLinks,
            },
            templates: {
                extraProjectDescription: `
                    <div class="detail-description">
                        Report a flood: 
                        <a href="https://survey123.arcgis.com/share/b9b1d621d16543378b6d3a6b3e02b424" target="_blank" class="link">
                            Community Flood Watch Project Survey
                        </a>
                    </div>
                `,
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
        links: [],
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
            project: null,
        },
        templates: {},
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
