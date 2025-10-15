# governo
This application aims to provide generally usable way to process events from Google Cloud. It's designed to take in Cloud Audit / Cloud Asset Inventory events and then allows you to define rules based on the event payload. It uses CEL under the hood for the condition language. It currently reads the rules from `config.yaml` file and reloads the config when there is changes made to it.

## Actions


## Flow

```mermaid
graph TD
    A[Pub/Sub Event] --> B(HTTP Server);
    B -- JSON Payload --> C(Dispatcher);
    C --> D{What is the event type?};
    D -- Cloud Asset Inventory --> E[Find rules with matching asset_type];
    D -- Audit Log --> F[Find rules with matching service and methods];
    E --> G(Select Matching Rules);
    F --> G;
    G --> H{For each matching rule...};
    H --> I(Evaluate CEL Condition);
    I -- Condition is true --> J(Execute Actions);
    I -- Condition is false --> H;
    J --> K[Action Plugins];
    K --> L[e.g., Log to console];
```
