package processing

import "io"

type Decoder interface {
	DecodeNext() (interface{}, error)
}

type Encoder interface {
	Encode(v interface{}) error
}

type Worker func(v interface{}) interface{}

func Process(in Decoder, out Encoder, w Worker, workers uint) {
	processQueue := make(chan interface{}, workers*4)
	outputQueue := make(chan interface{}, workers*4)
	doneChan := make(chan int, workers)
	// Start the output encoder
	go func() {
		for result := range outputQueue {
			out.Encode(result)
		}
		doneChan <- 1
	}()
	// Start all the workers
	for i := uint(0); i < workers; i++ {
		go func() {
			for obj := range processQueue {
				result := w(obj)
				outputQueue <- result
			}
			doneChan <- 1
		}()
	}
	// Read the input, send to workers
	for {
		obj, err := in.DecodeNext()
		if err == io.EOF {
			break
		}
		processQueue <- obj
	}
	close(processQueue)
	for i := uint(1); i < workers; i++ {
		<-doneChan
	}
	<-doneChan
}
