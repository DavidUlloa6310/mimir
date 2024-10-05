import requests
from requests.auth import HTTPBasicAuth

# Replace these with your ServiceNow instance details and credentials
instance_url = "https://dev274800.service-now.com/"
username = "admin"
password = "r8RGnqYX=%m0"

# Endpoint to fetch data from the incident table
url = f"{instance_url}/api/now/table/incident"

# Define any query parameters you need (optional)
query_params = {
    'sysparm_limit': '10',  # limits the number of incidents returned
    'sysparm_fields': 'number,short_description,priority,state'  # fields to return
}

# Set the headers for the request
headers = {
    "Accept": "application/json",
    "Content-Type": "application/json"
}

# Make the request to ServiceNow API
response = requests.get(url, auth=HTTPBasicAuth(username, password), headers=headers, params=query_params)

# Check if the request was successful
if response.status_code == 200:
    # Parse the JSON response
    incidents = response.json()
    print("Incident Data:", incidents)
else:
    print(f"Failed to retrieve data. Status code: {response.status_code}, Response: {response.text}")
