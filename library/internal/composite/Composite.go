package composite

import "auth/internal/store/sqllite"

type Composites struct {
	Auth AuthComposite
}

func (c *Composites) NewComposites(store sqllite.Store) {

	c.Auth.New(store.Logger, store.UserStorage)

}
