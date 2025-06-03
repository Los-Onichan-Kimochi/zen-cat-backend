import requests
import uuid
import json

# Base URL of your API
BASE_URL = "http://localhost:8098"

# Endpoint path
ENDPOINT = "/community-plan/"

# JWT Token (replace with your actual token)
JWT_TOKEN = "XD"

# Sample request data
request_data = {
    "communityId": str(uuid.uuid4()),  # Replace with actual UUID
    "planId": "MYCASO DE PRUEBA",       # Replace with actual UUID
    # Add other required fields as per your schema
}

# Headers
headers = {
    "Authorization": f"Bearer {JWT_TOKEN}",
    "Content-Type": "application/json"
}

# Make the POST request
response = requests.post(
    f"{BASE_URL}{ENDPOINT}",
    headers=headers,
    data=json.dumps(request_data)
)

# Print the response
print(f"Status Code: {response.status_code}")
print(f"Response: {response.json()}")
