protoc delivery/grpc/movie.proto \
--go_out=. \
--go-grpc_out=require_unimplemented_servers=false:. \
--grpc-gateway_out=logtostderr=true:.