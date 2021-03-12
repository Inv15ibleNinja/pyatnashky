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

func (field *Field) Move(key string) {
	//двигаем нолик по полю
	for i, v := range field.field { //ищем где нолик
		for j, vv := range v {
			if vv == 0 {
				//нашли нолик, проверим, можно ли двигать - крайние положения где i==0||i==siza-1 или j ==0 ||j== size-1
				//лучше написать отдельную функцию првоерки, в которую передавать координаты начала и направление

			}
		}
	}

	switch key {
	case "ArrUp":
	//0 уходит вниз

	case "ArrDown":
	//0 уходит вверх
	case "ArrLeft":
	//0 уходит вправо
	case "ArrRight":
		//0 уходит влево

	}
}

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
				keyPressed = "ArrUp"
			case term.KeyArrowDown:
				reset()
				keyPressed = "ArrDown"
			case term.KeyArrowLeft:
				reset()
				keyPressed = "ArrLeft"
			case term.KeyArrowRight:
				reset()
				keyPressed = "ArrRight"
			case term.KeySpace:
				reset()
				keyPressed = "Space"
				os.Exit(0)

			default:
				// we only want to read a single character or one key pressed event
				reset()
				//fmt.Println("ASCII : ", ev.Ch)

			}
		case term.EventError:
			panic(ev.Err)
		}
		f1.Move(keyPressed)
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		f1.Show()
	}

}
