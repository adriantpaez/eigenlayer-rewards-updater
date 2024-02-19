// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	big "math/big"

	context "context"

	distribution "github.com/Layr-Labs/eigenlayer-payment-updater/common/distribution"

	mock "github.com/stretchr/testify/mock"
)

// PaymentCalculator is an autogenerated mock type for the PaymentCalculator type
type PaymentCalculator struct {
	mock.Mock
}

// CalculateDistributionUntilTimestamp provides a mock function with given fields: ctx, startTimestamp, endTimestamp
func (_m *PaymentCalculator) CalculateDistributionUntilTimestamp(ctx context.Context, startTimestamp *big.Int, endTimestamp *big.Int) (*big.Int, *distribution.Distribution, error) {
	ret := _m.Called(ctx, startTimestamp, endTimestamp)

	if len(ret) == 0 {
		panic("no return value specified for CalculateDistributionUntilTimestamp")
	}

	var r0 *big.Int
	var r1 *distribution.Distribution
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int, *big.Int) (*big.Int, *distribution.Distribution, error)); ok {
		return rf(ctx, startTimestamp, endTimestamp)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int, *big.Int) *big.Int); ok {
		r0 = rf(ctx, startTimestamp, endTimestamp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int, *big.Int) *distribution.Distribution); ok {
		r1 = rf(ctx, startTimestamp, endTimestamp)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*distribution.Distribution)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *big.Int, *big.Int) error); ok {
		r2 = rf(ctx, startTimestamp, endTimestamp)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// NewPaymentCalculator creates a new instance of PaymentCalculator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentCalculator(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentCalculator {
	mock := &PaymentCalculator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
