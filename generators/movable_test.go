package generators

import (
	adapters2 "github.com/bindasov/adapterGenerator/adapters"
	"github.com/bindasov/ioc/ioc"
	"github.com/bindasov/spaceBattle/adapters"
	"github.com/bindasov/spaceBattle/commands"
	"github.com/bindasov/spaceBattle/services/mocks"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestMovableGenerator_Generate(t *testing.T) {
	type deps struct {
		movableAdapter adapters.MovableAdapter
		IoC            *ioc.IoC
	}
	tests := []struct {
		name    string
		handler func(*testing.T, *deps)
	}{
		{
			name: "generation success",
			handler: func(t *testing.T, deps *deps) {
				expected := "package adapters\n\nimport (\n\t\"github.com/bindasov/ioc/ioc\"\n\t\"github.com/bindasov/spaceBattle/adapters\"\n\t\"github.com/bindasov/spaceBattle/models\"\n)\n\nfunc NewUObject(adapter adapters.MovableAdapter, ioc *ioc.IoC) *UObject {\n\tobj := &UObject{\n\t\tobj: adapter,\n\t\tioc: ioc,\n\t}\n\treturn obj\n}\n\ntype UObject struct {\n\tobj adapters.MovableAdapter\n\tioc *ioc.IoC\n}\n\nfunc (m *UObject) GetPosition() *models.Vector {\n\treturn m.ioc.Resolve(\"IMovable:Position.Get\", m.obj).(*models.Vector)\n}\nfunc (m *UObject) GetVelocity() *models.Vector {\n\treturn m.ioc.Resolve(\"IMovable:Velocity.Get\", m.obj).(*models.Vector)\n}\nfunc (m *UObject) SetPosition(value *models.Vector) {\n\tm.ioc.Resolve(\"IMovable:Position.Set\", m.obj, value)\n}\n"
				result := Generate(reflect.TypeOf((*adapters.MovableAdapter)(nil)).Elem())
				require.Equal(t, expected, result)
			},
		},
		{
			name: "getting adapter success",
			handler: func(t *testing.T, deps *deps) {
				expected := reflect.TypeOf((*adapters2.UObject)(nil))
				result := deps.IoC.Resolve("Adapter", reflect.TypeOf((*adapters.MovableAdapter)(nil)).Elem(), deps.movableAdapter, deps.IoC)
				require.Equal(t, expected, reflect.TypeOf(result))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			movableMock := mocks.NewMovable(t)
			rotableMock := mocks.NewRotable(t)
			adapter := adapters.NewMovable(movableMock, rotableMock)

			IoC := ioc.NewIoC()

			IoC.Resolve("IoC.Register", "Adapter", func(args ...interface{}) interface{} {
				Generate(args[0])
				return adapters2.NewUObject(args[1].(adapters.MovableAdapter), args[2].(*ioc.IoC))
			}).(commands.Command).Execute()

			deps := &deps{
				movableAdapter: adapter,
				IoC:            IoC,
			}

			tc.handler(t, deps)
		})
	}
}
