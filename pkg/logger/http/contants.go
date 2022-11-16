package http

import "github.com/evrone/go-clean-template/pkg/logger"

const ContextLoggerField = "logger"

const (
    RequestId logger.Field = "request_id"
    Method    logger.Field = "method"
    URL       logger.Field = "url"
)
