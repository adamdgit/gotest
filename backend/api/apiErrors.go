package api

type Errors struct {
	Error string `json:"error" example:"Internal Server Error"`
	Code  string `json:"code,omitempty" example:"INTERNAL_SERVER_ERROR"`
}

var (
	ErrUnauthorized = Errors{
		Error: "Unauthorized",
		Code:  "UNAUTHORIZED",
	}

	ErrSessionExpired = Errors{
		Error: "Session expired, please log in again.",
		Code:  "SESSION_EXPIRED",
	}

	ErrInternalServer = Errors{
		Error: "Internal Server Error",
		Code:  "INTERNAL_SERVER_ERROR",
	}

	ErrInvalidRequest = Errors{
		Error: "Invalid request",
		Code:  "INVALID_REQUEST",
	}

	ErrForbidden = Errors{
		Error: "Forbidden",
		Code:  "FORBIDDEN",
	}

	ErrNotFound = Errors{
		Error: "Resource not found",
		Code:  "NOT_FOUND",
	}

	ErrConflict = Errors{
		Error: "Conflict",
		Code:  "CONFLICT",
	}

	ErrRateLimited = Errors{
		Error: "Too many requests",
		Code:  "RATE_LIMITED",
	}

	ErrInvalidBody = Errors{
		Error: "Invalid request body",
		Code:  "INVALID_BODY",
	}

	ErrInvalidCredentials = Errors{
		Error: "Invalid login details",
		Code:  "INVALID_CREDENTIALS",
	}
)
