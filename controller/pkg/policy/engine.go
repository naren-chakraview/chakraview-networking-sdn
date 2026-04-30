package policy

import (
	"fmt"
	"sync"
)

type PolicyManager struct {
	mu     sync.RWMutex
	engine *PolicyEngine
}

func NewPolicyManager() *PolicyManager {
	return &PolicyManager{
		engine: NewPolicyEngine(),
	}
}

func (pm *PolicyManager) CreatePolicy(intent *PolicyIntent) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.engine.intents[intent.ID]; exists {
		return fmt.Errorf("policy %s already exists", intent.ID)
	}

	config, err := pm.engine.Translate(intent)
	if err != nil {
		return err
	}

	pm.engine.intents[intent.ID] = intent
	pm.engine.configs[intent.ID] = config

	fmt.Printf("Policy %s created and translated\n", intent.ID)
	return nil
}

func (pm *PolicyManager) DeletePolicy(id string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.engine.intents[id]; !exists {
		return fmt.Errorf("policy %s not found", id)
	}

	delete(pm.engine.intents, id)
	delete(pm.engine.configs, id)

	return nil
}

func (pm *PolicyManager) ListPolicies() []*PolicyIntent {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	policies := make([]*PolicyIntent, 0, len(pm.engine.intents))
	for _, p := range pm.engine.intents {
		policies = append(policies, p)
	}
	return policies
}

func (pm *PolicyManager) GetConfig(id string) *PolicyConfig {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.engine.configs[id]
}
