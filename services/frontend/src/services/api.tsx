import { API_BASE_URL } from "./constants";
import { APIErrorReason, ensureAPIError, networkAPIError, unmarshallAPIError } from "./errors";

interface APIRequestConfig<RequestT> {
    endpoint: string
    method: string
    body?: RequestT
    accessToken?: string
    onTokenExpired?: (() => Promise<string>)
}

export async function executeAPIRequest<RequestT, ResponseT>(request: APIRequestConfig<RequestT>): Promise<ResponseT> {
    const execute = async (token: string | undefined): Promise<Response> => {
        try {
            return await fetch(`${API_BASE_URL}/${request.endpoint}`, {
                method: request.method,
                headers: {
                    "Content-Type": "application/json",
                    "Accept": "application/json",
                    ...(token ? {"Authorization": `Bearer ${token}`} : {} ),
                },
                body: request ? JSON.stringify(request.body) : null,
            });
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
