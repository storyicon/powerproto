# Using GoGoFast

This example uses the following public libraries:
* [googleapis](https://github.com/googleapis/googleapis)
* [gogoproto](https://github.com/gogo/protobuf/tree/master/gogoproto)

The following plug-ins are used:
* [protoc-gen-gofast](https://github.com/gogo/protobuf/protoc-gen-gofast)
* [protoc-gen-grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

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
│   └── service.proto
└── powerproto.yaml

1 directory, 5 files
```