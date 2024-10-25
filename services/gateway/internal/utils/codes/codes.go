package codes

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

var grpcToHTTPStatus = map[codes.Code]int{
	codes.OK:                 http.StatusOK,                  // 0 -> 200 OK
	codes.InvalidArgument:    http.StatusBadRequest,          // 3 -> 400 Bad Request
	codes.FailedPrecondition: http.StatusBadRequest,          // 9 -> 400 Bad Request
	codes.OutOfRange:         http.StatusBadRequest,          // 11 -> 400 Bad Request
	codes.Unauthenticated:    http.StatusUnauthorized,        // 16 -> 401 Unauthorized
	codes.PermissionDenied:   http.StatusForbidden,           // 7 -> 403 Forbidden
	codes.NotFound:           http.StatusNotFound,            // 5 -> 404 Not Found
	codes.AlreadyExists:      http.StatusConflict,            // 6 -> 409 Conflict
	codes.Aborted:            http.StatusConflict,            // 10 -> 409 Conflict
	codes.ResourceExhausted:  http.StatusTooManyRequests,     // 8 -> 429 Too Many Requests
	codes.Canceled:           499,                            // 1 -> 499 Client Closed Request
	codes.Unknown:            http.StatusInternalServerError, // 2 -> 500 Internal Server Error
	codes.Internal:           http.StatusInternalServerError, // 13 -> 500 Internal Server Error
	codes.DataLoss:           http.StatusInternalServerError, // 15 -> 500 Internal Server Error
	codes.Unimplemented:      http.StatusNotImplemented,      // 12 -> 501 Not Implemented
	codes.Unavailable:        http.StatusServiceUnavailable,  // 14 -> 503 Service Unavailable
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,      // 4 -> 504 Gateway Timeout
}

func GRPCCodeToHTTPStatus(code codes.Code) int {
	if status, ok := grpcToHTTPStatus[code]; ok {
		return status
	}

	return http.StatusInternalServerError
}
