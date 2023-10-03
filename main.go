package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/nsf/termbox-go"
)

func reverse(arr []uuid.UUID) []uuid.UUID {
	// Crear un nuevo arreglo para almacenar el resultado invertido
	invertido := make([]uuid.UUID, len(arr))

	// Invertir el arreglo
	for i, j := 0, len(arr)-1; i < len(arr); i, j = i+1, j-1 {
		invertido[i] = arr[j]
	}

	return invertido
}

type entity struct {
	x, y     int
	memory   []uuid.UUID
	maze     *Maze
	acc      int
	lastMove Direction
}

func (e *entity) Memoize() {
	uuidString := e.maze.grid[e.x][e.y].uuid.String()

	if len(e.memory) > 0 && e.memory[len(e.memory)-1].String() != uuidString {
		e.memory = append(e.memory, e.maze.grid[e.x][e.y].uuid)
	} else if len(e.memory) == 0 {
		e.memory = append(e.memory, e.maze.grid[e.x][e.y].uuid)
	}

}

func (e *entity) Render() {
	e.maze.Render(e)
	message := "acc:" + fmt.Sprint(e.acc)
	for i, char := range message {
		termbox.SetCell(i, e.maze.height, char, termbox.ColorBlack, termbox.ColorLightMagenta)
	}
	termbox.Flush()
	e.RenderMemory()
}

