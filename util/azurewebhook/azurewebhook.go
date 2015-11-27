package azurewebhook

type Payload struct {
	Status string `json:"status"`
	Context Context `json:"context"`
	Properties map[string]string `json:"properties"`
}

type Context struct {
	Timestamp string `json:"timestamp"`
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	ConditionType string `json:"conditionType"`
	Condition Condition `json:"condition"`
	SubscriptionID string `json:"subscriptionId"`
	ResourceGroupName string `json:"resourceGroupName"`
	ResourceName string `json:"resourceName"`
	ResourceType string `json:"resourceType"`
	ResourceId string `json:"resourceId"`
	ResourceRegion string `json:"resourceRegion"`
	PortalLink string `json:"portalLink"`
}

type Condition struct {
	MetricName string `json:"metricName"`
	MetricUnit string `json:"metricUnit"`
	MetricValue string `json:"metricValue"`
	Threshold string `json:"threshold"`
	WindowSize string `json:"windowSize"`
	TimeAggregation string `json:"timeAggregation"`
	Operator string `json:"operator"`
}