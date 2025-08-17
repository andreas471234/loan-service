package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSMNew(t *testing.T) {
	fsm := NewFSM()
	assert.NotNil(t, fsm)
	assert.Equal(t, StatusProposed, fsm.GetCurrentState())
}

func TestFSMValidTransitions(t *testing.T) {
	fsm := NewFSM()
	fsm.SetCurrentState(StatusProposed)

	// Test valid transition
	assert.True(t, fsm.CanTransition(StatusApproved))
	
	err := fsm.Transition(StatusApproved)
	assert.NoError(t, err)
	assert.Equal(t, StatusApproved, fsm.GetCurrentState())

	// Test invalid transition
	assert.False(t, fsm.CanTransition(StatusProposed))
	
	err = fsm.Transition(StatusProposed)
	assert.Error(t, err)
}

func TestFSMGetValidTransitions(t *testing.T) {
	fsm := NewFSM()
	fsm.SetCurrentState(StatusProposed)

	transitions := fsm.GetValidTransitions()
	assert.Len(t, transitions, 1)
	assert.Equal(t, StatusApproved, transitions[0].To)
	assert.Equal(t, "approve", transitions[0].Action)

	fsm.SetCurrentState(StatusApproved)
	transitions = fsm.GetValidTransitions()
	// Approved state has transition to invested
	assert.Len(t, transitions, 1)
	assert.Equal(t, StatusInvested, transitions[0].To)
	assert.Equal(t, "invest", transitions[0].Action)

	fsm.SetCurrentState(StatusInvested)
	transitions = fsm.GetValidTransitions()
	// Invested state has transition to disbursed
	assert.Len(t, transitions, 1)
	assert.Equal(t, StatusDisbursed, transitions[0].To)
	assert.Equal(t, "disburse", transitions[0].Action)

	fsm.SetCurrentState(StatusDisbursed)
	transitions = fsm.GetValidTransitions()
	// Disbursed state has no transitions
	assert.Len(t, transitions, 0)
}

func TestFSMCompleteLifecycle(t *testing.T) {
	fsm := NewFSM()
	
	// Proposed -> Approved
	assert.True(t, fsm.CanTransition(StatusApproved))
	err := fsm.Transition(StatusApproved)
	assert.NoError(t, err)
	assert.Equal(t, StatusApproved, fsm.GetCurrentState())

	// Approved -> Invested
	assert.True(t, fsm.CanTransition(StatusInvested))
	err = fsm.Transition(StatusInvested)
	assert.NoError(t, err)
	assert.Equal(t, StatusInvested, fsm.GetCurrentState())

	// Invested -> Disbursed
	assert.True(t, fsm.CanTransition(StatusDisbursed))
	err = fsm.Transition(StatusDisbursed)
	assert.NoError(t, err)
	assert.Equal(t, StatusDisbursed, fsm.GetCurrentState())
}
