package lifetime

import (
	"github.com/ndodanli/backend-api/pkg/infrastructure/db/sqldb/pg"
	"github.com/ndodanli/backend-api/pkg/infrastructure/services"
)

var (
	TxSessionManagerType = &pg.TxSessionManager{}
	AppServicesType      = &services.AppServices{}
)
