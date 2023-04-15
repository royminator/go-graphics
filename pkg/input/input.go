package input

import (
	"github.com/veandco/go-sdl2/sdl"
)

type (
	InputMapper struct {
		activeCtx InputContext
		entities  map[Action][]ExecsInput
		contexts  []InputContext
	}

	InputState struct {
		triggeredButtons []ButtonState
	}

	ButtonState struct {
		Id   ButtonId
		Mode ButtonMode
	}

	InputContext struct {
		actions map[Action][]ButtonState
	}

	ExecsInput interface {
		Id() int
		ExecInput(Action)
	}

	Action     int
	ButtonId   byte
	ButtonMode int
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

func ReadAndExecInputs(mapper InputMapper) {
	inputState := readInputs()
	mapper.mapAndExec(inputState)
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
		state.triggeredButtons = append(state.triggeredButtons, btn)
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

func (m *InputMapper) mapAndExec(state InputState) {
	actions := m.activeCtx.mapToDomainInputs(state)
	m.execInputs(actions)
}

func (ctx InputContext) mapToDomainInputs(state InputState) []Action {
	var triggeredActions []Action
	for action, btns := range ctx.actions {
		if isActionTriggered(btns, state.triggeredButtons) {
			triggeredActions = append(triggeredActions, action)
		}
	}
	return triggeredActions
}

func isActionTriggered(btns, triggered []ButtonState) bool {
	for _, btn := range btns {
		if !containsBtn(btn, triggered) {
			return false
		}
	}
	return true
}

func containsBtn(btn ButtonState, btns []ButtonState) bool {
	for _, b := range btns {
		if b == btn {
			return true
		}
	}
	return false
}

func (m InputMapper) execInputs(actions []Action) {
	for _, action := range actions {
		for _, entity := range m.entities[action] {
			entity.ExecInput(action)
		}
	}
}

func MakeMapper() InputMapper {
	ctx := InputContext{
		actions: map[Action][]ButtonState{
			MOVE_NORTH: {{W, PRESSED}},
		},
	}

	return InputMapper{
		activeCtx: ctx,
		entities:  map[Action][]ExecsInput{},
		contexts:  []InputContext{ctx},
	}
}

func (m *InputMapper) RegisterEntity(action Action, entity ExecsInput) {
	entities, found := m.entities[action]
	if found {
		entities = append(entities, entity)
		m.entities[action] = entities
		return
	}
	m.entities[action] = []ExecsInput{entity}
}

func (m *InputMapper) RemoveEntity(id int) {
	for action, entities := range m.entities {
		removed := removeEntityWithId(entities, id)
		m.entities[action] = removed
	}
}

func removeEntityWithId(entities []ExecsInput, id int) []ExecsInput {
	for i, entity := range entities {
		if entity.Id() == id {
			return removeEntityIndex(entities, i)
		}
	}
	return entities
}

func removeEntityIndex(entities []ExecsInput, i int) []ExecsInput {
	return append(entities[:i], entities[i+1:]...)
}
