// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	context "context"
	big "math/big"

	mock "github.com/stretchr/testify/mock"
)

// Transactor is an autogenerated mock type for the Transactor type
type Transactor struct {
	mock.Mock
}

// CurrRewardsCalculationEndTimestamp provides a mock function with given fields:
func (_m *Transactor) CurrRewardsCalculationEndTimestamp() (uint32, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CurrRewardsCalculationEndTimestamp")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func() (uint32, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNumberOfPublishedRoots provides a mock function with given fields:
func (_m *Transactor) GetNumberOfPublishedRoots() (*big.Int, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetNumberOfPublishedRoots")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func() (*big.Int, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *big.Int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRootIndex provides a mock function with given fields: root
func (_m *Transactor) GetRootIndex(root [32]byte) (uint32, error) {
	ret := _m.Called(root)

	if len(ret) == 0 {
		panic("no return value specified for GetRootIndex")
	}

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func([32]byte) (uint32, error)); ok {
		return rf(root)
	}
	if rf, ok := ret.Get(0).(func([32]byte) uint32); ok {
		r0 = rf(root)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func([32]byte) error); ok {
		r1 = rf(root)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitRoot provides a mock function with given fields: ctx, root, rewardsUnixTimestamp
func (_m *Transactor) SubmitRoot(ctx context.Context, root [32]byte, rewardsUnixTimestamp uint32) error {
	ret := _m.Called(ctx, root, rewardsUnixTimestamp)

	if len(ret) == 0 {
		panic("no return value specified for SubmitRoot")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, [32]byte, uint32) error); ok {
		r0 = rf(ctx, root, rewardsUnixTimestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransactor creates a new instance of Transactor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactor(t interface {
	mock.TestingT
	Cleanup(func())
}) *Transactor {
	mock := &Transactor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
