package backend

import (
	"context"
	"fmt"

	// "github.com/bgentry/que-go"
	"github.com/govau/que-go"
)

func exampleJob(ctx context.Context, j *que.Job) error {
	fmt.Printf("Example\n")
	return nil
}
