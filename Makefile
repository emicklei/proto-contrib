install:
	cd cmd/proto2avro && go install
	cd cmd/proto2gql && go install
	cd cmd/proto2openapi && go install
	cd cmd/proto2xsd && go install
	cd cmd/protocompose && go install
	cd cmd/protodep && go install
	cd cmd/protofmt && go install