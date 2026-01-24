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

	ErrAccessExpired = Errors{
		Error: "Access token expired, please refresh",
		Code:  "ACCESS_EXPIRED",
	}

	ErrSessionExpired = Errors{
		Error: "Session expired, please log in again",
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
		Error: "Forbidden to access this resource",
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

	ErrNoRowsReturned = Errors{
		Error: "No rows returned from database",
		Code:  "NO_ROWS",
	}

	ErrEmailInUse = Errors{
		Error: "Email already in use",
		Code:  "EMAIL_INUSE",
	}

	ErrInvalidBody = Errors{
		Error: "Invalid request body",
		Code:  "INVALID_BODY",
	}

	ErrInvalidForm = Errors{
		Error: "Missing required form fields",
		Code:  "INVALID_FORM",
	}

	ErrInvalidCredentials = Errors{
		Error: "Invalid login details",
		Code:  "INVALID_CREDENTIALS",
	}
)
