# General
01. When getting a char field if the value has a spacebar at the end it will result in weird printing where
    closing curly brace will be moved to the next line.
02. Rows are stored in table in order of their editing. If record with `id` of 1 is edited when fetching
    table this record will be at the bottom.
03. When inserting db management system will assign index by using counter which value equals 
    to last assigned index. The record with this index may be long deleted at that moment.

# PGX
01. `(r *connRow) Scan()` returns the same `rows.Err()` that we can get from struct field.
02. *WARNING*: `(r *connRows) Close()` may overwrite `rows.Err()` with error of it's own. Beware when closing rows manually.
03. Some functions like `(pgx.Tx) Rollback()` will do nothing under some conditions. 
    So, they can be defered without fear. 

# Complete 
02. Learn how to do transactions **DONE?**
     -- After transaction the row was added with `id` of `8` while maximum `id` at the moment was 6. **done**
        Investigate.

# TODO 
01. Try scanning into struct like you did with `UnmarshalJSON()` at `bscrap`
03. Learn how to handle event stream from database
     -- Kafka or something native?
04. Learn how to create TTL indexes.

# Extra
01. Login to postgres using password via ipv4 connection.