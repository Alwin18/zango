package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Alwin18/zango/internal/service"
)

func Start() {
	a := app.New()
	w := a.NewWindow("Zango Control Panel")

	// Widget status (pakai canvas.Text agar bisa ubah warna)
	httpStatus := newStatusLabel()
	dbStatus := newStatusLabel()

	// Update semua status
	updateStatuses := func() {
		httpActive, _ := service.Status("http")
		dbActive, _ := service.Status("db")

		httpStatus.set(httpActive)
		dbStatus.set(dbActive)
	}

	// Tombol HTTP
	startHttp := widget.NewButton("Start", func() {
		_ = service.Start("http")
		updateStatuses()
	})
	stopHttp := widget.NewButton("Stop", func() {
		_ = service.Stop("http")
		updateStatuses()
	})

	// Tombol DB
	startDb := widget.NewButton("Start", func() {
		_ = service.Start("db")
		updateStatuses()
	})
	stopDb := widget.NewButton("Stop", func() {
		_ = service.Stop("db")
		updateStatuses()
	})

	// Tabel tampilan
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

	content := container.NewVBox(
		widget.NewLabelWithStyle("Zango Control Panel", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		httpRow,
		dbRow,
	)

	w.SetContent(content)
	updateStatuses()
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
