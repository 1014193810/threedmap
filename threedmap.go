// threedmap project threedmap.go
package threedmap

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	//	"github.com/golang/image/colornames"
)

type Smap struct {
	Nodes [][]node

	Size Ssize
}

type node struct {
	Color   color.Color
	X, Y, Z float64
}

type Ssize struct {
	X0, Y0, XE, YE int
}

func (s Ssize) Width() int {
	return s.XE - s.X0
}
func (s Ssize) Center() (float64, float64) {
	return float64(s.XE+s.X0) / 2, float64(s.YE+s.Y0) / 2
}
func (s Ssize) Highth() int {
	return s.YE - s.Y0
}
func (M *Smap) Get(X, Y int) node {

	return M.Nodes[Y][X]
}
func (M *Smap) SetZ(X, Y int, z float64) {

	//	n := M.Size.Width()*Y + X

	//	M.Nodes[n].Z = z

	/////
	M.Nodes[Y][X].Z = z

}
func (M *Smap) SetColor(X, Y int, c color.Color) {
	//	n := M.Size.Width()*Y + X
	//	M.Nodes[n].Color = c

	/////
	M.Nodes[Y][X].Color = c

}
func (M *Smap) SetLight(x, y, z float64, c color.Color) {
	r := math.Sqrt(x*x + y*y + z*z)
	x = x / r
	y = y / r
	z = z / r

}
func NewSmap(s Ssize) *Smap {
	m := new(Smap)
	m.Size = s
	//	n := (s.Width() * s.Highth())
	m.Nodes = make([][]node, s.Highth())
	for i := 0; i < s.Highth(); i++ {
		m.Nodes[i] = make([]node, s.Width())
	}

	for i := 0; i < s.Width(); i++ {
		for j := 0; j < s.Highth(); j++ {
			//num := s.Width()*j + i
			x, y := s.Center()
			m.Nodes[j][i].X = float64(i) - x
			m.Nodes[j][i].Y = float64(j) - y

		}
	}
	return m
}

var Img *image.NRGBA

var Flags []float64

func Initmap(s Ssize) *Smap {

	Img = image.NewNRGBA(image.Rect(-s.Width()/2, -s.Highth()/2, s.Width()/2-1, s.Highth()/2-1))

	M := NewSmap(s)
	Flags = make([]float64, s.Highth()*s.Width())

	return M

}

var (
	sins, coss, sinf, cosf, k1, k2, k3, k4, k5, k6, k7, k8, k9, X, Y, Z float64
)

func (m *Smap) Drow2(s, f float64, r, d float64, dx, dy int) *image.NRGBA {
	sins = math.Sin(s)
	coss = math.Cos(s)
	sinf = math.Sin(f)
	cosf = math.Cos(f)
	k1 = -sins
	k2 = -cosf * coss
	k3 = -sinf * coss / d
	k4 = coss
	k5 = -cosf * sins
	k6 = -sinf * sins / d
	k7 = sinf
	k8 = -cosf / d
	k9 = r / d
	X = r * sinf * coss
	Y = r * sinf * sins
	Z = r * cosf

	//Clearimg()
	Img = image.NewNRGBA(Img.Rect)
	Flags = make([]float64, len(Flags))

	for i := 0; i < m.Size.Width(); i += dx {
		for j := 0; j < m.Size.Highth(); j += dy {
			node := m.Nodes[j][i]
			f := node.X*k3 + node.Y*k6 + node.Z*k8 + k9

			x := (node.X*k1 + node.Y*k4) / f
			y := -(node.X*k2 + node.Y*k5 + node.Z*k7) / f
			dis := (X-node.X)*(X-node.X) + (Y-node.Y)*(Y-node.Y) + (Z-node.Z)*(Z-node.Z)
			X0 := int(x+1000) - 1000
			Y0 := int(y+1000) - 1000

			n := (X0 + m.Size.Width()/2) + (Y0+m.Size.Highth()/2)*m.Size.Width()

			if n < 0 || n >= m.Size.Highth()*m.Size.Width() {
				continue
			}
			if Flags[n] == 0 || Flags[n] > dis {
				Flags[n] = dis
				Img.Set(X0, Y0, node.Color)

			}
		}
	}
	return Img

}
func GenerateRunfunc(m *Smap, cfg pixelgl.WindowConfig, background color.Color) func() {
	return func() {
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}

		var (
			camPos       = pixel.ZV
			camSpeed     = 500.0
			camZoom      = 1.0
			camZoomSpeed = 1.2
		)

		var (
			frames = 0
			second = time.Tick(time.Second)
		)
		s, f := 0.0, 0.0

		last := time.Now()
		for !win.Closed() {
			dt := time.Since(last).Seconds()
			last = time.Now()

			cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
			win.SetMatrix(cam)

			img := m.Drow2(s, f, 10000, 2000, 2, 2)
			pic := pixel.PictureDataFromImage(img)
			sprite := pixel.NewSprite(pic, pic.Bounds())

			win.Clear(background)

			sprite.Draw(win, pixel.IM)
			if win.Pressed(pixelgl.KeyA) {
				if s > 0 {
					s -= 1 * dt
				}
			}
			if win.Pressed(pixelgl.KeyD) {
				if s < math.Pi*2 {
					s += 1 * dt
				}
			}
			if win.Pressed(pixelgl.KeyW) {
				if f > 0 {
					f -= 1 * dt
				}
			}
			if win.Pressed(pixelgl.KeyS) {
				if f < math.Pi {
					f += 1 * dt
				}
			}

			if win.Pressed(pixelgl.KeyLeft) {
				camPos.X -= camSpeed * dt
			}
			if win.Pressed(pixelgl.KeyRight) {
				camPos.X += camSpeed * dt
			}
			if win.Pressed(pixelgl.KeyDown) {
				camPos.Y -= camSpeed * dt
			}
			if win.Pressed(pixelgl.KeyUp) {
				camPos.Y += camSpeed * dt
			}
			camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

			win.Update()

			frames++
			select {
			case <-second:
				win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
				frames = 0
			default:
			}
		}
	}

}
