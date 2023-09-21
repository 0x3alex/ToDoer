package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"io"
	"os"
)

type Tasks struct {
	Tasks []string
}

var data Tasks

var currentSelectedTask = -1

func remove(slice []string, s int) []string {
	if s == -1 {
		return slice
	}
	return append(slice[:s], slice[s+1:]...)
}

func handleFileNotExists() {
	if _, err := os.Stat("todos.json"); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create("todos.json")
		if err != nil {
			panic(err.Error())
		}
		return
	}
}

func saveToDos() {
	handleFileNotExists()
	file, _ := json.MarshalIndent(data, "", " ")
	_ = os.WriteFile("todos.json", file, 0644)
}

func loadToDos() {
	handleFileNotExists()
	jsonFile, err := os.Open("todos.json")
	if err != nil {
		fmt.Println(err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return
	}
}

func newTaskWindow(ToDoer fyne.App, List fyne.Widget) {
	newTaskWindow := ToDoer.NewWindow("Create new ToDo")
	newTaskWindow.Resize(fyne.Size{Width: 500, Height: 100})
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter ToDo")
	content := container.NewVBox(input, widget.NewButton("Save", func() {
		data.Tasks = append(data.Tasks, input.Text)
		List.Refresh()
		saveToDos()
		newTaskWindow.Close()
		input.Text = ""
	}))
	newTaskWindow.SetContent(content)
	newTaskWindow.Show()
}

func main() {
	ToDoer := app.New()

	//ToDoer.Settings().SetTheme(&builtinTheme{variant: VariantDark})
	mainWindow := ToDoer.NewWindow("ToDoer")
	mainWindow.Resize(fyne.Size{Width: 500, Height: 500})
	mainWindow.SetFixedSize(true)

	loadToDos()
	list := widget.NewList(
		func() int {
			return len(data.Tasks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data.Tasks[i])
		})

	list.OnSelected = func(id widget.ListItemID) {
		currentSelectedTask = id
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			newTaskWindow(ToDoer, list)
			//ToDoer.SendNotification(fyne.NewNotification("Goofy", "ahh"))
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			data.Tasks = remove(data.Tasks, currentSelectedTask)
			list.Refresh()
			saveToDos()
		}),
		widget.NewToolbarAction(theme.ConfirmIcon(), func() {
			data.Tasks = remove(data.Tasks, currentSelectedTask)
			list.Refresh()
			saveToDos()
		}),
	)
	content := container.NewBorder(toolbar, nil, nil, nil, list)
	mainWindow.SetContent(content)
	mainWindow.ShowAndRun()
	ToDoer.Quit()
}
