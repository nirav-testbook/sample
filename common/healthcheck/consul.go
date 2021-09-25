package healthcheck

import (
	"time"

	kitlog "github.com/go-kit/kit/log"
	consul "github.com/hashicorp/consul/api"
)

func InitConsulHealthCheck(agent *consul.Agent, logger kitlog.Logger, svcId string, ttl time.Duration) {
	for {
		err := agent.UpdateTTL(svcId, "OK", consul.HealthPassing)
		if err != nil {
			logger.Log("check failed", err)
		}
		time.Sleep(ttl)
	}
}
