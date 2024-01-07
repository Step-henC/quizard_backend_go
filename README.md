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
  For users, Redis employs a bloom-filter data structure. Briefly, Bloom filters are non-deterministic data structures that have higher false-positive rates,
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
    build the image at the project's root directory

    `docker build -t quizard-be:multistage -f Dockerfile.multistage .`

    Then run the container

    ` docker run --publish 8080:8080 quizard-be:multistage`

    You should see the graphql playground at `localhost:8080` in your browser.

    ![gql playground](https://blog.logrocket.com/wp-content/uploads/2020/06/graphql-playground-send-http-headers-1.png)

    Please refer to the schema.graphqls in this project to craft queries. Start with

    `query GetAuth {
        getAuth
    }`

    to get JWT for subsequent queries.

    To create user, copy that Jwt (myJwt) from get auth. Then at the bottom in Headers section, add
    `{
        "Authorization": "Bearer {myJwt}"
    }`
    

    `mutation CreateUser($input: UserInput!) {
      createUser(input: $input) {
        email,
        password
      }
    }`

### Considerations
  - JWT Authentication: 
