package Repossection

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dlvhdr/gh-dash/config"
	"github.com/dlvhdr/gh-dash/data"
	"github.com/dlvhdr/gh-dash/ui/components/repo"
	"github.com/dlvhdr/gh-dash/ui/components/section"
	"github.com/dlvhdr/gh-dash/ui/components/table"
	"github.com/dlvhdr/gh-dash/ui/constants"
	"github.com/dlvhdr/gh-dash/ui/context"
	"github.com/dlvhdr/gh-dash/ui/keys"
	"github.com/dlvhdr/gh-dash/utils"
)

const SectionType = "repo"

type Model struct {
	section.Model
	Repos []data.RepositoryData
}

func NewModel(id int, ctx *context.ProgramContext, cfg config.ReposSectionConfig, lastUpdated time.Time) Model {
	m := Model{
		section.NewModel(
			id,
			ctx,
			cfg.ToSectionConfig(),
			SectionType,
			GetSectionColumns(cfg, ctx),
			"Owner",
			"Name",
			lastUpdated,
		),
		[]data.RepositoryData{},
	}

	return m
}

func (m Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmd tea.Cmd
	var err error

	switch msg := msg.(type) {

	case tea.KeyMsg:

		if m.IsSearchFocused() {
			switch {

			case msg.Type == tea.KeyCtrlC, msg.Type == tea.KeyEsc:
				m.SearchBar.SetValue(m.SearchValue)
				blinkCmd := m.SetIsSearching(false)
				return &m, blinkCmd

			case msg.Type == tea.KeyEnter:
				m.SearchValue = m.SearchBar.Value()
				m.SetIsSearching(false)
				m.ResetRows()
				return &m, tea.Batch(m.FetchNextPageSectionRows()...)
			}

			break
		}

		// switch {

		// case key.Matches(msg, keys.PRKeys.Diff):
		// 	cmd = m.diff()

		// case key.Matches(msg, keys.PRKeys.Checkout):
		// 	cmd, err = m.checkout()
		// 	if err != nil {
		// 		m.Ctx.Error = err
		// 	}

		// case key.Matches(msg, keys.PRKeys.Close):
		// 	cmd = m.close()

		// case key.Matches(msg, keys.PRKeys.Ready):
		// 	cmd = m.ready()

		// case key.Matches(msg, keys.PRKeys.Merge):
		// 	cmd = m.merge()

		// case key.Matches(msg, keys.PRKeys.Reopen):
		// 	cmd = m.reopen()

		// }

	// case UpdatePRMsg:
	// 	for i, currPr := range m.Repos {
	// 		if currPr.Number == msg.PrNumber {
	// 			if msg.IsClosed != nil {
	// 				if *msg.IsClosed == true {
	// 					currPr.State = "CLOSED"
	// 				} else {
	// 					currPr.State = "OPEN"
	// 				}
	// 			}
	// 			if msg.NewComment != nil {
	// 				currPr.Comments.Nodes = append(currPr.Comments.Nodes, *msg.NewComment)
	// 			}
	// 			if msg.ReadyForReview != nil && *msg.ReadyForReview {
	// 				currPr.IsDraft = false
	// 			}
	// 			if msg.IsMerged != nil && *msg.IsMerged {
	// 				currPr.State = "MERGED"
	// 				currPr.Mergeable = ""
	// 			}
	// 			m.Repos[i] = currPr
	// 			m.Table.SetRows(m.BuildRows())
	// 			break
	// 		}
	// 	}

	case SectionRepositoriesFetchedMsg:
		if m.PageInfo != nil {
			m.Repos = append(m.Repos, msg.Repos...)
		} else {
			m.Repos = msg.Repos
		}
		m.TotalCount = msg.TotalCount
		m.PageInfo = &msg.PageInfo
		m.Table.SetRows(m.BuildRows())
		m.UpdateLastUpdated(time.Now())
		m.UpdateTotalItemsCount(m.TotalCount)
	}

	search, searchCmd := m.SearchBar.Update(msg)
	m.SearchBar = search
	return &m, tea.Batch(cmd, searchCmd)
}

