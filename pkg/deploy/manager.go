// manager.go
package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/cheyang/scloud/pkg/msg"
)

type Manager struct {
	PlannerObserver msg.Receiver

	lock sync.RWMutex

	ToDeploy Deployment // Operate on map, no need pointer

	latestDeploy Deployment

	FinishReport chan interface{}

	worker Deployer

	Workspace string
}

func (m *Manager) Deploy() {
	m.lock.Lock()
	defer m.lock.Unlock()

	current_deployment := m.setDeployment()

	if current_deployment.IsEmpty() {
		return
	}

	m.worker.Deploy(current_deployment, m.createWorkerDir())

}

func (m *Manager) setDeployment() Deployment {

	deployment := m.ToDeploy

	m.latestDeploy = deployment

	m.ToDeploy = Deployment{}

	return deployment
}

func (m *Manager) createWorkerDir() string {

	currentTime := time.Now().Unix()

	tm := time.Unix(currentTime, 0)

	timestamp := tm.Format("20060102150405")

	baseName := fmt.Sprintf("scloud_%s", timestamp)

	workingDir := filepath.Join(m.Workspace, baseName)

	err := os.MkdirAll(workingDir, 0744)

	if err != nil {
		fmt.Fprintf(os.stdErr, "create work dir error: %v", err)
	}

	return workingDir
}
