module crocodile

require (
	github.com/SAP/go-hdb v0.14.1 // indirect
	github.com/StackExchange/wmi v0.0.0-20181212234831-e0a55b97c705 // indirect
	github.com/coredns/coredns v1.4.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/go-control-plane v0.6.9 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/hashicorp/consul v1.5.1 // indirect
	github.com/hashicorp/go-gcp-common v0.5.0 // indirect
	github.com/hashicorp/go-memdb v1.0.0 // indirect
	github.com/hashicorp/go-plugin v1.0.0 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/raft-boltdb v0.0.0-20171010151810-6e5ba93211ea // indirect
	github.com/hashicorp/vault v1.1.0 // indirect
	github.com/hashicorp/vault-plugin-auth-alicloud v0.0.0-20190320211238-36e70c54375f // indirect
	github.com/hashicorp/vault-plugin-auth-azure v0.0.0-20190320211138-f34b96803f04 // indirect
	github.com/hashicorp/vault-plugin-auth-centrify v0.0.0-20190320211357-44eb061bdfd8 // indirect
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.0.0-20190328163920-79931ee7aad5 // indirect
	github.com/hashicorp/vault-plugin-secrets-ad v0.0.0-20190327182327-ed2c3d4c3d95 // indirect
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.0.0-20190320213517-3307bdf683cb // indirect
	github.com/hashicorp/vault-plugin-secrets-azure v0.0.0-20190320211922-2dc8a8a5e490 // indirect
	github.com/hashicorp/vault-plugin-secrets-gcp v0.0.0-20190320211452-71903323ecb4 // indirect
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.0.0-20190320213325-9e326a9e802d // indirect
	github.com/influxdata/influxdb v1.7.5 // indirect
	github.com/labulaka521/logging v0.0.0-20190526092138-ef7e2414b576
	github.com/lyft/protoc-gen-validate v0.0.14 // indirect
	github.com/micro/go-micro v1.5.0
	github.com/micro/go-plugins v1.1.0
	github.com/micro/micro v1.5.0 // indirect
	github.com/micro/protoc-gen-micro v0.8.0 // indirect
	github.com/micro/util v0.2.0
	github.com/protocolbuffers/protobuf v3.8.0+incompatible // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible // indirect
	github.com/ugorji/go/codec v0.0.0-20190320090025-2dc34c0b8780 // indirect
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5
	golang.org/x/tools v0.0.0-20190606050223-4d9ae51c2468 // indirect
	google.golang.org/genproto v0.0.0-20190605220351-eb0b1bdb6ae6 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.0
	gotest.tools v2.3.0+incompatible // indirect
	honnef.co/go/tools v0.0.0-20190605142022-0a11fc526260 // indirect
	layeh.com/radius v0.0.0-20190322222518-890bc1058917 // indirect
)

replace (
	github.com/golang/lint => golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
	github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.2
	github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
)
