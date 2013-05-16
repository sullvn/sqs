package sqs_test

import (
	"../sqs"
	"launchpad.net/goamz/aws"
	. "launchpad.net/gocheck"
)

var _ = Suite(&S{})

type S struct {
	//testServer.PrepareResponse(200, nil, TestChangeMessageVisibilityXmlOK)
	HTTPSuite
	sqs *sqs.SQS
}

func (s *S) SetUpSuite(c *C) {
	s.HTTPSuite.SetUpSuite(c)
	auth := aws.Auth{"abc", "123"}
	s.sqs = sqs.New(auth, aws.Region{SQSEndpoint: testServer.URL})
}

func (s *S) TestCreateQueue(c *C) {
	testServer.PrepareResponse(200, nil, TestCreateQueueXmlOK)

	timeOutAttribute := sqs.Attribute{"VisibilityTimeout", "60"}
	maxMessageSizeAttribute := sqs.Attribute{"MaximumMessageSize", "65536"}
	messageRetentionAttribute := sqs.Attribute{"MessageRetentionPeriod", "345600"}
	q, err := s.sqs.CreateQueue("testQueue", []sqs.Attribute{timeOutAttribute, maxMessageSizeAttribute, messageRetentionAttribute})
	req := testServer.WaitRequest()

	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Header["Date"], Not(Equals), "")
	c.Assert(q.Url, Equals, "http://sqs.us-east-1.amazonaws.com/123456789012/testQueue")
	c.Assert(err, IsNil)
}

func (s *S) TestListQueues(c *C) {
	testServer.PrepareResponse(200, nil, TestListQueuesXmlOK)

	resp, err := s.sqs.ListQueues()
	req := testServer.WaitRequest()

	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Header["Date"], Not(Equals), "")
	c.Assert(len(resp.QueueUrl), Not(Equals), 0)
	c.Assert(resp.QueueUrl[0], Equals, "http://sqs.us-east-1.amazonaws.com/123456789012/testQueue")
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "725275ae-0b9b-4762-b238-436d7c65a1ac")
	c.Assert(err, IsNil)
}

func (s *S) TestGetQueueUrl(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	resp, err := s.sqs.GetQueueUrl("testQueue")
	req := testServer.WaitRequest()

	c.Assert(req.Method, Equals, "GET")
	c.Assert(req.URL.Path, Equals, "/")
	c.Assert(req.Header["Date"], Not(Equals), "")
	c.Assert(resp.QueueUrl, Equals, "http://sqs.us-east-1.amazonaws.com/123456789012/testQueue")
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "470a6f13-2ed9-4181-ad8a-2fdea142988e")
	c.Assert(err, IsNil)
}

func (s *S) TestChangeMessageVisibility(c *C) {
	//testServer.PrepareResponse(200, nil, TestChangeMessageVisibilityXmlOK)
	testServer.PrepareResponse(200, nil, TestCreateQueueXmlOK)

	timeOutAttribute := sqs.Attribute{"VisibilityTimeout", "60"}
	maxMessageSizeAttribute := sqs.Attribute{"MaximumMessageSize", "65536"}
	messageRetentionAttribute := sqs.Attribute{"MessageRetentionPeriod", "345600"}
	q, err := s.sqs.CreateQueue("testQueue", []sqs.Attribute{timeOutAttribute, maxMessageSizeAttribute, messageRetentionAttribute})
	req := testServer.WaitRequest()

	c.Assert(req.Method, Equals, "GET")
	testServer.PrepareResponse(200, nil, TestChangeMessageVisibilityXmlOK)
	resp, err := q.ChangeMessageVisibility("MbZj6wDWli%2BJvwwJaBV%2B3dcjk2YW2vA3%2BSTFFljT", 0)
	testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "6a7a282a-d013-4a59-aba9-335b0fa48bed")
}

