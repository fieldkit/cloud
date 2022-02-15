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
