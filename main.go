package main

import (
    "net/http"
    "os"
    "encoding/json"
    "io"
    "io/ioutil"
    "strings"

    "github.com/gorilla/mux"
    "fmt"
)

type AzureWebhookPayload struct {
    Status string `json:"status"`
    Context struct {
        Timestamp string `json:"timestamp"`
        ID string `json:"id"`
        Name string `json:"name"`
        Description string `json:"description"`
        ConditionType string `json:"conditionType"`
        Condition struct {
            MetricName string `json:"metricName"`
            MetricUnit string `json:"metricUnit"`
            MetricValue string `json:"metricValue"`
            Threshold string `json:"threshold"`
            WindowSize string `json:"windowSize"`
            TimeAggregation string `json:"timeAggregation"`
            Operator string `json:"operator"`
        } `json:"condition"`
        SubscriptionID string `json:"subscriptionId"`
        ResourceGroupName string `json:"resourceGroupName"`
        ResourceName string `json:"resourceName"`
        ResourceType string `json:"resourceType"`
        ResourceId string `json:"resourceId"`
        ResourceRegion string `json:"resourceRegion"`
        PortalLink string `json:"portalLink"`
    } `json:"context"`
    Properties map[string]string `json:"properties"`
}

type SlackWebhookPayload struct {
    Text string `json:"text"`
}

func WebhookHandler (w http.ResponseWriter, r *http.Request) {
    var err error
    var body []byte
    var azurePayload AzureWebhookPayload

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

func SendSlackMessage (azurePayload AzureWebhookPayload) {
    var slackPayload SlackWebhookPayload

    // Construct a slack webhook payload
    // TODO Constructing an actual message (containing event message, etc...)
    slackPayload = SlackWebhookPayload{
        Text: `
            Hello World
        `,
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
