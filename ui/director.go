package ui

import (
	"log"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/paked/nes/nes"
)

type View interface {
	Enter()
	Exit()
	Update(t, dt float64)

	Console() *nes.Console
}

type Director struct {
	window             *glfw.Window
	audio              *Audio
	controllerAdapter1 ControllerAdapter
	controllerAdapter2 ControllerAdapter
	view               View
	menuView           View
	timestamp          float64
}

func NewDirector(window *glfw.Window, audio *Audio, ca1, ca2 ControllerAdapter) *Director {
	director := Director{}
	director.window = window
	director.audio = audio
	director.controllerAdapter1 = ca1
	director.controllerAdapter2 = ca2

	return &director
}

func (d *Director) SetTitle(title string) {
	d.window.SetTitle(title)
}

func (d *Director) SetView(view View) {
	if d.view != nil {
		d.view.Exit()
	}
	d.view = view
	if d.view != nil {
		d.view.Enter()
	}
	d.timestamp = glfw.GetTime()
}

func (d *Director) Step() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	timestamp := glfw.GetTime()
	dt := timestamp - d.timestamp
	d.timestamp = timestamp
	if d.view != nil {
		d.view.Update(timestamp, dt)
	}
}

func (d *Director) Start(paths []string) {
	d.menuView = NewMenuView(d, paths)
	if len(paths) == 1 {
		d.PlayGame(paths[0])
	} else {
		d.ShowMenu()
	}
	d.Run()
}

func (d *Director) Run() {
	for !d.window.ShouldClose() {
		d.Step()
		d.window.SwapBuffers()
		glfw.PollEvents()
	}
	d.SetView(nil)
}

func (d *Director) PlayGame(path string) {
	hash, err := hashFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}
	d.SetView(NewGameView(d, console, path, hash))
}

func (d *Director) ShowMenu() {
	d.SetView(d.menuView)
}

func (d *Director) Window() *glfw.Window {
	return d.window
}

func (d *Director) Console() *nes.Console {
	return d.view.Console()
}
