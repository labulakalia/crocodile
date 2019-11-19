proto:
	protoc --go_out=plugins=grpc:. core/proto/core.proto

tlskey:
	# 生成服务端私钥
	openssl ecparam -genkey -name secp384r1 -out core/tls/server.key
	# 生成自签公钥
	openssl req -new -x509 -sha256 -key core/tls/server.key -out core/tls/server.pem -days 3650