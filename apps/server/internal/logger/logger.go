package logger

import (
	"log/slog"
	"os"
	"strings"
)

func New(appName, level string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	})).With("app", appName)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Sanitize(value string) string {
	replacer := strings.NewReplacer("\r", "", "\n", "", "\t", "")
	return replacer.Replace(value)
}

func SysStartup(log *slog.Logger, port string) {
	log.Info("system startup", slog.String("event_type", "sys_startup"), slog.String("port", port))
}

func SysShutdown(log *slog.Logger) {
	log.Info("system shutdown", slog.String("event_type", "sys_shutdown"))
}

func SysCrash(log *slog.Logger, cause error) {
	log.Error("system crash", slog.String("event_type", "sys_crash"), slog.Any("error", cause))
}

func InputValidationFail(log *slog.Logger, requestID, detail string) {
	log.Warn("input validation failed",
		slog.String("event_type", "input_validation_fail"),
		slog.String("request_id", requestID),
		slog.String("detail", Sanitize(detail)),
	)
}

func AuthnLoginFail(log *slog.Logger, requestID, userID string) {
	log.Warn("authentication failed",
		slog.String("event_type", "authn_login_fail"),
		slog.String("request_id", requestID),
		slog.String("user_id", Sanitize(userID)),
	)
}

func AuthzFail(log *slog.Logger, requestID, userID, resource string) {
	log.Warn("authorization failed",
		slog.String("event_type", "authz_fail"),
		slog.String("request_id", requestID),
		slog.String("user_id", Sanitize(userID)),
		slog.String("resource", Sanitize(resource)),
	)
}
