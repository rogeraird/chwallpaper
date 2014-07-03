package main

import (
	"fmt"
	ch "github.com/rogeraird/chwallpaper"
	"time"
)

func main() {

	command := ch.GetCommand()

	if command == "" {
		panic("Could not get wallpaper property - Please use a recognised DE/WM")
	}

	nitrogen := ch.Nitrogen()
	wl := ch.FromJson()

	wp := make(map[int]*ch.Wallpaper)

	for i := range wl.Data {
		wp[wl.Data[i].Key] = &wl.Data[i]
	}

	time.Sleep(5 * time.Second)

	//fmt.Println("Woke up")

	startWs := ch.GetCWs()

	var currentWs int
	var counter int

	ch.SetWallpaper(*wp[startWs].Current(), command, nitrogen)

	for {
		currentWs = ch.GetCWs()

		fmt.Println(currentWs, " ", startWs)
		if currentWs != startWs {
			tempWs := currentWs
			if currentWs > len(wp) {
				tempWs = (currentWs % len(wp))
			}
			counter = 0
			ch.SetWallpaper(*wp[tempWs].Current(), command, nitrogen)
			startWs = currentWs
		} else {
			counter++

			// Multiply by 2 to account for 500ms sleep
			if counter > (30*2) && len(wp[currentWs].Wallpapers) > 1 {
				ch.SetWallpaper(*wp[currentWs].Next(), command, nitrogen)
				counter = 0
			}
		}

		time.Sleep(500 * time.Millisecond)

	}

}
