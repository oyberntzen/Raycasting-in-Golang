package main

import (
	"encoding/gob"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"

	"github.com/oyberntzen/Raycasting-in-Golang/game"
)

var env game.Enviroment
var players []game.Player

func main() {
	env = game.Enviroment{}
	path, _ := os.Getwd()
	imagesPath := filepath.Dir(filepath.Dir(path)) + "/images/"
	env.Init(game.Level01, imagesPath)
	gob.Register(game.Enviroment{})

	l, _ := net.Listen("tcp", ":8000")
	defer l.Close()

	for {
		c, _ := l.Accept()
		go playerConnection(c)
	}

}

func playerConnection(c net.Conn) {
	enc := gob.NewEncoder(c)
	dec := gob.NewDecoder(c)

	var width, height int
	handleError(dec.Decode(&width))
	handleError(dec.Decode(&height))

	p := game.Player{}
	p.Init(22.5, 10.5, -math.Pi/2, width, height)
	//num := len(players)
	num := len(players)
	players = append(players, p)
	player := &players[num]

	handleError(enc.Encode(env))
	handleError(enc.Encode(player))
	handleError(enc.Encode(num))

	for {
		err := dec.Decode(&player)
		if err != nil {
			players = append(players[:num], players[num+1:]...)
			break
		}
		handleError(enc.Encode(players))
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
