# protodep

For a given proto file, report all the required filenames recursively.

    $ protodep v1/myservice.proto

    v1/myservice.proto
    v1/shared.proto
    google/api/annotations.proto
    google/api/http.proto
    google/protobuf/descriptor.proto
    google/protobuf/timestamp.proto