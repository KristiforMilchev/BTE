package importcontract

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"bos/di"
	"bos/enums"
	"bos/types"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
)

type interactionsLoadedMsg struct {
	interactions []types.ContractInteraction
	err          error
}

type simulationCompletedMsg struct {
	report *types.SimulationReport
	err    error
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case interactionsLoadedMsg:
		m.handleInteractionsLoaded(msg)
		return m, nil

	case simulationCompletedMsg:
		m.simulating = false
		if msg.err != nil {
			m.status = "Simulation failed: " + msg.err.Error()
			return m, nil
		}
		return m, func() tea.Msg {
			return types.NavigateMsg{
				Screen: enums.ScreenSimulationReport,
				Payload: types.SimulationReportPayload{
					Report: reportValue(msg.report),
					Return: "import-contract",
				},
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc", "n", "N":
			return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenDashboard} }
		case "enter":
			return m, m.loadContract()
		case "C":
			m.clearContract()
			return m, textinput.Blink
		case "s", "S":
			return m, m.simulateContract()
		case "y", "Y":
			m.status = "Saving imported token will be wired later"
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) loadContract() tea.Cmd {
	address := strings.TrimSpace(m.input.Value())
	if !common.IsHexAddress(address) {
		m.loaded = false
		m.address = ""
		m.callable = nil
		m.readable = nil
		m.interactions = nil
		m.loadingInteractions = false
		m.status = "Enter a valid contract address"
		return nil
	}

	m.address = common.HexToAddress(address).Hex()
	m.input.SetValue(m.address)
	m.input.Blur()
	m.loaded = true
	m.callable = fakeCallableMethods()
	m.readable = fakeReadableMethods()
	m.interactions = nil
	m.loadingInteractions = false
	m.simulating = false
	if m.reader == nil {
		m.status = "Contract loaded • interaction reader not configured"
		return nil
	}

	m.loadingInteractions = true
	m.status = "Contract loaded • loading 24h interactions"
	contract := common.HexToAddress(m.address)
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		interactions, err := m.reader.RecentInteractions(ctx, contract, time.Now().Add(-24*time.Hour))
		return interactionsLoadedMsg{interactions: interactions, err: err}
	}
}

func (m *Model) clearContract() {
	m.input.SetValue("")
	m.input.Focus()
	m.address = ""
	m.loaded = false
	m.callable = nil
	m.readable = nil
	m.interactions = nil
	m.loadingInteractions = false
	m.simulating = false
	m.status = "No contract address selected"
}

func (m *Model) handleInteractionsLoaded(msg interactionsLoadedMsg) {
	m.loadingInteractions = false
	if msg.err != nil {
		m.interactions = nil
		m.status = "Contract loaded • failed to load 24h interactions: " + msg.err.Error()
		return
	}
	m.interactions = msg.interactions
	if len(m.interactions) == 0 {
		m.status = "Contract loaded • no 24h interactions found"
		return
	}
	m.status = "Contract loaded • 24h interactions loaded"
}

func (m *Model) simulateContract() tea.Cmd {
	if m.simulating {
		return nil
	}
	if !m.loaded || !common.IsHexAddress(m.address) {
		m.status = "Load a valid contract address before simulating"
		return nil
	}

	m.simulating = true
	m.status = "Starting fork simulation • confirm on Ledger"
	address := common.HexToAddress(m.address)
	methods := append([]Method(nil), m.callable...)
	return func() tea.Msg {
		report, err := runContractSimulation(address, methods)
		return simulationCompletedMsg{report: report, err: err}
	}
}

