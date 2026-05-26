package simulationreport

import (
	"bos/di"
	"bos/enums"
	"bos/types"
	"bos/utils"
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func (m *Model) Init() tea.Cmd {
	if m.payload.Draft == nil {
		return nil
	}
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg { return startMsg{} })
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startMsg:
		if m.payload.Draft == nil || m.running {
			return m, nil
		}
		m.running = true
		draft := *m.payload.Draft
		return m, func() tea.Msg {
			report, err := runSimulation(draft)
			return completedMsg{report: reportValue(report), err: err}
		}

	case completedMsg:
		m.running = false
		if msg.err != nil {
			return m, func() tea.Msg {
				return types.NavigateMsg{
					Screen: enums.ScreenError,
					Payload: types.ErrorPayload{
						Title:   "Simulation Failed",
						Message: msg.err.Error(),
						Return:  enums.ScreenDashboard,
					},
				}
			}
		}
		m.payload.Report = msg.report
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenDashboard} }
		case "esc":
			return m, func() tea.Msg { return types.NavigateMsg{Screen: enums.ScreenDashboard} }
		}
	}

	return m, nil
}

func runSimulation(draft types.TxDraft) (*types.SimulationReport, error) {
	simulator := di.GetSimulator()
	if simulator == nil {
		return nil, fmt.Errorf("simulator is not configured")
	}

	network := di.GetNetwork()
	if network == nil {
		return nil, fmt.Errorf("network provider is not configured")
	}

	caller := common.HexToAddress(draft.FromAddress)
	receiver := common.HexToAddress(draft.RecipientAddress)

	session, err := simulator.BeginSimulation(network.Network(), receiver, caller)
	if err != nil {
		return nil, err
	}
	enrichSession(session, draft)

	signedTx, err := signSimulationTransaction(*session, draft, caller, receiver)
	if err != nil {
		return nil, err
	}

	rawTx, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("encode signed transaction: %w", err)
	}

	report, err := simulator.ExecuteSignedTransaction(*session, rawTx)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func signSimulationTransaction(session types.SimulationSession, draft types.TxDraft, caller common.Address, receiver common.Address) (*coretypes.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, session.RPC)
	if err != nil {
		return nil, fmt.Errorf("connect to simulator RPC: %w", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("read simulator chain ID: %w", err)
	}
	if session.ChainID != "" {
		chainID, err = parseChainID(session.ChainID)
		if err != nil {
			return nil, fmt.Errorf("parse simulator chain ID %q: %w", session.ChainID, err)
		}
	}

	value, err := utils.ParseEtherToWei(draft.Amount)
	if err != nil {
		return nil, fmt.Errorf("parse simulation amount: %w", err)
	}

	nonce, err := client.PendingNonceAt(ctx, caller)
	if err != nil {
		return nil, fmt.Errorf("read simulator nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("read simulator gas price: %w", err)
	}

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  caller,
		To:    &receiver,
		Value: value,
	})
	if err != nil || gasLimit == 0 {
		gasLimit = 21000
	}

	tx := coretypes.NewTransaction(
		nonce,
		receiver,
		value,
		gasLimit,
		gasPrice,
		nil,
	)

	wallet := di.GetWallet()
	if wallet == nil {
		return nil, fmt.Errorf("wallet is not configured")
	}

	signedTx, err := wallet.SignTransaction(ctx, tx, chainID)
	if err != nil {
		return nil, fmt.Errorf("sign simulation transaction: %w", err)
	}
	if signedTx == nil {
		return nil, fmt.Errorf("wallet returned an empty signed transaction")
	}

	return signedTx, nil
}

func enrichSession(session *types.SimulationSession, draft types.TxDraft) {
	session.Amount = draft.Amount
	session.Asset = draft.Asset.Symbol
	session.BalanceBefore = draft.Asset.Balance
	session.BalanceAfter = pendingBalanceAfter(draft)
	session.AddressContract = false
}

func reportValue(report *types.SimulationReport) types.SimulationReport {
	if report == nil {
		return types.SimulationReport{}
	}
	return *report
}

func parseChainID(value string) (*big.Int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, fmt.Errorf("empty chain ID")
	}

	base := 10
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		base = 16
		value = value[2:]
	}

	chainID, ok := new(big.Int).SetString(value, base)
	if !ok {
		return nil, fmt.Errorf("invalid chain ID")
	}
	return chainID, nil
}
