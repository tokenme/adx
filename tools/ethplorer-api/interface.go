package ethplorer

import "fmt"

type Error struct {
	Code    int    `json:"code"`    // error code (integer)
	Message string `json:"message"` // error message
}

func (e Error) Error() string {
	return fmt.Sprintf("Code:%d, Msg:%s", e.Code, e.Message)
}

type TokenPrice struct {
	Rate            float64 `json:"rate"`     // current rate
	Currency        string  `json:"currency"` // token price currency (USD)
	Diff            float64 `json:"diff"`     // 24 hour rate difference (in percent)
	Diff7d          float64 `json:"diff7d"`
	Ts              int64   `json:"ts"` // last rate update timestamp
	MarketCapUsd    float64 `json:"marketCapUsd"`
	AvailableSupply float64 `json:"availableSupply"`
	Volume24h       float64 `json:"volume24h"`
}

type Token struct {
	Address        string     `json:"address"`        // token address
	TotalSupply    string     `json:"totalSupply"`    // total token supply
	Name           string     `json:"name"`           // token name
	Symbol         string     `json:"symbol"`         // token symbol
	Decimals       string     `json:"decimals"`       // number of significant digits
	Price          TokenPrice `json:"token_price"`    // token price (false, if not available)
	Owner          string     `json:"owner"`          // token owner address
	CountOps       uint64     `json:"countOps"`       // total count of token operations
	TotalIn        uint64     `json:"totalIn"`        // total amount of incoming tokens
	TotalOut       uint64     `json:"totalOut"`       // total amount of outgoing tokens
	HoldersCount   uint64     `json:"holdersCount"`   // total numnber of token holders
	IssuancesCount uint64     `json:"issuancesCount"` // total count of token issuances
}

type Contract struct {
	CreatorAddress  string `json:"creatorAddress"`  //contract creator address
	TransactionHash string `json:"transactionHash"` //contract creation transaction hash
	Timestamp       int64  `json:"timestamp"`       //contract creation timestamp
}

type ETH struct {
	Balance  float64 `json:"balance"`  // ETH balance
	TotalIn  float64 `json:"totalIn"`  // Total incoming ETH value
	TotalOut float64 `json:"totalOut"` // Total outgoing ETH value
}

type AddressToken struct {
	Token    Token  `json:"tokenInfo"` // token data (same format as token info)
	Balance  uint64 `json:"balance"`   // token balance (as is, not reduced to a floating point value)
	TotalIn  uint64 `json:"totalIn"`   // total incoming token value
	TotalOut uint64 `json:"totalOut"`  // total outgoing token value
}

type Address struct {
	Address  string          `json:"address"`      // address
	ETH      ETH             `json:"ETH"`          // ETH specific information
	Contract *Contract       `json:"contractInfo"` // exists if specified address is a contract
	Token    *Token          `json:"tokenInfo"`    // exists if specified address is a token contract address (same format as token info)
	Tokens   []*AddressToken `json:"tokens"`       // exists if specified address has any token balances
	CountTx  uint64          `json:"countTxs"`     // Total count of incoming and outcoming transactions (including creation one)
}

type Log struct {
	Address string   `json:"address"` // log record address
	Topics  []string `json:"topics"`  // log record topics
	Data    string   `json:"data"`    // log record data
}

type Operation struct {
	Timestamp       int64  `json:"timestamp"`       // operation timestamp
	TransactionHash string `json:"transactionHash"` // transaction hash
	Token           Token  `json:"tokenInfo"`       // token data (same format as token info)
	Type            string `json:"type"`            // operation type (transfer, approve, issuance, mint, burn, etc)
	Address         string `json:"address"`         // operation target address (if one)
	From            string `json:"from"`            // source address (if two addresses involved)
	To              string `json:"to"`              // destination address (if two addresses involved)
	Value           uint64 `json:"value"`           // operation value (as is, not reduced to a floating point value)
}

type Tx struct {
	Hash          string      `json:"hash"`          // transaction hash
	Timestamp     int64       `json:"timestamp"`     // transaction block time
	BlockNumber   uint64      `json:"blockNumber"`   // transaction block number
	Confirmations uint        `json:"confirmations"` // number of confirmations
	Success       bool        `json:"success"`       // true if there were no errors during tx execution
	From          string      `json:"from"`          // source address
	To            string      `json:"to"`            // destination address
	Value         uint64      `json:"value"`         // ETH send value
	Input         string      `json:"input"`         // transaction input data (hex)
	GasLimit      uint64      `json:"gasLimit"`      // gas limit set to this transaction
	GasUsed       uint64      `json:"gasUsed"`       // gas used for this transaction
	Logs          []Log       `json:"logs"`          // event logs
	Operations    []Operation `json:"operations"`    // token operations list for this transaction
}
