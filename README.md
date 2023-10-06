# The maze

## Descripcion

El programa consiste en un laberinto en el cual se debe encontrar la salida, el programa esta hecho en go y se ejecuta en la consola de comandos.

## Como funciona?

![alt](/laberintSolve.png)

*Se observa que hay tres tipos de bloques en el laberinto:*

* Solido
* Espacio vacio
* Salida

**La entidad se mueve aleatoriamente en una de las cuatro direcciones**, si la entidad llega a la salida se detiene el hilo de ejecución.

**Se toma prioriza ir hacia adelante en la direccion actual**,  por ejemplo si la entidad se mueve hacia arriba, se prioriza ir hacia arriba, si no puede ir hacia arriba, se mueve aleatoriamente en una de las otras tres direcciones.

Es importante que la entidad no tiene una memoria establecida de por donde ha pasado, por lo que puede pasar varias veces por el mismo lugar. Lo que hace que el programa sea ineficiente, pero aun asi se puede llegar a la salida.

## Codigo

### struct entity

```go
//El enum Direction es para saber hacia donde se mueve la entidad
type Direction int

const (
 Up Direction = iota
 Down
 Left
 Right
)

// La struct entity es la instancia que se movera en el laberinto

type entity struct {
 x, y     int
 memory   []uuid.UUID
 maze     *Maze
 acc      int
 lastMove Direction
}

// Memoize guarda la posición actual en la memoria
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

// RenderMemory muestra la memoria de la entidad

func (e *entity) RenderMemory() {
 historyLast := []uuid.UUID{}

 if len(e.memory) >= 10 {
  historyLast = e.memory[len(e.memory)-10:]
 } else {
  historyLast = e.memory
 }

 historyReverse := reverse(historyLast)

 for i, v := range historyReverse {
  uuidString := v.String()
  for j, ch := range ">" + uuidString {
   termbox.SetCell(j, e.maze.height+i+1, ch, termbox.ColorDefault, termbox.ColorDefault)
  }
  termbox.Flush()
 }
}

// moveLeft mueve la entidad hacia la izquierda

func (e *entity) moveLeft() {
 e.Memoize()
 e.lastMove = Right
 if e.isWin() {
  return
 }

 if e.y < len(e.maze.grid[0])-1 && e.maze.grid[e.x][e.y-1].val != so {
  e.y--
 }
}

// moveRight mueve la entidad hacia la derecha

func (e *entity) moveRight() {
 e.Memoize()
 e.lastMove = Left
 if e.isWin() {
  return
 }

 if e.y < len(e.maze.grid[0])-1 && e.maze.grid[e.x][e.y+1].val != so {
  e.y++
 }
}

// moveUp mueve la entidad hacia arriba

func (e *entity) moveUp() {
 e.Memoize()
 e.lastMove = Down
 if e.isWin() {
  return
 }

 if e.x < len(e.maze.grid)-1 && e.maze.grid[e.x-1][e.y].val != so {
  e.x--
 }
}

// moveDown mueve la entidad hacia abajo

func (e *entity) moveDown() {
 e.Memoize()
 e.lastMove = Up
 if e.isWin() {
  return
 }

 if e.x < len(e.maze.grid)-1 && e.maze.grid[e.x+1][e.y].val != so {
  e.x++
 }
}

// isWin verifica si la entidad llego a la salida

func (e *entity) isWin() bool {
 win := e.maze.grid[e.x][e.y].val == ex

 if win {
  message := "You win!"
  for i, char := range message {
   termbox.SetCell(i, e.maze.height, char, termbox.ColorBlack, termbox.ColorLightMagenta)
  }
  termbox.Flush()
 }

 return win
}

// checkWall verifica si hay paredes alrededor de la entidad

func (e *entity) checkWall(maze *Maze) []Direction {
 list := []Direction{}

 if maze.grid[e.x][e.y+1].val != so {
  list = append(list, Down)
 }
 if maze.grid[e.x][e.y-1].val != so {
  list = append(list, Up)
 }
 if maze.grid[e.x+1][e.y].val != so {
  list = append(list, Right)
 }
 if maze.grid[e.x-1][e.y].val != so {
  list = append(list, Left)
 }

 return list
}
```

### struct Maze

