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
	// Start the output encoder
	go func() {
		result := <-outputQueue
		out.Encode(result)
	}()
	// Start all the workers
	for i := uint(0); i < workers; i++ {
		go func() {
			for obj := range processQueue {
				result := w(obj)
				outputQueue <- result
			}
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
}
