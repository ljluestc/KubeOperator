package prepare

import (
	"fmt"
	"io"

	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/adm/phases"
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
)

const (
	prepareLoadBalancer = "04-load-balancer.yml"

	// vipTroubleshootingHint is written to the task log when the load balancer
	// phase fails. Cloud networks commonly drop VRRP and unassigned floating
	// IPs, which leaves the VIP reachable only from the node that holds it.
	vipTroubleshootingHint = "Load balancer phase failed. If the log shows \"no route to host\" for the VIP, " +
		"the network probably does not carry VRRP or floating IPs between master nodes (typical on public clouds). " +
		"Check that the VIP is in the masters' subnet and unused, that VRRP (IP protocol 112) and ARP are allowed, " +
		"or use an external load balancer (lbMode=external) provided by the cloud."
)

type LoadBalancerPhase struct {
}

func (s LoadBalancerPhase) Name() string {
	return "Install Load Balancer"
}

func (s LoadBalancerPhase) Run(b kobe.Interface, writer io.Writer) error {
	err := phases.RunPlaybookAndGetResult(b, prepareLoadBalancer, "", writer)
	if err != nil && writer != nil {
		_, _ = fmt.Fprintln(writer, vipTroubleshootingHint)
	}
	return err
}
