package main

import "testing"

var testToken string = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Mjc1LCJ1c2VybmFtZSI6InBhYmxvY2FycmFuemEiLCJlbWFpbCI6InNvbWVvbmVAZ21haWwuY29tIiwic29jaWFsX3NlcnZpY2VfYWNjb3VudCI6W10sImhhc19kaXNhYmxlZF9uZXdzbGV0dGVyIjpmYWxzZSwiand0X3NlY3JldCI6InNlY3JldHNlY3JldCIsImhhc1Bhc3N3b3JkU2V0Ijp0cnVlLCJuZWVkc1Bhc3N3b3JkUmVzZXQiOmZhbHNlLCJwdWJsaWNfa2V5IjpmYWxzZSwiZmVhdHVyZXMiOltdLCJpbnRlcmNvbVVzZXJOYW1lIjoiW1NUQUdJTkddIHBhYmxvY2FycmFuemEiLCJpbnRlcmNvbVVzZXJIYXNoIjoic29tZWhhc2giLCJwZXJtaXNzaW9ucyI6W10sImlhdCI6MTQ2NjI3MTc5MywiZXhwIjoxNDY2ODc2NTkzfQ.uXIMNA2OTLPmlkYc_fE6iECAy6dF6c_yYwxs8yvB5eU`

func TestGetUserId(t *testing.T) {
	if id, err := getUserId(testToken); err != nil {
		t.Error(err)
	} else if id != "275" {
		t.Errorf("Got wrong userId: %d", id)
	}
}
