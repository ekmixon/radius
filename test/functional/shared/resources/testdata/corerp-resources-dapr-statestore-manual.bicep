import radius as radius

param magpieimage string
param environment string
param namespace string = 'default'

resource app 'Applications.Core/applications@2022-03-15-privatepreview' = {
  name: 'corerp-resources-dsstm-old'
  properties: {
    environment: environment
  }
}

resource myapp 'Applications.Core/containers@2022-03-15-privatepreview' = {
  name: 'dapr-sts-manual-ctnr-old'
  properties: {
    application: app.id
    connections: {
      daprstatestore: {
        source: statestore.id
      }
    }
    container: {
      image: magpieimage
      readinessProbe: {
        kind: 'httpGet'
        containerPort: 3000
        path: '/healthz'
      }
    }
    extensions: [
      {
        kind: 'daprSidecar'
        appId: 'gnrc-sts-ctnr-old'
        appPort: 3000
      }
    ]
  }
}


module redis 'modules/redis-selfhost.bicep' = {
  name: 'dapr-sts-manual-redis-old-deployment'
  params: {
    name: 'dapr-sts-manual-redis-old'
    namespace: namespace
    application: app.name
  }
}


resource statestore 'Applications.Link/daprStateStores@2022-03-15-privatepreview' = {
  name: 'dapr-sts-manual-old'
  properties: {
    application: app.id
    environment: environment
    resourceProvisioning: 'manual'
    type: 'state.redis'
    metadata: {
      redisHost: '${redis.outputs.host}:${redis.outputs.port}'
      redisPassword: ''
    }
    version: 'v1'
  }
}
