package rspace

import (
	"context"
	"fmt"
	"strconv"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type SPDocInfo struct {
	GlobalId, Name, Created, LastModified, OwnerUsername, Tags string
	Signed                                                     bool
}

func tableRSpaceDocumentListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "name", Require: plugin.Optional},
		{Name: "tags", Require: plugin.Optional},
		{Name: "owner_username", Require: plugin.Optional},
	}
}

func tableRSpaceDocument() *plugin.Table {
	return &plugin.Table{
		Name:        "rspace_document",
		Description: "RSpace documents listing",
		List: &plugin.ListConfig{
			Hydrate:    listDocument,
			KeyColumns: tableRSpaceDocumentListKeyColumns(),
		},
		Get: &plugin.GetConfig{
			Hydrate:    getDocument,
			KeyColumns: plugin.SingleColumn("global_id"),
		},
		Columns: []*plugin.Column{
			{Name: "global_id", Transform: transform.FromCamel(), Type: proto.ColumnType_STRING, Description: "Global Id"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of Document"},
			{Name: "created", Transform: transform.FromCamel(), Type: proto.ColumnType_TIMESTAMP,
				Description: "Creation time of document"},
			{Name: "last_modified", Transform: transform.FromCamel(),
				Type: proto.ColumnType_TIMESTAMP, Description: "Last modified time of document"},
			{Name: "owner_username", Transform: transform.FromCamel(), Type: proto.ColumnType_STRING, Description: "Full name of owner"},
			{Name: "tags", Type: proto.ColumnType_STRING, Description: "Comma separated list of tags"},
			{Name: "signed", Type: proto.ColumnType_BOOL, Description: "Whether the document is signed or not"},
		},
	}
}

// getDocument retrieves a single document by its global ID
func getDocument(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx)
	if err != nil {
		logger.Warn("couldn't connect to RSpace")
		return nil, err
	}
	idStr := d.KeyColumnQuals["global_id"].GetStringValue()
	id, err := getIdFromGlobalId(idStr)
	if err != nil {
		return nil, err
	}
	logger.Debug("Parsed id as ", "id", id)

	doc, err := conn.DocumentById(id)

	if err != nil {
		return nil, err
	}
	logger.Debug("retrieved doc as ", "id", doc.Id)

	mappedDoc := SPDocInfo{doc.GlobalId, doc.Name, doc.Created,
		doc.LastModified, doc.UserInfo.Username, doc.Tags, doc.Signed}

	return mappedDoc, nil

}

func listDocument(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx)
	if err != nil {
		logger.Warn("couldn't connect to RSpace")
		return nil, err
	}
	q := buildQuery(d)
	logger.Warn("", "query", q)
	cfg := rspace.NewRecordListingConfig()
	cfg.PageSize = 100
	limit := getLimit(d)
	paginations, _ := calculatePageSizes(limit, HARD_LIMIT, 100)
	for i, v := range paginations {
		logger.Info("Retrieving pages", "page", i, "pageSize", v)
		cfg.PageSize = v
		docList, err := conn.AdvancedSearchDocuments(cfg, q)
		if err != nil {
			return nil, err
		}
		logger.Warn("There are " + strconv.Itoa(len(docList.Documents)) + " documents")
		for _, t := range docList.Documents {
			logger.Warn(fmt.Sprintf("id=%s and name=%s", t.GlobalId, t.Name))
		}
		for _, t := range docList.Documents {
			mappedDoc := SPDocInfo{t.GlobalId, t.Name, t.Created,
				t.LastModified, t.UserInfo.Username, t.Tags, t.Signed}
			d.StreamListItem(ctx, mappedDoc)
		}
		links := docList.Links
		if listingHasNextPage(links) {
			cfg.PageNumber = cfg.PageNumber + 1
		} else {
			break
		}
	}
	return nil, nil
}

func buildQuery(d *plugin.QueryData) *rspace.SearchQuery {
	builder := &rspace.SearchQueryBuilder{}
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		val := equalQuals["name"].GetStringValue()
		builder.AddTerm(val, rspace.NAME)
	}
	if equalQuals["tags"] != nil {
		val := equalQuals["tags"].GetStringValue()
		builder.AddTerm(val, rspace.TAG)
	}
	if equalQuals["owner_username"] != nil {
		val := equalQuals["owner_username"].GetStringValue()
		builder.AddTerm(val, rspace.OWNER)
	}
	q := builder.Build()
	return q
}
