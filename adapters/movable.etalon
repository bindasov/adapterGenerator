package adapters

import (
	"github.com/bindasov/ioc/ioc"
	"github.com/bindasov/spaceBattle/adapters"
	"github.com/bindasov/spaceBattle/models"
)

func NewUObject(adapter adapters.MovableAdapter, ioc *ioc.IoC) *UObject {
	obj := &UObject{
		obj: adapter,
		ioc: ioc,
	}
	return obj
}

type UObject struct {
	obj adapters.MovableAdapter
	ioc *ioc.IoC
}

func (m *UObject) GetPosition() *models.Vector {
	return m.ioc.Resolve("IMovable:Position.Get", m.obj).(*models.Vector)
}
func (m *UObject) GetVelocity() *models.Vector {
	return m.ioc.Resolve("IMovable:Velocity.Get", m.obj).(*models.Vector)
}
func (m *UObject) SetPosition(value *models.Vector) {
	m.ioc.Resolve("IMovable:Position.Set", m.obj, value)
}
