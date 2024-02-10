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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

var (
	frameOX = 0
	frameOY = 0
)

var (
	runnerImage  *ebiten.Image
	runnerImage2 *ebiten.Image
)

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
	for _, p := range g.keys {
		if ebiten.KeyName(p) == "n" {
			frameOY += 1
		}
		if ebiten.KeyName(p) == "p" {
			frameOY -= 1
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	i := (g.count / 10) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
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
	imageData, err := ioutil.ReadFile("asdf.jpg")
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
