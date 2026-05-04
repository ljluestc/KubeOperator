package prepare

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"io"
	"time"
)

const (
	prepareLoadBalancer      = "04-load-balancer.yml"
	loadBalancerMaxAttempts  = 3
	loadBalancerRetryBackoff = 10 * time.Second
)

type LoadBalancerPhase struct {
}

func (s LoadBalancerPhase) Name() string {
	return "Install Load Balancer"
}

func (s LoadBalancerPhase) Run(b kobe.Interface, writer io.Writer) error {
	var lastErr error
	for i := 1; i <= loadBalancerMaxAttempts; i++ {
		if err := phases.RunPlaybookAndGetResult(b, prepareLoadBalancer, "", writer); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < loadBalancerMaxAttempts {
			if writer != nil {
				_, _ = fmt.Fprintf(
					writer,
					"load balancer phase failed (attempt %d/%d): %v, retrying in %s...\n",
					i,
					loadBalancerMaxAttempts,
					lastErr,
					loadBalancerRetryBackoff,
				)
			}
			time.Sleep(loadBalancerRetryBackoff)
		}
	}

	return lastErr
}
