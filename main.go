package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nsf/termbox-go"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Maze struct {
	width, height int
	grid          [][]int
}

var defaultMaze = [][]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 1, 1, 1, 1, 0, 1},
	{1, 0, 1, 0, 0, 0, 1, 0, 1},
	{1, 0, 1, 0, 1, 0, 1, 0, 1},
	{1, 0, 1, 0, 1, 3, 1, 0, 1},
	{1, 0, 1, 0, 1, 1, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 1, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1},
}

// NewMaze crea un nuevo objeto Maze con el tamaño especificado.
func NewMaze() (*Maze, error) {
	maze := &Maze{
		width:  len(defaultMaze[0]),
		height: len(defaultMaze),
		grid:   defaultMaze,
	}

	return maze, nil
}

// Render dibuja el laberinto en la terminal.
func (m *Maze) Render(pos *map[string]int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i, row := range m.grid {

		for j, value := range row {
			char := ' '
			if value == 0 {
				char = ' '
			}

			if value == 1 {
				char = 'W'
			}

			if value == 3 {
				char = 'S'
			}

			if (*pos)["x"] == i && (*pos)["y"] == j {
				char = 'T'
			}

			termbox.SetCell(j, i, char, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

func moveLeft(maze *Maze, pos *map[string]int) {
	if isWin(maze, pos) {
		return
	}

	if (*pos)["y"] < len(maze.grid[0])-1 && maze.grid[(*pos)["x"]][(*pos)["y"]-1] != 1 {
		(*pos)["y"]--
	}
}

func moveRight(maze *Maze, pos *map[string]int) {
	if isWin(maze, pos) {
		return
	}

	if (*pos)["y"] < len(maze.grid[0])-1 && maze.grid[(*pos)["x"]][(*pos)["y"]+1] != 1 {
		(*pos)["y"]++
	}
}

func moveUp(maze *Maze, pos *map[string]int) {
	if isWin(maze, pos) {
		return
	}

	if (*pos)["x"] < len(maze.grid)-1 && maze.grid[(*pos)["x"]-1][(*pos)["y"]] != 1 {
		(*pos)["x"]--
	}
}

func moveDown(maze *Maze, pos *map[string]int) {
	if isWin(maze, pos) {
		return
	}

	if (*pos)["x"] < len(maze.grid)-1 && maze.grid[(*pos)["x"]+1][(*pos)["y"]] != 1 {
		(*pos)["x"]++
	}
}

func isWin(maze *Maze, pos *map[string]int) bool {
	message := "You win!"

	if maze.grid[(*pos)["x"]][(*pos)["y"]] == 3 {
		for _, char := range message {
			// fmt.Println(string(char))
			termbox.SetCell(0, len(maze.grid)+2, char, termbox.ColorDefault, termbox.ColorDefault)
		}
		termbox.Flush()
		return true
	}

	return false
}

// function for check if up, down, left or right is a wall
func checkWall(maze *Maze, pos *map[string]int) []Direction {
	list := []Direction{}

	if maze.grid[(*pos)["x"]][(*pos)["y"]+1] != 1 {
		list = append(list, Down)
	}
	if maze.grid[(*pos)["x"]][(*pos)["y"]-1] != 1 {
		list = append(list, Up)
	}
	if maze.grid[(*pos)["x"]+1][(*pos)["y"]] != 1 {
		list = append(list, Right)
	}
	if maze.grid[(*pos)["x"]-1][(*pos)["y"]] != 1 {
		list = append(list, Left)
	}

	return list
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}

	defer termbox.Close()
	pos := map[string]int{"x": 1, "y": 1}

	maze, _ := NewMaze()
	maze.Render(&pos)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Second / 20)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if isWin(maze, &pos) {
					break
				}
				check := checkWall(maze, &pos)
				// fmt.Println(check)
				selCheck := check[rand.Intn(len(check))]
				// Mover aleatoriamente en una de las cuatro direcciones

				switch selCheck {
				case 0:
					moveLeft(maze, &pos)
				case 1:
					moveRight(maze, &pos)
				case 2:
					moveUp(maze, &pos)
				case 3:
					moveDown(maze, &pos)
				}

				maze.Render(&pos)
			}
		}
	}()

	go func() {
		defer wg.Done()

		for {
			ev := termbox.PollEvent()
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					return
				case termbox.KeyArrowDown:
					moveDown(maze, &pos)
				case termbox.KeyArrowUp:
					moveUp(maze, &pos)
				case termbox.KeyArrowLeft:
					moveLeft(maze, &pos)
				case termbox.KeyArrowRight:
					moveRight(maze, &pos)
				}

				maze.Render(&pos)
			}
		}
	}()

	// Espera hasta que se presiona Esc o se recibe una señal de interrupción
	wg.Wait()
}
