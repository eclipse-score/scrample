package config

import (
    "scorex/internal/model"
    "scorex/internal/service/knowngood"
)

// Deprecated: use scorex/internal/model.KnownGood and scorex/internal/service/knowngood.Load.
type KnownGood = model.KnownGood

// Deprecated: use scorex/internal/service/knowngood.Load.
func LoadKnownGood(urlOrPath string) (*KnownGood, error) {
    return knowngood.Load(urlOrPath)
}
