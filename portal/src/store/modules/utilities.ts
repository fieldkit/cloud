export interface Services {}

export type ServicesFactory = () => Services;

export class ServiceRef {
    constructor(private readonly services: ServicesFactory | null = null) {}

    private verify(): Services {
        if (!this.services) {
            throw new Error(`Services unfilled`);
        }
        return this.services();
    }
}
