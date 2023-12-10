package lifetime

import (
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/db/sqldb/postgresql"
	"github.com/ndodanli/go-clean-architecture/pkg/infrastructure/services"
)

var (
	TxSessionManagerType = &postgresql.TxSessionManager{}
	AppServicesType      = &services.AppServices{}
)
