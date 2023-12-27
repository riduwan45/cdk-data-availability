// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	cdkvalidium "github.com/0xPolygon/cdk-data-availability/etherman/smartcontracts/cdkvalidium"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"

	common "github.com/ethereum/go-ethereum/common"

	context "context"

	etherman "github.com/0xPolygon/cdk-data-availability/etherman"

	event "github.com/ethereum/go-ethereum/event"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// Etherman is an autogenerated mock type for the Etherman type
type Etherman struct {
	mock.Mock
}

// FilterSequenceBatches provides a mock function with given fields: opts, numBatch
func (_m *Etherman) FilterSequenceBatches(opts *bind.FilterOpts, numBatch []uint64) (*cdkvalidium.CdkvalidiumSequenceBatchesIterator, error) {
	ret := _m.Called(opts, numBatch)

	if len(ret) == 0 {
		panic("no return value specified for FilterSequenceBatches")
	}

	var r0 *cdkvalidium.CdkvalidiumSequenceBatchesIterator
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts, []uint64) (*cdkvalidium.CdkvalidiumSequenceBatchesIterator, error)); ok {
		return rf(opts, numBatch)
	}
	if rf, ok := ret.Get(0).(func(*bind.FilterOpts, []uint64) *cdkvalidium.CdkvalidiumSequenceBatchesIterator); ok {
		r0 = rf(opts, numBatch)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cdkvalidium.CdkvalidiumSequenceBatchesIterator)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.FilterOpts, []uint64) error); ok {
		r1 = rf(opts, numBatch)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCurrentDataCommittee provides a mock function with given fields:
func (_m *Etherman) GetCurrentDataCommittee() (*etherman.DataCommittee, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCurrentDataCommittee")
	}

	var r0 *etherman.DataCommittee
	var r1 error
	if rf, ok := ret.Get(0).(func() (*etherman.DataCommittee, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *etherman.DataCommittee); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*etherman.DataCommittee)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCurrentDataCommitteeMembers provides a mock function with given fields:
func (_m *Etherman) GetCurrentDataCommitteeMembers() ([]etherman.DataCommitteeMember, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetCurrentDataCommitteeMembers")
	}

	var r0 []etherman.DataCommitteeMember
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]etherman.DataCommitteeMember, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []etherman.DataCommitteeMember); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]etherman.DataCommitteeMember)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTx provides a mock function with given fields: ctx, txHash
func (_m *Etherman) GetTx(ctx context.Context, txHash common.Hash) (*types.Transaction, bool, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for GetTx")
	}

	var r0 *types.Transaction
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Transaction, bool, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Transaction); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) bool); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, common.Hash) error); ok {
		r2 = rf(ctx, txHash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// HeaderByNumber provides a mock function with given fields: ctx, number
func (_m *Etherman) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	ret := _m.Called(ctx, number)

	if len(ret) == 0 {
		panic("no return value specified for HeaderByNumber")
	}

	var r0 *types.Header
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (*types.Header, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *types.Header); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Header)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TrustedSequencer provides a mock function with given fields:
func (_m *Etherman) TrustedSequencer() (common.Address, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TrustedSequencer")
	}

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func() (common.Address, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TrustedSequencerURL provides a mock function with given fields:
func (_m *Etherman) TrustedSequencerURL() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for TrustedSequencerURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchSetTrustedSequencer provides a mock function with given fields: opts, sink
func (_m *Etherman) WatchSetTrustedSequencer(opts *bind.WatchOpts, sink chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencer) (event.Subscription, error) {
	ret := _m.Called(opts, sink)

	if len(ret) == 0 {
		panic("no return value specified for WatchSetTrustedSequencer")
	}

	var r0 event.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencer) (event.Subscription, error)); ok {
		return rf(opts, sink)
	}
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencer) event.Subscription); ok {
		r0 = rf(opts, sink)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.WatchOpts, chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencer) error); ok {
		r1 = rf(opts, sink)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchSetTrustedSequencerURL provides a mock function with given fields: opts, sink
func (_m *Etherman) WatchSetTrustedSequencerURL(opts *bind.WatchOpts, sink chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencerURL) (event.Subscription, error) {
	ret := _m.Called(opts, sink)

	if len(ret) == 0 {
		panic("no return value specified for WatchSetTrustedSequencerURL")
	}

	var r0 event.Subscription
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencerURL) (event.Subscription, error)); ok {
		return rf(opts, sink)
	}
	if rf, ok := ret.Get(0).(func(*bind.WatchOpts, chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencerURL) event.Subscription); ok {
		r0 = rf(opts, sink)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(event.Subscription)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.WatchOpts, chan<- *cdkvalidium.CdkvalidiumSetTrustedSequencerURL) error); ok {
		r1 = rf(opts, sink)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEtherman creates a new instance of Etherman. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEtherman(t interface {
	mock.TestingT
	Cleanup(func())
}) *Etherman {
	mock := &Etherman{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}