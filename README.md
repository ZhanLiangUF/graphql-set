# graphql-set
GraphQL service for retrieval and storage of integer sets

## Installation
Git clone
```bash
https://github.com/ZhanLiangUF/graphql-set.git
```

Download docker - https://www.docker.com/products/docker-desktop and make sure to run it 

Change into directory

```bash
docker-compose up
```

Run locally on port 8080

Here are a sample query and mutation to test:

```graphql
mutation CreateNewSet {
  createSet(input:{members:[1,2,3,4,5,6,7,8,9]}) {
    members
    intersectingSets {
      members
    }
  }
}
```

```graphql
query GetAllSet {
	sets {
    members
    intersectingSets {
      members
    }
  }  
}
```