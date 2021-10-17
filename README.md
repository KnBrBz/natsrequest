# Nats request timeout snippet

## Problem

Program fails with message like `Source.NNNN.N request: *Request.Publish: nats: timeout`.

Number of iteraction that succeeds may vary.

## Program structure

- Events - imitates events every 100ms.
- Sources - imitates requests from up to 2000 sources, starts from 200 with step 200. Each source sends a request every 500ms.
- Destinations - every destination subscribes for all 100 events and one source request. One goroutine for events processing, other for requests.

There is two Nats connection instances for Events, Sources and one encoded connection for Destinations.
