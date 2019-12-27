package sync

import (
	"context"

	"github.com/nono/cozy-fuse/config"
	"github.com/nono/cozy-fuse/local"
	"github.com/nono/cozy-fuse/remote"
)

// Run start running the synchronization between the local file system and the
// remote cozy instance.
func Run(ctx context.Context, cfg *config.Config) error {
	remote.Do(cfg)
	return local.Mount(ctx, cfg)
}
