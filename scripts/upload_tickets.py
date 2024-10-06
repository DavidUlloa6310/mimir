import requests
from requests.auth import HTTPBasicAuth
import weaviate
import json

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
    print("Incident Data fetched successfully.")
else:
    print(f"Failed to retrieve data. Status code: {response.status_code}, Response: {response.text}")
    incidents = None

# Connect to Weaviate
client = weaviate.Client("http://localhost:8080")  # Change this to your Weaviate instance URL

# Define the Weaviate schema for a Ticket if it's not already set up
ticket_schema = {
    "classes": [
        {
            "class": "Ticket",
            "description": "A ServiceNow incident ticket",
            "properties": [
                {
                    "name": "number",
                    "dataType": ["string"],
                    "description": "Incident number"
                },
                {
                    "name": "shortDescription",
                    "dataType": ["string"],
                    "description": "Short description of the incident"
                },
                {
                    "name": "priority",
                    "dataType": ["string"],
                    "description": "Incident priority"
                },
                {
                    "name": "state",
                    "dataType": ["string"],
                    "description": "Current state of the incident"
                }
            ]
        }
    ]
}

# Check if the schema exists, if not create it
try:
    client.schema.create(ticket_schema)
except weaviate.exceptions.UnexpectedStatusCodeException as e:
    print("Schema already exists or an error occurred:", e)

# Insert incidents into the Weaviate database
if incidents:
    for incident in incidents['result']:
        # Create a ticket object for each incident
        ticket_data = {
            "number": incident.get('number'),
            "shortDescription": incident.get('short_description'),
            "priority": incident.get('priority'),
            "state": incident.get('state')
        }
        
        # Insert ticket into Weaviate
        try:
            client.data_object.create(
                data_object=ticket_data,
                class_name="Ticket"
            )
            print(f"Inserted ticket {ticket_data['number']} into Weaviate.")
        except weaviate.exceptions.UnexpectedStatusCodeException as e:
            print(f"Failed to insert ticket {ticket_data['number']}: {e}")

else:
    print("No incidents were retrieved, so nothing to insert.")
