package errs

var Success = New(0, "success")
var Unknown = New(1000, "system error [%s]")
var BasArgs = New(1001, "bad args [%s]")
var RpcError = New(1002, "call remote error [%s]")
var RedisError = New(1003, "call cache error [%s]")
var AttrKeyLimit = New(1004, "session attribute key[%s] limit")
var SessionInvalid = New(1005, "session invalid")
