package nbastats

//go:generate paramgen -output=paramgen_src.go SourceConfig

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"golang.org/x/time/rate"
)

type Source struct {
	sdk.UnimplementedSource

	config                  SourceConfig
	lastPositionRead        sdk.Position //nolint:unused // this is just an example
	limiter                 *rate.Limiter
	cachedSpeedDistanceData []byte
}

type SourceConfig struct {
	// Config includes parameters that are the same in the source and destination.
	Config
	// SourceConfigParam is named foo and must be provided by the user.
	SourceConfigParam string `json:"foo" validate:"required"`
}

func NewSource() sdk.Source {
	// Create Source and wrap it in the default middleware.
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

func (s *Source) Parameters() map[string]sdk.Parameter {
	// Parameters is a map of named Parameters that describe how to configure
	// the Source. Parameters can be generated from SourceConfig with paramgen.
	return s.config.Parameters()
}

func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	// Configure is the first function to be called in a connector. It provides
	// the connector with the configuration that can be validated and stored.
	// In case the configuration is not valid it should return an error.
	// Testing if your connector can reach the configured data source should be
	// done in Open, not in Configure.
	// The SDK will validate the configuration and populate default values
	// before calling Configure. If you need to do more complex validations you
	// can do them manually here.

	sdk.Logger(ctx).Info().Msg("Configuring Source...")
	err := sdk.Util.ParseConfig(cfg, &s.config)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	return nil
}

func (s *Source) Open(ctx context.Context, pos sdk.Position) error {
	// Open is called after Configure to signal the plugin it can prepare to
	// start producing records. If needed, the plugin should open connections in
	// this function. The position parameter will contain the position of the
	// last record that was successfully processed, Source should therefore
	// start producing records after this position. The context passed to Open
	// will be cancelled once the plugin receives a stop signal from Conduit.
	s.limiter = rate.NewLimiter(rate.Every(s.config.PollingPeriod), 1)
	return nil
}

func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	// Read returns a new Record and is supposed to block until there is either
	// a new record or the context gets cancelled. It can also return the error
	// ErrBackoffRetry to signal to the SDK it should call Read again with a
	// backoff retry.
	// If Read receives a cancelled context or the context is cancelled while
	// Read is running it must stop retrieving new records from the source
	// system and start returning records that have already been buffered. If
	// there are no buffered records left Read must return the context error to
	// signal a graceful stop. If Read returns ErrBackoffRetry while the context
	// is cancelled it will also signal that there are no records left and Read
	// won't be called again.
	// After Read returns an error the function won't be called again (except if
	// the error is ErrBackoffRetry, as mentioned above).
	// Read can be called concurrently with Ack.
	err := s.limiter.Wait(ctx)
	if err != nil {
		return sdk.Record{}, err
	} else {
		sdk.Logger(ctx).Info().Msgf("Waiting for %s before next request for data", s.config.PollingPeriod)
	}
	rec, err := s.getRecord(ctx)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("error getting the weather data: %w", err)
	}
	return rec, nil
}

func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	// Ack signals to the implementation that the record with the supplied
	// position was successfully processed. This method might be called after
	// the context of Read is already cancelled, since there might be
	// outstanding acks that need to be delivered. When Teardown is called it is
	// guaranteed there won't be any more calls to Ack.
	// Ack can be called concurrently with Read.
	sdk.Logger(ctx).Debug().Str("position", string(position)).Msg("got ack")
	return nil
}

func (s *Source) Teardown(ctx context.Context) error {
	// Teardown signals to the plugin that there will be no more calls to any
	// other function. After Teardown returns, the plugin should be ready for a
	// graceful shutdown.
	return nil
}

func (s *Source) getRecord(ctx context.Context) (sdk.Record, error) {
	speedDistanceData, err := fetchNBASpeedDistanceStats(s.config.PerMode)
	if err != nil {
		return sdk.Record{}, err
	}

	sdk.Logger(ctx).Info().Msg("Successfully fetched the NBA Speed and Distance data...")
	// if s.cachedSpeedDistanceData == nil || bytes.Equal(speedDistanceData, s.cachedSpeedDistanceData) == false {
	// 	s.cachedSpeedDistanceData = speedDistanceData
	// 	sdk.Logger(ctx).Info().Msg("Successfully fetched the NBA Speed and Distance data...")
	// } else {
	// 	sdk.Logger(ctx).Info().Msg("Fetched stats data is same as cached data....")
	// 	return sdk.Record{}, nil
	// }
	// Get current timestamp
	currentTime := time.Now()

	// Format the timestamp as a string (You can customize the format as needed)
	timestampStr := currentTime.Format("2006-01-02-1504")

	// Create the final string using the pattern with the formatted timestamp
	key := fmt.Sprintf("%s_%s", timestampStr, s.config.PerMode)
	recordKey := sdk.RawData(key)
	recordValue := sdk.RawData(speedDistanceData)
	return sdk.Util.Source.NewRecordCreate(
		sdk.Position(recordKey),
		nil,
		recordKey,
		recordValue,
	), nil
}
