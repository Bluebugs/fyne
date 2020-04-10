// +build wasm js

package dialog


import (
        "fyne.io/fyne"
)

func (f *fileDialog) loadPlaces() []fyne.CanvasObject {
        return nil
}

func isHidden(file, _ string) bool {
        return false
}

func fileOpenOSOverride(f *FileDialog) bool {
        return true
}

func fileSaveOSOverride(f *FileDialog) bool {
        return true
}
