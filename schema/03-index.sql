UPDATE messages_raw SET time = '2016-01-01' WHERE time IS NULL;

CREATE UNIQUE INDEX messages_raw_sqs_id ON messages_raw (sqs_id);
