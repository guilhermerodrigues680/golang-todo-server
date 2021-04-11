package storage

import "errors"

// ErrNotFound é o erro returnado pelo storage quando uma busca no storage não retorna resultados.
var ErrNotFound = errors.New("no items found")
