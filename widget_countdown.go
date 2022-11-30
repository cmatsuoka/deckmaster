package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/muesli/deckmaster/countdown"
)

var (
	countingColor = color.RGBA{64, 255, 64, 255}
	expiredColor  = color.RGBA{255, 64, 64, 255}
)

type CountdownWidgetState struct {
	timer *countdown.Countdown
}

func NewCountdownWidgetState(seconds, rate int) CountdownWidgetState {
	timer := countdown.New(seconds, rate)
	return CountdownWidgetState{
		timer: timer,
	}
}

// CountdownWidget is a widget displaying the current time/date.
type CountdownWidget struct {
	*StatefulWidget

	font   string
	start  time.Time
	layout *Layout
	state  CountdownWidgetState
}

// NewCountdownWidget returns a new CountdownWidget.
func NewCountdownWidget(sw *StatefulWidget, opts WidgetConfig) *CountdownWidget {
	sw.setInterval(time.Duration(opts.Interval)*time.Millisecond, time.Second/4)

	var state_id string
	var seconds int64

	_ = ConfigValue(opts.Config["seconds"], &seconds)
	_ = ConfigValue(opts.Config["state"], &state_id)

	// Set defaults
	if state_id == "" {
		state_id = "countdown"
	}
	if seconds == 0 {
		seconds = 300
	}

	verbosef("New countdown widget (%ds), state id = %q", seconds, state_id)

	sw.state_id = state_id
	layout := NewLayout(int(sw.dev.Pixels))

	var cstate CountdownWidgetState
	state, err := stateRegistry.Recover(state_id)
	if err != nil {
		cstate = NewCountdownWidgetState(int(seconds), 4)
		stateRegistry.Register(state_id, cstate)
		verbosef("Create new state")
	} else {
		cstate = state.(CountdownWidgetState)
		verbosef("Recover existing state")
	}

	return &CountdownWidget{
		StatefulWidget: sw,
		start:          time.Now(),
		layout:         layout,
		state:          cstate,
	}
}

// Update renders the widget.
func (w *CountdownWidget) Update() error {
	size := int(w.dev.Pixels)
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	timer := w.state.timer
	remaining := int(timer.Remaining().Seconds())

	var str string
	if timer.Expired() {
		str = "0:00"
	} else {
		str = fmt.Sprintf("%d:%02d", remaining/60, remaining%60)
	}

	bounds := img.Bounds()

	var color color.Color
	switch {
	case timer.Expired():
		color = expiredColor
	case timer.Paused():
		color = DefaultColor
	default:
		if remaining < 10 {
			color = expiredColor
		} else {
			color = countingColor
		}
	}

	drawString(img,
		bounds,
		ttfFont,
		str,
		w.dev.DPI,
		-1,
		color,
		image.Pt(-1, -1),
	)

	return w.render(w.dev, img)
}

// TriggerAction gets called when a button is pressed.
func (w *CountdownWidget) TriggerAction(hold bool) {
	timer := w.state.timer
	switch {
	case hold:
		verbosef("Reset counter %q", w.state_id)
		timer.Reset()
	case timer.Paused():
		if !timer.Expired() {
			verbosef("Start counter %q", w.state_id)
			timer.Resume()
		}
	default:
		if !timer.Expired() {
			verbosef("Pause counter %q", w.state_id)
			timer.Pause()
		}
	}
	w.Update()
}
