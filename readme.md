# terriyaki - ⚠️🚧 WIP 🚧⚠️

Terriyaki combines the power of Terraform with the Developer Experience of Go.

## TODO: Documentation

1. Overview & diagram of components of terriyaki
2. How to create a RootModule
   3. Create your struct
   4. Embed `tki.RootModule`
   5. Create your backend struct (implement the `tki.Backend` interface)
   6. Add your provider(s) to your root module
   7. Add/embed your resources:
      8. Embedding a struct of resources is akin to calling a Terraform module, in a way...
   9. Nothing prevents you from defining your own approach. But nested structs are not supported.
8. 
2. Type system: `StringValue`, `List[StringValue]`, etc. Actual value (cty.Value) vs reference. Traversals. 
3. Missing functionality: provider aliases, `for_each`, terraform functions

## TODO: bugs

1. [SOLVED] `aws_route_table` in-line `route{}` block is empty in generated code... Uses attributes-as-block: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route_table#route

## IDEAS:

### Error Handling

Create an error type for HCL encoding, that has source location.

### GOJEN

1. Make resource/data interface a struct
```go
type IamRole struct {
	Name  string `validate:"required"`
	Args  IamRoleArgs // TODO: how to hcl encode this??
	state *IamRoleOut
}
//...
func (i *IamRole) State() (*IamRoleState, bool) {
return i.state, i.state != nil
}

```

### tki exported things

Move everything that doesn't need to be in `tki` to `core`.
Update code generation also.
Rename `Terriyaki*` functions to `Tki*` or `Internal*`