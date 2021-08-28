# Using Multi Config Items

This example uses two different configurations for "using-gogo" and "using-googleapis" through the "scope" field.

It uses the following public libraries:
* [googleapis](https://github.com/googleapis/googleapis)

The following plug-ins are used:
* [protoc-gen-go](https://google.golang.org/protobuf/cmd/protoc-gen-go)
* [protoc-gen-go-grpc](https://google.golang.org/grpc/cmd/protoc-gen-go-grpc)
* [protoc-gen-grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway)
* [protoc-gen-gogo](https://github.com/gogo/protobuf/protoc-gen-gogo)


You can compile the proto file in this directory by executing the following command:
```
powerproto build -r .
```

Not surprisingly, you will get the following output:
```
➜ tree
.
├── README.md
├── powerproto.yaml
├── using-gogo
│   ├── service.pb.go
│   ├── service.pb.gw.go
│   └── service.proto
└── using-googleapis
    ├── service.pb.go
    ├── service.pb.gw.go
    ├── service.proto
    └── service_grpc.pb.go

2 directories, 9 files
```
