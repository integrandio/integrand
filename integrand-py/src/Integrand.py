import requests
from typing import Dict, Optional
import SSE

class Integrand:
    def __init__(self, integrand_base_endpoint: str, integrand_api_key: str):
        self.integrand_base_endpoint = integrand_base_endpoint
        self.api_key = integrand_api_key

    def EndpointRequest(self, route: str, security_key: str, data: Dict):
        url = f'{self.integrand_base_endpoint}/api/v1/connector/f/{route}'
        headers = {
            'Content-Type': 'application/json',
        }
        params = {
            "apikey": security_key
        }
        response = requests.post(url, headers=headers, json=data, params=params)
        response.raise_for_status()
        response_body = response.json()
        return response_body
    
    # Glue API
    def GetAllConnectors(self):
        url = f'{self.integrand_base_endpoint}/api/v1/connector'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def GetConnector(self, connector_id: str):
        url = f'{self.integrand_base_endpoint}/api/v1/connector/{connector_id}'
        token = 'Bearer ' + self.api_key
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def CreateConnector(self, connector_id: str, topicName: str):
        url = f'{self.integrand_base_endpoint}/api/v1/connector'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        body = {
            'id': connector_id,
            'topicName': topicName
        }
        response = requests.post(url, headers=headers, json=body)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def DeleteConnector(self, connector_id: str):
        url = f'{self.integrand_base_endpoint}/api/v1/connector/{connector_id}'
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
        url = f'{self.integrand_base_endpoint}/api/v1/topic'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        body = {
            "topicName": topicName 
        }
        response = requests.post(url, headers=headers, json=body)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def DeleteTopic(self, topicName: str):
        url = f'{self.integrand_base_endpoint}/api/v1/topic/{topicName}'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.delete(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body

    def GetEventsFromTopic(self, topicName: str, offset: Optional[int] = None, limit: Optional[int] = None):
        params: Dict[str, Union[str, int]] = dict()
        url = f'{self.integrand_base_endpoint}/api/v1/topic/{topicName}/events'
        if offset is not None:
            params["offset"] = offset
        if limit is not None:
            params["limit"] = limit
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, stream=True, headers=headers, params=params)
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
    
    def CreateWorkflow(self, topicName: str, functionName: str, sinkURL: str):
        url = f'{self.integrand_base_endpoint}/api/v1/workflow'
        headers = {
            'Content-Type': 'application/json',
            'Authorization': f'Bearer {self.api_key}',
        }
        body = {
            "topicName": topicName ,
            "functionName": functionName ,
            "sinkURL": sinkURL 
        }
        response = requests.post(url, headers=headers, json=body)
        response.raise_for_status()
        response_body = response.json()
        return response_body
    
    def DeleteWorkflow(self, id: str):
        url = f'{self.integrand_base_endpoint}/api/v1/workflow/{id}'
        headers = {
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.delete(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body
    
    def UpdateWorkflow(self, id: str):
        url = f'{self.integrand_base_endpoint}/api/v1/workflow/{id}'
        headers = {
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.put(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body
    
    def GetWorkflow(self, id: str):
        url = f'{self.integrand_base_endpoint}/api/v1/workflow/{id}'
        headers = {
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body
    
    def GetWorkflows(self):
        url = f'{self.integrand_base_endpoint}/api/v1/workflow/'
        headers = {
            'Authorization': f'Bearer {self.api_key}',
        }
        response = requests.get(url, headers=headers)
        response.raise_for_status()
        response_body = response.json()
        return response_body