func runContractSimulation(address common.Address, methods []Method) (*types.SimulationReport, error) {
	simulator := di.GetSimulator()
	if simulator == nil {
		return nil, fmt.Errorf("simulator is not configured")
	}

	network := di.GetNetwork()
	if network == nil {
		return nil, fmt.Errorf("network provider is not configured")
	}

	wallet := di.GetWallet()
	if wallet == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	caller, err := wallet.Account()
	if err != nil {
		return nil, fmt.Errorf("read wallet account: %w", err)
	}
	if caller == nil {
		return nil, fmt.Errorf("wallet returned an empty account")
	}

	functions := contractFunctions(methods)
	session, err := simulator.BeginContractSimulation(network.Network(), address, *caller, functions)
	if err != nil {
		return nil, err
	}

	signedTxs, err := signAPITransactions(*session, wallet)
	if err != nil {
		return nil, err
	}

	report, err := simulator.ExecuteSignedTransactions(*session, signedTxs)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func contractFunctions(methods []Method) []types.ContractFunction {
	functions := make([]types.ContractFunction, 0, len(methods))
	for _, method := range methods {
		functions = append(functions, types.ContractFunction{
			Name:       method.Name,
			Signature:  method.Signature,
			Mutability: "nonpayable",
		})
	}
	return functions
}

func signAPITransactions(session types.SimulationSession, wallet interface {
	SignTransaction(context.Context, *coretypes.Transaction, *big.Int) (*coretypes.Transaction, error)
}) ([]types.SignedSimulationTransaction, error) {
	transactions := session.Transactions
	if len(transactions) == 0 && session.Transaction.To != "" {
		transactions = []types.LedgerTransaction{session.Transaction}
	}
	if len(transactions) == 0 {
		return nil, fmt.Errorf("simulator did not return any transactions to sign")
	}

	signed := make([]types.SignedSimulationTransaction, 0, len(transactions))
	for _, txData := range transactions {
		txSession := session
		txSession.Transaction = txData
		signedTx, err := signAPITransaction(txSession, wallet)
		if err != nil {
			return nil, err
		}
		rawTx, err := signedTx.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("encode signed transaction: %w", err)
		}
		signed = append(signed, types.SignedSimulationTransaction{
			Function:          txData.Function,
			Signature:         txData.Signature,
			SignedTransaction: "0x" + hex.EncodeToString(rawTx),
		})
	}
	return signed, nil
}

func signAPITransaction(session types.SimulationSession, wallet interface {
	SignTransaction(context.Context, *coretypes.Transaction, *big.Int) (*coretypes.Transaction, error)
}) (*coretypes.Transaction, error) {
	txData := session.Transaction
	if txData.To == "" {
		return nil, fmt.Errorf("simulator did not return a transaction recipient")
	}
	if !common.IsHexAddress(txData.To) {
		return nil, fmt.Errorf("simulator returned invalid transaction recipient %q", txData.To)
	}

	chainIDValue := firstNonEmpty(txData.ChainID, session.ChainID)
	chainID, err := parseQuantityBig(chainIDValue)
	if err != nil {
		return nil, fmt.Errorf("parse simulator chain ID %q: %w", chainIDValue, err)
	}

	nonce, err := parseQuantityUint64(txData.Nonce)
	if err != nil {
		return nil, fmt.Errorf("parse simulator nonce %q: %w", txData.Nonce, err)
	}

	value, err := parseQuantityBig(firstNonEmpty(txData.Value, "0x0"))
	if err != nil {
		return nil, fmt.Errorf("parse simulator value %q: %w", txData.Value, err)
	}

	gasLimit, err := parseQuantityUint64(firstNonEmpty(txData.Gas, "0x5208"))
	if err != nil {
		return nil, fmt.Errorf("parse simulator gas %q: %w", txData.Gas, err)
	}

	gasPrice, err := parseQuantityBig(firstNonEmpty(txData.GasPrice, "0x0"))
	if err != nil {
		return nil, fmt.Errorf("parse simulator gas price %q: %w", txData.GasPrice, err)
	}

	tx := coretypes.NewTransaction(
		nonce,
		common.HexToAddress(txData.To),
		value,
		gasLimit,
		gasPrice,
		common.FromHex(firstNonEmpty(txData.Data, "0x")),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	signedTx, err := wallet.SignTransaction(ctx, tx, chainID)
	if err != nil {
		return nil, fmt.Errorf("sign simulation transaction: %w", err)
	}
	if signedTx == nil {
		return nil, fmt.Errorf("wallet returned an empty signed transaction")
	}
	return signedTx, nil
}

func parseQuantityBig(value string) (*big.Int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("empty value")
	}

	base := 10
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		base = 16
		value = value[2:]
	}
	if value == "" {
		return big.NewInt(0), nil
	}

	parsed, ok := new(big.Int).SetString(value, base)
	if !ok {
		return nil, fmt.Errorf("invalid quantity")
	}
	return parsed, nil
}

func parseQuantityUint64(value string) (uint64, error) {
	parsed, err := parseQuantityBig(value)
	if err != nil {
		return 0, err
	}
	if !parsed.IsUint64() {
		return 0, fmt.Errorf("quantity overflows uint64")
	}
	return parsed.Uint64(), nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func reportValue(report *types.SimulationReport) types.SimulationReport {
	if report == nil {
		return types.SimulationReport{}
	}
	return *report
}
