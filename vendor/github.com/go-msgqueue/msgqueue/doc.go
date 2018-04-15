/*
Package msgqueue implements task/job queue with in-memory, SQS, IronMQ backends.

go-msgqueue is a thin wrapper for SQS and IronMQ clients that uses Redis to implement rate limiting and call once semantic.

go-msgqueue consists of following components:
 - memqueue - in memory queue that can be used for local unit testing.
 - azsqs - Amazon SQS backend.
 - ironmq - IronMQ backend.
 - Manager - provides common interface for creating new queues.
 - Processor - queue processor that works with memqueue, azsqs, and ironmq.

rate limiting is implemented in the processor package using https://github.com/go-redis/redis_rate. Call once is implemented in clients by checking if message name exists in Redis database.
*/
package msgqueue
