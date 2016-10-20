package storage

import (
	"io"
	"math/rand"
	"sync"

	"github.com/denbeigh2000/jfsi"
)

func parallelReaderGroup(r io.Reader, n int) ([]io.Reader, <-chan error) {
	writers := make([]io.Writer, n)
	readers := make([]io.Reader, n)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		readers[i], writers[i] = io.Pipe()
	}

	writer := io.MultiWriter(writers...)
	go func() {
		_, err := io.Copy(writer, r)
		if err != nil {
			errs <- err
		}

		for _, pipe := range writers {
			wc := pipe.(io.WriteCloser)
			wc.Close()
		}
		close(errs)
	}()

	return readers, errs
}

func ParallelCreate(clients []Store, id jfsi.ID, r io.Reader) error {
	readers, copyErrs := parallelReaderGroup(r, len(clients))
	errs := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(len(clients) + 1)

	go func() {
		for err := range copyErrs {
			errs <- err
		}
		wg.Done()
	}()

	for i, read := range readers {
		go func(s Store, read io.Reader, i int) {
			err := s.Create(id, read)
			if err != nil {
				errs <- err
			}
			wg.Done()
		}(clients[i], read, i)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	// If we see any errors at all before the whole thing finishes,
	// return them in place of a success because we have no idea
	// if the whole thing will finish or not
	for err := range errs {
		return err
	}

	return nil
}

func ParallelUpdate(clients []Store, id jfsi.ID, r io.Reader) error {
	readers, copyErrs := parallelReaderGroup(r, len(clients))
	errs := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(len(clients) + 1)

	go func() {
		for err := range copyErrs {
			errs <- err
		}

		wg.Done()
	}()

	for i, read := range readers {
		go func(s Store, read io.Reader) {
			err := s.Update(id, read)
			if err != nil {
				errs <- err
			}
			wg.Done()
		}(clients[i], read)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	// If we see any errors at all before the whole thing finishes,
	// return them in place of a success because we have no idea
	// if the whole thing will finish or not
	for err := range errs {
		return err
	}

	return nil
}

func ParallelDelete(clients []Store, id jfsi.ID) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(clients))

	errs := make(chan error)

	for _, client := range clients {
		go func(s Store) {
			err := s.Delete(id)
			if err != nil {
				errs <- err
			}
			wg.Done()
		}(client)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		return err
	}

	return nil
}

func SelectiveRetrieve(clients []Store, id jfsi.ID) (io.Reader, error) {
	// Shuffling clients gives some non-determinism so load is distributed
	// Fisher-Yates shuffle implementation shamelessly stolen from
	// https://stackoverflow.com/a/12267471
	for i := range clients {
		j := rand.Intn(i + 1)
		clients[i], clients[j] = clients[j], clients[i]
	}

	var r io.Reader
	var err error

	for _, client := range clients {
		r, err = client.Retrieve(id)
		if err == nil {
			break
		} else {
			// ensure we don't erroneously pass back a non-nil reader
			// in any possible circumstance
			r = nil
		}
	}

	return r, err
}
