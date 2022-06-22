package rspace

import (
    "context"
    "os"
"strconv"
"fmt"
"net/url"

    "github.com/richarda23/rspace-client-go/rspace"

    "github.com/turbot/steampipe-plugin-sdk/grpc/proto"
    "github.com/turbot/steampipe-plugin-sdk/plugin"
)
const (
	APIKEY_ENV_NAME         = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME       = "RSPACE_URL"
)

func tableRSpaceEvent() *plugin.Table {
    return &plugin.Table{
        Name:        "rspace_event",
        Description: "RSpace events listing",
        List: &plugin.ListConfig{
            Hydrate: listEvent,
        },
        Columns: []*plugin.Column{
            {Name: "username", Type: proto.ColumnType_STRING, Description: "Username of person who performed event"},
            {Name: "domain", Type: proto.ColumnType_STRING, Description: "Event domain"},
            {Name: "action", Type: proto.ColumnType_STRING, Description: "Event action"},
            {Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp of event"},
            {Name: "payload", Type: proto.ColumnType_JSON, Description: "Values of custom fields in the event."},
        },
    }
}
func connect(ctx context.Context) (*rspace.RsWebClient, error) {
logger := plugin.Logger(ctx)
logger.Warn("Querying events API")
url, _ := url.Parse("https://pangolin8086.researchspace.com/api/v1")
	fmt.Println("url is " + url.String())
	apikey := "abcdefghijklmnop12"
	fmt.Println("api is " + apikey)
	webClient := rspace.NewWebClientCustomTimeout(url, apikey, 30)
        return webClient, nil
}
func listEvent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
logger := plugin.Logger(ctx)
    conn, err := connect(ctx)
    if err != nil {
       
        logger.Warn("couldn't connect to RSpace")
        return nil, err
    }
builder := rspace.ActivityQueryBuilder{}
builder.Domain("RECORD")
q,_ := builder.Build()

cfg := rspace.NewRecordListingConfig()
cfg.PageSize = 100
        docList, err := conn.Activities(q, cfg)
        if err != nil {
            return nil, err
        }
        logger.Warn("There are " + strconv.Itoa(len(docList.Activities)) + " activites")
        for _, t := range docList.Activities {
            d.StreamListItem(ctx, t)
        }
    return nil, nil
}

func getenv(envname string) string {
	return os.Getenv(envname)
}
