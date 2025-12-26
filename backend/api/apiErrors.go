package api

type ErrorResponse struct {
	Error string `json:"error" example:"Internal Server Error"`
	Code  string `json:"code,omitempty" example:"INTERNAL_SERVER_ERROR"`
}

var (
	ErrUnauthorized = ErrorResponse{
		Error: "Unauthorized",
		Code:  "UNAUTHORIZED",
	}

	ErrSessionExpired = ErrorResponse{
		Error: "Session expired, please log in again.",
		Code:  "SESSION_EXPIRED",
	}

	ErrInternalServer = ErrorResponse{
		Error: "Internal Server Error",
		Code:  "INTERNAL_SERVER_ERROR",
	}

	ErrInvalidRequest = ErrorResponse{
		Error: "Invalid request",
		Code:  "INVALID_REQUEST",
	}

	ErrForbidden = ErrorResponse{
		Error: "Forbidden",
		Code:  "FORBIDDEN",
	}

	ErrNotFound = ErrorResponse{
		Error: "Resource not found",
		Code:  "NOT_FOUND",
	}

	ErrConflict = ErrorResponse{
		Error: "Conflict",
		Code:  "CONFLICT",
	}

	ErrRateLimited = ErrorResponse{
		Error: "Too many requests",
		Code:  "RATE_LIMITED",
	}
)
