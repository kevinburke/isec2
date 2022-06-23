package isec2

import (
	"bytes"
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var testPrefix = ""

// 8 to 17 characters in an instance id
var instanceRe = regexp.MustCompile(`^i-[0-9a-f]{8,20}$`)

// occasionally overridden in tests
var ec2APIHost = net.JoinHostPort("169.254.169.254", "80")

var ErrCouldNotDetermine = errors.New("usual methods of checking EC2 availability all failed")

// IsEC2 reports whether you are running in EC2. We attempt to do this in
// a timely manner, ie. we may set a shorter timeout than is provided.
//
// IsEC2 returns an error if the status could not be determined.
func IsEC2(ctx context.Context) (bool, error) {
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/identify_ec2_instances.html
	// https://serverfault.com/a/903599

	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	f, err := os.Open(filepath.Join(testPrefix, "/sys/devices/virtual/dmi/id/board_asset_tag"))
	if err == nil {
		deadline, _ := ctx.Deadline()
		f.SetDeadline(deadline)
		p := make([]byte, 25)
		n, err := f.Read(p)
		if err == nil {
			p = bytes.TrimSpace(p[:n])
			if instanceRe.Match(p) {
				f.Close()
				return true, nil
			}
		}
	}
	f.Close()

	f, err = os.Open(filepath.Join(testPrefix, "/sys/hypervisor/uuid"))
	if err == nil {
		deadline, _ := ctx.Deadline()
		f.SetDeadline(deadline)
		var p [3]byte
		if _, err := f.Read(p[:]); err == nil && bytes.EqualFold(p[:], []byte("ec2")) {
			f.Close()
			return true, nil
		}
	}
	f.Close()

	f, err = os.Open(filepath.Join(testPrefix, "/sys/devices/virtual/dmi/id/product_uuid"))
	if err == nil {
		deadline, _ := ctx.Deadline()
		f.SetDeadline(deadline)
		var p [3]byte
		if _, err := f.Read(p[:]); err == nil && bytes.EqualFold(p[:], []byte("ec2")) {
			f.Close()
			return true, nil
		}
	}
	f.Close()

	// try to connect to the EC2 metadata API
	_, err = (&net.Dialer{}).DialContext(ctx, "tcp", ec2APIHost)
	if err == nil {
		return true, nil
	}
	operr, ok := err.(*net.OpError)
	if !ok {
		return false, ErrCouldNotDetermine
	}
	if operr.Timeout() {
		// can't hit EC2 API
		return false, nil
	}
	return false, ErrCouldNotDetermine
}
