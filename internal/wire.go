package internal

import "reflect"

type Wire struct{}

var WireType = reflect.TypeOf(Wire{})
