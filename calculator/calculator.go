package calculator

import (
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-payment-updater/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var BIPS_DENOMINATOR = big.NewInt(10000)

type PaymentCalculator interface {
	// CalculateDistributionsOverInterval returns the distributions of given tokens until the given end timestamp
	CalculateDistributionsUntilTimestamp(endTimestamp *big.Int) (*big.Int, map[gethcommon.Address]*common.Distribution, error)
}

type PaymentCalculatorImpl struct {
	intervalLength *big.Int
	dataService    PaymentCalculatorDataService
}

func NewPaymentCalculator(dataService PaymentCalculatorDataService) PaymentCalculator {
	return &PaymentCalculatorImpl{}
}

func (p *PaymentCalculatorImpl) CalculateDistributionsUntilTimestamp(endTimestamp *big.Int) (*big.Int, map[gethcommon.Address]*common.Distribution, error) {
	startTimestamp, err := p.dataService.GetPaymentsCalculatedUntilTimestamp()
	if err != nil {
		return nil, nil, err
	}

	// make sure the start timestamp is rounded to the nearest interval granularity
	if new(big.Int).Mod(startTimestamp, p.intervalLength).Cmp(big.NewInt(0)) != 0 {
		return nil, nil, fmt.Errorf("start timestamp must be rounded to the nearest interval granularity")
	}

	// round the end timestamp to the nearest interval granularity. the start is assumed to be rounded already
	endTimestamp.Sub(endTimestamp, new(big.Int).Mod(endTimestamp, p.intervalLength))

	// make sure the end timestamp is after the start timestamp
	if endTimestamp.Cmp(startTimestamp) <= 0 {
		return nil, nil, fmt.Errorf("end timestamp must be after start timestamp")
	}

	numIntervals := new(big.Int).Div(new(big.Int).Sub(endTimestamp, startTimestamp), p.intervalLength).Int64()

	// get all range payments that overlap with the given interval
	rangePayments, err := p.dataService.GetRangePaymentsWithOverlappingRange(startTimestamp, endTimestamp)
	if err != nil {
		return nil, nil, err
	}

	// get all distributions at the start timestamp
	distributions, err := p.dataService.GetDistributionsAtTimestamp(startTimestamp)
	if err != nil {
		return nil, nil, err
	}

	// loop through all range payments
	for _, rangePayment := range rangePayments {
		distribution := common.NewDistribution()
		// if we haven't seen this token before, add it to the map
		if distribution, ok := distributions[rangePayment.Token]; !ok {
			distribution = common.NewDistribution()
			distributions[rangePayment.Token] = distribution
		}

		// calculate the payment per second
		paymentPerSecond := new(big.Int).Div(rangePayment.Amount, new(big.Int).Sub(rangePayment.EndRangeTimestamp, rangePayment.StartRangeTimestamp))

		// loop through all intervals
		for i := int64(0); i < numIntervals; i++ {
			// calculate the start and end of the interval
			intervalStart := new(big.Int).Add(startTimestamp, new(big.Int).Mul(p.intervalLength, big.NewInt(i)))
			intervalEnd := new(big.Int).Add(startTimestamp, new(big.Int).Mul(p.intervalLength, big.NewInt(i+1)))

			// calculate overlap between the interval and the range payment
			overlapStart := max(intervalStart, rangePayment.StartRangeTimestamp)
			overlapEnd := min(intervalEnd, rangePayment.EndRangeTimestamp)

			// calculate the payment to distribute
			paymentToDistribute := new(big.Int).Mul(paymentPerSecond, new(big.Int).Sub(overlapEnd, overlapStart))

			// get the operator set at the interval start
			operatorSet, err := p.dataService.GetOperatorSetAtTimestamp(rangePayment.Avs, intervalStart)
			if err != nil {
				return nil, nil, err
			}

			// loop through all operators
			for _, operator := range operatorSet.Operators {
				// totalPaymentToOperatorAndStakers = paymentToDistribute * operatorDelegatedStrategyShares / totalStrategyShares
				totalPaymentToOperatorAndStakers := new(big.Int).Div(new(big.Int).Mul(paymentToDistribute, operator.TotalDelegatedStrategyShares[rangePayment.Strategy]), operatorSet.TotalStakedStrategyShares[rangePayment.Strategy])

				// increment token balance according to the operator's commission
				operatorAmt := distribution.Get(operator.Address)
				if operatorAmt == nil {
					operatorAmt = big.NewInt(0)
				}
				// operatorBalance += totalPaymentToOperatorAndStakers * operatorCommissions
				distribution.Set(operator.Address, operatorAmt.Add(operatorAmt, new(big.Int).Mul(totalPaymentToOperatorAndStakers, operator.Commissions[rangePayment.Avs])))

				// loop through all stakers
				for _, staker := range operator.Stakers {
					// increment token balance according to the staker's proportion of the strategy shares
					stakerAmt := distribution.Get(staker.Address)
					if stakerAmt == nil {
						stakerAmt = big.NewInt(0)
					}

					// stakerBalance += totalPaymentToOperatorAndStakers * (1 - operatorCommissions) * stakerShares / operatorDelegatedStrategyShares
					distribution.Set(staker.Address, stakerAmt.Add(stakerAmt, new(big.Int).Div(new(big.Int).Mul(new(big.Int).Mul(new(big.Int).Sub(BIPS_DENOMINATOR, operator.Commissions[rangePayment.Avs]), totalPaymentToOperatorAndStakers), staker.Shares[rangePayment.Avs]), operator.TotalDelegatedStrategyShares[rangePayment.Strategy])))
				}
			}
		}

	}

	return endTimestamp, distributions, nil
}

func max(a, b *big.Int) *big.Int {
	if a.Cmp(b) == 1 {
		return a
	}
	return b
}

func min(a, b *big.Int) *big.Int {
	if a.Cmp(b) == -1 {
		return a
	}
	return b
}
