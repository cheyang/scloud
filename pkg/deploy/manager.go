// manager.go
package deploy

import (
	"sync"

	"github.com/cheyang/scloud/pkg/msg"
)

type Manager struct {
	PlannerObserver msg.Receiver

	lock sync.RWMutex

	ToDeploy Deployment // Operate on map, no need pointer

	latestDeploy Deployment

	working bool

	FinishReport chan interface{}

	worker Deployer
}

func (m *Manager) Deploy() {
	if !m.isWorking() {
		return
	}
}

func (m *Manager) getDeployment() Deployment {
	m.lock.Lock()
	defer m.lock.Unlock()

	deployment := m.ToDeploy

	m.ToDeploy = Deployment{}

	return deployment
}

func (m *Manager) isWorking() bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.working

}
