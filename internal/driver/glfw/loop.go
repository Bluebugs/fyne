package glfw

import (
	"runtime"
	"sync"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/driver"
)

type funcData struct {
	f    func()
	done chan bool
}

// channel for queuing functions on the main thread
var funcQueue = make(chan funcData)
var runFlag = false
var runMutex = &sync.Mutex{}

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

func running() bool {
	runMutex.Lock()
	defer runMutex.Unlock()
	return runFlag
}

// force a function f to run on the main thread
func runOnMain(f func()) {
	// If we are on main just execute - otherwise add it to the main queue and wait.
	// The "running" variable is normally false when we are on the main thread.
	if !running() {
		f()
	} else {
		done := make(chan bool)

		funcQueue <- funcData{f: f, done: done}
		<-done
	}
}

func (d *gLDriver) repaintWindow(w *window) {
	canvas := w.canvas
	w.RunWithContext(func() {
		d.freeDirtyTextures(canvas)

		updateGLContext(w)
		if canvas.ensureMinSize() {
			w.fitContent()
		}
		canvas.paint(canvas.Size())

		w.viewport.SwapBuffers()
	})
}

func (d *gLDriver) freeDirtyTextures(canvas *glCanvas) {
	for {
		select {
		case object := <-canvas.refreshQueue:
			freeWalked := func(obj fyne.CanvasObject, _ fyne.Position, _ fyne.Position, _ fyne.Size) bool {
				canvas.painter.Free(obj)
				return false
			}
			driver.WalkCompleteObjectTree(object, freeWalked, nil)
		default:
			return
		}
	}
}
