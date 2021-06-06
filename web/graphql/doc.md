1) After modifying `/schemas/*.graphqls`, remember to run `make gqlgen` for new implemented code to be generated.

2) After running `make gqlgen`, if there is a file named `prelude.resolvers.go`.

**NOTE**: all resolver files **MUST BE** prefixed with `resolvers_`, followed by model name. <br />
*E.g*: `resolver_csv.go`, `resolvers_payment.go`
