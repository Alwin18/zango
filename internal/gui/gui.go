package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Alwin18/zango/internal/service"
	"github.com/Alwin18/zango/internal/storage"
)

func Start() {
	a := app.New()
	w := a.NewWindow("Zango Control Panel")

	tabs := container.NewAppTabs(
		container.NewTabItem("Control Panel", buildControlPanel()),
		container.NewTabItem("Logs", buildLogViewer()),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1200, 720))
	w.ShowAndRun()
}

type statusLabel struct {
	text *canvas.Text
}

func newStatusLabel() *statusLabel {
	txt := canvas.NewText("-", color.Black)
	txt.TextSize = 14
	return &statusLabel{text: txt}
}

func (s *statusLabel) set(active bool) {
	if active {
		s.text.Text = "🟢 Active"
		s.text.Color = color.RGBA{0, 200, 0, 255} // hijau
	} else {
		s.text.Text = "🔴 Inactive"
		s.text.Color = color.RGBA{200, 0, 0, 255} // merah
	}
	s.text.Refresh()
}

func buildControlPanel() fyne.CanvasObject {
	httpStatus := newStatusLabel()
	dbStatus := newStatusLabel()

	updateStatuses := func() {
		httpActive, _ := service.Status("http")
		dbActive, _ := service.Status("db")

		httpStatus.set(httpActive)
		dbStatus.set(dbActive)
	}

	startHttp := widget.NewButton("Start", func() {
		_ = service.Start("http")
		updateStatuses()
	})
	stopHttp := widget.NewButton("Stop", func() {
		_ = service.Stop("http")
		updateStatuses()
	})

	startDb := widget.NewButton("Start", func() {
		_ = service.Start("db")
		updateStatuses()
	})
	stopDb := widget.NewButton("Stop", func() {
		_ = service.Stop("db")
		updateStatuses()
	})

	httpRow := container.NewHBox(
		widget.NewLabel("HTTP"),
		httpStatus.text,
		startHttp,
		stopHttp,
	)

	dbRow := container.NewHBox(
		widget.NewLabel("Database"),
		dbStatus.text,
		startDb,
		stopDb,
	)

	updateStatuses()

	return container.NewVBox(
		widget.NewLabelWithStyle("GO-XAMPP Control Panel", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		httpRow,
		dbRow,
	)
}

func buildLogViewer() fyne.CanvasObject {
	logBox := container.NewVBox()

	// Header tetap
	header := container.NewHBox(
		makeBoldLabel("Service", 100),
		makeBoldLabel("Action", 80),
		makeBoldLabel("Timestamp", 220),
	)
	logBox.Add(header)

	scroll := container.NewScroll(logBox)
	scroll.SetMinSize(fyne.NewSize(500, 700))

	// Tombol Refresh
	refreshBtn := widget.NewButton("Refresh Logs", func() {
		refreshLogBox(logBox)
	})

	refreshBtn.Resize(fyne.NewSize(220, 40))

	refreshContainer := container.NewHBox(
		refreshBtn,
		layout.NewSpacer(), // agar tetap di kiri
	)

	// Pertama kali isi logBox
	refreshLogBox(logBox)

	return container.NewVBox(
		widget.NewLabelWithStyle("Service Logs", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		refreshContainer,
		scroll,
	)
}

func refreshLogBox(logBox *fyne.Container) {
	logs, err := storage.GetLatestLogs(100)
	if err != nil {
		return
	}

	logBox.Objects = logBox.Objects[:1] // pertahankan header

	for _, l := range logs {
		row := container.NewHBox(
			makeMonoLabel(l.ServiceName, 100),
			makeActionLabel(l.Action, 80),
			makeMonoLabel(l.Timestamp.Format("2006-01-02 15:04:05"), 220),
		)
		logBox.Add(row)
	}

	logBox.Refresh()

}

func makeMonoLabel(text string, width float32) fyne.CanvasObject {
	lbl := canvas.NewText(text, color.White)
	lbl.TextSize = 14
	lbl.TextStyle = fyne.TextStyle{Monospace: true}
	lbl.Alignment = fyne.TextAlignLeading

	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(width, 20))

	return container.NewStack(rect, lbl)
}

func makeBoldLabel(text string, width float32) fyne.CanvasObject {
	lbl := canvas.NewText(text, color.White)
	lbl.TextSize = 18
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	lbl.Alignment = fyne.TextAlignLeading

	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(width, 20))

	return container.NewStack(rect, lbl)
}

func makeActionLabel(action string, width float32) fyne.CanvasObject {
	var c color.Color
	switch action {
	case "start":
		c = color.RGBA{0, 180, 0, 255} // Hijau
	case "stop":
		c = color.RGBA{200, 0, 0, 255} // Merah
	default:
		c = color.Black
	}

	lbl := canvas.NewText(strings.ToUpper(action), c)
	lbl.TextSize = 14
	lbl.TextStyle = fyne.TextStyle{Monospace: true}
	lbl.Alignment = fyne.TextAlignLeading

	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(width, 20))
	return container.NewStack(rect, lbl)
}
