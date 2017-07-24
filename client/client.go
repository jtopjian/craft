package client

import (
	"github.com/sirupsen/logrus"
)

// Client represents a system client. At the moment, this is only
// to configure a global logger.
type Client struct {
	Logger *logrus.Logger
}
