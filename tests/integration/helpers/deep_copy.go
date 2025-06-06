package helpers

import (
	"reflect"

	didcli "github.com/pointidentity/pointidentity-node/x/did/client/cli"
	didtypes "github.com/pointidentity/pointidentity-node/x/did/types"
	resourcetypes "github.com/pointidentity/pointidentity-node/x/resource/types"
)

// IDeepCopy is an interface for deep copy in the decorator pattern.
type IDeepCopy interface {
	DeepCopy(src interface{}) interface{}
}

// TDeepCopy is a decorator for deep copy.
type TDeepCopy struct{}

func (d *TDeepCopy) DeepCopy(src interface{}) interface{} {
	return deepCopy(src)
}

// DeepCopyCreateDid is a decorator for deep copy of type MsgCreateDidPayload.
type DeepCopyCreateDid struct {
	TDeepCopy
}

func (d *DeepCopyCreateDid) DeepCopy(src didtypes.MsgCreateDidDocPayload) didtypes.MsgCreateDidDocPayload {
	return deepCopy(src).(didtypes.MsgCreateDidDocPayload)
}

// DeepCopyUpdateDid is a decorator for deep copy of type MsgUpdateDidPayload.
type DeepCopyUpdateDid struct {
	TDeepCopy
}

func (d *DeepCopyUpdateDid) DeepCopy(src didtypes.MsgUpdateDidDocPayload) didtypes.MsgUpdateDidDocPayload {
	return deepCopy(src).(didtypes.MsgUpdateDidDocPayload)
}

type DeepCopyDIDDocument struct {
	TDeepCopy
}

func (d *DeepCopyDIDDocument) DeepCopy(src didcli.DIDDocument) didcli.DIDDocument {
	return deepCopy(src).(didcli.DIDDocument)
}

// DeepCopyCreateResource is a decorator for deep copy of type MsgCreateResource.
type DeepCopyCreateResource struct {
	TDeepCopy
}

func (d *DeepCopyCreateResource) DeepCopy(src resourcetypes.MsgCreateResource) resourcetypes.MsgCreateResource {
	return deepCopy(src).(resourcetypes.MsgCreateResource)
}

func deepCopy(src interface{}) interface{} {
	var reflection interface{}
	var dst reflect.Value

	switch actualSrc := (src).(type) {
	case didtypes.MsgCreateDidDocPayload:
		// Create a reflection slice of the same length as the source slice
		reflection = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(actualSrc)), 1, 1).Interface()
		// Extract destination value as definition
		dst = reflect.ValueOf(reflection)
		// Define source value as slice
		slc := []didtypes.MsgCreateDidDocPayload{actualSrc}
		// Copy the source value into the destination
		reflect.Copy(dst, reflect.ValueOf(slc))
		// Return the destination value from the reflection slice
		return dst.Index(0).Interface().(didtypes.MsgCreateDidDocPayload)
	case didtypes.MsgUpdateDidDocPayload:
		// Create a reflection slice of the same length as the source slice
		reflection = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(actualSrc)), 1, 1).Interface()
		// Extract destination value as definition
		dst = reflect.ValueOf(reflection)
		// Define source value as slice
		slc := []didtypes.MsgUpdateDidDocPayload{actualSrc}
		// Copy the source value into the destination
		reflect.Copy(dst, reflect.ValueOf(slc))
		// Return the destination value from the reflection slice
		return dst.Index(0).Interface().(didtypes.MsgUpdateDidDocPayload)
	case resourcetypes.MsgCreateResource:
		// Create a reflection slice of the same length as the source slice
		reflection = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(actualSrc)), 1, 1).Interface()
		// Extract destination value as definition
		dst = reflect.ValueOf(reflection)
		// Define source value as slice
		slc := []resourcetypes.MsgCreateResource{actualSrc}
		// Copy the source value into the destination
		reflect.Copy(dst, reflect.ValueOf(slc))
		// Return the destination value from the reflection slice
		return dst.Index(0).Interface().(resourcetypes.MsgCreateResource)
	case didcli.DIDDocument:
		// Create a reflection slice of the same length as the source slice
		reflection = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(actualSrc)), 1, 1).Interface()
		// Extract destination value as definition
		dst = reflect.ValueOf(reflection)
		// Define source value as slice
		slc := []didcli.DIDDocument{actualSrc}
		// Copy the source value into the destination
		reflect.Copy(dst, reflect.ValueOf(slc))
		// Return the destination value from the reflection slice
		return dst.Index(0).Interface().(didcli.DIDDocument)
	default:
		panic("Unsupported type")
	}
}
