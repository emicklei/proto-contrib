# protodecode

Protodecode can decode a ProtocolBuffer marshaled message into a map[string]interface{} using its .proto definition.
This package was created to transform such messages directly to JSON without using the generated (Go) code to do the unmarshaling and marshaling. It can be used to inspect messages for quick view and debugging. Do not use this if performance is of importance.