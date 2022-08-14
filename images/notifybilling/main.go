package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/aws/jsii-runtime-go"
)

const lineNotifyURL = "https://notify-api.line.me/api/notify"

type Response struct {
	Code float64 `json:"code"`
	Message string `json:"message"`
}

type CostPerService struct {
	Service string
	Cost float64
}

// LINE通知を行う関数
func notifyToLINE(message *string) error {
	// アクセストークンを環境変数から取得する
	lineNotifyToken, ok := os.LookupEnv("LINE_NOTIFY_TOKEN")
	if !ok {
		return fmt.Errorf("notifyToLINE: NOT FOUND LINE_NOTIFY_TOKEN")
	}

	// LINE通知のリクエストを作成する
	// 詳しくは[LINE Notify API Document](https://notify-bot.line.me/doc/ja/)をご参照ください
	body := strings.NewReader(url.Values{
		"message": []string{*message},
	}.Encode())
	req, err := http.NewRequest(http.MethodPost, lineNotifyURL, body)
	if err != nil {
		return fmt.Errorf("notifyToLINE: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", lineNotifyToken))

	// リクエストを実行する
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
    	return fmt.Errorf("notifyToLINE: %w", err)
	}	
	defer resp.Body.Close() // resp.BodyをクローズしないとTCPコネクションが張られたままになる
	if resp.StatusCode != 200 {
    	return fmt.Errorf("notifyToLINE: %s", resp.Status)
	}

	return nil
}

// LINE通知で送るメッセージを整形する関数
func toMessage(CostsPerService *[]CostPerService) *string {
	// コストの大きい順に並べ替える
	sort.Slice(*CostsPerService, func(i, j int) bool {
		return (*CostsPerService)[i].Cost > (*CostsPerService)[j].Cost
	})

	// 合計コストを算出しつつ、Taxをスライスから除外する
	var totalCost float64
	var tax float64
	for i, v := range *CostsPerService {
		totalCost += v.Cost
		if v.Service == "Tax" {
			tax = v.Cost
			*CostsPerService = append((*CostsPerService)[:i], (*CostsPerService)[i+1:]...) // Taxをsliceから除外する
		}
	}

	// メッセージを整形する
	message := "\n***********************"
	message += fmt.Sprintf("\nTotal: $%.1f (Tax: $%.1f)", totalCost, tax)
	message += "\n***********************"
	for _, v := range *CostsPerService {
		message += fmt.Sprintf("\n%s: $%.1f", v.Service, v.Cost)
	}

	return &message
}

// 指定されて期間におけるAWS利用額を取得する関数
func getCostAndUsagePerService(CostExplorer *costexplorer.CostExplorer, startDate *string, endDate *string) (*[]CostPerService, error) {
	// 指定されて期間におけるAWS利用額をサービス単位で取得する
	output, err := CostExplorer.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: startDate,
			End: endDate,
		},
		GroupBy: []*costexplorer.GroupDefinition{
			{Key: jsii.String("SERVICE"), Type: jsii.String("DIMENSION"),},
		},
		Granularity: jsii.String("MONTHLY"),
		Metrics: []*string{jsii.String("BlendedCost"),},
	})
	if err != nil {
		return nil, fmt.Errorf("getTotalCostAndUsage: %w", err)
	}

	// 取得した結果をパースする
	var CostsPerService []CostPerService
	for _, result := range output.ResultsByTime {
		for _, group := range result.Groups {
			if *group.Metrics["BlendedCost"].Amount != "0" {
				CostFloat64, err := strconv.ParseFloat(*group.Metrics["BlendedCost"].Amount, 64)
				if err != nil {
					return nil, fmt.Errorf("getCostAndUsagePerService: %w", err)
				}
				CostsPerService = append(CostsPerService, CostPerService{Service: *group.Keys[0], Cost: CostFloat64})
			}
		}
	}

	return &CostsPerService, nil
}

// ハンドラ関数 = Lambda関数として実行されるロジック
func handleRequest() (Response, error) {
	session := session.Must(session.NewSession())
	CostExplorer := costexplorer.New(session, aws.NewConfig().WithRegion("ap-northeast-1"))

	// 月初から現在日におけるAWS利用額をサービス単位で取得する
	currentTime := time.Now()
	firstDayOfMonth := time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	today := currentTime.Format("2006-01-02")
	CostsPerService, err := getCostAndUsagePerService(CostExplorer, &firstDayOfMonth, &today)
	if err != nil {
		return Response{Code: 500, Message: "Internal Server Error"}, err
	}

	// 取得したAWS利用額を整形してLINE通知する
	message := toMessage(CostsPerService)
	err = notifyToLINE(message)
	if err != nil {
		return Response{Code: 500, Message: "Internal Server Error"}, err
	}

	return Response{Code: 200, Message: "OK"}, nil
}

// Lambda関数のコードが実行されるエントリポイント
func main() {
	lambda.Start(handleRequest)
}
