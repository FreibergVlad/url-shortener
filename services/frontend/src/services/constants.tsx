export const API_BASE_URL: string = requireEnv("VITE_API_BASE_URL")

function requireEnv(key: string): string {
    const val = import.meta.env[key];
    if (!val) {
        throw new Error(`env variable '${key}' should be set`);
    }
    return val;
}