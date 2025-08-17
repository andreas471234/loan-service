package domain

import (
	"errors"
)

// StateTransition represents a valid state transition
type StateTransition struct {
	From   LoanStatus
	To     LoanStatus
	Action string
}

// FSM defines the finite state machine for loan states
type FSM struct {
	CurrentState LoanStatus
	Transitions  []StateTransition
}

// NewFSM creates a new FSM instance
func NewFSM() *FSM {
	return &FSM{
		CurrentState: StatusProposed,
		Transitions: []StateTransition{
			{From: StatusProposed, To: StatusApproved, Action: "approve"},
			{From: StatusApproved, To: StatusInvested, Action: "invest"},
			{From: StatusInvested, To: StatusDisbursed, Action: "disburse"},
		},
	}
}

// CanTransition checks if a transition is valid
func (fsm *FSM) CanTransition(to LoanStatus) bool {
	for _, transition := range fsm.Transitions {
		if transition.From == fsm.CurrentState && transition.To == to {
			return true
		}
	}
	return false
}

// Transition performs a state transition
func (fsm *FSM) Transition(to LoanStatus) error {
	if !fsm.CanTransition(to) {
		return errors.New("invalid state transition")
	}
	fsm.CurrentState = to
	return nil
}

// GetCurrentState returns the current state
func (fsm *FSM) GetCurrentState() LoanStatus {
	return fsm.CurrentState
}

// SetCurrentState sets the current state (used when loading from database)
func (fsm *FSM) SetCurrentState(state LoanStatus) {
	fsm.CurrentState = state
}

// GetValidTransitions returns all valid transitions from current state
func (fsm *FSM) GetValidTransitions() []StateTransition {
	var valid []StateTransition
	for _, transition := range fsm.Transitions {
		if transition.From == fsm.CurrentState {
			valid = append(valid, transition)
		}
	}
	return valid
}
