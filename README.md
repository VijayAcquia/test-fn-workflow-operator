# Workflows Operator
This is the operator responsible for orchestrating workflows over NGC infrastructure.  It is a consumer of the `fnresources.acquia.io` api.

## Building

## Deployment
1. Generate a local Helm chart:

    ```bash
    cd helm
    ./package.sh
    ```

1. Deploy the Helm chart:

    ```bash
    helm install --name fn-drupal-operator --namespace services ./fn-drupal-operator \
      --set image.tag=${your_image_tag}
    ```

    If you want the Operator to only watch a specific Namespace, that can be specified by adding the command line option:

    ```bash
    --set watchNamespace=${namespace_to_watch}
    ```

## Local Testing
