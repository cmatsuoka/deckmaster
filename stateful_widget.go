package main

import (
	"fmt"
)

type WidgetState interface{}

// StatefulWidget is a widget that keeps track of its internal state.
type StatefulWidget struct {
	*BaseWidget

	state_id string
	state    WidgetState
}

// WidgetStateRegistry maintains a list of registered widget states.
type WidgetStateRegistry struct {
	widgetStates map[string]WidgetState
}

func NewWidgetStateRegistry() *WidgetStateRegistry {
	return &WidgetStateRegistry{widgetStates: map[string]WidgetState{}}
}

func (r *WidgetStateRegistry) Register(state_id string, state WidgetState) {
	verbosef("Register widget state id %q", state_id)
	r.widgetStates[state_id] = state
}

func (r *WidgetStateRegistry) Recover(state_id string) (WidgetState, error) {
	state, ok := r.widgetStates[state_id]
	if !ok {
		return nil, fmt.Errorf("invalid widget state id %q", state_id)
	}
	verbosef("Recover widget state %q", state_id)
	return state, nil
}

// NewStatefulWidget creates a widget with a registered state.
func NewStatefulWidget(bw *BaseWidget) *StatefulWidget {
	return &StatefulWidget{BaseWidget: bw}
}
