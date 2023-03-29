package repo

import (
	"encoding/json"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/gtg/ui/common"
)

type RepoModel struct {
	common common.Common
	client *http.Client
}

func New(c common.Common, client *http.Client) *RepoModel {
	return &RepoModel{
		common: c,
		client: client,
	}
}

func (m *RepoModel) Init() tea.Cmd {
	return nil
}

func (m *RepoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *RepoModel) View() string {
	return "ui model"
}

// func (m *Model) getRepos(user string) ([]Repo, error) {
// 	resp, err := m.client.Get(fmt.Sprintf("https://api.github.com/users/%s/repos", user))

// 	var repos []Repo

// 	if err != nil {
// 		return repos, err
// 	}

// 	err = json.NewDecoder(resp.Body).Decode(&repos)

// 	return repos, err
// }

func (m *RepoModel) getRepoLangs(user string, repo string) (map[string]int, error) {
	resp, err := m.client.Get(fmt.Sprintf("https://api.github.com/repos/%s/%s/languages", user, repo))

	langs := make(map[string]int)

	if err != nil {
		return langs, err
	}

	err = json.NewDecoder(resp.Body).Decode(&langs)

	return langs, err
}

// type Repo struct {
// 	Name        string
// 	Description string
// 	Language    string
// 	Updated_at  time.Time
// }
