import { ShortUser } from "./user";

export interface ShortURL {
    organizationId: string;
    domain: string;
    alias: string;
    shortUrl: string;
    expiresAt: string;
    createdAt: string;
    createdBy: ShortUser;
    description: string;
    tags: Array<string>;
}

export interface ListShortURLResponse {
    data: Array<ShortURL>;
    total: number;
}