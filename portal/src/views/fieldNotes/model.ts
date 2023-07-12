export interface PortalStationFieldNotes {
    id: number;
    author: { id: number; name: string; photo: object };
    body: string;
    createdAt: number;
    updatedAt: number;
    userId?: number;
}
