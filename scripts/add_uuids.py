import weaviate
import os
import uuid

from pathlib import Path
from dotenv import load_dotenv

load_dotenv(dotenv_path=Path("../.env"))

OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
WEAVIATE_API_KEY = os.getenv("WEAVIATE_API_KEY")
WEAVIATE_URL = os.getenv("WEAVIATE_URL")

# Initialize the client
client = weaviate.Client(
    url=f"https://{WEAVIATE_URL}",
    additional_headers={
        "X-OpenAI-Api-Key": OPENAI_API_KEY,
        # Uncomment the following line if authentication is required
        "Authorization": f"Bearer {WEAVIATE_API_KEY}"
    }
)

# Define the new property to add to your schema
new_property = {
    "dataType": ["string"],  # UUIDs will be strings
    "name": "ID",
}

# Add the new property to your class 'Accelerator'
client.schema.property.create(schema_class_name="Accelerator", schema_property=new_property)

# Fetch all the existing objects in the 'Accelerator' class
response = client.query.get("Accelerator", ["_additional { id }"]).do()

# Extract object IDs from the response
objects = response["data"]["Get"]["Accelerator"]

# Loop through each object and add a new UUID to the `ID` field
for obj in objects:
    object_id = obj["_additional"]["id"]
    
    # Generate a new UUID
    new_uuid = str(uuid.uuid4())

    # Update the object with the new UUID
    client.data_object.update(
        data_object={"ID": new_uuid},
        class_name="Accelerator",
        uuid=object_id
    )

print("All objects updated with new UUIDs!")