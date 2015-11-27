package main

import (
    "net/http"
    "os"
    "encoding/json"
    "io"
    "io/ioutil"
    "fmt"
    "strings"
    "time"

    "github.com/gorilla/mux"
    "github.com/jacobmarshall/go-azure-slack/util/slackwebhook"
    "github.com/jacobmarshall/go-azure-slack/util/azurewebhook"
)

func WebhookHandler (w http.ResponseWriter, r *http.Request) {
    var err error
    var body []byte
    var azurePayload azurewebhook.Payload

    // Attempt to read the entire body into memory
    body, err = ioutil.ReadAll(r.Body)

    if err != nil {
        // Send a 500 Internal Server Error (unable to read body)
        // Note: This could also be a Bad Request - but since it's unknown, 500 works fine
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        return
    }

    // Attempt to parse the body as JSON (use AzureWebhookPayload struct)
    err = json.Unmarshal(body, &azurePayload)

    if err != nil {
        // Send a 400 Bad Request as the payload was obviously not valid JSON
        http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
        return
    }

    // Send a slack message with the associated azure payload infomation
    go SendSlackMessage(azurePayload)

    // Write a 200 OK response, as payload parsed properly
    io.WriteString(w, http.StatusText(http.StatusOK))
}

func SendSlackMessage (azurePayload azurewebhook.Payload) {
    var slackPayload slackwebhook.Payload

    // Construct a slack webhook payload
    var color, prefix string
    var occurence time.Time
    if azurePayload.Status == "Activated" {
        color = "danger"
        prefix = "Alert: "
    } else {
        color = "good"
        prefix = "Resolved: "
    }
    occurence, err := time.Parse(time.RFC3339Nano, azurePayload.Context.Timestamp)
    slackPayload = slackwebhook.Payload{
        Attachments: []slackwebhook.Attachment{
            slackwebhook.Attachment{
                Color: color,
                Title: prefix + azurePayload.Context.Name,
                TitleLink: azurePayload.Context.PortalLink,
                Text: azurePayload.Context.Description +
                    "\n\n" +
                    "No. of " + strings.ToLower(azurePayload.Context.Condition.MetricName) +
                    " exceded " + azurePayload.Context.Condition.Threshold +
                    " over a period of " + azurePayload.Context.Condition.WindowSize +
                    " minutes (" + azurePayload.Context.Condition.MetricValue +
                    " " + azurePayload.Context.Condition.MetricName + ")",
                Fields: []slackwebhook.Field{
                    slackwebhook.Field{
                        Title: "Site",
                        Value: azurePayload.Context.ResourceName,
                        Short: true,
                    },
                    slackwebhook.Field{
                        Title: "Region",
                        Value: azurePayload.Context.ResourceRegion,
                        Short: true,
                    },
                    slackwebhook.Field{
                        Title: "Time",
                        Value: occurence.Format(time.RFC850),
                    },
                },
            },
        },
    }

    // Create a byte array of the JSON stringified slack payload
    json, _ := json.Marshal(slackPayload)
    reader := strings.NewReader(string(json))

    // Create a new HTTP client & dispatch a request to the webhook URL
    client := http.Client{}
    req, err := http.NewRequest("POST", os.Getenv("WEBHOOK_URL"), reader)
    client.Do(req)

    // If there was an error panic!
    if err != nil {
        panic(err)
    }
}

func main () {
    router := mux.NewRouter()

    // Register a webhook endpoint for receiving azure webhooks
    router.
        HandleFunc("/webhook", WebhookHandler).
        Methods("POST")

    // Register the mux router against the root of the webserver
    http.Handle("/", router)

    // Listen on the port provided by the caller (environment variable)
    fmt.Printf("Listening for connections on :%s\n", os.Getenv("PORT"))
    http.ListenAndServe(":" + os.Getenv("PORT"), nil)
}