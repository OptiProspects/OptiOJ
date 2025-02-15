package judge_grpc_service

// Proto 文件的原始定义，用于参考和维护
const ProtoDefinition = `
syntax = "proto3";

package judge_grpc_service;

service JudgeGrpcService {
    rpc Submit(SubmitRequest) returns (SubmitResponse);
}

message SubmitRequest {
    string language = 1;
    string source_code = 2;
    string input = 3;
    string expected_output = 4;
    int64 time_limit = 5; // 以秒为单位
    uint64 memory_limit = 6; // 以字节为单位
}

message SubmitResponse {
    int32 status = 1; // JudgeStatus
    double time_used = 2; // 以秒为单位
    double memory_used = 3; // 以 KB 为单位，保留两位小数
    string error_message = 4; // 错误信息
}
`
