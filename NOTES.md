# General
01. When getting a char field if the value has a spacebar at the end it will result in weird printing where
    closing curly brace will be moved to the next line.
02. Rows are stored in table in order of their editing. If record with `id` of 1 is edited when fetching
    table this record will be at the bottom.
03. When inserting, db management system will assign index by using counter whose value equals 
    to last assigned index.

# PGX
01. `(r *connRow) Scan()` returns the same `rows.Err()` that we can get from struct field.
02. *WARNING*: `(r *connRows) Close()` may overwrite `rows.Err()` with error of it's own. Beware when closing rows manually.
03. Some functions like `(pgx.Tx) Rollback()` will do nothing under some conditions. 
    So, they can be defered without fear. 
04. `pgconn.Conn` is not thread-safe. Hence in most practical scenarious pgxpool.Pool should be used.
05. Aquiring connection to `pgpool.Pool` can be disallowed by cancelling certain `context`
06. We can disallow receiving notification by cancelling context that was passed 
    into `(*pgxConn)WaitForNotification()`.

# Complete 
01. Try scanning into struct like you did with `UnmarshalJSON()` at `bscrap` **DONE**
     -- Consider using `scany` instead
02. Learn how to do transactions **DONE**
     -- After transaction the row was added with `id` of `8` while maximum `id` at the moment was 6. **done**
        Investigate.
06. Figure out how to unmarshal data from notification **DONE**
07. Add timestamp field to `oplog` **DONE**

# TODO 
03. Learn how to handle event stream from database
     -- Kafka or something native?
     -- Try native first
04. Learn how to create TTL indexes.
05. Refactor notification receiver
     -- Does cancelling context that was passed into `(*pgxConn)WaitForNotification()` cancels wait 
        for new notifications? **YES**
08. Figure out why program doesn't shutdown on `^C`
     -- why does it require to adidtionally call `os.Exit()`?
09. Put SQL code into repo as well.
10. Figure out how to handle case when both error and context cancellation occur in changestream
     -- are context cancellation and error occuring in changestream mutually exclusive?

# Extra
01. Login to postgres using password via ipv4 connection.