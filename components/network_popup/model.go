package networksPopup

import (
	"log"
	"strconv"
	"strings"

	"bos/components/search"
	"bos/di"
	"bos/types"

	tea "github.com/charmbracelet/bubbletea"
)

type focusArea int

const (
	focusSearch focusArea = iota
	focusTable
)

type Model struct {
	Visible bool

	search *search.Model
	focus  focusArea

	networks *[]types.Network
	filtered *[]types.Network
	selected int
	offset   int
}

func (m *Model) Init() tea.Cmd {
	if m.search == nil {
		return nil
	}
	return m.search.Init()
}

func New() *Model {
	networks, err := di.GetNetwork().Networks()
	if err != nil {
		log.Printf("no networks: %s", err)
	}

	return &Model{
		Visible:  false,
		search:   search.New("Search networks", 64),
		focus:    focusSearch,
		networks: networks,
		filtered: networks,
	}
}

func (m *Model) focusSearch() {
	m.focus = focusSearch
	m.search.Focus()
}

func (m *Model) focusTable() {
	m.focus = focusTable
	m.search.Blur()
}

func (m *Model) applySearch(query string) {
	m.filterNetworks(query)
	m.focusTable()
}

func (m *Model) filterNetworks(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		m.filtered = m.networks
	} else {
		filtered := make([]types.Network, 0, len(*m.networks))
		for _, network := range *m.networks {
			chainId := strconv.FormatInt(network.Chain.Int64(), 10)
			if strings.Contains(strings.ToLower(*network.Name), query) ||
				strings.Contains(strings.ToLower(*network.Rpc), query) ||
				strings.Contains(strings.ToLower(*network.Symbol), query) ||
				strings.Contains(strings.ToLower(chainId), query) {
				filtered = append(filtered, network)
			}
		}
		m.filtered = &filtered
	}

	m.selected = 0
	m.offset = 0
}

func (m *Model) submit() tea.Cmd {
	return func() tea.Msg {
		if m.filtered == nil || len(*m.filtered) == 0 {
			return nil
		}

		selected := (*m.filtered)[m.selected]
		di.GetNetwork().Change(&selected)
		return SubmittedMsg{Network: selected}
	}
}
