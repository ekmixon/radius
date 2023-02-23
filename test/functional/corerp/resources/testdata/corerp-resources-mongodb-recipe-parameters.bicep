import radius as radius

param rg string = resourceGroup().name

param sub string = subscription().subscriptionId

param magpieimage string 

resource env 'Applications.Core/environments@2022-03-15-privatepreview' = {
  name: 'corerp-resources-environment-recipe-parameters-env'
  location: 'global'
  properties: {
    compute: {
      kind: 'kubernetes'
      resourceId: 'self'
      namespace: 'corerp-resources-environment-recipe-parameters-env'
    }
    providers: {
      azure: {
        scope: '/subscriptions/${sub}/resourceGroups/${rg}'
      }
    }
    recipes: {
      mongodb: {
          linkType: 'Applications.Link/mongoDatabases' 
          templatePath: 'radiusdev.azurecr.io/recipes/functionaltest/parameters/mongodatabases/kubernetes:1.0'
          parameters: {
            mongodbName: 'mdb-operator-${rg}'
          }
      }
    }
  }
}

resource app 'Applications.Core/applications@2022-03-15-privatepreview' = {
  name: 'corerp-resources-mongodb-recipe-parameters'
  location: 'global'
  properties: {
    environment: env.id
    extensions: [
      {
          kind: 'kubernetesNamespace'
          namespace: 'corerp-resources-mongodb-recipe-param-app'
      }
    ]
  }
}

resource webapp 'Applications.Core/containers@2022-03-15-privatepreview' = {
  name: 'mdb-param-ctnr'
  location: 'global'
  properties: {
    application: app.id
    connections: {
      mongodb: {
        source: recipedb.id
      }
    }
    container: {
      image: magpieimage
      env: {
        DBCONNECTION: recipedb.connectionString()
      }
      readinessProbe:{
        kind:'httpGet'
        containerPort:3000
        path: '/healthz'
      }
    }
  }
}

resource recipedb 'Applications.Link/mongoDatabases@2022-03-15-privatepreview' = {
  name: 'mdb-recipe-param-db'
  location: 'global'
  properties: {
    application: app.id
    environment: env.id
    mode: 'recipe'
    recipe: {
      name: 'mongodb'
      parameters: {
        mongodbName: 'mdb-developer-${rg}'
      }
    }
  }
}