func (s *S) TestChangeMessageVisibilityBatch(c *C) {
	testServer.PrepareResponse(200, nil, TestCreateQueueXmlOK)

	timeOutAttribute := sqs.Attribute{"VisibilityTimeout", "60"}
	maxMessageSizeAttribute := sqs.Attribute{"MaximumMessageSize", "65536"}
	messageRetentionAttribute := sqs.Attribute{"MessageRetentionPeriod", "345600"}
	q, err := s.sqs.CreateQueue("testQueue", []sqs.Attribute{timeOutAttribute, maxMessageSizeAttribute, messageRetentionAttribute})
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestChangeMessaveVisibilityBatchXmlOK)

	messageVisibilityBatch := []sqs.ChangeMessageVisibilityBatchEntry{sqs.ChangeMessageVisibilityBatchEntry{"change_visibility_msg_2", "gfk0T0R0waama4fVFffkjKzmhMCymjQvfTFk2LxT33G4ms5subrE0deLKWSscPU1oD3J9zgeS4PQQ3U30qOumIE6AdAv3w%2F%2Fa1IXW6AqaWhGsEPaLm3Vf6IiWqdM8u5imB%2BNTwj3tQRzOWdTOePjOjPcTpRxBtXix%2BEvwJOZUma9wabv%2BSw6ZHjwmNcVDx8dZXJhVp16Bksiox%2FGrUvrVTCJRTWTLc59oHLLF8sEkKzRmGNzTDGTiV%2BYjHfQj60FD3rVaXmzTsoNxRhKJ72uIHVMGVQiAGgBX6HGv9LDmYhPXw4hy%2FNgIg%3D%3D", 45}, sqs.ChangeMessageVisibilityBatchEntry{"change_visibility_msg_3", "gfk0T0R0waama4fVFffkjKzmhMCymjQvfTFk2LxT33FUgBz3%2BnougdeLKWSscPU1%2FXgx%2BxcNnjnQQ3U30qOumIE6AdAv3w%2F%2Fa1IXW6AqaWhGsEPaLm3Vf6IiWqdM8u5imB%2BNTwj3tQRzOWdTOePjOsogjZM%2F7kzn4Ew27XLU9I%2FYaWYmKvDbq%2Fk3HKVB9HfB43kE49atP2aWrzNL4yunG41Q4cfRRtfJdcGQGNHQ2%2Byd0Usf5qR1dZr1iDo5xk946eQat83AxTRP%2BY4Qi0V7FAeSLH9su9xpX6HGv9LDmYhPXw4hy%2FNgIg%3D%3D", 45}}
	resp, err := q.ChangeMessageVisibilityBatch(messageVisibilityBatch)
	testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "ca9668f7-ab1b-4f7a-8859-f15747ab17a7")
	c.Assert(resp.Id[0], Equals, "change_visibility_msg_2")
	c.Assert(resp.Id[1], Equals, "change_visibility_msg_3")
}

func (s *S) TestReceiveMessage(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestReceiveMessageXmlOK)

	resp, err := q.ReceiveMessage([]string{"All"}, 5, 15)
	testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(len(resp.Messages), Not(Equals), 0)
}

func (s *S) TestDeleteMessage(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestDeleteMessageXmlOK)

	resp, err := q.DeleteMessage("MbZj6wDWli%2BJvwwJaBV%2B3dcjk2YW2vA3%2BSTFFljTM8tJJg6HRG6PYSasuWXPJB%2BCwLj1FjgXUv1uSj1gUPAWV66FU/WeR4mq2OKpEGYWbnLmpRCJVAyeMjeU5ZBdtcQ%2BQEauMZc8ZRv37sIW2iJKq3M9MFx1YvV11A2x/KSbkJ0=")
	testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "b5293cb5-d306-4a17-9048-b263635abe42")
}

func (s *S) TestDeleteMessageBatch(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestDeleteMessageBatchXmlOK)

	deleteMessageBatch := []sqs.DeleteMessageBatch{sqs.DeleteMessageBatch{Id: "msg1", ReceiptHandle: "gfk0T0R0waama4fVFffkjPQrrvzMrOg0fTFk2LxT33EuB8wR0ZCFgKWyXGWFoqqpCIiprQUEhir%2F5LeGPpYTLzjqLQxyQYaQALeSNHb0us3uE84uujxpBhsDkZUQkjFFkNqBXn48xlMcVhTcI3YLH%2Bd%2BIqetIOHgBCZAPx6r%2B09dWaBXei6nbK5Ygih21DCDdAwFV68Jo8DXhb3ErEfoDqx7vyvC5nCpdwqv%2BJhU%2FTNGjNN8t51v5c%2FAXvQsAzyZVNapxUrHIt4NxRhKJ72uICcxruyE8eRXlxIVNgeNP8ZEDcw7zZU1Zw%3D%3D"}, sqs.DeleteMessageBatch{Id: "msg2", ReceiptHandle: "gfk0T0R0waama4fVFffkjKzmhMCymjQvfTFk2LxT33G4ms5subrE0deLKWSscPU1oD3J9zgeS4PQQ3U30qOumIE6AdAv3w%2F%2Fa1IXW6AqaWhGsEPaLm3Vf6IiWqdM8u5imB%2BNTwj3tQRzOWdTOePjOjPcTpRxBtXix%2BEvwJOZUma9wabv%2BSw6ZHjwmNcVDx8dZXJhVp16Bksiox%2FGrUvrVTCJRTWTLc59oHLLF8sEkKzRmGNzTDGTiV%2BYjHfQj60FD3rVaXmzTsoNxRhKJ72uIHVMGVQiAGgB%2BqAbSqfKHDQtVOmJJgkHug%3D%3D"}}
	resp, err := q.DeleteMessageBatch(deleteMessageBatch)
	testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(len(resp.DeleteMessageBatchResult.Ids), Equals, 2)
	c.Assert(resp.DeleteMessageBatchResult.Ids[0], Equals, "msg1")
	c.Assert(resp.DeleteMessageBatchResult.Ids[1], Equals, "msg2")
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "d6f86b7a-74d1-4439-b43f-196a1e29cd85")
}

