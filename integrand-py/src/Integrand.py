import requests
from typing import Dict, Optional
import SSE

class Integrand:
    def __init__(self, integrand_base_endpoint: str, integrand_api_key: str):
        self.integrand_base_endpoint = integrand_base_endpoint
        self.api_key = integrand_api_key

    def EndpointRequest(self, route: str):
        url = f'{self.integrand_base_endpoint}/api/v1/glue/f/{route}'
        headers = {
            'Content-Type': 'application/json',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body
    
    # Glue API
    def GetAllGlueHandlers(self):
        url = f'{self.integrand_base_endpoint}/api/v1/glue'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def GetGlueHandler(self, glueHandlerId: str):
        url = f'{self.integrand_base_endpoint}/api/v1/glue/{glueHandlerId}'
        token = 'Bearer ' + self.api_key
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def CreateGlueHandler(self, id: str, topicName: str):
        url = f'{self.integrand_base_endpoint}/api/v1/glue'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        body = {
            'id': id,
            'topicName': topicName
        }
        response = requests.post(url, headers=headers, json=body)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def DeleteGlueHandler(self, glueHandlerId: str):
        url = f'{self.integrand_base_endpoint}/api/v1/glue/{glueHandlerId}'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.delete(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    # Topics API
    def GetAllTopics(self):
        url = f'{self.integrand_base_endpoint}/api/v1/topic'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def GetTopic(self, topicName: str):
        url = f'{self.integrand_base_endpoint}/api/v1/topic/{topicName}'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def CreateTopic(self, topicName: str):
        url = f'{self.integrand_base_endpoint}/api/v1/topic/{topicName}'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        body = {
            "topicName": topicName 
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def DeleteTopic(self, topicName: str):
        url = f'{self.integrand_base_endpoint}/api/v1/topic/{topicName}'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.delet(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def ConsumeTopic(self, topicName: str, offset: Optional[int] = None):
        params: Dict[str, Union[str, int]] = dict()
        url = f'{self.integrand_base_endpoint}/api/v1/topic/{topicName}/consume'
        if offset is not None:
            params["offset"] = offset
        headers = {
            'Accept': 'text/event-stream',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, stream=True, headers=headers, params=params)
        response.raise_for_status()
        sseClient = SSE.SSEClient(response)
        return sseClient.events()