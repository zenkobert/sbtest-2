protoc delivery/grpc/movie/movie.proto \
--go_out=. \
--go-grpc_out=require_unimplemented_servers=false:. 