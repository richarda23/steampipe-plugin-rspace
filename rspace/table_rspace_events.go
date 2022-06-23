package rspace

import (
	"context"
	"strconv"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)

func tableRSpaceEventListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "domain", Require: plugin.Optional},
		{Name: "action", Require: plugin.Optional},
		{Name: "timestamp", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
	}
}

func tableRSpaceEvent() *plugin.Table {
	return &plugin.Table{
		Name:        "rspace_event",
		Description: "RSpace events listing",
		List: &plugin.ListConfig{
			Hydrate:    listEvent,
			KeyColumns: tableRSpaceEventListKeyColumns(),
		},
		Columns: []*plugin.Column{
			{Name: "username", Type: proto.ColumnType_STRING, Description: "Username of person who performed event"},
			{Name: "full_name", Transform: transform.FromCamel(), Type: proto.ColumnType_STRING, Description: "Full name of person who performed event"},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Event domain"},
			{Name: "action", Type: proto.ColumnType_STRING, Description: "Event action"},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of event"},
			{Name: "payload", Type: proto.ColumnType_JSON, Description: "Values of custom fields in the event."},
		},
	}
}

func listEvent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx)
	if err != nil {
		logger.Warn("couldn't connect to RSpace")
		return nil, err
	}
	builder := rspace.ActivityQueryBuilder{}
	equalQuals := d.KeyColumnQuals

	if equalQuals["domain"] != nil {
		val := equalQuals["domain"].GetStringValue()
		builder.Domain(val)
	}

	q, _ := builder.Build()

	cfg := rspace.NewRecordListingConfig()
	cfg.PageSize = 100
	docList, err := conn.Activities(q, cfg)
	if err != nil {
		return nil, err
	}
	logger.Warn("There are " + strconv.Itoa(len(docList.Activities)) + " activities")

	for _, t := range docList.Activities {
		d.StreamListItem(ctx, t)
	}
	return nil, nil
}
