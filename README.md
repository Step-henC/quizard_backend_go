# Quizard Backend

### What Is This?

This repository houses the backend service for the Quizard App with the UI found [at this repository.](https://github.com/Step-henC/quizard_ui)
The quizard backend service is written in Go with incoming quiz and user information with JWT headers routed to an [ElasticSearch database](https://pkg.go.dev/github.com/elastic/go-elasticsearch/v7) 
and Redis for the cache/ bloom filter.

### How Does This Service work?
- GraphQL: This service is bootstrapped from GqlGen's grapql server. As a result, the main files changed are the schema files to identify the types, the schema.resolver files
  to implement logic on what to do with data, and any handlers. GqlGen's graphql uses net/http package in Go. However the Gin web server framework is regarded as [faster](https://veryfirstfact.com/comparing-gorilla-mux-gin-net-http-for-http-web-framework/#:~:text=Gin%20makes%20use%20of%20httprouter%2C%20which%20performs%20operations%20more%20quickly)
  than net/http package. So, I wrapped this [gqlgen graphql server in the gin framework](https://gqlgen.com/recipes/gin/) for faster performance. The http context can then be extracted from each endpoint and
  evaluated for JWTs.
- JWT Authentication: JWTs are handled most popularly in either application (in-memory) storage or Http-only. Http-only cookies are generally considered more safe, however
  this application uses GraphQL servers that have been bootstrapped from Apollo (in front end) and [GqlGen](https://gqlgen.com/). As a result the middle ground was to use an in-memory token serving
  as a refresh token. This means that the Apollo graphql server requests a JWT token with every request. The JWT request is made with a 'GetAuth' query. The resulting token in stored in-memory of the React
  app in the frontend.
- Databases: ElasticSearch is growing in popularity for the benefits it affords in [high performance and fast delivery times for developing products](https://aws.amazon.com/what-is/elasticsearch/#:~:text=Elasticsearch%20benefits,-Fast%20time%2Dto&text=Elasticsearch%20offers%20simple%20REST%2Dbased,applications%20for%20various%20use%20cases.)
  This application uses a Docker-Compose file for the elasticSearch database. The input data is resolved into json where is sent to elasticsearch for storage.
- Cache: Redis is a popular key-value storage for caching data. This application uses the [cache-aside pattern for Redis caching](https://codedamn.com/news/backend/advanced-redis-caching-techniques#:~:text=use%20in%20Redis.-,Cache%2DAside%20Pattern,-The%20cache%2Daside)
  For users, Redis employs a [bloom-filter data structure](https://redis.io/docs/data-types/probabilistic/bloom-filter/). Briefly, Bloom filters are non-deterministic data structures that have higher false-positive rates,
  but zero false-negatives. Meaning, bloom filters cannot tell you whats in the database, but can tell you what is not in the database. The filter
  is implemented before a timely, json serializing db read to elasticsearch is made for users. The result is less reads and hopefully high performance.
  Redis is also a docker compose, with the UI portion provided by a [redis-commander docker container](https://migueldoctor.medium.com/run-redis-redis-commander-in-3-steps-using-docker-195fc6fa7076).

### How to run this services? 

Make sure you have [Docker Engine installed](https://docs.docker.com/engine/install/). We need them for the Databases

##### Run Databases
 - Redis and Elasticsearch are housed in a multi-container docker compose file at the project's root directory. Start DB's run the following command

   `docker compose -d up`

    You should see Elasticsearch in your browser at `localhost:9200` and Redis Commander at `localhost:8081`

   To view ElasticSearch users, visit your browser at `localhost:9200/users/_search?pretty`
   To view Elasticsearch quizzes, visit browse at `localhost:9200/quizzes/_search?pretty`

##### Run Service
   
###### No Go Installed?
  - This app is containerized with [Docker, for those without Go installed](https://docs.docker.com/language/golang/build-images/). To run, open up 
    docker desktop and make sure engine is running. Then, simply 

    run docker compose at the project's root directory

    `docker-compose up`

    You should see the graphql playground at `localhost:8080` in your browser.

    ![gql playground](https://blog.logrocket.com/wp-content/uploads/2020/06/graphql-playground-send-http-headers-1.png)

    Then, open Docker desktop and execute the following curl commands in the elasticsearch container 

    `curl -X PUT http://localhost:9200/users && curl -X PUT http://localhost:9200/quizzes`

    This is until I refactor code to create indices in api call.

    We should be able to CRUD data now.

    Please refer to the schema.graphqls in this project to craft queries. Start with

    `query GetAuth {
        getAuth
    }`

    to get JWT for subsequent queries.

    To create user, copy that Jwt (myJwt) from get auth. Then at the bottom in Headers section, add
    `{
        "Authorization": "Bearer {myJwt}"
    }`
    

    Then, run the following mutation:
    
    `mutation CreateUser($input: UserInput!) {
      createUser(input: $input) {
        email,
        password
      }
    }`

    In the variables section at the bottom, add
    {
      "email": "anyemail",
      "password": "anypassword"
    }

    You should see the info returned. Then check elasticsearch for users in the browser at `localhost:9200/users/_search?pretty`

    This will get you on your way to creating quizzes!

###### Go Installed?

Pull down as you would any go project. Be sure to initiate the gqlgen graphql server with the following commands

`go get github.com/99designs/gqlgen` 
`go run github.com/99designs/gqlgen init ` 
`go run github.com/99designs/gqlgen generate`

For databases, go to the `dbcontainers` directory and `docker-compose up` with just the databases.
Continue to CRUD data as mentioned in the previous section


### Considerations
  - JWT Authentication: JWTs are regarded as safer as Http-only, however with the graphql server being bootstrapped from gqlgen, I used a refresh token 
    with a token request regenerating a new token. Additionally, the gql server wrapped in a gin framework made context tricky for http-only jwts. 
    Nonetheless, for development purposes, I did not flush out token secrets (or hide them). Honestly JWTs were not a 
    main focus in this project. I focused on implementing JWTs in the most simplest code imaginable. However, future projects will employ JWTs with a bit 
    more intentionality.

 - Data Architecture: The gqlgen generates the models. All mutations need an input type. However, mutations cannot return input types. Moreover, I need 
   the original model/type to be returned to help with writing to Apollo GraphQL's in memory cache in the front end. Mapping input types to their 
   corresponding query types made mapping data less ideal. Future projects will fine tune the shape of the data in relation to gqlgen model generation
