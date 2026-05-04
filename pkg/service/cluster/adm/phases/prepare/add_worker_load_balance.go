package prepare

import (
	"fmt"
	"io"
	"time"

	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareAddWorkerLoadBalancer = "91-add-worker-04-load-balancer.yml"
)

type AddWorkerLoadBalancerPhase struct {
}

func (s AddWorkerLoadBalancerPhase) Name() string {
	return "Install Load Balancer"
}

func (s AddWorkerLoadBalancerPhase) Run(b kobe.Interface, writer io.Writer) error {
	var lastErr error
	for i := 1; i <= loadBalancerMaxAttempts; i++ {
		if err := phases.RunPlaybookAndGetResult(b, prepareAddWorkerLoadBalancer, "", writer); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if i < loadBalancerMaxAttempts {
			if writer != nil {
				_, _ = fmt.Fprintf(
					writer,
					"add-worker load balancer phase failed (attempt %d/%d): %v, retrying in %s...\n",
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