func (e *entity) RenderMemory() {
	historyLast := []uuid.UUID{}

	if len(e.memory) >= 10 {
		historyLast = e.memory[len(e.memory)-10:]
	} else {
		historyLast = e.memory
	}

	// Invertir el arreglo
	historyReverse := reverse(historyLast)

	for i, v := range historyReverse {
		uuidString := v.String()
		for j, ch := range ">" + uuidString {
			termbox.SetCell(j, e.maze.height+i+1, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
		termbox.Flush()
	}
}

func (e *entity) moveLeft() {
	e.Memoize()
	e.lastMove = Right
	if e.isWin() {
		return
	}

	if e.y < len(e.maze.grid[0])-1 && e.maze.grid[e.x][e.y-1].val != solid {
		e.y--
	}
}

func (e *entity) moveRight() {
	e.Memoize()
	e.lastMove = Left
	if e.isWin() {
		return
	}

	if e.y < len(e.maze.grid[0])-1 && e.maze.grid[e.x][e.y+1].val != solid {
		e.y++
	}
}

func (e *entity) moveUp() {
	e.Memoize()
	e.lastMove = Down
	if e.isWin() {
		return
	}

	if e.x < len(e.maze.grid)-1 && e.maze.grid[e.x-1][e.y].val != solid {
		e.x--
	}
}

func (e *entity) moveDown() {
	e.Memoize()
	e.lastMove = Up
	if e.isWin() {
		return
	}

	if e.x < len(e.maze.grid)-1 && e.maze.grid[e.x+1][e.y].val != solid {
		e.x++
	}
}

func (e *entity) isWin() bool {
	win := e.maze.grid[e.x][e.y].val == exit

	if win {
		message := "You win!"
		for i, char := range message {
			termbox.SetCell(i, e.maze.height, char, termbox.ColorBlack, termbox.ColorLightMagenta)
		}
		termbox.Flush()
	}

	return win
}

func (e *entity) checkWall(maze *Maze) []Direction {
	list := []Direction{}

	if maze.grid[e.x][e.y+1].val != solid {
		list = append(list, Down)
	}
	if maze.grid[e.x][e.y-1].val != solid {
		list = append(list, Up)
	}
	if maze.grid[e.x+1][e.y].val != solid {
		list = append(list, Right)
	}
	if maze.grid[e.x-1][e.y].val != solid {
		list = append(list, Left)
	}

	found := false

	for _, dir := range list {
		if e.lastMove == dir {
			found = true
			break
		}
	}

	if found {
		return []Direction{e.lastMove}
	}

	return list
}

type block int

const (
	empty block = iota
	solid
	exit
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type node int

const (
	regular node = iota
	wall
	branch
	deadPoint
	merge
)

type Cell struct {
	uuid     uuid.UUID
	val      block
	nodeType node
}

func genCells(maze [][]block) [][]Cell {
	cells := make([][]Cell, len(maze))

	for i, row := range maze {
		cells[i] = make([]Cell, len(row))

		for j, val := range row {
			cells[i][j] = Cell{
				uuid:     uuid.New(),
				val:      block(val),
				nodeType: regular,
			}
		}
	}

	return cells
}

type Maze struct {
	width, height int
	grid          [][]Cell
}

var defaultMaze = [][]block{
	{solid, solid, solid, solid, solid, solid, solid, solid, solid},
	{solid, empty, empty, empty, empty, empty, empty, empty, solid},
	{solid, empty, solid, solid, solid, solid, solid, empty, solid},
	{solid, empty, solid, empty, empty, empty, solid, empty, solid},
	{solid, empty, solid, empty, solid, empty, solid, empty, solid},
	{solid, empty, solid, empty, solid, exit, solid, empty, solid},
	{solid, empty, solid, empty, solid, solid, solid, empty, solid},
	{solid, empty, empty, empty, empty, empty, solid, empty, solid},
	{solid, solid, solid, solid, solid, solid, solid, solid, solid},
}

// NewMaze crea un nuevo objeto Maze con el tamaño especificado.
func NewMaze() (*Maze, error) {
	maze := &Maze{
		width:  len(defaultMaze[0]),
		height: len(defaultMaze),
		grid:   genCells(defaultMaze),
	}

	return maze, nil
}

// Render dibuja el laberinto en la terminal.
func (m *Maze) Render(e *entity) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i, row := range m.grid {

		for j, cell := range row {
			char := ' '
			fg := termbox.ColorDefault
			bg := termbox.ColorLightGray
			if cell.val == empty {
				char = ' '
				fg = termbox.ColorDefault
				bg = termbox.ColorDefault
			}

			if cell.val == solid {
				char = 'W'
				fg = termbox.ColorBlack
				bg = termbox.ColorYellow
			}

			if cell.val == exit {
				char = 'S'
				fg = termbox.ColorRed
				bg = termbox.ColorLightGreen
			}

			if e.x == i && e.y == j {
				char = 'T'
				fg = termbox.ColorGreen
				bg = termbox.ColorDefault
			}

			termbox.SetCell(j, i, char, fg, bg)
		}
	}

	termbox.Flush()
}

// function for check if up, down, left or right is a wall

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}

	defer termbox.Close()

	maze, errMaze := NewMaze()

	if errMaze != nil {
		fmt.Println(errMaze)
		return
	}

	e := entity{
		x:        1,
		y:        1,
		memory:   []uuid.UUID{},
		maze:     maze,
		acc:      10,
		lastMove: Up,
	}

	check := e.checkWall(maze)
	e.lastMove = check[rand.Intn(len(check))]

	e.Render()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	ticker := time.NewTicker(time.Second / time.Duration(e.acc))

	go func() {
		defer wg.Done()

		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if e.isWin() {
					break
				}
				check := e.checkWall(maze)
				// fmt.Println(check)
				selCheck := check[rand.Intn(len(check))]
				// Mover aleatoriamente en una de las cuatro direcciones

				switch selCheck {
				case 0:
					e.moveLeft()
					e.Render()
				case 1:
					e.moveRight()
					e.Render()
				case 2:
					e.moveUp()
					e.Render()
				case 3:
					e.moveDown()
					e.Render()
				}

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
					e.moveDown()
					e.Render()
				case termbox.KeyArrowUp:
					e.moveUp()
					e.Render()
				case termbox.KeyArrowLeft:
					e.moveLeft()
					e.Render()
				case termbox.KeyArrowRight:
					e.moveRight()
					e.Render()
				case termbox.KeyCtrlA:
					if e.acc-1 > 0 {
						e.acc--
					}
					ticker = time.NewTicker(time.Second / time.Duration(e.acc))
				case termbox.KeyCtrlS:
					e.acc++
					ticker = time.NewTicker(time.Second / time.Duration(e.acc))
				}

				if e.isWin() {
					break
				}

			}
		}
	}()

	// Espera hasta que se presiona Esc o se recibe una señal de interrupción
	wg.Wait()
}
