package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/spf13/viper"
	"image/color"
	"math/rand"
	"time"
)

type Engine struct {
	width                 int
	height                int
	fields                [][]Field
	stack                 FieldStack
	cursor                Field
	finishedDrawing       bool
	numberOfVisitedFields int
	currentField          *Field
	lastUpdateTime        time.Time
	updateInterval        time.Duration
}

func (e *Engine) chooseDirRandomly(directions []Direction) Direction {
	randomIndex := rand.Intn(len(directions))
	return directions[randomIndex]
}

func (e *Engine) getPossibleDirections(field *Field) []Direction {
	var possibleDirections []Direction
	if field.row > 0 {
		if !e.fields[field.row-1][field.col].visited {
			possibleDirections = append(possibleDirections, Left)
		}
	}
	if field.row < e.width-1 {
		if !e.fields[field.row+1][field.col].visited {
			possibleDirections = append(possibleDirections, Right)
		}
	}
	if field.col > 0 {
		if !e.fields[field.row][field.col-1].visited {
			possibleDirections = append(possibleDirections, Top)
		}
	}
	if field.col < e.height-1 {
		if !e.fields[field.row][field.col+1].visited {
			possibleDirections = append(possibleDirections, Bot)
		}
	}
	return possibleDirections
}

func (e *Engine) InitEngine() {
	e.width = viper.GetInt("maze.width")
	e.height = viper.GetInt("maze.height")
	e.fields = make([][]Field, e.width)
	for i := range e.fields {
		e.fields[i] = make([]Field, e.height)
		for k := range e.fields[i] {
			e.fields[i][k] = Field{
				x:   float32(i * 25),
				y:   float32(k * 25),
				row: i,
				col: k,
			}
		}
	}
	e.currentField = &e.fields[0][0]
	e.currentField.visited = true
	e.stack.Add(e.currentField)
	e.numberOfVisitedFields = 1
	e.lastUpdateTime = time.Now()
	e.updateInterval = time.Duration(viper.GetInt("maze.speed")) * time.Millisecond
}

func (e *Engine) Update() error {
	if e.finishedDrawing {
		return nil
	}

	if time.Since(e.lastUpdateTime) < e.updateInterval {
		return nil
	}
	e.lastUpdateTime = time.Now()
	possibleDirections := e.getPossibleDirections(e.currentField)
	if len(possibleDirections) == 0 {
		if poppedField, ok := e.stack.Pop(); ok {
			e.currentField = &e.fields[poppedField.row][poppedField.col]
			e.cursor = *e.currentField
		} else {
			e.finishedDrawing = true
		}
		return nil
	}
	newDir := e.chooseDirRandomly(possibleDirections)
	switch newDir {
	case Top:
		e.currentField.isDirTop = true
		e.currentField = &e.fields[e.currentField.row][e.currentField.col-1]
	case Bot:
		e.currentField.isDirBot = true
		e.currentField = &e.fields[e.currentField.row][e.currentField.col+1]
	case Left:
		e.currentField.isDirLeft = true
		e.currentField = &e.fields[e.currentField.row-1][e.currentField.col]
	case Right:
		e.currentField.isDirRight = true
		e.currentField = &e.fields[e.currentField.row+1][e.currentField.col]
	}

	e.currentField.visited = true
	e.cursor = *e.currentField
	e.stack.Add(e.currentField)
	e.numberOfVisitedFields++

	if e.numberOfVisitedFields >= e.height*e.width {
		e.finishedDrawing = true
	}
	return nil
}

func (e *Engine) Layout(int, int) (screenWidth int, screenHeight int) {
	return e.width * 25, e.height * 25
}

func (e *Engine) Draw(screen *ebiten.Image) {
	background := color.RGBA{R: 95, G: 95, B: 95, A: 255}
	backgroundVisited := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	green := color.RGBA{G: 255, A: 255}
	blue := color.RGBA{B: 255, A: 255}
	vector.DrawFilledRect(screen, 0, 0, float32(e.width*25), float32(e.height*25), background, false)
	for _, i := range e.fields {
		for _, k := range i {
			if !k.visited {
				vector.DrawFilledRect(screen, k.x, k.y, 25, 25, backgroundVisited, false)
				continue
			}
			if !e.finishedDrawing {
				vector.DrawFilledRect(screen, e.cursor.x, e.cursor.y, 25, 25, blue, false)
			}
			if !k.isDirRight {
				if k.row < e.height-1 {
					if !e.fields[k.row+1][k.col].isDirLeft {
						vector.StrokeLine(screen, k.x+25, k.y, k.x+25, k.y+25, 1, green, false)
					}
				}
			}
			if !k.isDirLeft {
				if k.row > 0 {
					if !e.fields[k.row-1][k.col].isDirRight {
						vector.StrokeLine(screen, k.x, k.y, k.x, k.y+25, 1, green, false)
					}
				}
			}
			if !k.isDirTop {
				if k.col > 0 {
					if !e.fields[k.row][k.col-1].isDirBot {
						vector.StrokeLine(screen, k.x, k.y, k.x+25, k.y, 1, green, false)
					}
				}

			}
			if !k.isDirBot {
				if k.col < e.height-1 {
					if !e.fields[k.row][k.col+1].isDirTop {
						vector.StrokeLine(screen, k.x, k.y+25, k.x+25, k.y+25, 1, green, false)
					}
				}
			}
		}
	}
}
