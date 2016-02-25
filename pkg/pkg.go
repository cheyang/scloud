// pkg.go
package pkg

import (
	"path/filepath"

	"github.com/cheyang/scloud/pkg/utils"
)

func GetDefaultStore(clusterName string) *persist.Filestore {
	homeDir := utils.GetHomeDir()
	clusterDir := filepath.Join(homeDir, ".scloud", clusterName)
	return &persist.Filestore{
		Path: filepath.Join(homeDir),
	}
}
