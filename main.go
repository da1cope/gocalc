package main

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type calc struct {
	display *widget.Label
	current string
	op      string
	prev    float64
}

func (c *calc) input(num string) {
	c.current += num
	c.updateDisplay()
}

func (c *calc) operator(op string) {
	if c.current != "" {
		c.prev, _ = strconv.ParseFloat(c.current, 64)
		c.current = ""
	}
	c.op = op
	c.updateDisplay()
}

func (c *calc) clear() {
	c.current = ""
	c.op = ""
	c.prev = 0
	c.updateDisplay()
}

func (c *calc) equals() {
	if c.current == "" || c.op == "" {
		return
	}
	currentNum, _ := strconv.ParseFloat(c.current, 64)
	var result float64

	switch c.op {
	case "+":
		result = c.prev + currentNum
	case "−":
		result = c.prev - currentNum
	case "×":
		result = c.prev * currentNum
	case "÷":
		if currentNum != 0 {
			result = c.prev / currentNum
		} else {
			c.display.SetText("Error")
			c.clear()
			return
		}
	default:
		return
	}

	c.current = strconv.FormatFloat(result, 'f', -1, 64)
	c.updateDisplay()
	c.op = ""
	c.prev = 0
}

func (c *calc) toggleSign() {
	if c.current == "" {
		return
	}
	num, _ := strconv.ParseFloat(c.current, 64)
	num = -num
	c.current = strconv.FormatFloat(num, 'f', -1, 64)
	c.updateDisplay()
}

func (c *calc) percentage() {
	if c.current == "" {
		return
	}
	num, _ := strconv.ParseFloat(c.current, 64)
	num /= 100
	c.current = strconv.FormatFloat(num, 'f', -1, 64)
	c.updateDisplay()
}

func (c *calc) updateDisplay() {
	if c.current == "" && c.op != "" {
		c.display.SetText(strconv.FormatFloat(c.prev, 'f', -1, 64) + " " + c.op)
	} else if c.current != "" {
		c.display.SetText(c.current)
	} else {
		c.display.SetText("0")
	}
}

func hasDecimal(s string) bool {
	parts := strings.Split(s, ".")
	return len(parts) > 1
}

func main() {
	a := app.New()
	w := a.NewWindow("Simple Calculator")
	w.Resize(fyne.NewSize(320, 460))

	calc := &calc{
		display: widget.NewLabel("0"),
	}

	calc.display.TextStyle = fyne.TextStyle{Bold: true}
	display := widget.NewLabel("0")
	display.Alignment = fyne.TextAlignTrailing
	display.TextStyle = fyne.TextStyle{Bold: true}
	calc.display = display

	// Helper to create number buttons
	numBtn := func(label string, action func()) *widget.Button {
		b := widget.NewButton(label, action)
		b.Importance = widget.MediumImportance
		return b
	}

	// Helper to create operator buttons (high importance for styling)
	opBtn := func(label string, action func()) *widget.Button {
		b := widget.NewButton(label, action)
		b.Importance = widget.HighImportance
		return b
	}

	row1 := container.NewGridWithColumns(4,
		widget.NewButton("C", calc.clear),
		widget.NewButton("±", calc.toggleSign),
		widget.NewButton("%", calc.percentage),
		opBtn("÷", func() { calc.operator("÷") }),
	)
	row2 := container.NewGridWithColumns(4,
		numBtn("7", func() { calc.input("7") }),
		numBtn("8", func() { calc.input("8") }),
		numBtn("9", func() { calc.input("9") }),
		opBtn("×", func() { calc.operator("×") }),
	)
	row3 := container.NewGridWithColumns(4,
		numBtn("4", func() { calc.input("4") }),
		numBtn("5", func() { calc.input("5") }),
		numBtn("6", func() { calc.input("6") }),
		opBtn("−", func() { calc.operator("−") }),
	)
	row4 := container.NewGridWithColumns(4,
		numBtn("1", func() { calc.input("1") }),
		numBtn("2", func() { calc.input("2") }),
		numBtn("3", func() { calc.input("3") }),
		opBtn("+", func() { calc.operator("+") }),
	)
	row5 := container.NewGridWithColumns(4,
		container.NewGridWithColumns(2, numBtn("0", func() { calc.input("0") }), widget.NewButton(".", func() {
			if !hasDecimal(calc.current) {
				if calc.current == "" {
					calc.current = "0"
				}
				calc.input(".")
			}
		})),
		widget.NewButton("=", calc.equals),
	)

	content := container.NewVBox(
		container.NewPadded(display),
		container.NewPadded(row1),
		container.NewPadded(row2),
		container.NewPadded(row3),
		container.NewPadded(row4),
		container.NewPadded(row5),
	)

	w.SetContent(content)
	w.ShowAndRun()
}
