# Dear Julien

> WARNING: these notes were written in a bit of a hurry. Proceed with caution

## Making it "work"

1. Start the NATS test server in [natsserver](./hack/natsserver/)
    - this outputs some files in an `out` directory that you need to pass to the controller
    - for convenience there is a `out/.env` file you can source for running the controller
2. Start the controller
    - I was using KinD for local testing and all good

    ```bash
    # Setup the env based on the NATS test server you started
    source $(pwd)/hack/natsserver/out/.env

    # Install CRDs and run the operator
    make install run

    # Apply the custom resources, which will create an account, user, stream and consumer.
    # They are all called "sample" and will be created without a specific namespace...
    kubectl apply -f config/samples/out/
    ```

    - Check the logs from the controller. There might be some errors because of ordering, but that's ok. As long as after a little while everything looks ok :)

3. Start the [ingest container builder](./hack/ingest/cmd/main.go)
    - it wil listen to a KV bucket and run [ko build](./hack/ingest/kobuild.go) based on the [ingestion template](./hack/ingest/templates/ingestion/)
    - NOTE: you need to run this with a NATS user for an account that has JetStream enabled
        - the NATS system account DOES NOT have Jetstream enabled
        - so just use the `sample` user/account that we created

    ```bash
    # Go to correct directory
    cd hack/ingest

    # Get the secret for the NATS user
    kubectl get secrets sample-creds -o json | jq '.data."user.creds"' --raw-output | base64 -d > user.creds

    # Setup env
    export NATS_URL=nats://0.0.0.0:4222
    export NATS_CREDS=user.creds

    # Start the service
    go run cmd/main.go
    ```

4. Send a schema to NATS using [schemastore](./hack/schemastore/)
    - it will put some proto files as txtar to a bucket
    - this needs to use the same account as `ingest` because buckets are per-account

    ```bash
    # Go to correct directory
    cd hack/schemastore

    # Get the secret for the NATS user
    kubectl get secrets sample-creds -o json | jq '.data."user.creds"' --raw-output | base64 -d > user.creds

    # Setup env
    export NATS_URL=nats://0.0.0.0:4222
    export NATS_CREDS=user.creds

    # Put the event.proto file as a schema.
    # The schema name does not matter, `ingest` listens to everything in the hardcoded bucket name.
    # You can pass more -in files, but there has to be one top-level message called "Event".
    go run main.go -schema=whatever -dir ../ingest/testdata/event/ -in=event.proto
    ```
