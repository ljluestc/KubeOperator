package model

import (
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
)

func TestResolveLBRole(t *testing.T) {
	t.Run("worker node has empty lb role and does not consume master slot", func(t *testing.T) {
		role, next := resolveLBRole(constant.NodeRoleNameWorker, 0)
		if role != "" {
			t.Fatalf("expected role %q, got %q", "", role)
		}
		if next != 0 {
			t.Fatalf("expected master count 0, got %d", next)
		}
	})

	t.Run("first master gets master lb role", func(t *testing.T) {
		role, next := resolveLBRole(constant.NodeRoleNameMaster, 0)
		if role != "master" {
			t.Fatalf("expected role %q, got %q", "master", role)
		}
		if next != 1 {
			t.Fatalf("expected master count 1, got %d", next)
		}
	})

	t.Run("subsequent masters get backup lb role", func(t *testing.T) {
		role, next := resolveLBRole(constant.NodeRoleNameMaster, 1)
		if role != "backup" {
			t.Fatalf("expected role %q, got %q", "backup", role)
		}
		if next != 2 {
			t.Fatalf("expected master count 2, got %d", next)
		}
	})

	t.Run("worker before master still leaves first master as master", func(t *testing.T) {
		count := 0
		role, next := resolveLBRole(constant.NodeRoleNameWorker, count)
		if role != "" || next != 0 {
			t.Fatalf("expected worker result (%q,0), got (%q,%d)", "", role, next)
		}

		role, next = resolveLBRole(constant.NodeRoleNameMaster, next)
		if role != "master" || next != 1 {
			t.Fatalf("expected first master result (%q,1), got (%q,%d)", "master", role, next)
		}
	})
}
