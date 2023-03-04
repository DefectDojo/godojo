package distros

import (
	"fmt"
	"strings"

	c "github.com/mtesauro/commandeer"
)

func CmdsForTarget(cp *c.CmdPkg, t string) ([]c.SingleCmd, error) {
	// Cycle through Ubuntu install targets
	for k := range cp.Targets {
		if strings.Compare(
			strings.ToLower(cp.Targets[k].ID),
			strings.ToLower(t)) == 0 {
			// Return the commands matching that target
			return cp.Targets[k].PkgCmds, nil
		}
	}

	return make([]c.SingleCmd, 1), fmt.Errorf("Unable to find commands for OS target %s\n", t)
}
