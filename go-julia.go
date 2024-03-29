package main

import (
	"log"
	"math"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func init_tcell() tcell.Screen {
    s, err := tcell.NewScreen()
    if err != nil {
        log.Fatalf("%+v", err)
    }
    err = s.Init()
    if err != nil {
        log.Fatalf("%+v", err)
    }
    return s
}

func quit(s tcell.Screen) {
    s.Fini()
    os.Exit(0)
}

func get_color_lerp(idx int) (int32, int32, int32) {
    br, bg, bb := BEGIN_COLOR[0], BEGIN_COLOR[1], BEGIN_COLOR[2]
    er, eg, eb := END_COLOR[0], END_COLOR[1], END_COLOR[2]
    return  int32(float32(br) + FACTOR*float32(idx)*float32(er-br)),
            int32(float32(bg) + FACTOR*float32(idx)*float32(eg-bg)),
            int32(float32(bb) + FACTOR*float32(idx)*float32(eb-bb))
}

func make_styles() []tcell.Style {
    styles := make([]tcell.Style, MAX_IT)

    for i := range styles {
        styles[i] = tcell.StyleDefault.Background(tcell.NewRGBColor(get_color_lerp(i))).Foreground(tcell.ColorReset)
    }

    return styles
}

func magnitude(re, im float32) float32 {
    return float32(math.Sqrt(float64(re*re + im*im)))
}

func draw_frame(s tcell.Screen, styles []tcell.Style, zoom, pos_x, pos_y float32) {
    sz_x, sz_y := s.Size()
    esc := RADIUS/zoom

    for i := 0; i < sz_x; i++ {
        for j := 0; j < sz_y; j++ {
            re := ((float32(i)/zoom) - pos_x)/(float32(sz_x-1)/zoom)*(2.0*esc)-esc
            im := ((float32(j)/zoom) - pos_y)/(float32(sz_y-1)/zoom)*(2.0*esc)-esc
            
            it := 0
            for magnitude(re, im) < RADIUS*RADIUS && it < MAX_IT {
                re_tmp := re*re - im*im
                im = 2*re*im + JULIA_CONST_IM
                re = re_tmp + JULIA_CONST_RE
                it++
            }
            s.SetContent(i, j, ' ', nil, styles[it])
        }
    }
}

const MAX_IT int = 100
const FACTOR float32 = 1/(float32(MAX_IT)-1)
var BEGIN_COLOR = [3]int32{0, 255, 0}
var END_COLOR = [3]int32{0, 0, 0}
var JULIA_CONST_RE float32 = -0.8
var JULIA_CONST_IM float32 = 0.156
//var RADIUS float32 = (1.0+float32(math.Sqrt(1.0-4.0*math.Sqrt(float64(JULIA_CONST_RE*JULIA_CONST_RE+JULIA_CONST_IM*JULIA_CONST_IM)))))/2.0
const RADIUS float32 = 2.0
func main() {
    s := init_tcell()
    default_style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
    s.SetStyle(default_style)
    s.Clear()

    args := os.Args[1:]
    if len(args) == 2 {
        re, _ := strconv.ParseFloat(args[0], 32)
        im, _ := strconv.ParseFloat(args[1], 32)
        JULIA_CONST_RE = float32(re)
        JULIA_CONST_IM = float32(im)
    }

    styles := append(make_styles(), tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorReset))
    var zoom float32 = 1.0
    const d_zoom float32 = 0.5
    var pos_x float32 = 0.0
    var pos_y float32 = 0.0
    var d_pos float32 = 1.0
    draw_frame(s, styles, zoom, pos_x, pos_y)

    for {
        s.Show()

        ev := s.PollEvent()
        switch ev := ev.(type) {
        case *tcell.EventResize:
            s.Sync()
            s.Show()
        case *tcell.EventKey:
            if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
                quit(s)
            }
            if ev.Key() == tcell.KeyRight {
                pos_x -= d_pos
                draw_frame(s, styles, zoom, pos_x, pos_y)
            }
            if ev.Key() == tcell.KeyLeft {
                pos_x += d_pos
                draw_frame(s, styles, zoom, pos_x, pos_y)
            }
            if ev.Key() == tcell.KeyUp {
                pos_y += d_pos
                draw_frame(s, styles, zoom, pos_x, pos_y)
            }
            if ev.Key() == tcell.KeyDown {
                pos_y -= d_pos
                draw_frame(s, styles, zoom, pos_x, pos_y)
            }
            if ev.Key() == tcell.KeyRune {
                if ev.Rune() == '+' || ev.Rune() == '=' {
                    zoom += d_zoom
                    d_pos = 1.0/zoom
                    draw_frame(s, styles, zoom, pos_x, pos_y)
                }
                if (ev.Rune() == '-' || ev.Rune() == '_') && zoom > 0.0 {
                    zoom -= d_zoom
                    d_pos = 1.0/zoom
                    draw_frame(s, styles, zoom, pos_x, pos_y)
                }
                if ev.Rune() == ' ' {
                    pos_x = 0.0
                    pos_y = 0.0
                    zoom = 1.0
                    draw_frame(s, styles, zoom, pos_x, pos_y)
                }
            }
        }
    }
}
