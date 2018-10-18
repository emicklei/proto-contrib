# protodep

For a given proto file, report all the required filenames recursively.

    $ protodep v1/myservice.proto

    v1/myservice.proto
    v1/shared.proto
    google/api/annotations.proto
    google/api/http.proto
    google/protobuf/descriptor.proto
    google/protobuf/timestamp.proto

It is also possible to add multiple proto files

    $ protodep v1/myservice.proto v1/myotherservice.proto

    v1/myservice.proto
    v1/shared.proto
    google/api/annotations.proto
    google/api/http.proto
    google/protobuf/descriptor.proto
    google/protobuf/timestamp.proto
    v1/myotherservice.proto
    v1/othershared.proto

Use the format flag to get JSON output

    $ protodep -format json v1/myservice.proto v1/myotherservice.proto

    ["v1/myservice.proto", "v1/shared.proto", ...]