package policy

type PolicyIntent struct {
	ID       string
	Name     string
	Type     string /* "routing", "isolation", "qos" */
	Source   string
	Dest     string
	Action   string /* "allow", "deny", "redirect" */
	Priority int
}

type PolicyConfig struct {
	Rules []PolicyRule
	ACLs  []ACLEntry
}

type PolicyRule struct {
	ID   string
	From string
	To   string
	Action string
	Tags map[string]string
}

type ACLEntry struct {
	ID       string
	SourceIP string
	DestIP   string
	Action   string
}

type PolicyEngine struct {
	intents map[string]*PolicyIntent
	configs map[string]*PolicyConfig
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{
		intents: make(map[string]*PolicyIntent),
		configs: make(map[string]*PolicyConfig),
	}
}

func (pe *PolicyEngine) AddIntent(intent *PolicyIntent) error {
	pe.intents[intent.ID] = intent
	return nil
}

func (pe *PolicyEngine) Translate(intent *PolicyIntent) (*PolicyConfig, error) {
	config := &PolicyConfig{
		Rules: make([]PolicyRule, 0),
		ACLs:  make([]ACLEntry, 0),
	}

	rule := PolicyRule{
		ID:     intent.ID,
		From:   intent.Source,
		To:     intent.Dest,
		Action: intent.Action,
		Tags: map[string]string{
			"priority": string(rune(intent.Priority)),
		},
	}
	config.Rules = append(config.Rules, rule)

	return config, nil
}

func (pe *PolicyEngine) GetConfig(id string) *PolicyConfig {
	return pe.configs[id]
}
