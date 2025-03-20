export enum ErrorReason {
  UNKNOWN = 'UNKNOWN',
  NETWORK = 'NETWORK',
  INTERNAL_SERVER_ERROR = 'INTERNAL_SERVER_ERROR',
  BAD_REQUEST = 'BAD_REQUEST',
  ALREADY_EXISTS = 'ALREADY_EXISTS',
  INVALID_CREDENTIALS = 'INVALID_CREDENTIALS',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
}

export enum FriendlyErrorMessage {
  UNKNOWN = 'An unexpected error occurred. Please try again later.',
  NETWORK = 'Network error. Are you offline?',
  INTERNAL_SERVER_ERROR = 'Service unavailable right now. Try again later.',
  BAD_REQUEST = 'Invalid request parameters.',
  ALREADY_EXISTS = 'Resource already exists.',
  INVALID_CREDENTIALS = 'Email or password is not correct.',
  TOKEN_EXPIRED = 'Token expired. Please login again.',
}

export interface APIError extends Error {
  code: number
  message: string
  details: APIErrorInfoDetail[]
}

export interface APIErrorInfoDetail {
  '@type': string
  'reason'?: ErrorReason
  'fieldViolations'?: FieldViolation[]
}

interface FieldViolation {
  description: string
  field: string
  localizedMessage: string
  reason: string
}

interface AppErrorParams {
  reason: ErrorReason
  message: string
  friendlyMessage: string
  details?: FieldViolation[]
}

export class AppError extends Error {
  friendlyMessage: string
  reason: ErrorReason
  details: FieldViolation[]

  constructor({ reason, message, friendlyMessage, details }: AppErrorParams) {
    super(message)
    this.name = 'AppError'
    this.friendlyMessage = friendlyMessage
    this.reason = reason
    this.details = details ?? []
  }
}

export function ensureAppError(originalError: unknown): AppError {
  if (originalError instanceof AppError) {
    return originalError
  }
  if (originalError instanceof Error) {
    return unknownAppError(originalError.message)
  }
  return unknownAppError(JSON.stringify(originalError))
}

export function unknownAppError(message: string): AppError {
  return new AppError({
    reason: ErrorReason.UNKNOWN,
    friendlyMessage: FriendlyErrorMessage.UNKNOWN,
    message: message,
  })
}

export function networkAppError(message: string): AppError {
  return new AppError({
    reason: ErrorReason.NETWORK,
    friendlyMessage: FriendlyErrorMessage.NETWORK,
    message: message,
  })
}

export async function unmarshallAPIError(response: Response): Promise<AppError> {
  const { message, details } = await response.json() as APIError

  const error = unknownAppError(message)

  if (!Array.isArray(details)) {
    return error
  }

  details.forEach((detail) => {
    switch (detail['@type']) {
      case 'type.googleapis.com/google.rpc.ErrorInfo':
        error.reason = ErrorReason[detail.reason as keyof typeof ErrorReason]
        error.friendlyMessage = FriendlyErrorMessage[error.reason]
        break
      case 'type.googleapis.com/google.rpc.BadRequest':
        for (const fieldViolation of detail.fieldViolations ?? []) {
          error.details.push(fieldViolation)
        }
        break
    }
  })

  return error
}
