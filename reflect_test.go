package ref

import (
  "time"
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/google/uuid"
  "flag"
)

type FAddress struct {
  Country        string     `json:"country"          yaml:"country"`
  Index          string     `json:"index"            yaml:"index"`
  City           string     `json:"city"             yaml:"city"`
  Street         string     `json:"street"           yaml:"street"`
  House          string     `json:"house"            yaml:"house"`
  Room           string     `json:"room"             yaml:"room"`
}

type FBankAccount struct {
  BIK                    string     `json:"bik"                         yaml:"bik"`
  BankName               string     `json:"bank_name"                   yaml:"bank_name"`
  CorrespondentAccount   string     `json:"correspondent_account"       yaml:"correspondent_account"`
  Account                string     `json:"account"                     yaml:"account"`
  Currency_CODE          string     `json:"currency_code"               yaml:"currency_code"`
  Default                bool       `json:"default_account"             yaml:"default_account"`
}

type FBankAccounts []FBankAccount

////////////////////////////////
// Organization
///////////////////////////////

type Organization struct {
  ID             uuid.UUID     `db:"id"                         json:"id"            yaml:"id"               gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()"`
  CreatedAt      time.Time     `db:"created_at;default: now()"  json:"created_at"    sql:"default: now()"    gorm:"type:timestamp with time zone"`
  UpdatedAt      time.Time     `db:"updated_at;default: null"   json:"updated_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`
  DeletedAt     *time.Time     `db:"deleted_at;default: null"   json:"deleted_at"    sql:"default: null"     gorm:"type:timestamp with time zone"`

  CODE           string    `db:"code"         json:"code"          yaml:"code"             gorm:"type:varchar(96);default: null"`

  Name           string    `db:"name"         json:"name"          yaml:"name"             sql:"column:name"        gorm:"column:name;type:varchar(256)"`
  Description    string    `db:"description"  json:"description"   yaml:"description"      gorm:"default: null"`
  
  Bank           FBankAccounts      `db:"bank"           json:"bank,ommitempty"                 gorm:"type:jsonb;"`
  AddressLegal      FAddress         `db:"address_legal"             json:"address_legal,omitempty"                 gorm:"type:jsonb;"`
  AddressBilling    FAddress         `db:"address_billing"           json:"address_billing,omitempty"               gorm:"type:jsonb;"`
  AddressShipping   FAddress         `db:"address_shipping"          json:"address_shipping,omitempty"              gorm:"type:jsonb;"`
}

////////////////////////////////
// User
///////////////////////////////

type User struct {
  ID            uuid.UUID       `json:"ID"`
  Login         string          `json:"login"`
  EMail         string          `json:"email"`
  DisplayName   string          `json:"display_name"`
  Avatar        string          `json:"avatar"`
  Language      string          `json:"lang"`
  Group         string          `json:"group"`
  Groups      []string          `json:"groups"`
  TimeLogin     time.Time       `json:"-"`
  AuthCode      string          `json:"-"`
  Disable       bool            `json:"disable"`
}

/////////////////////////
// TESTS
/////////////////////////
func TestReflect(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  var adr FAddress
  adr.Country = "Russia"
  adr.Index   = "127888"
  adr.City    = "Moscow"
  
  uid1, _ := uuid.Parse("00000002-0003-0004-0005-000000000004")
  
  m_need := map[string]interface{}{"city":"Moscow", "country":"Russia", "index":"127888"}
  m := ConvertToMap(adr)
  assert.Equal(t, m_need, m)


  var org1 Organization
  org1.ID = uid1
  org1.Name = "OOO `Org`"
  org1.AddressLegal.Country = "Russia"
  org1.AddressLegal.Index   = "127888"
  org1.AddressLegal.City    = "Moscow"
  org1.Bank = make(FBankAccounts, 2)
  org1.Bank[0].BIK     = "1111111"
  org1.Bank[0].Account = "111111134583459834573279"
  org1.Bank[1].BIK     = "21111111"
  org1.Bank[1].Account = "2111111134583459834573279"
  
  o1_need := map[string]interface{}{"id": uid1, "address_legal.city":"Moscow", "address_legal.country":"Russia", "address_legal.index":"127888", "bank:0.account":"111111134583459834573279", "bank:0.bik":"1111111", "bank:1.account":"2111111134583459834573279", "bank:1.bik":"21111111", "name":"OOO `Org`"}
  o1 := ConvertToMap(org1)
  assert.Equal(t, o1_need, o1)

  org2 := Organization{}
  
  ConvertFromMap(&org2, &o1)
  assert.Equal(t, org1, org2)
  
  
  user1 := User{Login: "user", EMail: "user@mail.user", Group: "user_group", Groups: []string{"mail_user_group", "storage_user_group"}}

  u1_need := map[string]interface{}{"email":"user@mail.user", "group":"user_group", "groups:0":"mail_user_group", "groups:1":"storage_user_group", "login":"user"}
  u1 := ConvertToMap(user1)
  assert.Equal(t, u1_need, u1)
  
  assert.Equal(t, true, FieldExists(user1, "EMail"))
  assert.Equal(t, false, FieldExists(user1, "e-mail"))
}

