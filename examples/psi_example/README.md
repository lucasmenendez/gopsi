# Basic PSI protocol example

The following diagram explains the basic example of the library using [Client](/pkg/client/client.go) (alice) and [Server](./pkg/server/server.go) (bob) structs. This example is implemented into [psi_example/main.go](./main.go.go) file.

However you can use [SRA](/pkg/sra/sra.go) or [BloomFilters](/pkg/bloomfilter/bloomfilter.go) isolated and also design your own protocol using it.

## Iteration flow

```mermaid
sequenceDiagram
    participant alice
    participant bob

    Note over alice,bob: Request intersection
    alice->>+bob: send alice RSA public key
    bob-->>bob: generate prime number 
    bob-->>-alice: send prime number encrypted w/ alice public key

    Note over alice,bob: perform intersection
    par both encrypts its own data with the common prime
        alice-->>alice: encrypt its data with SRA
    and
        bob-->>bob: encrypt its data with SRA
    end

    bob->>alice: send encrypted data
    alice-->>alice: re-encrypt bob data with SRA
    alice->>bob: send its encrypted data and bobs re-encrypted data
    bob-->>bob: initialize intersection with re-encrypted data
    bob-->>bob: re-encrypt alice data with SRA
    bob-->>alice: send the intersection between re-encrypted data sets
```

## Run example

```sh
go run main.go
```