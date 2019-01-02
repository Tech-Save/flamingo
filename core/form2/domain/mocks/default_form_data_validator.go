// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import domain "flamingo.me/flamingo/core/form2/domain"
import mock "github.com/stretchr/testify/mock"
import web "flamingo.me/flamingo/framework/web"

// DefaultFormDataValidator is an autogenerated mock type for the DefaultFormDataValidator type
type DefaultFormDataValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: ctx, req, validatorProvider, formData
func (_m *DefaultFormDataValidator) Validate(ctx context.Context, req *web.Request, validatorProvider domain.ValidatorProvider, formData interface{}) (*domain.ValidationInfo, error) {
	ret := _m.Called(ctx, req, validatorProvider, formData)

	var r0 *domain.ValidationInfo
	if rf, ok := ret.Get(0).(func(context.Context, *web.Request, domain.ValidatorProvider, interface{}) *domain.ValidationInfo); ok {
		r0 = rf(ctx, req, validatorProvider, formData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ValidationInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *web.Request, domain.ValidatorProvider, interface{}) error); ok {
		r1 = rf(ctx, req, validatorProvider, formData)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}