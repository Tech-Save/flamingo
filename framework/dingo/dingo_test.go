package dingo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	Interface interface {
		Test() int
	}

	InterfaceSub Interface

	InterfaceImpl1 struct {
		i   int
		foo string
	}

	InterfaceImpl2 struct {
		i int
	}

	IfaceProvider func() Interface

	DepTest struct {
		Iface  Interface `inject:""`
		Iface2 Interface `inject:"test"`

		IfaceProvider      IfaceProvider `inject:""`
		IfaceProvided      Interface     `inject:"provider"`
		IfaceImpl1Provided Interface     `inject:"providerimpl1"`
		IfaceInstance      Interface     `inject:"instance"`
	}

	TestSingleton struct {
		i int
	}

	TestModule struct{}

	PreTestModule struct{}
)

func InterfaceProvider(str string) Interface {
	return &InterfaceImpl1{foo: str}
}

func InterfaceImpl1Provider(str string) *InterfaceImpl1 {
	return &InterfaceImpl1{foo: str}
}

func (ptm *PreTestModule) Configure(injector *Injector) {
	injector.Bind((*string)(nil)).ToInstance("Hello World")
}

func (tm *TestModule) Configure(injector *Injector) {
	injector.Bind((*Interface)(nil)).To((*InterfaceSub)(nil))
	injector.Bind((*InterfaceSub)(nil)).To(InterfaceImpl1{})
	injector.Bind((*Interface)(nil)).AnnotatedWith("test").To(InterfaceImpl2{})

	injector.Bind((*Interface)(nil)).AnnotatedWith("provider").ToProvider(InterfaceProvider)
	injector.Bind((*Interface)(nil)).AnnotatedWith("providerimpl1").ToProvider(InterfaceImpl1Provider)
	injector.Bind((*Interface)(nil)).AnnotatedWith("instance").ToInstance(new(InterfaceImpl2))

	injector.Bind(TestSingleton{}).AsEagerSingleton()
}

func (if1 *InterfaceImpl1) Test() int {
	return 1
}

func (if2 *InterfaceImpl2) Test() int {
	return 2
}

// AOP
type (
	AopInterface interface {
		Test() string
	}

	AopImpl struct{}

	AopDep struct {
		A AopInterface `inject:""`
	}

	AopInterceptor1 struct {
		AopInterface
	}

	AopInterceptor2 struct {
		AopInterface
	}

	AopModule struct{}
)

func (m *AopModule) Configure(injector *Injector) {
	injector.Bind((*AopInterface)(nil)).To(AopImpl{})

	injector.BindInterceptor((*AopInterface)(nil), AopInterceptor1{})
	injector.BindInterceptor((*AopInterface)(nil), AopInterceptor2{})
}

func (a *AopImpl) Test() string {
	return "Test"
}

func (a *AopInterceptor1) Test() string {
	return a.AopInterface.Test() + " 1"
}

func (a *AopInterceptor2) Test() string {
	return a.AopInterface.Test() + " 2"
}

func TestDingo(t *testing.T) {
	t.Run("Dingo Test", func(t *testing.T) {
		t.Run("Simple resolve", func(t *testing.T) {
			t.Run("Should resolve dependencies on request", func(t *testing.T) {
				injector := NewInjector(new(PreTestModule), new(TestModule))

				var iface Interface
				iface = injector.GetInstance(new(Interface)).(Interface)

				assert.Equal(t, 1, iface.Test())

				dt := *injector.GetInstance(new(DepTest)).(*DepTest)

				assert.Equal(t, 1, dt.Iface.Test())
				assert.Equal(t, 2, dt.Iface2.Test())

				var dt2 DepTest
				injector.requestInjection(&dt2)

				assert.Equal(t, 1, dt2.Iface.Test())
				assert.Equal(t, 2, dt2.Iface2.Test())

				assert.Equal(t, 1, dt.IfaceProvided.Test())
				assert.Equal(t, 1, dt.IfaceImpl1Provided.Test())
				assert.Equal(t, 2, dt.IfaceInstance.Test())

				assert.Equal(t, 1, dt.IfaceProvider().Test())
				assert.Equal(t, "Hello World", dt.IfaceProvided.(*InterfaceImpl1).foo)
				assert.Equal(t, "Hello World", dt.IfaceImpl1Provided.(*InterfaceImpl1).foo)
			})

			t.Run("Should resolve scopes", func(t *testing.T) {
				injector := NewInjector(new(TestModule))

				assert.Equal(t, injector.GetInstance(TestSingleton{}), injector.GetInstance(TestSingleton{}))
			})
		})

		t.Run("AOP", func(t *testing.T) {
			t.Run("Should intercept interfaces", func(t *testing.T) {
				injector := NewInjector(new(AopModule))

				var dep AopDep
				injector.requestInjection(&dep)

				assert.Equal(t, "Test 1 2", dep.A.Test())
			})
		})
	})
}