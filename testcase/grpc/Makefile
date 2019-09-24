search:
	protoc --go_out=plugins=grpc:. proto/search.proto

stream:
	protoc --go_out=plugins=grpc:. proto/stream.proto

# 不使用ca证书验证双方
key:
	# 生成服务端私钥
	openssl ecparam -genkey -name secp384r1 -out conf/server.key
	# 生成自签公钥
	openssl req -new -x509 -sha256 -key conf/server.key -out server.pem -days 3650

# 使用ca来验证双方key
cakey:
	# 生成key
	openssl genrsa -out conf1/ca.key 1024
	# 生成密钥
	openssl req -new -x509 -days 7200 -key conf1/ca.key -out conf1/ca.pem

caserver:
	# 生成服务端私钥
	openssl ecparam -genkey -name secp384r1 -out conf1/server/server.key
	# 生成CSR CSR 是 Cerificate Signing Request 的英文缩写，为证书请求文件。主要作用是 CA 会利用 CSR 文件进行签名使得攻击者无法伪装或篡改原有证书
	openssl req -new -key conf1/server/server.key -out conf1/server/server.csr
	# 基于CA签发
	openssl x509 -req -sha256 -CA conf1/ca.pem -CAkey conf1/ca.key -CAcreateserial -days 3650 -in conf1/server/server.csr -out conf1/server/server.pem

caclient:
	# 生成key
	openssl ecparam -genkey -name secp384r1 -out conf1/client/client.key
	# 生成csr
	openssl req -new -key conf1/client/client.key -out conf1/client/client.csr
	#
	openssl x509 -req -sha256 -CA conf1/ca.pem -CAkey conf1/ca.key -CAcreateserial -days 3650 -in conf1/client/client.csr -out conf1/client/client.pem