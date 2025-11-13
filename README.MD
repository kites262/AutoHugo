# AutoHugo

AutoHugo presents a SaaS solution for automating the deployment of Hugo static sites.

By specifying a **repository URL**, AutoHugo automatically pulls Hugo source code from GitHub, builds and deploys it as a static website, while also providing webhook functionality to automatically trigger rebuilds and deployments when the source code is updated.

## Features

- **Automated Deployment**: Seamlessly deploy your Hugo sites with minimal effort.
- **Webhook**: Integrate with GitHub to trigger deployments on push events.

## Configuration

Use the provided `docker-compose.yml` to set up the service. Ensure to configure the GitHub SSH key path in the secrets section.
