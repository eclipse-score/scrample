package projectinit

import serviceprojectinit "scorex/internal/service/projectinit"

// Deprecated: use scorex/internal/service/projectinit.Options.
type Options = serviceprojectinit.Options

// Deprecated: use scorex/internal/service/projectinit.Result.
type Result = serviceprojectinit.Result

// Deprecated: use scorex/internal/service/projectinit.Run.
func Run(opts Options) (*Result, error) {
    return serviceprojectinit.Run(opts)
}
