export default function interpolatePartner(baseString) {
    const floodnetClass = "floodnet";

    if (document.body.classList.contains(floodnetClass)) {
        return baseString + floodnetClass;
    }

    return baseString + "fieldkit";
}
