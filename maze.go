package main

import (
	"fmt"
	"os"
	"time"
	"os/exec"
	"math"
	"math/rand"
)

type stage struct {
	world []string
	wallColor string
	crashWallColor string
	floorColor string
	width int
	height int
}

var s stage

type player struct {
	x int
	y int
	character rune
	color string
	movements int
	crashes int
	crashed bool
}

func checkCrash(x int, y int) bool {
	return string([]rune(s.world[y])[x]) == "█"
}

var p player

func (p *player) move(key string) {
	if (gameOver) {
		return
	}

	p.crashed = false

	switch key {
		case "w", "W":
			nextY := int(math.Max(float64(p.y - 1), 0.0))
			if (!checkCrash(p.x, nextY)) {
				p.y = nextY
				p.movements++
			} else {
				p.crash()
			}
			break
		case "a", "A":		
			nextX := int(math.Max(float64(p.x - 1), 0.0))
			if (!checkCrash(nextX, p.y)) {
				p.x = nextX
				p.movements++
			} else {
				p.crash()
			}
			break
		case "s", "S":		
			nextY := int(math.Min(float64(p.y + 1), float64(s.height - 1)))
			if (!checkCrash(p.x, nextY)) {
				p.y = nextY
				p.movements++
			} else {
				p.crash()
			}
			break
		case "d", "D":
			nextX := int(math.Min(float64(p.x + 1), float64(s.width - 1)))
			if (!checkCrash(nextX, p.y)) {
				p.x = nextX
				p.movements++
			} else {
				p.crash()
			}
			break
	}

	if (p.x == t.x && p.y == t.y) {
		gameOver = true
	}
}

var gameOver bool

func (p *player) crash() {
	fmt.Print("\a")
	p.crashes++
	p.crashed = true
}

type target struct {
	x int
	y int
	character rune
	color string
}

var t target

func (t *target) move() {
	nextX := t.x + 1 - rand.Intn(3)
	nextY := t.y + 1 - rand.Intn(3)
	
	if (string([]rune(s.world[t.y])[nextX]) != "█") {
		t.x = nextX
	}
	
	if (string([]rune(s.world[nextY])[t.x]) != "█") {
		t.y = nextY
	}
}

func draw() {
	fmt.Print("[H[2J")
	if (gameOver) {
		fmt.Println("You win. Game over.")
		fmt.Printf("[38;5;39m[1mScore: %d[m", p.movements - p.crashes)
	} else {
		wallColor := s.wallColor

		if (p.crashed) {
			wallColor = s.crashWallColor
		}

		for y, xLen := 0, len([]rune(s.world[0])); y < len(s.world); y++ {
			fmt.Printf("%s%s", wallColor, s.floorColor)
			for x := 0; x < xLen; x++ {
				if (x == p.x && y == p.y) {
					fmt.Printf("%s%c[m%s%s", p.color, p.character, wallColor, s.floorColor)
				} else if (x == t.x && y == t.y) {
					fmt.Printf("%s%c[m%s%s", t.color, t.character, wallColor, s.floorColor)
				} else {
					fmt.Print(string([]rune(s.world[y])[x]))
				}
			}
			fmt.Println("[m")
		}
	}

	fmt.Printf("\nMovements: %d Crashes: %d\n", p.movements, p.crashes)
}

func main() {
	fmt.Println("Loading...")

	rand.Seed(time.Now().UnixNano())

	gameOver = false

	s = stage { world: []string{
		"██████████████████████████████████████████████████",
		"██  ██                       ██       ██        ██",
		"██  ██          ███████████             ██  ██  ██",
		"██  ██████████  ██            ██████    ██      ██",
		"██          ██  ██  ██  ██    ██                ██",
		"██    ████  ██  ██  ██  ██    ██  ██    ██  ██  ██",
		"██    ██        ██  ██  ██    ██████        ██  ██",
		"██    ████████████  ██      ██                  ██",
		"██    ██            ██████              ██  ██  ██",
		"██    ███████████       ██████   █████████      ██",
		"██                 ███      ██          ██  ██  ██",
		"██   ████████████      ██   ██      ██  ██████  ██",
		"██             ██  ██████     ██    ██      ██  ██",
		"██                          ██          ██  ██  ██",
		"██████████████████████████████████████████████████"},
		wallColor: "[38;5;231m",
		crashWallColor: "[38;5;1m",
		floorColor: "[48;5;248m",
		width: 50,
		height: 15}

	p = player{ x: 2, y: 2, character: '#', color: "[38;5;9m", 
				movements: 0, crashes: 0, crashed: false }

	t = target { x: 46, y: 2, character: '@', color: "[38;5;17m" }
	
	draw()

	ch := make(chan string)
	go func(ch chan string) {
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		var b []byte = make([]byte, 1)

		for {
			os.Stdin.Read(b)
			ch <- string(b)
		}
	}(ch)

	for {
		select {
			case key, _ := <-ch:		
				p.move(key)		
				draw()
				t.move()
		}

		time.Sleep(time.Millisecond * 100)
	}
}