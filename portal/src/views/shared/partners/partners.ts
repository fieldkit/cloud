import { DisplayStation } from "@/store";
import _ from "lodash";

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
    viz: {
        groupStation: (station: unknown) => string | null;
    };
}

function getNeighborhood(station: DisplayStation): string | null {
    if (station.attributes) {
        const maybeAttribute = station.attributes.find((attr) => attr.name === "Neighborhood");
        if (maybeAttribute) return maybeAttribute.stringValue;
    }
    return null;
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
                    url: "https://www.floodnet.nyc/",
                },
            ],
            stationLocationName: (station: DisplayStation) => {
                return getNeighborhood(station) || station.locationName;
            },
            viz: {
                groupStation: (station: DisplayStation): string | null => {
                    return getNeighborhood(station) || null;
                },
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
        viz: {
            groupStation: (station: DisplayStation): string | null => {
                return null;
            },
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
