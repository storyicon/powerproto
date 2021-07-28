# Greeter

This is the simplest grpc example, which does not reference any external grpc libraries.

The following plug-ins are used:
* [protoc-gen-go](google.golang.org/protobuf/cmd/protoc-gen-go)
* [protoc-gen-go-grpc](google.golang.org/grpc/cmd/protoc-gen-go-grpc)

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
│   ├── greeter.pb.go
│   ├── greeter.proto
│   └── greeter_grpc.pb.go
└── powerproto.yaml

1 directory, 5 files
```