# Batch Log Server

## Configuration
The server retrieves its configuration through environment variables

BATCH_SIZE - Specifies the max size of the cache

BATCH_INTERVAL - Specifies the time period before flushing the cache to the configured endpoint

BATCH_ENDPOINT - Specifies the endpoint to make an http post request to when flushing the cache