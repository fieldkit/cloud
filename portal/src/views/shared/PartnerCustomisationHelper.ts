export enum FKPartnersEnum {
    floodnet = "floodnet",
}

export function isCustomisationEnabledFor(partner: FKPartnersEnum) {
    return document.body.classList.contains(partner);
}

export function isCustomisationEnabled() {
    return document.body.classList.contains(FKPartnersEnum.floodnet);
}

export function interpolatePartner(baseString) {
    if (isCustomisationEnabledFor(FKPartnersEnum.floodnet)) {
        return baseString + FKPartnersEnum.floodnet;
    }

    return baseString + "fieldkit";
}

interface Route {
    name: string;
    params: { id: number };
}

export function getPartnerForcedLandingPage(): Route | null {
    if (isCustomisationEnabled()) {
        // Terrible. Only remotely acceptable because it's isolated to this
        // partner related file. I'm thinking of possible better approaches to
        // this, not totally sure what I like just yet.
        return {
            name: "viewProjectBigMap",
            params: { id: 174 },
        };
    }
    return null;
}
