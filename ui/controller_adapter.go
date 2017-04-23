package ui

import (
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/paked/nes/nes"
)

type ControllerAdapter interface {
	Buttons() [8]bool

	SetWindow(*glfw.Window)

	Trigger(int, bool)
}

type KeyboardControllerAdapter struct {
	*BasicControllerAdapter
}

func NewKeyboardControllerAdapter() *KeyboardControllerAdapter {
	return &KeyboardControllerAdapter{
		&BasicControllerAdapter{},
	}
}

func NewKeyboardControllerAdapterWithWindow(window *glfw.Window) *KeyboardControllerAdapter {
	kba := NewKeyboardControllerAdapter()

	kba.SetWindow(window)

	return kba
}

func (kba *KeyboardControllerAdapter) SetWindow(window *glfw.Window) {
	kba.BasicControllerAdapter.SetWindow(window)

	bindings := map[glfw.Key]int{
		glfw.KeyZ:          nes.ButtonA,
		glfw.KeyX:          nes.ButtonB,
		glfw.KeyRightShift: nes.ButtonSelect,
		glfw.KeyEnter:      nes.ButtonStart,
		glfw.KeyUp:         nes.ButtonUp,
		glfw.KeyDown:       nes.ButtonDown,
		glfw.KeyLeft:       nes.ButtonLeft,
		glfw.KeyRight:      nes.ButtonRight,
	}

	kba.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var v bool

		if action == glfw.Press {
			v = true
		} else if action == glfw.Release {
			v = false
		} else {
			return
		}

		button, ok := bindings[key]
		if !ok {
			return
		}

		kba.Trigger(button, v)
	})
}

func (kba *KeyboardControllerAdapter) Trigger(button int, newState bool) {
	kba.state[button] = newState
}

func (kba *KeyboardControllerAdapter) Buttons() [8]bool {
	return kba.state
}

type DummyControllerAdapter struct{}

func NewDummyControllerAdapter() *DummyControllerAdapter {
	return &DummyControllerAdapter{}
}

func (*DummyControllerAdapter) Buttons() [8]bool {
	return [8]bool{}
}

func (*DummyControllerAdapter) Trigger(_ int, _ bool)    {}
func (*DummyControllerAdapter) SetWindow(_ *glfw.Window) {}

type BasicControllerAdapter struct {
	state  [8]bool
	window *glfw.Window
}

func (bca *BasicControllerAdapter) Trigger(button int, newState bool) {
	bca.state[button] = newState
}

func (bca *BasicControllerAdapter) SetWindow(window *glfw.Window) {
	bca.window = window
}

func (bca *BasicControllerAdapter) Buttons() [8]bool {
	return bca.state
}
