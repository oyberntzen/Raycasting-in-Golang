package main

import (
	"encoding/gob"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/oyberntzen/Raycasting-in-Golang/game"
)

var env game.Enviroment
var players sync.Map

func main() {
	env = game.Enviroment{}
	path, _ := os.Getwd()
	imagesPath := filepath.Dir(filepath.Dir(path)) + "/images/"
	env.Init(game.Level01, imagesPath)
	gob.Register(game.Enviroment{})

	l, _ := net.Listen("tcp", ":8000")
	defer l.Close()

	for i := 0; true; i++ {
		c, _ := l.Accept()
		go playerConnection(c, i)
	}

}

func playerConnection(c net.Conn, id int) {
	enc := gob.NewEncoder(c)
	dec := gob.NewDecoder(c)

	var width, height int
	handleError(dec.Decode(&width))
	handleError(dec.Decode(&height))

	p := game.Player{}
	p.Init(22.5, 10.5, -math.Pi/2, width, height)
	//num := len(players)
	players.Store(id, p)

	handleError(enc.Encode(env))
	handleError(enc.Encode(p))
	handleError(enc.Encode(id))

	for {
		err := dec.Decode(&p)
		players.Store(id, p)
		if err != nil {
			players.Delete(id)
			break
		}
		newMap := make(map[int]game.Player)
		players.Range(func(key, value interface{}) bool {
			i, player := key.(int), value.(game.Player)
			newMap[i] = player
			return true
		})
		handleError(enc.Encode(newMap))
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
