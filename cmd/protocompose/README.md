# protocompose

This tool can process `@compose` "annotations" in comments of Protobuf message definitions.

## how it works

Given the source file containing:
```
message Source {
    // the identifier
    string id = 1;    

    // date of issue
    google.type.Date issue_date = 2;

    // unused
    string external_code = 3;
}

// Dimension is for objects that have a size: W x H
message Dimension {
  // width in pixels
  int32 width = 1;

  // height in pixels
  int32 height = 2;
}

// FileReference represents a store media object
message FileReference {
  string file_name = 1;
  string mime_type = 2;
}

```

and the any target file containing an `@compose` annotation:

```
// @compose somepackage.v2.Source.id
// @compose somepackage.v2.Source.issue_date
// @compose ...somepackage.v2.Dimension
// @compose #somepackage.v2.FileReference
message Composed {
}
```

after processing all proto files with `protocompose`, you will get:

```
// @compose somepackage.v2.Source.id
// @compose somepackage.v2.Source.issue_date
message Composed {
  
  // the identifier
  string id = 1;
  
  // date of issue
  google.type.Date issue_date = 2;

  // width in pixels
  int32 width = 3;

  // height in pixels
  int32 height = 4;

  // FileReference represents a store media object
  somepackage.v2.FileReference filereference = 5;
}
```
which contains copies of the fields as specified by each annotation.

Existing fields in a composed message are removed.
Field numbers start at 1 and will follow the order as described.