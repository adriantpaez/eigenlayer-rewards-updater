package updater

import (
	"context"
	"time"

	calculator "github.com/Layr-Labs/eigenlayer-payment-updater/calculator"
	"github.com/Layr-Labs/eigenlayer-payment-updater/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type Updater struct {
	UpdateInterval time.Duration
	dataService    PaymentDataService
	calculator     calculator.PaymentCalculator
	transactor     UpdaterTransactor
}

func NewUpdater(
	updateInterval time.Duration,
	calculator calculator.PaymentCalculator,
	chainClient *common.ChainClient,
	claimingManagerAddress gethcommon.Address,
) (*Updater, error) {
	dataService := NewPaymentDataService(chainClient)

	transactor, err := NewUpdaterTransactor(chainClient, claimingManagerAddress)
	if err != nil {
		log.Fatal().Msgf("failed to create transactor: %s", err)
	}

	return &Updater{
		UpdateInterval: updateInterval,
		dataService:    dataService,
		calculator:     calculator,
		transactor:     transactor,
	}, nil
}

func (u *Updater) Start() error {
	// run a loop unning once every u.UpdateInterval that calls u.update()
	log.Info().Msg("service started")
	ctx := context.Background()
	ticker := time.NewTicker(u.UpdateInterval)
	for range ticker.C {
		log.Info().Msg("running update")
		if err := u.update(ctx); err != nil {
			log.Error().Msgf("failed to update: %s", err)
		}
	}
	return nil
}

func (u *Updater) update(ctx context.Context) error {
	// get the interval of time that we need to update payments for
	log.Info().Msg("getting latest finalized timestamp")
	latestFinalizedTimestamp, err := u.dataService.GetLatestFinalizedTimestamp(ctx)
	if err != nil {
		return err
	}

	// give the interval to the distribution calculator, get the map from address => token => amount
	log.Info().Msg("calculating distribution")
	paymentsCalculatedUntilTimestamp, newDistributions, err := u.calculator.CalculateDistributionsUntilTimestamp(ctx, latestFinalizedTimestamp)
	if err != nil {
		return err
	}

	// add the pending distribution to the previous distribution
	log.Info().Msg("adding pending distribution to previous distribution and merklizing")
	distributionRoots := make([][]byte, len(newDistributions))
	for token, distribution := range newDistributions {
		newDistributions[token].Add(distribution)

		distributionRoot, err := newDistributions[token].Merklize(token)
		if err != nil {
			return err
		}

		distributionRoots = append(distributionRoots, distributionRoot[:])
	}

	// merklize the distribution roots
	log.Info().Msg("merklizing distribution roots")
	root, err := common.Merklize(distributionRoots)

	// send the merkle root to the smart contract
	log.Info().Msg("updating payments")
	if err := u.transactor.SubmitRoot(ctx, root, paymentsCalculatedUntilTimestamp); err != nil {
		return err
	}

	return nil
}
