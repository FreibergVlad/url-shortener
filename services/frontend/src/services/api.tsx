import { API_BASE_URL } from "./constants";
import { APIErrorReason, ensureAPIError, networkAPIError, unmarshallAPIError } from "./errors";

interface APIRequestConfig<RequestT> {
    endpoint: string
    method: string
    body?: RequestT
    queryParams?: Record<string, string>
    accessToken?: string
    onTokenExpired?: (() => Promise<string>)
}

export async function executeAPIRequest<RequestT, ResponseT>(request: APIRequestConfig<RequestT>): Promise<ResponseT> {
    const execute = async (token: string | undefined): Promise<Response> => {
        let url = `${API_BASE_URL}/${request.endpoint}`;
        if (request.queryParams) {
            url += `?${new URLSearchParams(request.queryParams)}`;
        }
        const body = request.body ? JSON.stringify(request.body) : null;
        const headers = {
            "Content-Type": "application/json",
            "Accept": "application/json",
            ...(token ? {"Authorization": `Bearer ${token}`} : {} ),
        };

        try {
            return await fetch(url, {method: request.method, headers, body});
        } catch (e: unknown) {
            throw networkAPIError(ensureAPIError(e).message);
        }
    };

    let response = await execute(request.accessToken);
    if (response.ok) {
        return await response.json();
    }

    const error = await unmarshallAPIError(response);
    if (error.reason !== APIErrorReason.TOKEN_EXPIRED || !request.onTokenExpired) {
        throw error;
    }

    const refreshedToken = await request.onTokenExpired();
    response = await execute(refreshedToken);
    if (response.ok) {
        return await response.json();
    }

    throw unmarshallAPIError(response);
}
