# graphql-relay-go [![Build Status](https://semaphoreci.com/api/v1/stratumn/graphql-pagination-go/branches/master/badge.svg)](https://semaphoreci.com/stratumn/graphql-pagination-go)

This repository is a fork of https://github.com/graphql-go/relay

At stratumn, we wanted an hybrid system between react-relay implementation of pagination and flat lists

So we made the choice to declare pagination as above:

```
  type StuffList {
    # list of segments
    items: [Stuff!]!
    # information about the pagination
    pageInfo: PageInfo!
    # total count of entities in DB
    totalCount: Int!
  }
```

A Go/Golang library to help construct a [graphql-go](https://github.com/graphql-go/graphql) server supporting react-relay.

Source code for demo can be found at https://github.com/graphql-go/playground

### Test

```bash
$ go get github.com/stratumn/graphql-pagination-go
$ go build && go test ./...
```
