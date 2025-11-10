package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/redis/go-redis/v9"
	"time"
)

func initRedisClient(ctx context.Context, svc *checkoutService) {

	redisAddr := "redis-13783.crce179.ap-south-1-1.ec2.redns.redis-cloud.com:13783"
	// Optionally fetch credentials from AWS Secrets Manager if configured
	//secretName := os.Getenv("arn:aws:kms:ap-south-1:038184794282:key/0f925495-fbf7-4937-a5b3-5124a5ab18b0")
	//awsRegion := os.Getenv("ap-south-1")
	//if secretName != "" && awsRegion != "" {
	//if _, err := getRedisCredentialsFromAWS(ctx, secretName, awsRegion); err != nil {
	//	log.Warnf("failed to fetch redis credentials from secrets manager: %v", err)
	//} else {
	username := "default"
	password := "d2Vvv7ORCsdwvWDAWNs4jOEHVTkUYxee"
	//	}
	//}

	fmt.Println("username:", username)
	fmt.Println("password:", password)

	rc := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: username,
		Password: password,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := rc.Ping(pingCtx).Err(); err != nil {
		log.Warnf("redis not available at %q: %v", redisAddr, err)
		return
	}

	log.Infof("connected to redis at %q", redisAddr)
	svc.redisClient = rc
}

func getRedisCredentialsFromAWS(ctx context.Context, secretName, region string) (*RedisSecret, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := secretsmanager.NewFromConfig(cfg)
	resp, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	if resp.SecretString == nil {
		return nil, fmt.Errorf("secret %s has no SecretString", secretName)
	}
	var out RedisSecret
	if err := json.Unmarshal([]byte(*resp.SecretString), &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret JSON: %w", err)
	}
	return &out, nil
}
