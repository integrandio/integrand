# This is the integration test suite for LucidMQ
from requests import HTTPError
from Integrand import Integrand
import pytest
import string
import os
import json
import time
import random
from typing import Dict

INTEGRAND_URL = os.environ.get('INTEGRAND_URL', 'http://localhost:8000')  # The server's hostname or IP address
INTEGRAND_API_KEY = os.environ.get('INTEGRAND_API_KEY', '11111')

def get_random_string(length: int):
    # choose from all lowercase letter
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str

@pytest.fixture(scope="class")
def clean_up_topics():
    print("Setting Up Class")
    yield
    print("Tearing Down Class")
    integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
    response = integrand.GetAllTopics()
    topics = response['data']
    for topic in topics:
        integrand.DeleteTopic(topic['topicName'])
        
@pytest.fixture(params=["ld_ld_sync", "calendly_sync"])
def workflow_function_name(request):
    return request.param

@pytest.mark.usefixtures("clean_up_topics")
class TestConnectorAPI:
    #TODO: Setup by deleting all the topics
    def test_get_all_connectors_empty(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.GetAllConnectors()
        assert response['status'] == 'success'
        assert len(response['data']) == 0

    def test_get_connector_empty(self):
        route = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        with pytest.raises(Exception) as httpException:
            response = integrand.GetConnector(route)
        #TODO: Check the error here
        print(httpException)

    def test_create_connector(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.CreateConnector(id, topicName)
        assert response['status'] == 'success'
        assert response['data']['id'] == id
        assert response['data']['topicName'] == topicName
        # Clean up
        integrand.DeleteConnector(id)
    
    def test_get_all_connectors(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateConnector(id, topicName)
        response = integrand.GetAllConnectors()
        assert response['status'] == 'success'
        assert len(response['data']) == 1
        # Clean up
        integrand.DeleteConnector(id)

    def test_get_connector(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateConnector(id, topicName)
        response = integrand.GetConnector(id)
        assert response['status'] == 'success'
        assert response['data']['id'] == id
        assert response['data']['topicName'] == topicName
        # Clean up
        integrand.DeleteConnector(id)
    
    def test_delete_connector(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateConnector(id, topicName)
        response = integrand.DeleteConnector(id)
        assert response['status'] == 'success'


class TestTopicAPI:
    def test_get_topics_empty(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.GetAllTopics()
        assert response['status'] == 'success'
        assert len(response['data']) == 0
    
    def test_get_topic_empty(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        with pytest.raises(Exception) as httpException:
            response = integrand.GetTopic(topicName)
        #TODO: Check the error here
        print(httpException)

    def test_create_topic(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.CreateTopic(topicName)
        assert response['status'] == 'success'
        assert response['data']['topicName'] == topicName
        assert response['data']['oldestOffset'] == 0
        assert response['data']['nextOffset'] == 0
        # Clean up
        integrand.DeleteTopic(topicName)

    def test_get_all_topics(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateTopic(topicName)
        response = integrand.GetAllTopics()
        assert response['status'] == 'success'
        assert len(response['data']) == 1
        # Clean up
        integrand.DeleteTopic(topicName)

    def test_get_topic(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateTopic(topicName)
        response = integrand.GetTopic(topicName)
        assert response['status'] == 'success'
        assert response['data']['topicName'] == topicName
        # Clean up
        integrand.DeleteTopic(topicName)
    
    def test_delete_topic(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateTopic(topicName)
        response = integrand.DeleteTopic(topicName)
        assert response['status'] == 'success'
        
    def test_delete_topic_used_by_workflow_raises_exception(self, workflow_function_name):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        connectorId = get_random_string(5)
        integrand.CreateTopic(topicName)
        integrand.CreateConnector(connectorId, topicName)
        response = integrand.CreateWorkflow(topicName, workflow_function_name, "http://example.com")
        workflowId = response['data']['id']
        try:
            integrand.DeleteTopic(topicName)
        except HTTPError as e:
            assert e.response.status_code == 500
        # Cleanup
        integrand.DeleteConnector(connectorId)
        integrand.DeleteWorkflow(workflowId)
        integrand.DeleteTopic(topicName)
        
    def test_delete_topic_used_by_connector_raises_exception(self):
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        integrand.CreateTopic(topicName)
        connectorId = get_random_string(5)
        integrand.CreateConnector(connectorId, topicName)
        try:
            integrand.DeleteTopic(topicName)
        except HTTPError as e:
            assert e.response.status_code == 500
        # Cleanup
        integrand.DeleteConnector(connectorId)
        integrand.DeleteTopic(topicName)
    
class TestMessages():
    def test_send_message(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        createResponse = integrand.CreateConnector(id, topicName)
        data = {'hello': 'world'}
        res = integrand.EndpointRequest(id, createResponse['data']['securityKey'], data)
        assert res['message'] == 'message sent successfully'
        response = integrand.GetEventsFromTopic(topicName, 0, 1)
        assert len(response['data']) == 1
        assert response['data'][0] == data
        # Clean up
        integrand.DeleteConnector(id)
        integrand.DeleteTopic(topicName)

    def test_send_multiple_messages(self):
        id = get_random_string(5)
        topicName = get_random_string(5)
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        createResponse = integrand.CreateConnector(id, topicName)
        # Generate Data
        messages: Dict = []
        for i in range(5):
            messages.append({'key'+ str(i): 'value'+ str(i)})
        for msg in messages:
            res = integrand.EndpointRequest(id, createResponse['data']['securityKey'], msg)
            assert res['message'] == 'message sent successfully'
        response = integrand.GetEventsFromTopic(topicName, 0, 5)
        assert len(response['data']) == len(messages)
        for ind, data in enumerate(messages):
            assert response['data'][ind] == data
        # Clean up
        integrand.DeleteConnector(id)
        integrand.DeleteTopic(topicName)

class TestsWorkflow():   
    def test_create_workflow(self, workflow_function_name):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        connectorId = get_random_string(5)
        topicName = get_random_string(5)
        
        response = integrand.CreateConnector(connectorId, topicName)
        connectorAPIKey = response['data']['securityKey']
        sinkURL = rf"{INTEGRAND_URL}/api/v1/connector/f/{connectorId}?apikey={connectorAPIKey}"
        
        response = integrand.CreateWorkflow(topicName, workflow_function_name, sinkURL)
        assert response['status'] == 'success'
        assert response['data']['topicName'] == topicName
        assert response['data']['functionName'] == workflow_function_name
        assert response['data']['enabled'] == True
        id = response['data']['id']
        # Cleanup
        integrand.DeleteConnector(connectorId)
        integrand.DeleteWorkflow(id)
        integrand.DeleteTopic(topicName)
        
    def test_delete_workflow(self, workflow_function_name):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        connectorId = get_random_string(5)
        topicName = get_random_string(5)
        
        response = integrand.CreateConnector(connectorId, topicName)
        connectorAPIKey = response['data']['securityKey']
        sinkURL = rf"{INTEGRAND_URL}/api/v1/connector/f/{connectorId}?apikey={connectorAPIKey}"
        
        response = integrand.CreateWorkflow(topicName, workflow_function_name, sinkURL)
        id = response['data']['id']
       
        response = integrand.DeleteWorkflow(id)
        assert response['status'] == 'success'
        # Cleanup
        integrand.DeleteConnector(connectorId)
        integrand.DeleteTopic(topicName)
        
    def test_update_workflow(self, workflow_function_name):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        connectorId = get_random_string(5)
        topicName = get_random_string(5)
        
        response = integrand.CreateConnector(connectorId, topicName)
        connectorAPIKey = response['data']['securityKey']
        sinkURL = rf"{INTEGRAND_URL}/api/v1/connector/f/{connectorId}?apikey={connectorAPIKey}"
        
        response = integrand.CreateWorkflow(topicName, workflow_function_name, sinkURL)
        id = response['data']['id']
        integrand.UpdateWorkflow(response['data']['id'])
        
        assert response['status'] == 'success'
        # Cleanup
        integrand.DeleteConnector(connectorId)
        integrand.DeleteWorkflow(id)
        integrand.DeleteTopic(topicName)

    def test_get_workflow(self, workflow_function_name):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        connectorId = get_random_string(5)
        topicName = get_random_string(5)
        
        response = integrand.CreateConnector(connectorId, topicName)
        connectorAPIKey = response['data']['securityKey']
        sinkURL = rf"{INTEGRAND_URL}/api/v1/connector/f/{connectorId}?apikey={connectorAPIKey}"
        
        response = integrand.CreateWorkflow(topicName, workflow_function_name, sinkURL)
        id = response['data']['id']
        
        response = integrand.GetWorkflow(response['data']['id'])
        assert response['status'] == 'success'
        assert response['data']['topicName'] == topicName
        assert response['data']['functionName'] == workflow_function_name
        assert response['data']['enabled'] == True
        # Cleanup
        integrand.DeleteConnector(connectorId)
        integrand.DeleteWorkflow(id)
        integrand.DeleteTopic(topicName)
        
    def test_get_workflows(self, workflow_function_name):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        connectorId = get_random_string(5)
        topicName = get_random_string(5)
        
        response = integrand.CreateConnector(connectorId, topicName)
        connectorAPIKey = response['data']['securityKey']
        sinkURL = rf"{INTEGRAND_URL}/api/v1/connector/f/{connectorId}?apikey={connectorAPIKey}"
        
        response = integrand.GetWorkflows()
        assert response['status'] == 'success'
        assert response['data'] == []
        
        response = integrand.CreateWorkflow(topicName, workflow_function_name, sinkURL)
        id = response['data']['id']
        
        response = integrand.GetWorkflows()
        assert response['status'] == 'success'
        assert response['data'][0]['topicName'] == topicName
        assert response['data'][0]['functionName'] == workflow_function_name
        assert response['data'][0]['enabled'] == True
        # Cleanup
        integrand.DeleteConnector(connectorId)
        integrand.DeleteWorkflow(id)
        integrand.DeleteTopic(topicName)
        
    def test_workflow_send_message_one_end_to_another(self, workflow_function_name):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        sourceConnectorId = get_random_string(5)
        sinkConnectorId = get_random_string(5)
        sourceTopicName = get_random_string(5)
        sinkTopicName = get_random_string(5)
        
        response = integrand.CreateConnector(sourceConnectorId, sourceTopicName)
        sourceConnectorAPIKey = response['data']['securityKey']
        response = integrand.CreateConnector(sinkConnectorId, sinkTopicName)
        sinkConnectorAPIKey = response['data']['securityKey']
        sinkURL = rf"{INTEGRAND_URL}/api/v1/connector/f/{sinkConnectorId}?apikey={sinkConnectorAPIKey}"
        
        response = integrand.CreateWorkflow(sourceTopicName, workflow_function_name, sinkURL)
        assert response['status'] == 'success'
        assert response['data']['topicName'] == sourceTopicName
        assert response['data']['functionName'] == workflow_function_name
        assert response['data']['enabled'] == True
        id = response['data']['id']
        
        with open(rf'tests/workflow_functions/{workflow_function_name}/message.json', 'r') as file:
            data = json.load(file)
            
        integrand.EndpointRequest(sourceConnectorId, sourceConnectorAPIKey, data)
        
        # Retry fetching event for 10 seconds every 200 ms since there's a retry delay for all workflows and they're do not run concurrently
        total_duration = 10  
        interval = 0.2 
        iterations = int(total_duration / interval)

        for i in range(iterations):
            offset = 0
            limit = 1
            sinkTopicEvents = integrand.GetEventsFromTopic(sinkTopicName, offset, limit)
            if sinkTopicEvents['data'] != None:
                with open(rf'tests/workflow_functions/{workflow_function_name}/message_output.json', 'r') as file:
                    data = json.load(file)
                assert sinkTopicEvents['data'][offset] == data
                break
            time.sleep(interval)        
            
        # If for loops exhausts itself, sink haven't gotten our new message
        if i == iterations: assert False
        
        # Cleanup
        integrand.DeleteConnector(sourceConnectorId)
        integrand.DeleteConnector(sinkConnectorId)
        integrand.DeleteWorkflow(id)
        integrand.DeleteTopic(sourceTopicName)
        integrand.DeleteTopic(sinkTopicName)
        
        
class TestsUser():
    def test_create_user(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        email = "test@gmail.com"
        password = "p@ssword!"
        response = integrand.CreateUser(email, password)
        id = response['data']['id']
        assert response['data']['email'] == email
        assert response['data']['authType'] == 'email'
        # Cleanup
        integrand.DeleteUser(id)

    def test_delete_user(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        email = "test@gmail.com"
        password = "p@ssword!"
        response = integrand.CreateUser(email, password)
        id = response['data']['id']
        integrand.DeleteUser(id)
        assert response['status'] == 'success'
        
    def test_delete_root_user_raises_exception(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        try:
            integrand.DeleteUser(1)
        except HTTPError as e:
            assert e.response.status_code == 500
        
    def test_update_user_password(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        email = "test@gmail.com"
        password = "p@ssword!"
        response = integrand.CreateUser(email, password)
        id = response['data']['id']
        newPassword = "newP@ssword!"
        response = integrand.UpdateUser(id, oldPassword=password, newPassword=newPassword)
        assert response['status'] == 'success'
        # Try to authenticate with old password
        try:
             integrand.Login(email, password)
        except HTTPError as e:
            assert e.response.status_code == 401
        # Authenticate with new password
        integrand.Login(email, newPassword)
        # Cleanup
        integrand.DeleteUser(id)
        
    def test_get_all_users(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.GetUsers()
        assert response['status'] == 'success'
        assert len(response['data']) == 1
        
    def test_get_single_user(self):
        ROOT_USER_ID = 1
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        response = integrand.GetUser(ROOT_USER_ID)
        assert response['status'] == 'success'
        assert response['data']['id'] == ROOT_USER_ID
        
        
class TestsAuth():
    def test_login(self):
        integrand = Integrand(INTEGRAND_URL, INTEGRAND_API_KEY)
        email = "admin"
        password = "admin"
        http_status_code = integrand.Login(email, password)
        assert http_status_code == 200
        