package belajargolangcontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)


func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}


func TestContextWithValue(t *testing.T) {
	contextA := context.Background() //! Parent bisa di sebut juga sebagai context background

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")

	contextG := context.WithValue(contextF, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)


	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))
	fmt.Println(contextA.Value("b"))
}

/*
	? menambahkan data dalam context bisa menggunakan function context.With.Value(), bukan merubah context, tapi membuat context baru atau child
*/


func CreateCounter(cxt context.Context) chan int {
	destination := make(chan int)

	go func() { //! go routine
		defer close(destination)
		counter := 1
		for {
			select {
			case <- cxt.Done():
				return
			default:
				destination <- counter
				counter ++
				time.Sleep(1 * time.Second) //! simulasi slow
			}
		}
	}()

	return destination
}


func TestContextWithCancel(t *testing.T) {
	fmt.Println("total goroutine before", runtime.NumGoroutine()) //! di awal ada 2 goroutine

	parent := context.Background() //! penambahan untuk sinyal cancel
	cxt, cancel := context.WithCancel(parent) //! penambahan untuk sinyal cancel

	destination := CreateCounter(cxt) //! penambahan parameter untuk sinyal cancel

	fmt.Println("total goroutine in proses", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break;
		}
	}

	
	cancel() //! mengirim sinyal cancel ke context (penambahan)
	time.Sleep(2 * time.Second)

	fmt.Println("total goroutine after", runtime.NumGoroutine()) //! di akhir bertambah goroutinenya yang di sebut goroutine leak, satu goroutine nyala terus dan tidak pernah mati.


	/*
		! berbahaya kalau misalkan ada satu goroutine yang jalan terus, case nya jika setiap request ada goroutine, maka setiap request tersebut kebentuk juga goroutine leak, makin lambat dan memori konsumsi makin tinggi dan worst case app mati.

		TODO: Untuk membatalkan terjadinya goroutine leak kita bisa menggunakan context.Context sebagai parameter di func. kemudian menggunakan select pada perulangan FOR LOOP, membuat sebuah parent dan menggunakan func context.WithCancel()
	*/ 
}


func TestContextWithTimeout(t *testing.T) {
	fmt.Println("total goroutine before", runtime.NumGoroutine())

	parent := context.Background()
	cxt, cancel := context.WithTimeout(parent, 5 * time.Second) //! Pembatalan secara otomatis menggunkana timeout, dan sinyal cancel otomatis di kirim
	defer cancel()

	destination := CreateCounter(cxt)
	fmt.Println("total goroutine in proses", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break;
		}
	}
	time.Sleep(2 * time.Second)

	fmt.Println("total goroutine after", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("total goroutine before", runtime.NumGoroutine())

	parent := context.Background()
	cxt, cancel := context.WithDeadline(parent, time.Now().Add(5 * time.Second)) //! Pembatalan secara otomatis menggunkana deadline, dan sinyal cancel otomatis di kirim
	defer cancel()

	destination := CreateCounter(cxt)
	fmt.Println("total goroutine in proses", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break;
		}
	}
	time.Sleep(2 * time.Second)

	fmt.Println("total goroutine after", runtime.NumGoroutine())
}