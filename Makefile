user:
	micro new --namespace=crocodile --type web --alias=user github.com/labulaka521/crocodile/web/user
	micro new --namespace=crocodile --type srv --alias=user github.com/labulaka521/crocodile/service/user

run_api:
	micro --api_namespace=crocodile.web --registry=etcdv3 api --handler=web

build_all:
	cd service/actuator&&make build

