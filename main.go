package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"text/tabwriter"
	"time"

	term "github.com/nsf/termbox-go"
)

//массив от 1 до 15, генерим случайную координату и записываем в поле(изначально заполнено 0) значение из массива, заменяя его в массиве на 0,
//когда кроме 0 не осталось- можно играть
type Field struct {
	field [][]int
	size  int
}

const (
	arrUp    = "ArrUp"
	arrDown  = "ArrDown"
	arrLeft  = "ArrLeft"
	arrRight = "ArrRight"
	space    = "Space"
)

func (f *Field) Init(size int, inline []int) {
	f.size = size
	rand.Seed(time.Now().UTC().UnixNano())

	f.field = make([][]int, f.size)
	for i := 0; i < f.size; i++ {
		f.field[i] = make([]int, f.size)
		for j := 0; j < f.size; j++ {
			f.field[i][j] = 0
		}
	}

	var tmp []int = inline //берем начальную последовательность
	//и рандомно пересталяем её 100 раз, получаем псевдослучайную расстановку
	for i := 0; i < 100; i++ {
		x := rand.Intn(15-1) + 1
		y := rand.Intn(15-1) + 1
		t := tmp[x]
		tmp[x] = tmp[y]
		tmp[y] = t
	}
	//заполняем поле
	for i := 0; i < f.size; i++ {
		f.field[i] = make([]int, f.size)
		for j := 0; j < f.size; j++ {
			//берем последнее значение из расстановки и всталяем его в поле, отрезаем
			if i+j < f.size*2-2 {
				f.field[i][j] = tmp[len(tmp)-1]
				tmp = tmp[:len(tmp)-1]
			}
		}
	}
}
func swap(x *int, y *int) {
	tmp := x
	x = y
	y = tmp
}

// func (field *Field) Move1(key string) {
// 	//двигаем нолик по полю
// 	var x, y int
// 	// i - row, j - column
// 	for i, v := range field.field { //ищем где нолик
// 		for j, vv := range v {
// 			if vv == 0 {
// 				x = i
// 				y = j
// 				switch key {
// 				case arrDown :
// 					if x != 0 {
// 						tmpfield[i][j]
// 				}
// 			}
// 			if vv == 15 {
// 				x = i
// 				y = j
// 				fmt.Println("X: " + strconv.Itoa(x) + " Y: " + strconv.Itoa(y) + " Key: " + "15")
// 			}
// 		}
// 	}
// 	fmt.Println("X: " + strconv.Itoa(x) + " Y: " + strconv.Itoa(y) + " Key: " + key)
// 	switch key {
// 	//0 уходит вниз
// 	case "ArrUp":
// 		if x != 0 {
// 			swap(&field.field[x][y], &field.field[x+1][y])
// 		}
// 	//0 уходит вверх
// 	case "ArrDown":
// 		if x != 0 {
// 			//field
// 			//swap(&field.field[x][y], &field.field[x+1][y])
// 			//field.field[x][y], field.field[x+1][y] = field.field[x+1][y], field.field[x][y]
// 			fmt.Println("нажал вниз")
// 		}
// 	//0 уходит вправо
// 	case "ArrLeft":
// 		if y != 0 {
// 			swap(&field.field[x][y], &field.field[x+1][y])
// 		}
// 	//0 уходит влево
// 	case "ArrRight":
// 		if y != 0 {
// 			swap(&field.field[x][y], &field.field[x+1][y])
// 		}
// 	}
// }
func (field *Field) Show() {
	wr := tabwriter.NewWriter(os.Stdout, 4, 0, 1, ' ', tabwriter.AlignRight)
	for _, w := range field.field {
		var str string
		for i, v := range w {
			//fmt.Print(strconv.Itoa(v) + "  ")
			if v == 0 {
				str += "#"
			} else {
				str += strconv.Itoa((v))
			}
			if i != len(w) {
				str += "\t"
			}
		}
		fmt.Fprintln(wr, str)
		wr.Flush()
	}
}

func (f *Field) Move(key string) {
	for i := 0; i < f.size; i++ {
		for j := 0; j < f.size; j++ {
			if f.field[i][j] == 0 {
				//нашли где пустая клетка, проверяем, можно ли двигать в указанном направлении
				switch key {
				case arrDown:
					{
						if i != 0 {
							tmp := f.field[i-1][j]
							f.field[i-1][j] = 0
							f.field[i][j] = tmp
						}
					}
				case arrUp:
					{
						if i != f.size-1 {
							tmp := f.field[i+1][j]
							f.field[i+1][j] = 0
							f.field[i][j] = tmp
						}
					}
				case arrLeft:
					{
						if j != f.size-1 {
							tmp := f.field[i][j+1]
							f.field[i][j+1] = 0
							f.field[i][j] = tmp
						}
					}
				case arrRight:
					{
						if j != 0 {
							tmp := f.field[i][j-1]
							f.field[i][j-1] = 0
							f.field[i][j] = tmp
						}
					}
				}
			}

		}

	}

}
func reset() {
	term.Sync() // cosmestic purpose
}

func main() {

	var f1 Field
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	initLine := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	f1.Init(4, initLine)

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()
	for { //в цикле ждем ввода сообщения

		var keyPressed string

	keyPressListenerLoop:
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyArrowUp:
				reset()
				keyPressed = arrUp
			case term.KeyArrowDown:
				reset()
				keyPressed = arrDown
			case term.KeyArrowLeft:
				reset()
				keyPressed = arrLeft
			case term.KeyArrowRight:
				reset()
				keyPressed = arrRight
			case term.KeySpace:
				reset()
				keyPressed = space
				os.Exit(0)

			default:
				reset()

			}
		case term.EventError:
			panic(ev.Err)
		}
		//двигаем фишки
		f1.Move(keyPressed)
		//чистим экран
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		//отображаем обновленное поле
		f1.Show()
	}

}
