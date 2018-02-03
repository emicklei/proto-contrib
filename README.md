# proto-contrib

Packages and tools on top of https://github.com/emicklei/proto

### usage of protofmt command

	> protofmt -help
		Usage of protofmt [flags] [path ...]
  		-w	write result to (source) files instead of stdout

See folder `cmd/protofmt/README.md` for more details.

### usage of proto2xsd command

	> proto2xsd -help
		Usage of proto2xsd [flags] [path ...]
  		-ns string
    		namespace of the target types (default "http://your.company.com/domain/version")

See folder `cmd/proto2xsd/README.md` for more details.

### usage of proto2gql command

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

See folder `cmd/proto2gql/README.md` for more details.