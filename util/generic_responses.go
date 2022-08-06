package util

import "encoding/json"

var BadRequestResponse, _ = json.Marshal(map[string]string{"status": "Error", "message": "Invalid parameters"})
var UnAuthenticatedRequestResponse, _ = json.Marshal(map[string]string{"status": "UnAuthenticated", "message": "You are not allowed to perform this action, please login first."})
var NotFoundResponse, _ = json.Marshal(map[string]string{"status": "Error", "message": "Resource not found"})
