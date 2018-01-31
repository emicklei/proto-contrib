# proto2gql

GraphQL Schema conversion tool for Google Protocol Buffers version 3

	> proto2gql -help
	    Usage of proto2gql [flags] [path ...]

        -filter string
            Regexp to filter out matched types
        -filterN string
            Regexp to filter out not matched types
        -go_out string
            Writes transformed files to .go file
        -js_out string
            Writes transformed files to .js file
        -no_prefix
            Disables package prefix for type names
        -package_alias value
            Renames packages using given aliases
        -resolve_import value
            Resolves given external packages
        -std_out
            Writes transformed files to stdout
        -txt_out string
            Writes transformed files to .graphql file

### build
	make