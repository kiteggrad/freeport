package freeport

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func Test_Get(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	require := require.New(t)

	const triesCount = 1000
	ports := make(map[Port]bool, triesCount)

	for i := 0; i < triesCount; i++ {
		port, err := Get()
		require.NoError(err)

		require.False(ports[port])
		ports[port] = true
	}
}

func Test_Retry(t *testing.T) {
	t.Cleanup(func() { goleak.VerifyNone(t) })

	require := require.New(t)

	const triesCount = 2
	ports := make(map[Port]bool, triesCount)
	var lastRequestedPort Port

	i := 0
	port, err := NewRetrier(nil, nil).Retry(t.Context(), func(port Port) error {
		defer func() { i++ }()

		require.False(ports[port])
		ports[port] = true
		lastRequestedPort = port

		if i == triesCount {
			return nil
		}

		return errors.New("some")
	})
	require.NoError(err)
	require.Equal(lastRequestedPort, port)
}
