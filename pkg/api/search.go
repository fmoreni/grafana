package api

import (
	"strconv"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/metrics"
	m "github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/services/search"
)

func Search(c *m.ReqContext) {
	query := c.Query("query")
	tags := c.QueryStrings("tag")
	starred := c.Query("starred")
	limit := c.QueryInt("limit")
	dashboardType := c.Query("type")
	permission := m.PERMISSION_VIEW

	if limit == 0 {
		limit = 1000
	}

	if c.Query("permission") == "Edit" {
		permission = m.PERMISSION_EDIT
	}

	dbids := make([]int64, 0)
	for _, id := range c.QueryStrings("dashboardIds") {
		dashboardId, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			dbids = append(dbids, dashboardId)
		}
	}

	folderIds := make([]int64, 0)
	for _, id := range c.QueryStrings("folderIds") {
		folderId, err := strconv.ParseInt(id, 10, 64)
		if err == nil {
			folderIds = append(folderIds, folderId)
		}
	}

	searchQuery := search.Query{
		Title:        query,
		Tags:         tags,
		SignedInUser: c.SignedInUser,
		Limit:        limit,
		IsStarred:    starred == "true",
		OrgId:        c.OrgId,
		DashboardIds: dbids,
		Type:         dashboardType,
		FolderIds:    folderIds,
		Permission:   permission,
	}

	err := bus.Dispatch(&searchQuery)
	if err != nil {
		c.JsonApiErr(500, "Search failed", err)
		return
	}

	c.TimeRequest(metrics.M_Api_Dashboard_Search)
	c.JSON(200, searchQuery.Result)
}
