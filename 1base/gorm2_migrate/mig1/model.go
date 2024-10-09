package mig1

import (
	"database/sql/driver"
	"encoding/json"
)

type Cafeteria struct {
	ID          int64              `db:"id"`
	SnowflakeID int64              `db:"snowflake_id" gorm:"index"`
	Name        string             `gorm:"size:32"`
	Settings    *CafeteriaSettings `gorm:"type:jsonb"`
	CityID      int64              `db:"cityID"`
	Subtitle    string             `gorm:"type:varchar(64)"`
}

func (c Cafeteria) TableName() string {
	return "my_cafeterias"
}

type CafeteriaSettings struct {
	Status  bool
	Address string `gorm:"size:128"`
	Code    string `json:"code,omitempty"`
}

type CafeteriaWithRestaurants struct {
	Cafeteria
	Restaurants []*SimpleRestaurant `gorm:"foreignKey:CafeteriaID;references:ID"`
}

func (c CafeteriaSettings) Value() (driver.Value, error) {
	return json.Marshal(c)
}
func (c *CafeteriaSettings) Scan(src interface{}) error {
	if err := json.Unmarshal(src.([]byte), c); err != nil {
		return err
	}
	return nil
}

type SimpleRestaurant struct {
	ID          int64              `json:"id" db:"id"`
	SnowflakeID int64              `json:"snowflakeID,omitempty" db:"snowflake_id"`
	RefID       int64              `json:"refID" db:"ref_id"`
	CafeteriaID int64              `json:"cafeteriaID" db:"cafeteria_id" gorm:"not null;"`
	Name        string             `json:"name" db:"name" gorm:"type:varchar(100);not null"`
	UUID        string             `gorm:"type:uuid;not null;default:gen_random_uuid();index"`
	Settings    RestaurantSettings `json:"settings,omitempty" db:"settings" gorm:"type:jsonb;not null"`
	Style       string             `json:"style" db:"style" gorm:"size:10"`
}

type RestaurantSettings struct {
	Code         string                 `json:"code,omitempty"`
	Keyword      []string               `json:"keyword,omitempty"`
	EnableImage  bool                   `json:"enableImage,omitempty"`
	ContractInfo RestaurantContractInfo `json:"contractInfo,omitempty"`
}

type RestaurantContractInfo struct {
	Contact         string `json:"contact,omitempty"`         // 联系方式
	ContactPerson   string `json:"contactPerson,omitempty"`   // 接口人
	TransferBank    string `json:"transferBank,omitempty"`    // 转帐银行
	TransferName    string `json:"transferName,omitempty"`    // 转帐名称
	TransferAccount string `json:"transferAccount,omitempty"` // 转帐帐号
	PaymentMethod   string `json:"paymentMethod,omitempty"`   // 结款方式
}

func (c RestaurantSettings) Value() (driver.Value, error) {
	return json.Marshal(c)
}
func (c *RestaurantSettings) Scan(src interface{}) error {
	if err := json.Unmarshal(src.([]byte), c); err != nil {
		return err
	}
	return nil
}
