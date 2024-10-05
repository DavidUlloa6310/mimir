import json
import openai
import weaviate
import os
import time
from pathlib import Path

from dotenv import load_dotenv

load_dotenv(dotenv_path=Path("../.env"))

# Load environment variables
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
WEAVIATE_API_KEY = os.getenv("WEAVIATE_API_KEY")
WEAVIATE_URL = os.getenv("WEAVIATE_URL")

# Check if the necessary environment variables are set
if not OPENAI_API_KEY:
    raise ValueError("Please set the OPENAI_API_KEY environment variable.")
if not WEAVIATE_URL:
    raise ValueError("Please set the WEAVIATE_URL environment variable.")

# Set OpenAI API key
openai.api_key = OPENAI_API_KEY

# Initialize Weaviate client
client = weaviate.Client(
    url=WEAVIATE_URL,
    additional_headers={
        "X-OpenAI-Api-Key": OPENAI_API_KEY,
        # Uncomment the following line if authentication is required
        "Authorization": f"Bearer {WEAVIATE_API_KEY}"
    }
)

# Define the class schema in Weaviate if it doesn't exist
class_name = "Accelerator"

schema = {
    "classes": [
        {
            "class": class_name,
            "description": "A class representing an accelerator with a description.",
            "properties": [
                {
                    "name": "name",
                    "dataType": ["text"]
                },
                {
                    "name": "url",
                    "dataType": ["text"]
                },
                {
                    "name": "category",
                    "dataType": ["text"]
                },
                {
                    "name": "description",
                    "dataType": ["text"]
                }
            ]
        }
    ]
}

# Check if the class already exists
existing_classes = client.schema.get()["classes"]
if not any(c["class"] == class_name for c in existing_classes):
    client.schema.create(schema)

# Load the JSON file
with open('accelerators.json', 'r') as f:
    data = json.load(f)

accelerators = data['accelerators']

for accelerator in accelerators:
    url = accelerator.get('url')
    name = accelerator.get('name')
    category = accelerator.get('category')

    # Skip if essential fields are missing
    if not url or not name:
        print(f"Skipping accelerator due to missing 'url' or 'name': {accelerator}")
        continue

    # Create the prompt
    prompt = f"""
You are a helpful assistant that generates detailed markdown descriptions for accelerators.
Please provide a comprehensive description for the accelerator below:

**Name:** {name}
**URL:** {url}
**Category:** {category}

The description should include key features, benefits, and any relevant information to help someone understand what this accelerator offers.

Please format the description in markdown.
"""

    # Call the GPT-4 API
    try:
        response = openai.chat.completions.create(
            model='gpt-4',
            messages=[
                {"role": "user", "content": prompt.strip()}
            ],
            max_tokens=500,
            n=1,
            stop=None,
            temperature=0.7,
        )

        description = response.choices[0].message.content.strip()

        # Prepare the data object for Weaviate
        data_object = {
            'name': name,
            'url': url,
            'category': category,
            'description': description
        }

        # Store the object in Weaviate
        client.data_object.create(
            data_object,
            class_name=class_name
        )

        print(f"Successfully processed accelerator: {name}")

    except Exception as e:
        print(f"Error processing accelerator '{name}': {str(e)}")

    # Optional: Delay to respect API rate limits
    time.sleep(1)
