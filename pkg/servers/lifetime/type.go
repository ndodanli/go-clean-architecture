package lifetime

import (
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/pg"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
)

var (
	TxSessionManagerType = &pg.TxSessionManager{}
	AppServicesType      = &services.AppServices{}
)
