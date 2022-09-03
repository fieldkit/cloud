package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/fieldkit/cloud/server/common/logging"
	"go.uber.org/zap"
)

const (
	TailLength = 20
)

type PileKey string

type pileEntry struct {
	Key   PileKey   `json:"key"`
	Added time.Time `json:"added"`
	Size  int64     `json:"size"`
}

type PileMeta struct {
	Size int64        `json:"size"`
	Tail []*pileEntry `json:"tail"`
}

type Pile struct {
	lock    sync.RWMutex
	metrics *logging.Metrics
	name    string
	path    string
	log     *zap.SugaredLogger
	meta    *PileMeta
}

func NewPile(metrics *logging.Metrics, name string) *Pile {
	return &Pile{
		metrics: metrics,
		name:    name,
		path:    path.Join("/tmp/pile", name),
		meta:    nil,
		log:     nil,
	}
}

func (pile *Pile) IsOpen() bool {
	return pile.meta != nil
}

func (pile *Pile) Open(ctx context.Context) error {
	pile.log = logging.Logger(ctx).Named("pile").Sugar().With("pile", pile.path)

	pile.lock.Lock()

	defer pile.lock.Unlock()

	if pile.meta != nil {
		return fmt.Errorf("already opened")
	}

	if err := os.MkdirAll(pile.path, 0755); err != nil {
		return err
	}

	pathMeta := path.Join(pile.path, "pile.json")

	if _, err := os.Stat(pathMeta); os.IsNotExist(err) {
		pile.meta = &PileMeta{}

		return pile.flush(ctx)
	}

	if file, err := os.OpenFile(pathMeta, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return err
	} else {
		meta := &PileMeta{}
		if err := json.NewDecoder(file).Decode(meta); err != nil {
			return err
		}
		pile.meta = meta
	}

	return nil
}

func (pile *Pile) flush(ctx context.Context) error {
	pathMeta := path.Join(pile.path, "pile.json")

	if file, err := os.OpenFile(pathMeta, os.O_RDWR|os.O_CREATE, 0666); err != nil {
		return err
	} else {
		if err := json.NewEncoder(file).Encode(pile.meta); err != nil {
			return err
		}
	}

	return nil
}

func (pile *Pile) getNeedlePath(key PileKey) string {
	return path.Join(pile.path, string(key))
}

func (pile *Pile) Find(ctx context.Context, key PileKey) (io.Reader, int64, error) {
	pile.lock.RLock()

	defer pile.lock.RUnlock()

	if pile.meta == nil {
		return nil, 0, fmt.Errorf("uninitialized")
	}

	needlePath := pile.getNeedlePath(key)
	if info, err := os.Stat(needlePath); os.IsNotExist(err) {
		pile.log.Infow("miss", "key", key)

		pile.metrics.PileMiss(pile.name)

		return nil, 0, nil
	} else {
		pile.log.Infow("hit", "key", key)

		pile.metrics.PileHit(pile.name)

		opened, err := os.OpenFile(needlePath, os.O_RDONLY, 0)
		if err != nil {
			return nil, 0, err
		}

		return opened, info.Size(), nil
	}
}

func (pile *Pile) AddBytes(ctx context.Context, key PileKey, data []byte) error {
	return pile.Add(ctx, key, bytes.NewReader(data))
}

func (pile *Pile) Add(ctx context.Context, key PileKey, reader io.Reader) error {
	pile.log.Infow("adding", "key", key)

	pile.lock.Lock()

	defer pile.lock.Unlock()

	if pile.meta == nil {
		return fmt.Errorf("uninitialized")
	}

	needlePath := pile.getNeedlePath(key)
	if _, err := os.Stat(needlePath); !os.IsNotExist(err) {
		return fmt.Errorf("key exists")
	}

	writing, err := os.OpenFile(needlePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	copied, err := io.Copy(writing, reader)
	if err != nil {
		return err
	}

	entry := &pileEntry{
		Key:   key,
		Added: time.Now().UTC(),
		Size:  copied,
	}

	/// Were we being deleted while we were waiting?
	if pile.meta == nil {
		return fmt.Errorf("canceled, pile deleted")
	}

	// Check for a duplicate just in case there was a race.
	l := len(pile.meta.Tail)
	if l > 0 && pile.meta.Tail[l-1].Key == key {
		return nil
	}

	// Ok, everything seems fine so append to the tail.
	if len(pile.meta.Tail) == TailLength {
		for i := 0; i < TailLength; i++ {
			pile.meta.Tail[i] = pile.meta.Tail[i+1]
		}
		pile.meta.Tail[TailLength-1] = entry
	} else {
		pile.meta.Tail = append(pile.meta.Tail, entry)
	}

	// Keep total size of the pile correct.
	pile.meta.Size += copied

	// Flush, while lock is held.
	if err := pile.flush(ctx); err != nil {
		return err
	}

	pile.metrics.PileBytes(pile.name, pile.meta.Size)

	pile.log.Infow("added", "pile", pile.path, "key", key, "size", copied)

	return nil
}

func (pile *Pile) Delete(ctx context.Context) error {
	pile.log.Infow("deleting")

	pile.lock.Lock()

	defer pile.lock.Unlock()

	if pile.meta != nil {
		return nil
	}

	if err := os.RemoveAll(pile.path); err != nil {
		return err
	}

	pile.meta = nil

	return nil
}

func (pile *Pile) Close() error {
	pile.lock.Lock()

	defer pile.lock.Unlock()

	pile.meta = nil

	return nil
}
