package configs

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/onflow/flow-go-sdk"
	log "github.com/sirupsen/logrus"
)

type Config struct {

	// -- Logger config --
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// -- Feature flags --

	DisableRawTransactions   bool `env:"DISABLE_RAWTX"`
	DisableFungibleTokens    bool `env:"DISABLE_FT"`
	DisableNonFungibleTokens bool `env:"DISABLE_NFT"`
	DisableChainEvents       bool `env:"DISABLE_CHAIN_EVENTS"`

	// -- Admin account --

	AdminAddress    string `env:"ADMIN_ADDRESS,notEmpty"`
	AdminKeyIndex   int    `env:"ADMIN_KEY_INDEX" envDefault:"0"`
	AdminKeyType    string `env:"ADMIN_KEY_TYPE" envDefault:"local"`
	AdminPrivateKey string `env:"ADMIN_PRIVATE_KEY,notEmpty"`
	// This sets the number of proposal keys to be used on the admin account.
	// You can increase transaction throughput by using multiple proposal keys for
	// parallel transaction execution.
	AdminProposalKeyCount uint16 `env:"ADMIN_PROPOSAL_KEY_COUNT" envDefault:"1"`

	// -- Keys --

	// When "DefaultKeyType" is set to "local", private keys are generated by the API
	// and stored as encrypted text in the database.
	// KMS key types:
	// - aws_kms
	// - google_kms
	DefaultKeyType  string `env:"DEFAULT_KEY_TYPE" envDefault:"local"`
	DefaultKeyIndex int    `env:"DEFAULT_KEY_INDEX" envDefault:"0"`
	// If the default of "-1" is used for "DefaultKeyWeight"
	// the service will use flow.AccountKeyWeightThreshold from the Flow SDK.
	DefaultKeyWeight int    `env:"DEFAULT_KEY_WEIGHT" envDefault:"-1"`
	DefaultSignAlgo  string `env:"DEFAULT_SIGN_ALGO" envDefault:"ECDSA_P256"`
	DefaultHashAlgo  string `env:"DEFAULT_HASH_ALGO" envDefault:"SHA3_256"`
	// This symmetrical key is used to encrypt private keys
	// that are stored in the database. Values per type:
	// - local: 32 bytes long encryption key
	// - aws_kms: key ARN, e.g. arn:aws:kms:us-west-1:123456789000:key/00000000-1111-2222-3333-444444444444
	// - google_kms: key resource name (without version info), e.g. projects/my-project/locations/europe-north1/keyRings/my-keyring/cryptoKeys/my-encryption-key
	EncryptionKey string `env:"ENCRYPTION_KEY,notEmpty"`
	// Encryption key type, one of: local, aws_kms, google_kms
	EncryptionKeyType string `env:"ENCRYPTION_KEY_TYPE,notEmpty" envDefault:"local"`
	// DefaultAccountKeyCount specifies how many times the account key will be duplicated upon account creation, does not affect existing accounts
	DefaultAccountKeyCount uint `env:"DEFAULT_ACCOUNT_KEY_COUNT" envDefault:"1"`

	// -- Database --

	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"wallet.db"`
	DatabaseType    string `env:"DATABASE_TYPE" envDefault:"sqlite"`
	DatabaseVersion string `env:"DATABASE_VERSION" envDefault:""`

	// -- Host and chain access --

	Host                 string        `env:"HOST"`
	Port                 int           `env:"PORT" envDefault:"3000"`
	ServerRequestTimeout time.Duration `env:"SERVER_REQUEST_TIMEOUT" envDefault:"60s"`
	AccessAPIHost        string        `env:"ACCESS_API_HOST,notEmpty"`
	ChainID              flow.ChainID  `env:"CHAIN_ID" envDefault:"flow-emulator"`

	// -- Templates --

	EnabledTokens                            []string `env:"ENABLED_TOKENS" envSeparator:","`
	ScriptPathCreateAccount                  string   `env:"SCRIPT_PATH_CREATE_ACCOUNT" envDefault:""`
	InitFungibleTokenVaultsOnAccountCreation bool     `env:"INIT_FUNGIBLE_TOKEN_VAULTS_ON_ACCOUNT_CREATION" envDefault:"false"`

	// -- Workerpool --

	// Defines the maximum number of active jobs that can be queued before
	// new jobs are rejected.
	WorkerQueueCapacity uint `env:"WORKER_QUEUE_CAPACITY" envDefault:"1000"`
	// Number of concurrent workers handling incoming jobs.
	// You can increase the number of workers if you're sending
	// too many transactions and find that the queue is often backlogged.
	WorkerCount uint `env:"WORKER_COUNT" envDefault:"1"`
	// Webhook endpoint to receive job status updates
	JobStatusWebhookUrl string `env:"JOB_STATUS_WEBHOOK" envDefault:""`
	// Duration for which to wait for a response, if 0 wait indefinitely. Default: 30s.
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// For more info: https://pkg.go.dev/time#ParseDuration
	JobStatusWebhookTimeout time.Duration `env:"JOB_STATUS_WEBHOOK_TIMEOUT" envDefault:"30s"`

	// -- Google KMS --

	GoogleKMSProjectID  string `env:"GOOGLE_KMS_PROJECT_ID"`
	GoogleKMSLocationID string `env:"GOOGLE_KMS_LOCATION_ID"`
	GoogleKMSKeyRingID  string `env:"GOOGLE_KMS_KEYRING_ID"`

	// -- Misc --

	// Duration for which to wait for a transaction seal, if 0 wait indefinitely. Default: 0.
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// For more info: https://pkg.go.dev/time#ParseDuration
	TransactionTimeout time.Duration `env:"TRANSACTION_TIMEOUT" envDefault:"0"`

	// Idempotency middleware configuration
	DisableIdempotencyMiddleware bool `env:"DISABLE_IDEMPOTENCY_MIDDLEWARE" envDefault:"false"`
	// Idempotency middleware database type;
	// - "local", in-memory w/ no multi-instance support
	// - "shared", sql (gorm) database shared with the app (DatabaseType)
	// - "redis"
	IdempotencyMiddlewareDatabaseType string `env:"IDEMPOTENCY_MIDDLEWARE_DATABASE_TYPE" envDefault:"local"`
	// Redis URL for idempotency key storage, e.g. "redis://walletapi:wallet-api-redis@localhost:6379/"
	IdempotencyMiddlewareRedisURL string `env:"IDEMPOTENCY_MIDDLEWARE_REDIS_URL" envDefault:""`

	// Set the starting height for event polling. This won't have any effect if the value in
	// database (chain_event_status[0].latest_height) is greater.
	// If 0 (default) use latest block height if starting fresh (no previous value in database).
	ChainListenerStartingHeight uint64 `env:"EVENTS_STARTING_HEIGHT" envDefault:"0"`
	// Maximum number of blocks to check at once.
	ChainListenerMaxBlocks uint64 `env:"EVENTS_MAX_BLOCKS" envDefault:"100"`
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// For more info: https://pkg.go.dev/time#ParseDuration
	ChainListenerInterval time.Duration `env:"EVENTS_INTERVAL" envDefault:"10s"`

	// Max transactions per second, rate at which the service can submit transactions to Flow (excluding ops)
	TransactionMaxSendRate int `env:"MAX_TPS" envDefault:"10"`

	// maxJobErrorCount is the maximum number of times a Job can be tried to
	// execute before considering it completely failed.
	MaxJobErrorCount int `env:"MAX_JOB_ERROR_COUNT" envDefault:"10"`

	// Poll DB for new schedulable jobs every 30s.
	DBJobPollInterval time.Duration `env:"DB_JOB_POLL_INTERVAL" envDefault:"30s"`

	// Grace time period before re-scheduling jobs that are in state INIT or
	// ACCEPTED. These are jobs where the executor processing has been
	// unexpectedly disrupted (such as bug, dead node, disconnected
	// networking etc.).
	AcceptedGracePeriod time.Duration `env:"ACCEPTED_GRACE_PERIOD" envDefault:"180s"`

	// Grace time period before re-scheduling jobs that are up for immediate
	// restart (such as NO_AVAILABLE_WORKERS or ERROR).
	ReSchedulableGracePeriod time.Duration `env:"RESCHEDULABLE_GRACE_PERIOD" envDefault:"60s"`

	// Sleep duration in case of service isHalted
	PauseDuration time.Duration `env:"PAUSE_DURATION" envDefault:"60s"`

	GrpcMaxCallRecvMsgSize int `env:"GRPC_MAX_CALL_RECV_MSG_SIZE" envDefault:"16777216"`

	// -- ops ---
	// WorkerCount for system jobs, max number of in-flight transactions
	OpsWorkerCount uint `env:"OPS_WORKER_COUNT" envDefault:"200"`
	// Capacity of buffered jobs queues for system jobs.
	OpsWorkerQueueCapacity uint `env:"OPS_WORKER_QUEUE_CAPACITY" envDefault:"300000"`
}

// Parse parses environment variables and flags to a valid Config.
func Parse(opts ...env.Options) (*Config, error) {
	cfg := Config{}
	opts = append(opts, env.Options{Prefix: "FLOW_WALLET_"})
	err := env.Parse(&cfg, opts...)
	return &cfg, err
}

func ConfigureLogger(logLevel string) {
	ll, err := log.ParseLevel(logLevel)
	if err != nil {
		ll = log.DebugLevel
	}

	log.SetLevel(ll)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}

func SetenvIfNotSet(key string, value string) {
	if _, isSet := os.LookupEnv(key); !isSet {
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}
	}
}

func ParseTestConfig(t *testing.T) *Config {
	t.Helper()

	SetenvIfNotSet("FLOW_WALLET_ADMIN_ADDRESS", "0xf8d6e0586b0a20c7")
	SetenvIfNotSet("FLOW_WALLET_ADMIN_PRIVATE_KEY", "91a22fbd87392b019fbe332c32695c14cf2ba5b6521476a8540228bdf1987068")
	SetenvIfNotSet("FLOW_WALLET_ACCESS_API_HOST", "localhost:3569")
	SetenvIfNotSet("FLOW_WALLET_ENCRYPTION_KEY", "faae4ed1c30f4e4555ee3a71f1044a8e")
	SetenvIfNotSet("FLOW_WALLET_ENCRYPTION_KEY_TYPE", "local")
	SetenvIfNotSet("FLOW_WALLET_ENABLED_TOKENS", "FlowToken:0x0ae53cb6e3f42a79:flowToken,FUSD:0xf8d6e0586b0a20c7:fusd")

	cfg, err := Parse()
	if err != nil {
		t.Fatal(err)
	}

	// Check if the db DSN contains "test", if it does not, force it to prevent accidents
	if !strings.Contains(cfg.DatabaseDSN, "test") {
		cfg.DatabaseDSN = path.Join(t.TempDir(), "test.db")
		cfg.DatabaseType = "sqlite"
	}

	// Force this to prevent accidents
	cfg.ChainID = flow.Emulator

	return cfg
}
