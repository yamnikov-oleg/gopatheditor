package main

import (
	"github.com/yamnikov-oleg/w32"
	"github.com/yamnikov-oleg/wingo"
	"strings"
)

var (
	WND_WIDTH  int = 480
	WND_HEIGHT int = 540
)

const (
	CTRL_INTERVAL = 10
	TE_HEIGHT     = 25
	BUTTON_WIDTH  = 90
)

var (
	regPathKey w32.HKEY

	wnd *wingo.Window

	listbox                                              *wingo.ListBox
	teUpdate, teAppend                                   *wingo.TextEdit
	appendButton, deleteButton, saveButton, reloadButton *wingo.Button

	paths     []string
	selection = -1
)

func updateMetrics() {
	clSize := wnd.GetClientSize()
	WND_WIDTH = clSize.X
	WND_HEIGHT = clSize.Y

	listbox.SetPosition(wingo.Vector{CTRL_INTERVAL, CTRL_INTERVAL})
	listbox.SetSize(wingo.Vector{WND_WIDTH - 2*CTRL_INTERVAL, WND_HEIGHT - TE_HEIGHT*2 - CTRL_INTERVAL*4})

	teUpdate.SetPosition(wingo.Vector{CTRL_INTERVAL, WND_HEIGHT - TE_HEIGHT*2 - CTRL_INTERVAL*2})
	teUpdate.SetSize(wingo.Vector{WND_WIDTH - BUTTON_WIDTH - CTRL_INTERVAL*3, TE_HEIGHT})

	deleteButton.SetPosition(wingo.Vector{WND_WIDTH - BUTTON_WIDTH - CTRL_INTERVAL, WND_HEIGHT - TE_HEIGHT*2 - CTRL_INTERVAL*2})
	deleteButton.SetSize(wingo.Vector{BUTTON_WIDTH, TE_HEIGHT})

	teAppend.SetPosition(wingo.Vector{CTRL_INTERVAL, WND_HEIGHT - TE_HEIGHT - CTRL_INTERVAL})
	teAppend.SetSize(wingo.Vector{WND_WIDTH - BUTTON_WIDTH - CTRL_INTERVAL*3, TE_HEIGHT})

	appendButton.SetPosition(wingo.Vector{WND_WIDTH - BUTTON_WIDTH - CTRL_INTERVAL, WND_HEIGHT - TE_HEIGHT - CTRL_INTERVAL})
	appendButton.SetSize(wingo.Vector{BUTTON_WIDTH, TE_HEIGHT})

}

func loadPaths() {
	path := w32.RegGetString(regPathKey, "", "Path")
	paths = strings.Split(path, ";")
	listbox.SetList(paths)
	selection = -1
	teUpdate.SetText("")
}

func main() {
	regPathKey = w32.RegOpenKeyEx(w32.HKEY_LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, w32.KEY_ALL_ACCESS)

	wnd = wingo.NewWindow(true, true)
	wnd.SetTitle("Редактор PATH")
	wnd.SetSize(wingo.Vector{WND_WIDTH, WND_HEIGHT})
	wnd.Show()
	wnd.OnSizeChanged = func(w *wingo.Window, size wingo.Vector) {
		updateMetrics()
	}

	listbox = wnd.NewListBox()
	listbox.MakeDraggable()
	listbox.OnSelected = func(lb *wingo.ListBox, text string, index int) {
		selection = index
		teUpdate.SetText(text)
	}

	teUpdate = wnd.NewTextEdit()
	teUpdate.OnChange = func(te *wingo.TextEdit, text string) {
		if selection < 0 {
			return
		}
		oldText, _ := listbox.Get(selection)
		if oldText == text {
			return
		}
		listbox.Set(selection, text)
	}

	deleteButton = wnd.NewButton()
	deleteButton.SetText("Удалить")
	deleteButton.OnClick = func(b *wingo.Button) {
		if selection < 0 {
			return
		}
		listbox.Delete(selection)
		selection = -1
		teUpdate.SetText("")
	}

	teAppend = wnd.NewTextEdit()

	appendButton = wnd.NewButton()
	appendButton.SetText("Добавить")
	appendButton.OnClick = func(b *wingo.Button) {
		listbox.Append(teAppend.GetText())
		teAppend.SetText("")
	}

	menu := wingo.NewMenu()
	pathMenu := menu.AppendPopup("Переменная среды")
	pathMenu.AppendItemText("Загрузить").OnClick = func(mi *wingo.MenuItem) {
		loadPaths()
	}
	pathMenu.AppendItemText("Сохранить").OnClick = func(mi *wingo.MenuItem) {
		joined := strings.Join(listbox.GetList(), ";")
		w32.RegSetExpandString(regPathKey, "Path", joined)
	}

	wnd.ApplyMenu(menu)

	updateMetrics()
	loadPaths()

	wingo.Start()

	w32.RegCloseKey(regPathKey)

}
