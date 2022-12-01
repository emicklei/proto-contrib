# protocompose

This tool can process `@compose` "annotations" in comments of Protobuf message definitions.

## how it works

Given the source file containing:
```
message Source {
    // the identifier
    string id = 10;    

    // date of issue
    google.type.Date issue_date = 11;

    // unused
    string external_code = 12;
}

// Dimension is for objects that have a size: W x H
message Dimension {
  // width in pixels
  int32 width = 10;

  // height in pixels
  int32 height = 20;
}
```

and the any target file containing `@compose` annotations:

```
// @compose somepackage.v2.Source.id
// @compose somepackage.v2.Source.issue_date
// @compose somepackage.v2.Dimension.width
// @compose somepackage.v2.Dimension.height
message Composed {
  // must be empty
}
```

after processing all proto files with `protocompose`, you will get:

```
// @compose somepackage.v2.Source.id
// @compose somepackage.v2.Source.issue_date
// @compose somepackage.v2.Dimension.width
// @compose somepackage.v2.Dimension.height
message Composed {
  
  // the identifier
  string id = 1;
  
  // date of issue
  google.type.Date issue_date = 2;

  // width in pixels
  int32 width = 3;

  // height in pixels
  int32 height = 4;
}
```
which contains copies of the fields as specified by each annotation.

Field numbers follow the order as described.

Supported fields: Normal and Map.