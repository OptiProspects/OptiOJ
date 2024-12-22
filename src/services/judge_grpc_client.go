package services

import (
	"OptiOJ/src/models"
	pb "OptiOJ/src/proto/judge_grpc_service"
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type JudgeGrpcClient struct {
	client pb.JudgeGrpcServiceClient
	conn   *grpc.ClientConn
}

var judgeGrpcClient *JudgeGrpcClient

func InitJudgeGrpcClient() error {
	logrus.Info("开始初始化判题客户端...")

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Errorf("连接判题服务失败: %v", err)
		return fmt.Errorf("无法连接到判题服务: %v", err)
	}

	judgeGrpcClient = &JudgeGrpcClient{
		client: pb.NewJudgeGrpcServiceClient(conn),
		conn:   conn,
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 尝试一个简单的ping请求来验证连接
	_, err = judgeGrpcClient.client.Submit(ctx, &pb.SubmitRequest{})
	if err != nil {
		logrus.Warnf("判题服务连接测试失败: %v", err)
	} else {
		logrus.Info("判题服务连接测试成功")
	}

	logrus.Info("判题客户端初始化完成")
	return nil
}

func GetJudgeClient() *JudgeGrpcClient {
	if judgeGrpcClient == nil {
		logrus.Error("判题客户端未初始化，尝试重新初始化...")
		if err := InitJudgeGrpcClient(); err != nil {
			logrus.Errorf("重新初始化判题客户端失败: %v", err)
			return nil
		}
	}
	return judgeGrpcClient
}

func (c *JudgeGrpcClient) Submit(config *models.JudgeConfig, testCases []models.TestCase) (*models.RunResult, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("判题客户端未初始化")
	}

	// 转换测试用例格式
	protoTestCases := make([]*pb.TestCase, len(testCases))
	for i, tc := range testCases {
		// 读取输入文件
		input, err := ioutil.ReadFile(tc.InputFile)
		if err != nil {
			return nil, fmt.Errorf("读取输入文件失败: %v", err)
		}

		// 读取期望输出文件
		expectedOutput, err := ioutil.ReadFile(tc.OutputFile)
		if err != nil {
			return nil, fmt.Errorf("读取输出文件失败: %v", err)
		}

		protoTestCases[i] = &pb.TestCase{
			Input:          string(input),
			ExpectedOutput: string(expectedOutput),
		}
	}

	// 创建请求
	req := &pb.SubmitRequest{
		Language:    config.Language,
		SourceCode:  config.Code,
		TimeLimit:   int32(config.TimeLimit),   // 毫秒
		MemoryLimit: int32(config.MemoryLimit), // MB
		TestCases:   protoTestCases,
	}

	// 使用带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 发送请求
	resp, err := c.client.Submit(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("提交判题请求失败: %v", err)
	}

	// 将 gRPC 响应转换为 RunResult
	result := &models.RunResult{
		Status:       convertGrpcStatus(resp.Status),
		TimeUsed:     int(resp.TimeUsed),   // 毫秒
		MemoryUsed:   int(resp.MemoryUsed), // KB
		ErrorMessage: resp.ErrorMessage,
	}

	// 转换每个测试点的结果
	result.TestCaseResults = make([]models.TestCaseResult, len(resp.TestCaseResults))
	for i, tcResult := range resp.TestCaseResults {
		result.TestCaseResults[i] = models.TestCaseResult{
			Status:       convertGrpcStatus(tcResult.Status),
			TimeUsed:     float64(tcResult.TimeUsed),
			MemoryUsed:   float64(tcResult.MemoryUsed),
			ActualOutput: tcResult.ActualOutput,
			TestCaseID:   int(tcResult.TestCaseId),
		}
	}

	return result, nil
}

// convertGrpcStatus 将 gRPC 状态码转换为系统状态
func convertGrpcStatus(status int32) string {
	switch status {
	case 0:
		return models.StatusAccepted
	case 1:
		return models.StatusWrongAnswer
	case 2:
		return models.StatusTimeLimitExceed
	case 3:
		return models.StatusMemoryLimitExceed
	case 4:
		return models.StatusRuntimeError
	default:
		return models.StatusSystemError
	}
}
