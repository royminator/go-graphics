package input

import (
	"github.com/veandco/go-sdl2/sdl"
)

type (
	AppState struct {
		InputContext Context
	}

	InputState struct {
		TriggeredButtons []ButtonState
	}

	ButtonState struct {
		Id   ButtonId
		Mode ButtonMode
	}

	ButtonId   byte
	ButtonMode int

	Context struct {
		Actions map[Action][]ButtonState
	}

	Action int
)

const (
	W ButtonId = iota
	A
	R
	S
	Q
	ESC

	PRESSED ButtonMode = iota
	RELEASED

	MOVE_NORTH Action = iota
	MOVE_WEST
	MOVE_SOUTH
	MOVE_EAST
	QUIT
)

func ReadAndExecInputs(state *AppState) {
	for {
		inputState := readInputs()
		state.InputContext.inputStateToDomainInputs(inputState)
		// executeInputs(domainInputs, state)
	}
}

func readInputs() InputState {
	var state InputState
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		state.recordSDLEvent(event)
	}
	return state
}

func (state *InputState) recordSDLEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		state.recordKeyboardEvent(t)
	}
}

func (state *InputState) recordKeyboardEvent(event *sdl.KeyboardEvent) {
	if btn, found := sdlEventToButtonState(event); found {
		state.TriggeredButtons = append(state.TriggeredButtons, btn)
	}
}

func sdlEventToButtonState(event *sdl.KeyboardEvent) (ButtonState, bool) {
	id, idFound := sdlKeycodeToButtonId(event.Keysym.Sym)
	mode := sdlKeyStateToButtonMode(event.State)

	if idFound {
		return ButtonState{id, mode}, true
	}
	return ButtonState{}, false
}

func sdlKeycodeToButtonId(code sdl.Keycode) (ButtonId, bool) {
	switch code {
	case sdl.K_w:
		return W, true
	case sdl.K_a:
		return A, true
	case sdl.K_r:
		return R, true
	case sdl.K_s:
		return S, true
	case sdl.K_q:
		return Q, true
	case sdl.K_ESCAPE:
		return ESC, true
	default:
		return 0, false
	}
}

func sdlKeyStateToButtonMode(state uint8) ButtonMode {
	if state == sdl.PRESSED {
		return PRESSED
	} else {
		return RELEASED
	}
}

func (ctx Context) inputStateToDomainInputs(state InputState) []Action {
	var triggeredActions []Action
	for action, btns := range ctx.Actions {
		if isActionTriggered(btns, state.TriggeredButtons) {
			triggeredActions = append(triggeredActions, action)
		}
	}
	return triggeredActions
}

func isActionTriggered(btns, triggered []ButtonState) bool {
	var count int = 0
	for _, btn := range btns {
		if containsBtn(btn, triggered) {
			count++
		}
	}

	if count == len(btns) {
		return true
	}
	return false
}

func containsBtn(btn ButtonState, btns []ButtonState) bool {
	for _, b := range btns {
		if b == btn {
			return true
		}
	}

	return false
}
