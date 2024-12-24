package app

import "github.com/viwet/GoBIP39/words"

type configConstraint interface {
	DepositConfig | BLSToExecutionConfig
}

type State[Config configConstraint] struct {
	cfg *Config

	mnemonic []string
	list     words.List
	password string
}

func NewState[Config configConstraint](cfg *Config) *State[Config] {
	return &State[Config]{cfg: cfg}
}

func (s *State[Config]) WithMnemonic(menmonic []string, list words.List) *State[Config] {
	s.mnemonic = menmonic
	s.list = list
	return s
}

func (s *State[Config]) WithPassword(password string) *State[Config] {
	s.password = password
	return s
}

func (s *State[Config]) Mnemonic() []string {
	return s.mnemonic
}

func (s *State[Config]) Words() words.List {
	return s.list
}

func (s *State[Config]) Config() *Config {
	return s.cfg
}
