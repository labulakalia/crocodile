package rbac

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestRabc(t *testing.T) {
	//// https://casbin.org/docs/zh-CN/supported-models
	//e := casbin.NewEnforcer("/Users/labulakalia/workerspace/golang/crocodile/rbac/conf/rbac.conf", "/Users/labulakalia/workerspace/golang/crocodile/rbac/conf/rbac.csv")
	//t.Log(e.Enforce("superAdmin2", "project", "read"))
	//t.Log(e.Enforce("quyuan", "asse", "read"))
	//t.Log(e.Enforce("wenyin", "asse", "read"))
	//t.Log(e.Enforce("wenyin", "project", "read"))
	//t.Log(e.GetAllActions())
	//t.Log(e.GetAllSubjects())
	//t.Log(e.GetAllRoles())
	//t.Log(e.GetAllObjects())
	//
	//m := model.Model{}
	//m.AddDef("r", "r", "sub, obj, act")
	//m.AddDef("p", "p", "sub, obj, act")
	//m.AddDef("g", "g", "_,_")
	//m.AddDef("e", "e", "some(where (p.eft == allow))")
	//m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")
	//
	//a := fileadapter.NewFilteredAdapter("/Users/labulakalia/workerspace/golang/crocodile/rbac/conf/rbac.csv")
	//
	//e = casbin.NewEnforcer(m, a)
	//t.Log(e.Enforce("superAdmin2", "project", "read"))
	//t.Log(e.Enforce("quyuan", "asse", "read"))
	//t.Log(e.Enforce("wenyin", "asse", "read"))
	//t.Log(e.Enforce("wenyin", "project", "read"))
	//
	n := 0x1234
	res := *(*byte)(unsafe.Pointer(&n))
	fmt.Println(res == 0x34)
	t.Log()
}
