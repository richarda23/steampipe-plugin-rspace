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

type SPFileInfo struct {
	GlobalId, Name, Created, Caption string
	Size                             int
}

func tableRSpaceFileListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{}
}

func tableRSpaceFile() *plugin.Table {
	return &plugin.Table{
		Name:        "rspace_file",
		Description: "RSpace Files listing",
		List: &plugin.ListConfig{
			Hydrate:    listFile,
			KeyColumns: tableRSpaceFileListKeyColumns(),
		},
		Get: &plugin.GetConfig{
			Hydrate:    getFile,
			KeyColumns: plugin.SingleColumn("global_id"),
		},
		Columns: []*plugin.Column{
			{Name: "global_id", Transform: transform.FromCamel(), Type: proto.ColumnType_STRING, Description: "Global Id"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of File"},
			{Name: "created", Transform: transform.FromCamel(), Type: proto.ColumnType_TIMESTAMP,
				Description: "Creation time of File"},
			{Name: "caption", Type: proto.ColumnType_STRING,
				Description: "Caption of file"},
			{Name: "size", Type: proto.ColumnType_INT, Description: "Full size in bytes"},
		},
	}
}

// getFile retrieves a single File by its global ID
func getFile(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	file, err := conn.FileById(id)

	if err != nil {
		return nil, err
	}
	logger.Debug("retrieved doc as ", "id", file.Id)

	mappedDoc := SPFileInfo{file.GlobalId, file.Name, file.Created,
		file.Caption, file.Size}

	return mappedDoc, nil

}

func listFile(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx)
	if err != nil {
		logger.Warn("couldn't connect to RSpace")
		return nil, err
	}
	cfg := rspace.NewRecordListingConfig()
	cfg.PageSize = 100
	limit := getLimit(d)
	paginations, _ := calculatePageSizes(limit, HARD_LIMIT, 100)
	for i, v := range paginations {
		logger.Info("Retrieving files", "page", i, "pageSize", v)
		cfg.PageSize = v
		docList, err := conn.Files(cfg, "")
		if err != nil {
			return nil, err
		}
		logger.Warn("There are " + strconv.Itoa(len(docList.Files)) + " Files")
		for _, t := range docList.Files {
			logger.Warn(fmt.Sprintf("id=%s and name=%s", t.GlobalId, t.Name))
		}
		for _, t := range docList.Files {
			mappedDoc := SPFileInfo{t.GlobalId, t.Name, t.Created,
				t.Caption, t.Size}
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
