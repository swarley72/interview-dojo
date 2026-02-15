package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("failed encode json", "err", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func handleGRPCError(w http.ResponseWriter, err error, logMsg string) {
	s, _ := status.FromError(err)
	httpStatus := grpcCodeToHTTPStatus(s.Code())

	if httpStatus == http.StatusInternalServerError {
		slog.Error(logMsg, "error", err)
		writeError(w, httpStatus, "internal server error")
		return
	}

	writeError(w, httpStatus, s.Message())
}

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
