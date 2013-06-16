package accounts

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/kobeld/duoerl/global"
	"labix.org/v2/mgo/bson"
	"time"
)

type Account struct {
	Id              bson.ObjectId `bson:"_id"`
	Name            string
	Email           string
	Password        string
	CreatedAt       time.Time
	ConfirmPassword string `bson:"-"`
	Profile         Profile
}

type Profile struct {
	Gender      bool
	Description string
	Location    string
	Birthday    time.Time
	HairTexture string
}

func (this *Account) MakeId() interface{} {
	if this.Id == "" {
		this.Id = bson.NewObjectId()
	}
	return this.Id
}

func NewAccount() *Account {
	return &Account{}
}

func (this *Account) IsPwdMatch(pwd string) bool {
	if bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(pwd)) != nil {
		return false
	}
	return true
}

func LoginWith(email string, pwd string) (account *Account) {
	a, _ := FindByEmail(email)
	if a == nil {
		return
	}
	if a.IsPwdMatch(pwd) {
		account = a
	}

	return
}

func (this *Account) Signup() (err error) {
	this.CreatedAt = time.Now()
	this.encryptPwd()
	return this.Save()
}

func (this *Account) encryptPwd() {
	hp, _ := bcrypt.GenerateFromPassword([]byte(this.Password), 0)
	this.Password = string(hp)
}

func (this *Account) Birthday() string {
	if this.Profile.Birthday.IsZero() {
		return global.TEXT_BIRTHDAY_SECRET
	}
	return this.Profile.Birthday.Format(global.DATE_BIRTHDAY)
}

func (this *Account) Gender() string {
	if !this.Profile.Gender {
		return global.TEXT_GENDER_FEMALE
	}
	return global.TEXT_GENDER_MALE
}
