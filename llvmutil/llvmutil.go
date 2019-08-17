// Package llvmutil provides several LLVM helper functions used across the
// compiler.
package llvmutil

import (
	"tinygo.org/x/go-llvm"
)

// GetUses returns a list of LLVM values where this value is used as an operand.
// It is legal to pass in a nil value: that will result in no values being returned.
func GetUses(value llvm.Value) []llvm.Value {
	if value.IsNil() {
		return nil
	}
	var uses []llvm.Value
	use := value.FirstUse()
	for !use.IsNil() {
		uses = append(uses, use.User())
		use = use.NextUse()
	}
	return uses
}

// GetZeroValue returns a zero LLVM value for any LLVM type. It is similar to
// the 'zeroinitializer' on a global.
// Unfortunately, I haven't found a way to do it directly with the Go API but
// this works just fine.
func GetZeroValue(typ llvm.Type) llvm.Value {
	switch typ.TypeKind() {
	case llvm.ArrayTypeKind:
		subTyp := typ.ElementType()
		subVal := GetZeroValue(subTyp)
		vals := make([]llvm.Value, typ.ArrayLength())
		for i := range vals {
			vals[i] = subVal
		}
		return llvm.ConstArray(subTyp, vals)
	case llvm.FloatTypeKind, llvm.DoubleTypeKind:
		return llvm.ConstFloat(typ, 0.0)
	case llvm.IntegerTypeKind:
		return llvm.ConstInt(typ, 0, false)
	case llvm.PointerTypeKind:
		return llvm.ConstPointerNull(typ)
	case llvm.StructTypeKind:
		types := typ.StructElementTypes()
		vals := make([]llvm.Value, len(types))
		for i, subTyp := range types {
			vals[i] = GetZeroValue(subTyp)
		}
		if typ.StructName() != "" {
			return llvm.ConstNamedStruct(typ, vals)
		} else {
			return typ.Context().ConstStruct(vals, false)
		}
	default:
		panic("unknown LLVM zero inititializer: " + typ.String())
	}
}
