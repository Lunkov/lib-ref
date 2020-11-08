package ref

import (
  "reflect"
  "github.com/golang/glog"
)

func RunMethodIfExists(nameInterface string, any interface{}, nameFunc string, args ...interface{}) ([]reflect.Value, bool) {
  v := reflect.ValueOf(any)
	method := v.MethodByName(nameFunc)
	if method.Kind() == reflect.Invalid {
    glog.Warningf("WRN: NOT FOUND runMethodIfExists(%s.%s)\n", nameInterface, nameFunc)
		return []reflect.Value{}, false
	}

	if method.Type().NumIn() != len(args) {
    glog.Errorf("ERR: runMethodIfExists(%s.%s): expected %d args, actually %d.\n", 
      nameInterface,
			nameFunc,
			len(args),
			method.Type().NumIn())
    return []reflect.Value{}, false
	}

	// Create a slice of reflect.Values to pass to the method. Simultaneously
	// check types.
	argVals := make([]reflect.Value, len(args))
	for i, arg := range args {
		argVal := reflect.ValueOf(arg)

		if argVal.Type() != method.Type().In(i) {
      glog.Errorf("ERR: runMethodIfExists(%s): expected arg %d to have type %v.\n", 
        nameFunc,
				i,
				argVal.Type())
		}

		argVals[i] = argVal
	}
  if glog.V(9) {
    glog.Infof("DBG: Call(%s.%s)\n", nameInterface, nameFunc)
  }
	return method.Call(argVals), true
}
