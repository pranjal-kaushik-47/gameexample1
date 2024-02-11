// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameWidth  = 32
	frameHeight = 32
)

var (
	frameOX    = 0
	frameOY    = 0
	frameCount = 8
	frameRate  = 10
	edit       = false
	facing     = "right"
	flip       = false
	xdirection = 1
	ydirection = 1
	xmovement  = 0
	ymovement  = 0
)

var (
	runnerImage  *ebiten.Image
	runnerImage2 *ebiten.Image
)

type animationConf struct {
	frameOX    int
	frameOY    int
	frameCount int
	frameRate  int
}

var idle = animationConf{
	frameOX:    0,
	frameOY:    0,
	frameCount: 8,
	frameRate:  10,
}
var running = animationConf{
	frameOX:    0,
	frameOY:    32,
	frameCount: 8,
	frameRate:  10,
}
var attack1 = animationConf{
	frameOX:    0,
	frameOY:    64,
	frameCount: 8,
	frameRate:  10,
}
var attack2 = animationConf{
	frameOX:    0,
	frameOY:    96,
	frameCount: 8,
	frameRate:  10,
}
var attack3 = animationConf{
	frameOX:    0,
	frameOY:    128,
	frameCount: 8,
	frameRate:  10,
}

var jump = animationConf{
	frameOX:    0,
	frameOY:    160,
	frameCount: 6,
	frameRate:  15,
}

var die = animationConf{
	frameOX:    0,
	frameOY:    192,
	frameCount: 7,
	frameRate:  20,
}

var animationMap = map[ebiten.Key]animationConf{
	ebiten.KeySpace: jump,
	ebiten.KeyF:     attack1,
	ebiten.KeyD:     running,
	ebiten.KeyA:     running,
	ebiten.KeyW:     running,
	ebiten.KeyS:     running,
	ebiten.KeyE:     idle,
	ebiten.KeyX:     die,
}

type Game struct {
	count int
	keys  []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	for _, p := range g.keys {
		if inpututil.KeyPressDuration(p) > 5 {
			conf := animationMap[p]
			defultConf := animationConf{}
			if conf != defultConf {
				frameOX, frameOY, frameCount, frameRate = conf.frameOX, conf.frameOY, conf.frameCount, conf.frameRate
			}

			switch p {
			case ebiten.KeyD:
				xdirection = 1
				xmovement = inpututil.KeyPressDuration(p)
				if facing == "right" {
					flip = false
				} else {
					flip = true
				}
			case ebiten.KeyA:
				xdirection = -1
				xmovement = inpututil.KeyPressDuration(p)
				if facing == "right" {
					flip = true
				} else {
					flip = false
				}
			case ebiten.KeyW:
				ydirection = -1
				ymovement = inpututil.KeyPressDuration(p)
			case ebiten.KeyS:
				ydirection = 1
				ymovement = inpututil.KeyPressDuration(p)
			case ebiten.KeyUp:
				frameOY--
			case ebiten.KeyDown:
				frameOY++
			case ebiten.KeyEnter:
				conf.frameOY = frameOY
			case ebiten.KeyZ:
				edit = true
			}
		}
	}

	if flip {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate((screenWidth/2+30)+float64(xmovement*xdirection), (screenHeight/2)+float64(ymovement*ydirection))
	} else {
		op.GeoM.Translate((screenWidth/2)+float64(xmovement*xdirection), (screenHeight/2)+float64(ymovement*ydirection))
	}

	i := (g.count / frameRate) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v : %v : %v : %v : %v, %v", frameOX, frameOY, frameCount, frameRate, xdirection, ydirection), 10, 20)
	if !edit {
		if i == 0 {
			frameOX, frameOY, frameCount, frameRate = idle.frameOX, idle.frameOY, idle.frameCount, idle.frameRate
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func getSpritSheet() *os.File {
	url := "https://kidscancode.org/godot_recipes/3.x/img/adventurer_sprite_sheet_v1.1.png"
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("asdf.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func main() {
	imageData, err := ioutil.ReadFile("asdf.png")
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return
	}
	imageBuffer := bytes.NewReader(imageData)

	decodedImage, _, err := image.Decode(imageBuffer)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	runnerImage = ebiten.NewImageFromImage(decodedImage)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
