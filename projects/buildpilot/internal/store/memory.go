package store

import (
	"sync"

	"github.com/ranjith-ka/buildpilot/internal/domain"
)

type Memory struct {
	mu     sync.RWMutex
	agents map[string]domain.Agent
	jobs   []domain.Job
}

func NewMemory() *Memory {
	return &Memory{agents: make(map[string]domain.Agent)}
}

func (m *Memory) UpsertAgent(agent domain.Agent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.agents[agent.ID] = agent
}

func (m *Memory) Agent(id string) (domain.Agent, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agent, ok := m.agents[id]
	return agent, ok
}

func (m *Memory) Enqueue(job domain.Job) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs = append(m.jobs, job)
}

func (m *Memory) NextJob(agentID string) (domain.Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for index := range m.jobs {
		if m.jobs[index].AgentID == agentID && m.jobs[index].Status == domain.JobQueued {
			m.jobs[index].Status = domain.JobPreparing
			return m.jobs[index], true
		}
	}
	return domain.Job{}, false
}
