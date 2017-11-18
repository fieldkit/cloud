// @flow weak

export type ActiveExpedition = {
    project: {
        name: string,
        slug: string
    },
    expedition: {
        name: string,
        slug: string
    }
}

export type Focus = {
    feature: ?GeoJSONFeature,
    time: number
}
