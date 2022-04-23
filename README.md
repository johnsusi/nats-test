Filtering streams by consumer test
==================================

This tests the performance for an interleaved stream (simulating for instance 
multiple IOT-devices reporting to the same stream with the ID in the subject.)

1M messages are produced for 100 entities, 10k messages each. Messages will be 
published as

`[entity 0 message 0] [entity 1 message 0] ... `

Server is run with `nats-server -c nats.conf`
Stream is created using `nats str create test`
with subjects `foo.*.bar.*` and the rest are the defaults.

```
$ nats str info test
Information for Stream test created 2022-04-23T17:35:46+02:00

Configuration:

             Subjects: foo.*.bar.*
     Acknowledgements: true
            Retention: File - Limits
             Replicas: 1
       Discard Policy: Old
     Duplicate Window: 2m0s
    Allows Msg Delete: true
         Allows Purge: true
       Allows Rollups: false
     Maximum Messages: unlimited
        Maximum Bytes: unlimited
          Maximum Age: unlimited
 Maximum Message Size: unlimited
    Maximum Consumers: unlimited


State:

             Messages: 1,000,000
                Bytes: 55 MiB
             FirstSeq: 1 @ 2022-04-23T15:37:03 UTC
              LastSeq: 1,000,000 @ 2022-04-23T15:37:14 UTC
     Active Consumers: 0
```

Producer is run with `go run producer.go`

Consumer 1 is run with `go run consumer.go`

Consumer 2 is run with `go run consumer.go filtered`

Expected outcome
----------------
A filtered consumer is faster or as fast as an unfiltered consumer.

Actual outcome
--------------

```
$ go run consumer/consumer.go 
2022/04/23 17:39:58 Processed 1000 messages in 0 seconds. 386503 msgs / sec.
2022/04/23 17:39:58 Processed 2000 messages in 0 seconds. 488436 msgs / sec.
...
2022/04/23 17:40:01 Processed 999000 messages in 2 seconds. 397756 msgs / sec.
2022/04/23 17:40:01 Processed 1000000 messages in 2 seconds. 397955 msgs / sec.
```

```
$ go run consumer/consumer.go filtered
2022/04/23 17:40:11 Processed 1000 messages in 4 seconds. 221 msgs / sec.
2022/04/23 17:40:14 Processed 2000 messages in 8 seconds. 245 msgs / sec.
...
2022/04/23 17:40:52 Processed 9000 messages in 46 seconds. 194 msgs / sec.
2022/04/23 17:40:55 Processed 10000 messages in 49 seconds. 202 msgs / sec.
```

Note: msgs / sec is expected to be lower for the filtered consumer since it skips 99 out of 100 messages in the stream. 

Final thoughts:
---------------

The structure of the subject seems to matter. Doing the same experiment but without encoding a second wildcard in the subject does not exhibit this behavior ('foo.*' vs 'foo.50').

