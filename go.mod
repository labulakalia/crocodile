module github.com/labulaka521/crocodile

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/casbin/casbin/v2 v2.1.2
	github.com/casbin/gorm-adapter/v2 v2.0.3
	github.com/denisenkom/go-mssqldb v0.0.0-20191001013358-cfbb681360f0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/go-bindata-assetfs v1.0.0
	github.com/gin-contrib/pprof v1.2.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.7 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.3.3
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/websocket v1.4.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/jinzhu/gorm v1.9.11 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/procfs v0.0.5 // indirect
	github.com/soheilhy/cmux v0.1.4
	github.com/spf13/cobra v0.0.5
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
	go.uber.org/zap v1.12.0
	golang.org/x/crypto v0.0.0-20191108234033-bd318be0434a
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	golang.org/x/tools v0.0.0-20200131000851-b4207ef49307 // indirect
	google.golang.org/genproto v0.0.0-20191115221424-83cc0476cb11 // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tucnak/telebot.v2 v2.0.0-20200120165535-b6c3367fed99
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace gopkg.in/jcmturner/rpc.v1 => gopkg.in/jcmturner/rpc.v1 v1.1.0

// +heroku goVersion go1.13
