package store

import (
	"github.com/David-Orson/Surf-Secrets-3-Api/store/model"
)

var store Store

type Store interface {
	Account() AccountStore
	Auth() AuthStore
	Finder() FinderStore
	Map() MapStore
	Match() MatchStore
	Token() TokenStore

	HelperStore
}

type AccountStore interface {
	Get(username string) (model.Account, error)
	GetAll() ([]model.Account, error)
	Create(account *model.Account) error
	Update(account *model.Account) error
	Delete(id int) error
	IncrementWin(id int) error
	IncrementLoss(id int) error
	IncrementDispute(id int) error
}

type AuthStore interface {
	Login(account *model.Account) (model.Token, error)
}

type FinderStore interface {
	GetPost(id int) (model.FinderPost, error)
	GetAvailablePost(id int) (model.FinderPost, error)
	GetAllPosts() ([]model.FinderPost, error)
	CreatePost(finderPost *model.FinderPost) error
	SetAccepted(id int) error
}

type HelperStore interface {
	IsDuplicate(table string, column string, field interface{}, id int) bool
}

type MapStore interface {
	GetAll() ([]model.Map, error)
}
type MatchStore interface {
	Get(id int) (model.Match, error)
	GetAll() ([]model.Match, error)
	GetByAccount(id int) ([]model.Match, error)
	GetDisputesByAccount(id int) ([]model.Match, error)
	Create(match *model.Match) error
	Update(match *model.Match) error
}

type TokenStore interface {
	GetAll() ([]model.Token, error)
	DeleteAllByAccountId(accountId int) error
}