func (s *S) TestAddPermission(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestAddPermissionXmlOK)
	resp, err := q.AddPermission("testLabel", []sqs.AccountPermission{sqs.AccountPermission{"125074342641", "SendMessage"}, sqs.AccountPermission{"125074342642", "ReceiveMessage"}})
	testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "9a285199-c8d6-47c2-bdb2-314cb47d599d")
}

func (s *S) TestRemovePermission(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestRemovePermissionXmlOK)
	resp, err := q.RemovePermission("testLabel")
	testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "f8bdb362-6616-42c0-977a-ce9a8bcce3bb")
}

func (s *S) TestSendMessage(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestSendMessageXmlOK)

	resp, err := q.SendMessage("This is a Message")
	testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(resp.SendMessageResult.MD5OfMessageBody, Equals, "fafb00f5732ab283681e124bf8747ed1")
	c.Assert(resp.SendMessageResult.MessageId, Equals, "5fea7756-0ea4-451a-a703-a558b933e274")
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "27daac76-34dd-47df-bd01-1f6e873584a0")
}

func (s *S) TestSendMessageWithDelay(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestSendMessageXmlOK)

	resp, err := q.SendMessageWithDelay("This is a Message", 60)
	testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(resp.SendMessageResult.MD5OfMessageBody, Equals, "fafb00f5732ab283681e124bf8747ed1")
	c.Assert(resp.SendMessageResult.MessageId, Equals, "5fea7756-0ea4-451a-a703-a558b933e274")
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "27daac76-34dd-47df-bd01-1f6e873584a0")
}

func (s *S) TestSendMessageBatch(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestSendMessageBatchXmlOK)

	sendMessageBatchRequests := []sqs.SendMessageBatchRequestEntry{sqs.SendMessageBatchRequestEntry{Id: "test_msg_001", MessageBody: "test message body 1", DelaySeconds: 30}}
	resp, err := q.SendMessageBatch(sendMessageBatchRequests)
	testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(len(resp.SendMessageBatchResult.Entries), Equals, 2)
	c.Assert(resp.SendMessageBatchResult.Entries[0].Id, Equals, "test_msg_001")
}

func (s *S) TestGetQueueAttributes(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestGetQueueAttributesAllXmlOK)

	resp, err := q.GetQueueAttributes([]string{"ALL"})
	testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(len(resp.Attributes), Equals, 8)
	c.Assert(resp.Attributes[0].Name, Equals, "VisibilityTimeout")
	c.Assert(resp.Attributes[0].Value, Equals, "30")
}

func (s *S) TestGetQueueAttributesSelective(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestGetQueueAttributesSelectiveXmlOK)

	resp, err := q.GetQueueAttributes([]string{"VisibilityTimeout", "DelaySeconds"})
	testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(len(resp.Attributes), Equals, 2)
	c.Assert(resp.Attributes[0].Name, Equals, "VisibilityTimeout")
	c.Assert(resp.Attributes[0].Value, Equals, "30")
	c.Assert(resp.Attributes[1].Name, Equals, "DelaySeconds")
	c.Assert(resp.Attributes[1].Value, Equals, "0")
}

func (s *S) TestDeleteQueue(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestDeleteQueueXmlOK)
	resp, err := q.Delete()
	testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "6fde8d1e-52cd-4581-8cd9-c512f4c64223")
}

func (s *S) TestSetQueueAttributes(c *C) {
	testServer.PrepareResponse(200, nil, TestGetQueueUrlXmlOK)

	q, err := s.sqs.GetQueue("testQueue")
	testServer.WaitRequest()

	testServer.PrepareResponse(200, nil, TestSetQueueAttributesXmlOK)
	var policyStr = `
  {
        "Version":"2008-10-17",
        "Id":"/123456789012/testQueue/SQSDefaultPolicy",
        "Statement":  [
             {
             "Sid":"Queue1ReceiveMessage",
             "Effect":"Allow",
             "Principal":{"AWS":"*"},
             "Action":"SQS:ReceiveMessage",
             "Resource":"arn:aws:sqs:us-east-1:123456789012:testQueue"
              }
         ]    
   }
  `
	resp, err := q.SetQueueAttributes(sqs.Attribute{"Policy", policyStr})
	testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp.ResponseMetadata.RequestId, Equals, "e5cca473-4fc0-4198-a451-8abb94d02c75")
}
