protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --validate_out="lang=go,paths=source_relative:." *.proto
protoc --plugin=protoc-gen-grpc-java --grpc-java_out=. --java_out=. --validate_out="lang=java:."  *.proto
