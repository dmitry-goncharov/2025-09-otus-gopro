package storagefactory

import (
	"fmt"

	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/app"
	"github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/config"
	memorystorage "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/dmitry-goncharov/2025-09-otus-gopro/hw12_13_14_15_calendar/internal/storage/sql"
)

const (
	IMStorage = "im"
	DBStorage = "db"
)

func NewStorage(storageConf *config.StorageConf) (app.Storage, error) {
	switch storageConf.Type {
	case IMStorage:
		return memorystorage.New(), nil
	case DBStorage:
		return sqlstorage.New(storageConf.DB.Dsn), nil
	default:
		return nil, fmt.Errorf("illegal storage type: %s", storageConf.Type)
	}
}
