# github-exporter
The github-exporter deployment makes github clones and views statistics available to prometheus. On startup it queries the specified GITHUB_USER and starts to monitor all its repositories. The statistics are gathered with the personal access token.

The values are made available for scraping by prometheus. The scrape url is http://<service_address>:<PORT>/metrics. If you use the prometheus-operator deployment in combination with our helm chart the scrape config is not needed. The helm chart contains a serviceMonitor definition.

# Environment variables
The following enviroment variables are supported:
| Variable | Description |
| -------- | -------- |
| REFRESH_SECONDS | The amount of seconds between successive polls of github, defaults to 300 |
| PORT_NUMBER | The port that is used for publishing the metrics, defaults to 9900 |
| GITHUB_USER | The github user account that owns the repositories that should be monitored |

# First deployment
The first deployment will fail because the secret is not yet created. After the first deployment the namespace exists so we can create the secret(in the next section). After this the pod should be rstarted to pickup the new secret.

# Secrets
The github-exporter needs a github personal access token to read the statistics. This token can be genarated as follows
[github](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#creating-a-fine-grained-personal-access-token)

Make sure that the token has the "Administration" permission.

This token should be put in a secret in the namespace of the github-exporter deployment with the following command f.e.:
```
kubectl create secret generic github-exporter -n github-exporter --from-literal=token=github_pat_YOUR_TOKEN
```
The name of the secret should be "github-exporter". 

Another approach would be to use external-secrets in combination with Vault, as long as the secret name stays the same.

# Available metrics
At this time 2 metrics are available, `github_clones` and `github_views`, both metrics have the `unique` label to differentiate between all or unique entries.
