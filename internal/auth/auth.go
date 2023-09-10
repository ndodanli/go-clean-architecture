package auth

import (
	"fmt"
	"github.com/ndodanli/go-clean-architecture/configs"
	"github.com/ndodanli/go-clean-architecture/pkg/casbin/v2"
	"github.com/ndodanli/go-clean-architecture/pkg/casbin/v2/model"
	"github.com/ndodanli/go-clean-architecture/pkg/xorm-adapter/v2"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Auth struct {
	enforcer *casbin.Enforcer
}

func (a *Auth) Enforcer() *casbin.Enforcer {
	return a.enforcer
}

func NewAuth(cfg *configs.Config) (*Auth, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=auth sslmode=disable", cfg.Postgresql.HOST, cfg.Postgresql.PORT, cfg.Postgresql.USER, cfg.Postgresql.PASS)

	a, err := xormadapter.NewAdapter("postgres", connStr)
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}
	// Read the rbac_model.conf file and add it to the model
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("error: currentPath: %s", err)
		return nil, err
	}

	m, err := model.NewModelFromFile(currentPath + "/internal/auth/rbac_model.conf")

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		log.Fatalf("error: enforcer: %s", err)
		return nil, err
	}
	//
	e.EnableAutoSave(true)
	var policyErr error
	var ok bool
	ok, policyErr = e.AddNamedPoliciesEx("p", [][]string{
		{"user1", "data1", "read", "tenant1"}, {"user2", "data2", "read", "tenant1"},
	})

	ok, policyErr = e.AddNamedPoliciesEx("p", [][]string{
		{"user1", "data1", "read", "tenant1"}, {"user2", "data2", "read", "tenant1"},
	})
	print(ok)

	users := e.GetAllUsersByDomain("tenant1")

	print(users)

	if policyErr != nil {
		log.Fatalf("error: add policy: %s", err)
		return nil, err
	}

	return &Auth{
		enforcer: e,
	}, nil
}
