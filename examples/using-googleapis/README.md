# Using Googleapis

This example uses the following public libraries:
* [googleapis](https://github.com/googleapis/googleapis)

The following plug-ins are used:
* [protoc-gen-go](https://google.golang.org/protobuf/cmd/protoc-gen-go)
* [protoc-gen-go-grpc](https://google.golang.org/grpc/cmd/protoc-gen-go-grpc)
* [protoc-gen-grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway)

You can compile the proto file in this directory by executing the following command:
```
powerproto build -r ./apis
```

Not surprisingly, you will get the following output:
```
➜ tree
.
├── README.md
├── apis
│   ├── service.pb.go
│   ├── service.pb.gw.go
│   ├── service.proto
│   └── service_grpc.pb.go
└── powerproto.yaml

1 directory, 6 files
```
