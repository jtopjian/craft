package testhelper

import (
	"github.com/jtopjian/craft/client"
	"github.com/sirupsen/logrus"
)

func TestClient() client.Client {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	c := client.Client{
		Logger: logger,
	}

	return c
}
