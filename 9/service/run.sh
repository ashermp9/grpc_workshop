

mkdir -p google/api
cd google/api
wget https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto
wget https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto

cd ../..
ls

protoc -I. -Igoogle/api --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
       greeting_service.proto