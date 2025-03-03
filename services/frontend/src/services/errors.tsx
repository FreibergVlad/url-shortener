
export enum APIErrorReason {
    UNKNOWN = "UNKNOWN",
    NETWORK = "NETWORK",
    INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR",
    BAD_REQUEST = "BAD_REQUEST",
    ALREADY_EXISTS = "ALREADY_EXISTS",
    INVALID_CREDENTIALS = "INVALID_CREDENTIALS",
    TOKEN_EXPIRED = "TOKEN_EXPIRED",
}

export enum FriendlyErrorMessage {
    UNKNOWN = "An unexpected error occurred. Please try again later.",
    NETWORK = "Network error. Are you offline?",
    INTERNAL_SERVER_ERROR = "Service unavailable right now. Try again later.",
    BAD_REQUEST = "Invalid request parameters.",
    ALREADY_EXISTS = "Resource already exists.",
    INVALID_CREDENTIALS = "Email or password is not correct.",
    TOKEN_EXPIRED = "Token expired. Please login again.",
}

interface APIErrorParams {
    reason: APIErrorReason
    message: string
    friendlyMessage: string
    details?: Record<string, string>
}

export class APIError extends Error {
    friendlyMessage: string;
    reason: APIErrorReason;
    details?: Record<string, string>

    constructor({reason, message, friendlyMessage, details}: APIErrorParams) {
        super(message);
        this.name = "APIError";
        this.friendlyMessage = friendlyMessage;
        this.reason = reason;
        this.details = details;
    }
}

export function ensureAPIError(originalError: unknown): APIError {
    if (originalError instanceof APIError) {
        return originalError;
    }
    if (originalError instanceof Error) {
        return unknownAPIError(originalError.message);
    }
    return unknownAPIError(JSON.stringify(originalError));
}

export function unknownAPIError(message: string): APIError {
    return new APIError({
        reason: APIErrorReason.UNKNOWN,
        friendlyMessage: FriendlyErrorMessage.UNKNOWN,
        message: message
    });
}

export function networkAPIError(message: string): APIError {
    return new APIError({
        reason: APIErrorReason.NETWORK,
        friendlyMessage: FriendlyErrorMessage.NETWORK,
        message: message
    });
}

export async function unmarshallAPIError(response: Response): Promise<APIError> {
    const {message, details} = await response.json();

    const error = unknownAPIError(message);

    if (!Array.isArray(details)) {
        return error;
    }

    details.forEach((detail) => {
        switch (detail["@type"]) {
        case "type.googleapis.com/google.rpc.ErrorInfo":
            error.reason = APIErrorReason[detail.reason as keyof typeof APIErrorReason] || APIErrorReason.UNKNOWN;
            error.friendlyMessage = FriendlyErrorMessage[error.reason]
            break;
        case "type.googleapis.com/google.rpc.BadRequest":
            break;
        }
    });

    return error;
}