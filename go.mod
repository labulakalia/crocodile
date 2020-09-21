module github.com/labulaka521/crocodile

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/casbin/casbin/v2 v2.11.3
	github.com/casbin/gorm-adapter/v2 v2.1.0
	github.com/coreos/go-etcd v2.0.0+incompatible // indirect
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20200910202707-1e08a3fab204 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/go-openapi/spec v0.19.9 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v7 v7.4.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/lib/pq v1.8.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/micro/go-micro v1.18.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/soheilhy/cmux v0.1.4
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.7
	github.com/ugorji/go v1.1.8 // indirect
	github.com/urfave/cli v1.22.4 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	golang.org/x/sys v0.0.0-20200917073148-efd3b9a0ff20 // indirect
	golang.org/x/text v0.3.3
	golang.org/x/tools v0.0.0-20200917221617-d56e4e40bc9d // indirect
	google.golang.org/genproto v0.0.0-20200917134801-bb4cff56e0d0 // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/tucnak/telebot.v2 v2.3.4
	gopkg.in/yaml.v2 v2.3.0 // indirect
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.20.1
)

replace gopkg.in/jcmturner/rpc.v1 => gopkg.in/jcmturner/rpc.v1 v1.1.0

// +heroku goVersion go1.13
