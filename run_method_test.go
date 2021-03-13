package ref

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "flag"
  "reflect"
)

func (u *User) GetLogin() string {
  return u.Login
}

/////////////////////////
// TESTS
/////////////////////////
func TestRunMethod(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  var user User
  user.Login = "UserLogin"

  res, ok := RunMethodIfExists(&user, "GetLogin111")
  assert.Equal(t, false, ok)
  assert.Equal(t, []reflect.Value{}, res)

  res, ok = RunMethodIfExists(&user, "GetLogin")
  str, oks := res[0].Interface().(string)
  assert.Equal(t, true, ok)
  assert.Equal(t, true, oks)
  assert.Equal(t, "UserLogin", str)
}

