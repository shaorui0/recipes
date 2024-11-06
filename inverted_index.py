from collections import defaultdict
import datetime

# Sample time series data with structured content
time_series_data = [
    {
        "timestamp": datetime.datetime(2024, 11, 1, 9, 0),
        "__name__": "http_requests_total",
        "pod": "example-pod",
        "job": "example-job",
        "path": "/api/v1/resource",
        "status": "200",
        "method": "GET"
    },
    {
        "timestamp": datetime.datetime(2024, 11, 1, 12, 30),
        "__name__": "http_requests_total",
        "pod": "example-pod-2",
        "job": "example-job",
        "path": "/api/v1/resource",
        "status": "500",
        "method": "POST"
    },
    {
        "timestamp": datetime.datetime(2024, 11, 2, 15, 15),
        "__name__": "http_requests_total",
        "pod": "example-pod",
        "job": "example-job-2",
        "path": "/api/v2/other",
        "status": "404",
        "method": "GET"
    }
]

def create_inverted_index_time_series(data):
    inverted_index = defaultdict(list)

    for record in data:
        timestamp = record["timestamp"]
        for key, value in record.items():
            if key != "timestamp":  # Exclude the timestamp field from indexing
                value = str(value).lower()  # Normalize the value to lowercase for consistency
                inverted_index[(key, value)].append(timestamp)

    return inverted_index

# Create the inverted index
inverted_index_time_series = create_inverted_index_time_series(time_series_data)

# Print the inverted index
for key, timestamps in inverted_index_time_series.items():
    print(f"{key}: {timestamps}")

# Example of how to query the index for a specific field and value
def query_inverted_index_time_series(field, value, inverted_index):
    return inverted_index.get((field, value.lower()), [])

# Query example
query_field = "status"
query_value = "200"
result = query_inverted_index_time_series(query_field, query_value, inverted_index_time_series)
print(f"\nTimestamps containing {query_field}='{query_value}': {result}")
