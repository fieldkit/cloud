package backend

import (
	"context"
	"io"
	"sync"

	"github.com/hashicorp/go-multierror"
)

type AsyncWriterFunc func(ctx context.Context, writer io.Writer) error

type AsyncReaderFunc func(ctx context.Context, reader io.Reader) error

type AsyncFileWriter struct {
	reader *io.PipeReader
	writer *io.PipeWriter
	read   AsyncReaderFunc
	write  AsyncWriterFunc
	wg     sync.WaitGroup
	errors *multierror.Error
}

func NewAsyncFileWriter(readFunc AsyncReaderFunc, writeFunc AsyncWriterFunc) *AsyncFileWriter {
	reader, writer := io.Pipe()

	return &AsyncFileWriter{
		reader: reader,
		writer: writer,
		read:   readFunc,
		write:  writeFunc,
	}
}

func (w *AsyncFileWriter) Start(ctx context.Context) error {
	w.wg.Add(2)

	go func() {
		defer w.reader.Close()

		err := w.read(ctx, w.reader)
		if err != nil {
			w.errors = multierror.Append(w.errors, err)
		}

		w.wg.Done()
	}()

	go func() {
		defer w.writer.Close()

		err := w.write(ctx, w.writer)
		if err != nil {
			w.errors = multierror.Append(w.errors, err)
		}

		w.wg.Done()
	}()

	return nil
}

func (w *AsyncFileWriter) Wait(ctx context.Context) error {
	w.wg.Wait()
	return w.errors.ErrorOrNil()
}