func GetSectionColumns(cfg config.ReposSectionConfig, ctx *context.ProgramContext) []table.Column {
	dLayout := ctx.Config.Defaults.Layout.Repos
	sLayout := cfg.Layout

	ownerLayout := config.MergeColumnConfigs(dLayout.Owner, sLayout.Owner)
	// updatedAtLayout := config.MergeColumnConfigs(dLayout.UpdatedAt, sLayout.UpdatedAt)
	repoLayout := config.MergeColumnConfigs(dLayout.Name, sLayout.Name)
	// titleLayout := config.MergeColumnConfigs(dLayout.Title, sLayout.Title)
	// authorLayout := config.MergeColumnConfigs(dLayout.Author, sLayout.Author)
	// assigneesLayout := config.MergeColumnConfigs(dLayout.Assignees, sLayout.Assignees)
	// reviewStatusLayout := config.MergeColumnConfigs(dLayout.ReviewStatus, sLayout.ReviewStatus)
	// stateLayout := config.MergeColumnConfigs(dLayout.State, sLayout.State)
	// ciLayout := config.MergeColumnConfigs(dLayout.Ci, sLayout.Ci)
	// linesLayout := config.MergeColumnConfigs(dLayout.Lines, sLayout.Lines)

	return []table.Column{
		{
			Title:  "Owner",
			Width:  ownerLayout.Width,
			Hidden: ownerLayout.Hidden,
		},
		{
			Title:  "Name",
			Width:  repoLayout.Width,
			Hidden: repoLayout.Hidden,
		},
		// {
		// 	Title:  "",
		// 	Hidden: stateLayout.Hidden,
		// },
		// {
		// 	Title:  "",
		// 	Width:  repoLayout.Width,
		// 	Hidden: repoLayout.Hidden,
		// },
		// {
		// 	Title:  "Title",
		// 	Grow:   utils.BoolPtr(true),
		// 	Hidden: titleLayout.Hidden,
		// },
		// {
		// 	Title:  "Author",
		// 	Width:  authorLayout.Width,
		// 	Hidden: authorLayout.Hidden,
		// },
		// {
		// 	Title:  "Assignees",
		// 	Width:  assigneesLayout.Width,
		// 	Hidden: assigneesLayout.Hidden,
		// },
		// {
		// 	Title:  "",
		// 	Width:  utils.IntPtr(4),
		// 	Hidden: reviewStatusLayout.Hidden,
		// },
		// {
		// 	Title:  "",
		// 	Width:  &ctx.Styles.Reposection.CiCellWidth,
		// 	Grow:   new(bool),
		// 	Hidden: ciLayout.Hidden,
		// },
		// {
		// 	Title:  "",
		// 	Width:  linesLayout.Width,
		// 	Hidden: linesLayout.Hidden,
		// },
	}
}

func (m *Model) BuildRows() []table.Row {
	var rows []table.Row
	for _, currPr := range m.Repos {
		prModel := repo.Repository{Ctx: m.Ctx, Data: currPr}
		rows = append(rows, prModel.ToTableRow())
	}

	if rows == nil {
		rows = []table.Row{}
	}

	return rows
}

func (m *Model) NumRows() int {
	return len(m.Repos)
}

type SectionRepositoriesFetchedMsg struct {
	Repos      []data.RepositoryData
	TotalCount int
	PageInfo   data.PageInfo
}

func (m *Model) GetCurrRow() data.RowData {
	if len(m.Repos) == 0 {
		return nil
	}
	repo := m.Repos[m.Table.GetCurrItem()]
	return &repo
}

func (m *Model) FetchNextPageSectionRows() []tea.Cmd {
	if m == nil {
		return nil
	}

	if m.PageInfo != nil && !m.PageInfo.HasNextPage {
		return nil
	}

	var cmds []tea.Cmd

	startCursor := time.Now().String()
	if m.PageInfo != nil {
		startCursor = m.PageInfo.StartCursor
	}
	taskId := fmt.Sprintf("fetching_Repos_%d_%s", m.Id, startCursor)
	task := context.Task{
		Id:           taskId,
		StartText:    fmt.Sprintf(`Fetching Repos for "%s"`, m.Config.Title),
		FinishedText: fmt.Sprintf(`Repos for "%s" have been fetched`, m.Config.Title),
		State:        context.TaskStart,
		Error:        nil,
	}
	startCmd := m.Ctx.StartTask(task)
	cmds = append(cmds, startCmd)

	fetchCmd := func() tea.Msg {
		limit := m.Config.Limit
		if limit == nil {
			limit = &m.Ctx.Config.Defaults.ReposLimit
		}
		res, err := data.FetchRepositories(m.GetFilters(), *limit, m.PageInfo)
		if err != nil {
			log.Printf("err %v", err)
			return constants.TaskFinishedMsg{
				SectionId:   m.Id,
				SectionType: m.Type,
				TaskId:      taskId,
				Err:         err,
			}
		}

		return constants.TaskFinishedMsg{
			SectionId:   m.Id,
			SectionType: m.Type,
			TaskId:      taskId,
			Msg: SectionRepositoriesFetchedMsg{
				Repos:      res.Repos,
				TotalCount: res.TotalCount,
				PageInfo:   res.PageInfo,
			},
		}
	}
	cmds = append(cmds, fetchCmd)

	return cmds
}

func (m *Model) ResetRows() {
	m.Repos = nil
	m.Table.Rows = nil
	m.ResetPageInfo()
	m.Table.ResetCurrItem()
}

func FetchAllSections(ctx context.ProgramContext) (sections []section.Section, fetchAllCmd tea.Cmd) {
	fetchReposCmds := make([]tea.Cmd, 0, len(ctx.Config.ReposSections))
	sections = make([]section.Section, 0, len(ctx.Config.ReposSections))
	for i, sectionConfig := range ctx.Config.ReposSections {
		sectionModel := NewModel(i+1, &ctx, sectionConfig, time.Now()) // 0 is the search section
		sections = append(sections, &sectionModel)
		fetchReposCmds = append(fetchReposCmds, sectionModel.FetchNextPageSectionRows()...)
	}
	return sections, tea.Batch(fetchReposCmds...)
}

type UpdatePRMsg struct {
	PrNumber       int
	IsClosed       *bool
	NewComment     *data.Comment
	ReadyForReview *bool
	IsMerged       *bool
}
