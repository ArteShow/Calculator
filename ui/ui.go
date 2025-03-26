package ui

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func SendExpression(expression string, resultArea *widgets.QTextEdit) {
	// Example body for the POST request
	body := `{"expression": "` + expression + `"}`

	// Send the expression to the internal server
	resp, err := http.Post("http://localhost:8082/api/v1/calculate", "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		widgets.QMessageBox_Critical(nil, "Error", "Failed to send expression", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	resultArea.SetPlainText(string(responseBody))

	widgets.QMessageBox_Information(nil, "Success", "Expression sent successfully", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func GetExpressions(resultArea *widgets.QTextEdit) {
	// Example GET request to retrieve expressions
	resp, err := http.Get("http://localhost:8082/api/v1/expressions")
	if err != nil {
		widgets.QMessageBox_Critical(nil, "Error", "Failed to get expressions", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	resultArea.SetPlainText(string(responseBody))

	widgets.QMessageBox_Information(nil, "Success", "Expressions retrieved successfully", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func GetExpressionByID(id string, resultArea *widgets.QTextEdit) {
	// Example GET request to retrieve an expression by ID
	resp, err := http.Get("http://localhost:8082/api/v1/expression/" + id)
	if err != nil {
		widgets.QMessageBox_Critical(nil, "Error", "Failed to get expression by ID", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := ioutil.ReadAll(resp.Body)
	resultArea.SetPlainText(string(responseBody))

	widgets.QMessageBox_Information(nil, "Success", "Expression retrieved successfully", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
}

func RunUI() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Calculator")
	window.SetMinimumSize2(300, 400)

	centralWidget := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQVBoxLayout()

	title := widgets.NewQLabel2("Simple Calculator", nil, 0)
	title.SetAlignment(core.Qt__AlignCenter)
	layout.AddWidget(title, 0, 0)

	selectBox := widgets.NewQComboBox(nil)
	selectBox.AddItems([]string{"Calculate", "Expressions", "Get Expression by ID"})
	layout.AddWidget(selectBox, 0, 0)

	label := widgets.NewQLabel2("Enter Expression:", nil, 0)
	label.SetAlignment(core.Qt__AlignCenter)
	layout.AddWidget(label, 0, 0)

	input := widgets.NewQLineEdit(nil)
	layout.AddWidget(input, 0, 0)

	sendButton := widgets.NewQPushButton2("Send", nil)
	sendButton.SetFixedSize2(100, 30) // Make the button smaller
	layout.AddWidget(sendButton, 0, core.Qt__AlignCenter)

	getExpressionsButton := widgets.NewQPushButton2("Get Expressions", nil)
	getExpressionsButton.SetFixedSize2(150, 30)
	layout.AddWidget(getExpressionsButton, 0, core.Qt__AlignCenter)

	getExpressionByIDButton := widgets.NewQPushButton2("Get Expression by ID", nil)
	getExpressionByIDButton.SetFixedSize2(200, 30)
	layout.AddWidget(getExpressionByIDButton, 0, core.Qt__AlignCenter)

	resultArea := widgets.NewQTextEdit(nil)
	resultArea.SetReadOnly(true)
	layout.AddWidget(resultArea, 0, 0)

	// Connect button actions
	sendButton.ConnectClicked(func(bool) {
		expression := input.Text()
		SendExpression(expression, resultArea)
	})

	getExpressionsButton.ConnectClicked(func(bool) {
		GetExpressions(resultArea)
	})

	getExpressionByIDButton.ConnectClicked(func(bool) {
		id := input.Text()
		GetExpressionByID(id, resultArea)
	})

	// Show/hide widgets based on selected value
	selectBox.ConnectCurrentIndexChanged(func(index int) {
		switch index {
		case 0: // Calculate
			label.SetText("Enter Expression:")
			label.Show()
			input.Show()
			sendButton.Show()
			getExpressionsButton.Hide()
			getExpressionByIDButton.Hide()
		case 1: // Expressions
			label.Hide()
			input.Hide()
			sendButton.Hide()
			getExpressionsButton.Show()
			getExpressionByIDButton.Hide()
		case 2: // Get Expression by ID
			label.SetText("Enter ID:")
			label.Show()
			input.Show()
			sendButton.Hide()
			getExpressionsButton.Hide()
			getExpressionByIDButton.Show()
		}
	})

	// Initialize visibility
	selectBox.SetCurrentIndex(0)

	centralWidget.SetLayout(layout)
	window.SetCentralWidget(centralWidget)

	window.Show()

	app.Exec()
}

func main() {
	RunUI()
}