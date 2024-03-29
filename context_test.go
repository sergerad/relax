package relax

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_Signals_CancelContext(t *testing.T) {
	for _, signal := range []string{"SIGINT", "SIGTERM"} {
		t.Run(signal, func(t *testing.T) {
			ctx := Context()
			go func() {
				cmd := exec.Command("kill", "-"+signal, fmt.Sprint(os.Getpid()))
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				require.NoError(t, cmd.Run())
			}()
			// Exit only if signal causes cancel
			<-ctx.Done()
			t.Log("Exiting gracefully")
		})
	}
}
