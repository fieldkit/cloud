interface Route {
    name: string;
    params: { id: number };
}

export interface PartnerCustomization {
    title: string; // TODO i18n
    class: string;
    icon: string;
    interpolate: (base: string) => string;
    forcedLandingPage: Route | null;
}

export function getPartnerCustomization(): PartnerCustomization | null {
    if (window.location.hostname.indexOf("floodnet.") === 0) {
        return {
            title: "Data Dashboard - FloodNet",
            class: "floodnet",
            icon: "/favicon-floodnet.ico",
            interpolate: (baseString: string) => {
                return baseString + "floodnet";
            },
            forcedLandingPage: {
                name: "viewProjectBigMap",
                params: { id: 174 },
            },
        };
    }
    return null;
}

export function getPartnerForcedLandingPage(): Route | null {
    const partnerCustomization = getPartnerCustomization();
    return partnerCustomization?.forcedLandingPage || null;
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
