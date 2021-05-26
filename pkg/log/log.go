package log

import (
	"flag"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Options is a simple wrapper around the controller runtime logging.
type Options struct {
	zap.Options
}

// AddFlagSet binds options to the flagset.
func (o *Options) AddFlagSet(f *flag.FlagSet) {
	o.Options.BindFlags(f)
}

// New creates a new Zap logr.
func New(o *Options) logr.Logger {
	return zap.New(zap.UseFlagOptions(&o.Options))
}
