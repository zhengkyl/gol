package repo

import (
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemDelegate struct{}

// Render renders the item's view.
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {

}

// Height is the height of the list item.
func (d ItemDelegate) Height() int {
	return 10
}

// Spacing is the size of the horizontal gap between list items in cells.
func (d ItemDelegate) Spacing() int {
	return 1
}

// Update is the update loop for items. All messages in the list's update
// loop will pass through here except when the user is setting a filter.
// Use this method to perform item-level updates appropriate to this
// delegate.
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {

	// var cmds []tea.Cmd

	// for _, item := range m.Items() {
	// 	repo := item.(RepoModel)
	// 	// _, cmd := repo.Update(msg)
	// }
	return nil
}

// ctx := context.Background()
// ts := oauth2.StaticTokenSource(
// 	&oauth2.Token{AccessToken: ""},
// )
// tc := oauth2.NewClient(ctx, ts)

// client := github.NewClient(tc)

// m.client = client

// repos, _, err := client.Repositories.List(ctx, "", nil)
// client.Repositories.ListLanguages()

// return func() tea.Msg {
// 	repos, err := m.getRepos("zhengkyl")
// 	if err != nil {
// 		return err
// 	}

// 	sort.Slice(repos, func(i, j int) bool {
// 		return repos[i].Updated_at.Before(repos[j].Updated_at)
// 	})

// 	return repos[:5]

// 	// for _, repo := range repos {
// 	// 	m.getRepoLangs("zhengkyl", repo.Name)
// 	// }

// 	// return repos
// }
