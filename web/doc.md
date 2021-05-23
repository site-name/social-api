1) After modifying /schemas/*.grapohqls, remember to run `make gqlgen`

2) After running `make gqlgen` <br />
If there is a file named `prelude.resolvers.go`, delete it <br />

**NOTE**: all resolver files must be prefixed with `resolvers_`, followed by model name <br />
*E.g*: `resolver_csv.go`, `resolvers_payment.go`
