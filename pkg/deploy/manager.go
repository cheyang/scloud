// manager.go
package deploy

import (
	"sync"
)

type Manager struct {
	PlannerObserver msg.Receiver

	lock sync.RWMutex

	ToDeploy Deployment // Operate on map, no need pointer

	latestDeploy Deployment

	isWorking bool

	FinishReport chan interface{}

	worker deployers.Deployer
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

}

func (m *Manager) IsWorking() bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	return isWorking

}
