package main

import (
	"fmt"
	"log"

	"github.com/gogpu/gg"
	_ "github.com/gogpu/gg/gpu"
	"github.com/gogpu/gg/integration/ggcanvas"
	"github.com/gogpu/gogpu"
	"github.com/gogpu/ui/app"
	"github.com/gogpu/ui/core/button"
	"github.com/gogpu/ui/primitives"
	"github.com/gogpu/ui/render"
	"github.com/gogpu/ui/theme/material3"
	"github.com/gogpu/ui/widget"
)

func main() {
	gogpuApp := gogpu.NewApp(gogpu.DefaultConfig().
		WithTitle("My First App").
		WithSize(600, 400).
		WithContinuousRender(false))

	m3 := material3.New(widget.Hex(0x6750A4))

	uiApp := app.New(
		app.WithWindowProvider(gogpuApp),
		app.WithPlatformProvider(gogpuApp),
		app.WithEventSource(gogpuApp.EventSource()),
	)
	uiApp.SetRoot(
		primitives.Box(
			primitives.Text("Hello, gogpu/ui!").FontSize(24).Bold(),
			button.New(
				button.TextOpt("Click Me"),
				button.OnClick(func() { fmt.Println("Clicked!") }),
				button.PainterOpt(material3.ButtonPainter{Theme: m3}),
			),
		).Padding(24).Gap(12),
	)

	var canvas *ggcanvas.Canvas
	gogpuApp.OnDraw(func(dc *gogpu.Context) {
		w, h := dc.Width(), dc.Height()
		if w <= 0 || h <= 0 {
			return
		}
		if canvas == nil {
			provider := gogpuApp.GPUContextProvider()
			if provider == nil {
				return
			}
			var err error
			canvas, err = ggcanvas.New(provider, w, h)
			if err != nil {
				log.Printf("ggcanvas: %v", err)
				return
			}
		}
		uiApp.Frame()
		cw, ch := canvas.Size()
		if cw != w || ch != h {
			if err := canvas.Resize(w, h); err != nil {
				log.Printf("resize: %v", err)
			}
			cw, ch = w, h
		}
		// sv := dc.SurfaceView()
		// sw, sh := dc.SurfaceSize()
		canvas.Draw(func(cc *gg.Context) {
			cc.SetRGBA(0.94, 0.94, 0.94, 1)
			cc.DrawRectangle(0, 0, float64(cw), float64(ch))
			cc.Fill()
			widgetCanvas := render.NewCanvas(cc, cw, ch)
			uiApp.Window().DrawTo(widgetCanvas)
		})
		// if err := canvas.RenderDirect(sv, sw, sh); err != nil {
		// 	log.Printf("render: %v", err)
		// }
		if err := canvas.Render(dc.RenderTarget()); err != nil {
			log.Printf("render: %v", err)
		}
	})
	gogpuApp.OnClose(func() { gg.CloseAccelerator() })

	if err := gogpuApp.Run(); err != nil {
		log.Fatal(err)
	}
}
