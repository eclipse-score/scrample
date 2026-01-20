package skeleton

import serviceskeleton "scorex/internal/service/skeleton"

// Deprecated: use scorex/internal/service/skeleton.Generate.
func Generate(props Properties) error {
    return serviceskeleton.Generate(props)
}
