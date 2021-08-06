package dupont

import (
	"github.com/lorenzosaino/go-sysctl"
	"go.uber.org/zap"
)

var (
	sysctlValues = map[string]string{
		"net.ipv4.ip_forward": "1",
	}
)

func ensureSysctl(log *zap.Logger) error {
	for k, v := range sysctlValues {
		err := sysctl.Set(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
