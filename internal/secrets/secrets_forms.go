package secrets

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type EncryptedPayloadForm struct {
	// don't use: uint8 --> it serializes to char
	Encrypted []uint `json:"encrypted"`
	Salt      []uint `json:"salt"`
	IV        []uint `json:"iv"`
}

func (f EncryptedPayloadForm) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Salt, validation.When(len(f.Encrypted) > 0, validation.Required)),
		validation.Field(&f.IV, validation.When(len(f.Encrypted) > 0, validation.Required)))
}

type UserSecretForm struct {
	CategoryId        int                  `json:"categoryId"`
	NewCategoryName   string               `json:"newCategoryName"`
	Description       string               `json:"description"`
	Username          string               `json:"username"`
	PasswordEncrypted EncryptedPayloadForm `json:"passwordEncrypted"`
	SafeNoteEncrypted EncryptedPayloadForm `json:"safeNoteEncrypted"`
	URLSite           string               `json:"urlSite"`
}

func (f UserSecretForm) ValidateFront() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.CategoryId, validation.When(f.NewCategoryName == "", validation.Required)),
		validation.Field(&f.NewCategoryName, validation.When(f.CategoryId == 0, validation.Required, validation.Length(3, 50))),
		validation.Field(&f.Description, validation.When(f.Description != "", validation.Length(0, 250))),
		validation.Field(&f.Username, validation.When(f.Username != "", validation.Length(0, 250))),
		validation.Field(&f.PasswordEncrypted),
		validation.Field(&f.SafeNoteEncrypted),
		validation.Field(&f.URLSite, validation.When(f.URLSite != "", is.URL, validation.Length(0, 250))))
}

func (f UserSecretForm) Save(userId int) (int, error) {
	bytesPasswordEncrypted, _ := json.Marshal(f.PasswordEncrypted)
	bytesSafeNoteEncrypted, _ := json.Marshal(f.SafeNoteEncrypted)

	newSecretId, err := SaveNewUserSecretDB(NewUserSecretModel{
		UserId:            userId,
		CategoryId:        f.CategoryId,
		NewCategoryName:   f.NewCategoryName,
		Description:       f.Description,
		Username:          f.Username,
		PasswordEncrypted: bytesPasswordEncrypted,
		SafeNoteEncrypted: bytesSafeNoteEncrypted,
		URLSite:           f.URLSite,
	})

	return newSecretId, err
}

type UpdateUserSecretForm struct {
	Id                int                  `json:"id"`
	CategoryId        int                  `json:"categoryId"`
	NewCategoryName   string               `json:"newCategoryName"`
	Description       string               `json:"description"`
	Username          string               `json:"username"`
	PasswordEncrypted EncryptedPayloadForm `json:"passwordEncrypted"`
	SafeNoteEncrypted EncryptedPayloadForm `json:"safeNoteEncrypted"`
	URLSite           string               `json:"urlSite"`
}

func (f UpdateUserSecretForm) ValidateFront() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Id, validation.Required),
		validation.Field(&f.CategoryId, validation.When(f.NewCategoryName == "", validation.Required)),
		validation.Field(&f.NewCategoryName, validation.When(f.CategoryId == 0, validation.Required, validation.Length(3, 50))),
		validation.Field(&f.Description, validation.When(f.Description != "", validation.Length(0, 250))),
		validation.Field(&f.Username, validation.When(f.Username != "", validation.Length(0, 250))),
		validation.Field(&f.PasswordEncrypted),
		validation.Field(&f.SafeNoteEncrypted),
		validation.Field(&f.URLSite, validation.When(f.URLSite != "", is.URL, validation.Length(0, 250))))
}

func (f UpdateUserSecretForm) Update(userId int) error {
	bytesPasswordEncrypted, _ := json.Marshal(f.PasswordEncrypted)
	bytesSafeNoteEncrypted, _ := json.Marshal(f.SafeNoteEncrypted)

	err := UpdateUserSecretDB(UpdateUserSecretModel{
		UserId:            userId,
		SecretId:          f.Id,
		CategoryId:        f.CategoryId,
		NewCategoryName:   f.NewCategoryName,
		Description:       f.Description,
		Username:          f.Username,
		PasswordEncrypted: bytesPasswordEncrypted,
		SafeNoteEncrypted: bytesSafeNoteEncrypted,
		URLSite:           f.URLSite,
	})
	return err
}
