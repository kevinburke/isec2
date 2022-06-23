package isec2

import (
	"context"
	"net"
	"testing"
)

func TestIsEC2(t *testing.T) {
	t.Run("HypervisorUUID", func(t *testing.T) {
		testPrefix = "testdata/test1"
		defer func() { testPrefix = "" }()
		yes, err := IsEC2(context.Background())
		if err != nil {
			t.Fatal("expected non-nil err, got nil")
		}
		if !yes {
			t.Fatal("should have gotten yes, did not get it")
		}
	})

	t.Run("ProductUUID", func(t *testing.T) {
		testPrefix = "testdata/test2"
		defer func() { testPrefix = "" }()
		yes, err := IsEC2(context.Background())
		if err != nil {
			t.Fatal("expected non-nil err, got nil")
		}
		if !yes {
			t.Fatal("should have gotten yes, did not get it")
		}
	})

	t.Run("BoardAssetTag", func(t *testing.T) {
		testPrefix = "testdata/test3"
		defer func() { testPrefix = "" }()
		yes, err := IsEC2(context.Background())
		if err != nil {
			t.Fatal("expected non-nil err, got nil")
		}
		if !yes {
			t.Fatal("should have gotten yes, did not get it")
		}
	})

	t.Run("Nothing", func(t *testing.T) {
		testPrefix = "testdata/test4"
		oldHost := ec2APIHost
		ec2APIHost = net.JoinHostPort("10.10.10.11", "80")
		defer func() { testPrefix = ""; ec2APIHost = oldHost }()
		yes, err := IsEC2(context.Background())
		if err != nil {
			t.Fatal("expected non-nil err, got nil")
		}
		if yes {
			t.Fatal("should not have gotten yes, but got it")
		}
	})
}
