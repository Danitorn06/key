package auth

import (
    "github.com/casbin/casbin/v2"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
    e, err := casbin.NewEnforcer("model.conf", "policy.csv")
    if err != nil {
        panic("failed to initialize Casbin: " + err.Error())
    }
    Enforcer = e
}