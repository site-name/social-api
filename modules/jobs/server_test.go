package jobs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartWorkers(t *testing.T) {
	t.Run("uninitialized", func(t *testing.T) {
		jobServer, _, _ := makeJobServer(t)
		err := jobServer.StartWorkers()
		require.Equal(t, ErrWorkersUninitialized, err)
	})
}