```go
// block es el tipo de bloque que se puede encontrar en el laberinto
type block int

const (
 em block = iota
 so
 ex
)

// Cell es una celda del laberinto diferenciada por una uuid, tipo de bloque y coordenadas
type Cell struct {
 uuid     uuid.UUID
 val      block
 nodeType node
}

// genCells genera una matriz de celdas a partir de un array de blocks

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

// Maze es el laberinto representado en una matriz de celdas

type Maze struct {
 width, height int
 grid          [][]Cell
}

// NewMaze construye el laberinto a partir de una matriz de celdas

func NewMaze(mazeArr [][]block) (*Maze, error) {
 maze := &Maze{
  width:  len(mazeArr[0]),
  height: len(mazeArr),
  grid:   genCells(mazeArr),
 }

 return maze, nil
}
```

### Estado inicial del laberinto

Se representa con una matriz de bloques

```go
var customMaze = [][]block{
 {so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so},
 {so, em, em, so, em, em, em, em, so, em, em, so, so, em, so, so, so, so, so, em, em, em, em, em, so},
 {so, em, so, so, em, so, so, em, so, em, so, so, so, em, so, em, so, so, em, so, so, so, so, em, so},
 {so, em, em, em, em, em, so, em, em, em, em, em, em, em, em, em, em, em, em, so, em, em, so, em, so},
 {so, em, so, em, so, em, so, so, so, so, so, so, so, so, em, so, em, so, so, so, em, so, so, em, so},
 {so, em, so, em, so, em, so, so, so, so, so, em, em, em, em, so, em, em, em, em, em, em, so, em, so},
 {so, em, so, em, so, em, em, em, em, em, em, em, so, so, so, so, em, so, so, em, so, em, so, em, so},
 {so, so, so, so, so, em, so, so, em, so, so, em, so, em, em, so, em, so, em, em, so, em, em, em, so},
 {so, so, em, em, em, em, so, so, em, so, so, em, so, so, so, so, em, so, em, so, so, so, so, so, so},
 {so, em, em, so, so, em, em, em, em, em, so, em, so, so, so, so, em, so, em, so, so, so, so, so, so},
 {so, em, so, so, so, em, so, so, so, em, so, em, em, em, em, so, em, so, em, em, so, so, so, so, so},
 {so, em, so, so, so, em, em, so, so, em, so, em, so, so, em, so, em, so, so, em, em, so, so, so, so},
 {so, em, em, em, em, so, em, so, so, em, so, so, so, so, em, so, em, em, em, so, em, so, so, so, so},
 {so, so, so, em, so, so, em, em, em, em, so, em, em, em, em, so, so, so, em, so, em, so, so, so, so},
 {so, so, so, em, so, em, em, so, so, so, so, em, so, so, em, so, so, so, em, so, em, so, so, so, so},
 {so, so, so, em, so, so, em, so, em, so, so, em, so, so, em, so, so, so, em, so, em, so, so, so, so},
 {so, em, em, em, so, so, em, so, em, em, em, em, so, em, so, so, em, em, em, so, em, so, so, so, so},
 {so, em, so, em, em, em, em, so, em, so, so, so, so, em, so, so, so, so, so, so, em, so, so, so, so},
 {so, em, so, so, so, so, so, so, em, so, em, em, em, em, so, em, em, em, em, em, em, so, so, so, so},
 {so, em, em, em, so, em, em, em, em, so, em, so, so, em, em, ex, so, so, so, so, so, so, so, so, so},
 {so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so, so},
}
```

### Estado inicial de la entidad

```go
e := entity{
  x:        1,
  y:        1,
  memory:   []uuid.UUID{},
  maze:     maze,
  acc:      10,
  lastMove: Up,
 }

```

### Hilo de ejecución

Se ejecuta una go routine que se encarga de mover la entidad en el laberinto, se mueve aleatoriamente en una de las cuatro direcciones, si la entidad llega a la salida se detiene el hilo de ejecución.

```go
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
```

### Ejecutar el programa

Se deben ejecutar los sigueintes coamndos:

```bash
go install # Para instalar las dependencias
go run main.go # Para ejecutar el programa
```

El programa inicia automaticamente con una velocidad por defecto, pero usando la conbinacion de teclas `ctrl + s` se puede aumantar la velocidad y con `ctrl + a` se puede disminuir la velocidad.
