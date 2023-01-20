package data

import (
	"fmt"

	"github.com/cli/go-gh"
	graphql "github.com/cli/shurcooL-graphql"
)

type RepositoryData struct {
	Repository Repository
	// Number int
	// Title  string
	// Body   string
	// Author struct {
	// 	Login string
	// }
	// UpdatedAt      time.Time
	// Url            string
	// State          string
	// Mergeable      string
	// ReviewDecision string
	// Additions      int
	// Deletions      int
	// HeadRefName    string
	// BaseRefName    string
	// HeadRepository struct {
	// 	Name string
	// }
	// HeadRef struct {
	// 	Name string
	// }
	// Repository    Repository
	// Assignees     Assignees `graphql:"assignees(first: 3)"`
	// Comments      Comments  `graphql:"comments(last: 5, orderBy: { field: UPDATED_AT, direction: DESC })"`
	// LatestReviews Reviews   `graphql:"latestReviews(last: 3)"`
	// IsDraft       bool
	// Commits       Commits `graphql:"commits(last: 1)"`
}

// type CheckRun struct {
// 	Name       graphql.String
// 	Status     graphql.String
// 	Conclusion graphql.String
// 	CheckSuite struct {
// 		Creator struct {
// 			Login graphql.String
// 		}
// 		WorkflowRun struct {
// 			Workflow struct {
// 				Name graphql.String
// 			}
// 		}
// 	}
// }

// type StatusContext struct {
// 	Context graphql.String
// 	State   graphql.String
// 	Creator struct {
// 		Login graphql.String
// 	}
// }

// type Commits struct {
// 	Nodes []struct {
// 		Commit struct {
// 			Deployments struct {
// 				Nodes []struct {
// 					Task        graphql.String
// 					Description graphql.String
// 				}
// 			} `graphql:"deployments(last: 10)"`
// 			StatusCheckRollup struct {
// 				Contexts struct {
// 					TotalCount graphql.Int
// 					Nodes      []struct {
// 						Typename      graphql.String `graphql:"__typename"`
// 						CheckRun      CheckRun       `graphql:"... on CheckRun"`
// 						StatusContext StatusContext  `graphql:"... on StatusContext"`
// 					}
// 				} `graphql:"contexts(last: 20)"`
// 			}
// 		}
// 	}
// }

// type Comment struct {
// 	Author struct {
// 		Login string
// 	}
// 	Body      string
// 	UpdatedAt time.Time
// }

// type Comments struct {
// 	Nodes      []Comment
// 	TotalCount int
// }

// type Review struct {
// 	Author struct {
// 		Login string
// 	}
// 	Body      string
// 	State     string
// 	UpdatedAt time.Time
// }

// type Reviews struct {
// 	Nodes []Review
// }

// type PageInfo struct {
// 	HasNextPage bool
// 	StartCursor string
// 	EndCursor   string
// }

func (data RepositoryData) GetRepoOwner() string {
	return data.Repository.Owner
}

// func (data RepositoryData) GetNumber() int {
// 	return data.Number
// }

// func (data RepositoryData) GetUrl() string {
// 	return data.Url
// }

// func (data RepositoryData) GetUpdatedAt() time.Time {
// 	return data.UpdatedAt
// }

func makeRepositoriesQuery(query string) string {
	return fmt.Sprintf("is:pr %s sort:updated", query)
}

type RepositoriesResponse struct {
	Prs        []RepositoryData
	TotalCount int
	PageInfo   PageInfo
}

func FetchRepositories(query string, limit int, pageInfo *PageInfo) (RepositoriesResponse, error) {
	var err error
	client, err := gh.GQLClient(nil)

	if err != nil {
		return RepositoriesResponse{}, err
	}

	var queryResult struct {
		Search struct {
			Nodes []struct {
				Repository RepositoryData `graphql:"... on Repository"`
			}
			IssueCount int
			PageInfo   PageInfo
		} `graphql:"search(type: ISSUE, first: $limit, after: $endCursor, query: $query)"`
	}
	var endCursor *string
	if pageInfo != nil {
		endCursor = &pageInfo.EndCursor
	}
	variables := map[string]interface{}{
		"query":     graphql.String(makeRepositoriesQuery(query)),
		"limit":     graphql.Int(limit),
		"endCursor": (*graphql.String)(endCursor),
	}
	err = client.Query("SearchRepositories", &queryResult, variables)
	if err != nil {
		return RepositoriesResponse{}, err
	}

	prs := make([]RepositoryData, 0, len(queryResult.Search.Nodes))
	for _, node := range queryResult.Search.Nodes {
		// if node.Repository.Repository.IsArchived {
		// 	continue
		// }
		prs = append(prs, node.Repository)
	}

	return RepositoriesResponse{
		Prs:        prs,
		TotalCount: queryResult.Search.IssueCount,
		PageInfo:   queryResult.Search.PageInfo,
	}, nil
}
