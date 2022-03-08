# grafana-action
This repo contains the source code for the /grafana-action GitHub Action.
This action pushes annotation to a Grafana dashboard.

**Required** Env variables: The following env variables MUST be defined in your workflow
```
env:
  PROJECT_NAME:           ${{ github.event.repository.name }}
  TEAM:                   <My_team>
  GRAFANA_API_TOKEN:      ${{ secrets.grafana_api_token }}
  GRAFANA_URL:            https://grafana.<org>/api
```


## Example usage 1 - send a pre defined template of started message
**A previous checkout step with fetch-depth of 0 is required!!!**

```
    - name: Check out code
      uses: actions/checkout@v2
    - name: Push Grafana annotation
        uses: online-applications/grafana-action@v1
```
