package store

import (
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

var store Store

type Store interface {
	Account() AccountStore
	Auth() AuthStore
	Token() TokenStore
}

type AccountStore interface {
	Get(id int) (model.Account, error)
	GetAll() ([]model.Account, error)
	Create(account *model.Account) error
	Update(account *model.Account) error
	Delete(id int) error
}

type AuthStore interface {
	Login(account *model.Account) (model.Token, error)
}

type TokenStore interface {
	GetAll() ([]model.Token, error)
	DeleteAllByAccountId(accountId int) error
}
