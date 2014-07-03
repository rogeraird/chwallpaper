package chwallpaper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

/* A wallpaper list struct is necessary so that it can be serialized into JSON. */
type WallpaperList struct {
	Data []Wallpaper
}

type Wallpaper struct {
	Key        int      // Workspace number for this Wallpaper
	currentPos int      // Index into Wallpapers
	Wallpapers []string // All wallpapers for this workspace
}

/* Return the current wallpaper */
func (w *Wallpaper) Current() *string {
	return &w.Wallpapers[w.currentPos]
}

/* Increment the currentPos and return the next wallpaper */
func (w *Wallpaper) Next() *string {
	w.currentPos = (w.currentPos + 1) % len(w.Wallpapers)
	return &w.Wallpapers[w.currentPos]
}

/* Check if nitrogen can be used to set the wallpaper, this is the case in OpenBox */
func Nitrogen() bool {

	cmd := exec.Command("nitrogen", "--help")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return false
	}

	return true
}

/* Return the necessary property for settting the wallpaper using gsettings */
func getWallpaperProperty() string {

	var cinnamon bool

	//Check for Cinnamon first due to complex cases
	cmd := exec.Command("echo", "$CINNAMON_VERSION")
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	if out.String() != "" {
		cinnamon = true
	} else {
		cinnamon = false
	}

	if cinnamon == true {

		cmd = exec.Command("gsettings", "get",
			"org.cinnamon.background", "picture-uri")

		cmd.Stdout = &out
		err = cmd.Run()

		if err != nil {
			if err.Error() == "exit status 1" {
				return "org.gnome.desktop.background"
			} else {
				log.Fatal(err)
			}
		}

		return "org.cinnamon.background"
	}
	// Now check for Ubuntu
	cmd = exec.Command("echo", "DESKTOP_SESSION")

	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err) // Continue on instead?
	}

	if out.String() == "gnome" || out.String() == "ubuntu" {
		return "org.gnome.desktop.background"
	}

	return ""

}

/* Return the command used to set the wallpaper */
func GetCommand() string {

	if Nitrogen() {
		return "nitrogen --set-zoom-fill"
	}

	property := getWallpaperProperty()

	return "gsettings set " + property + " picture-uri 'file://"

}

func NewWallpaper(key int, wallpapers []string) *Wallpaper {
	w := new(Wallpaper)
	w.Key = key
	w.Wallpapers = wallpapers
	return w
}

/* Get the number of workspaces, currently unused, will probably be removed in future
   due to the popularity of dynamic workspace creation */
func GetNumOfWs() int {
	fmt.Println("Warning: GetNumOfWs is deprecated - it will be removed in the future!\nDon't rely on this function!")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("xprop", "-root", "_NET_NUMBER_OF_DESKTOPS")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fi, err := os.Create(usr.HomeDir + "/go.log")
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		fi.Write([]byte("Error creating command: getNumOfWs"))
		log.Fatal(err)
	}
	s := out.String()[36:38]
	s = strings.Trim(s, " \n")
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		fi, err := os.Create(usr.HomeDir + "/go.log")
		if err != nil {
			panic(err)
		}
		fi.Write([]byte("Couldn't get number of workspaces"))
		log.Fatal(err, "\n", "Couldn't get number of workspaces")
	}
	return int(i)

}

/* Get the current workspace number. Workspaces are 1-indexed. */
func GetCWs() int {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("xprop", "-root", "_NET_CURRENT_DESKTOP")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fi, err := os.Create(usr.HomeDir + "/go.log")
		if err != nil {
			log.Fatal(err)
		}
		defer fi.Close()
		fi.Write([]byte(err.Error()))
		log.Fatal(err)
	}
	s := out.String()[33:34]
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		fi, err := os.Create(usr.HomeDir + "/go.log")
		if err != nil {
			log.Fatal(err)
		}
		defer fi.Close()
		fi.Write([]byte("Couldn't get current workspace"))
		log.Fatal("Couldn't get current workspace")
	}
	return 1 + int(i)
}

/* Change the wallpaper */
func SetWallpaper(wallpaper string, command string, nitrogen bool) {
	var cmd *exec.Cmd
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if nitrogen {
		command += " " + wallpaper
		cmdSlice := strings.Split(command, " ")
		//a = append(a[:i], a[i+1:]...)

		cmd = exec.Command(cmdSlice[0], cmdSlice[1:]...)
	} else {
		command += wallpaper + "'"
		cmdSlice := strings.Split(command, " ")
		cmd = exec.Command(cmdSlice[0], cmdSlice[1:]...)
	}
	err = cmd.Run()
	if err != nil {
		fi, err := os.Create(usr.HomeDir + "/go.log")
		if err != nil {
			log.Fatal("Couldn't write file")
		}
		defer fi.Close()
		fi.Write([]byte("Couldn't set wallpaper"))
		log.Fatal("Couldn't set wallpaper")
	}
}

/* Make a WallpaperList from JSON */
func FromJson() WallpaperList {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fi, err := os.Open(usr.HomeDir + "/.wallpapers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()

	wallpapers := new(WallpaperList)
	err = json.NewDecoder(fi).Decode(wallpapers)

	if err != nil {
		log.Fatal(err)
	}

	return *wallpapers
}
