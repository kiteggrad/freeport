package freeport

import (
	"fmt"
	"net"
	"strconv"
)

type Port int

func (p Port) String() string {
	return strconv.Itoa(int(p))
}

func (p Port) Int() int {
	return int(p)
}

// Get - returns free port that wasn't used before (from this package).
func Get() (port Port, err error) {
	return globGenerator.Get()
}

// MustGet - returns free port that wasn't used before (from this package) or panics.
func MustGet() (port Port) {
	return globGenerator.MustGet()
}

func tryGetFreeport() (port Port, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to net.ResolveTCPAddr: %w", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to net.ListenTCP: %w", err)
	}
	defer l.Close()

	return Port(l.Addr().(*net.TCPAddr).Port), nil
}
