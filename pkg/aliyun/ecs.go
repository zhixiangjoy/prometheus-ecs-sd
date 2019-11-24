package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/seanly/prometheus-ecs-sd/pkg/config"
	"strings"
)

type EcsClient struct {
	sdConfig *config.SDConfig
	client   ecs.Client
}

func (d *EcsClient)DescribeInstanceAttribute(instanceId string) (response *ecs.DescribeInstanceAttributeResponse, err error){
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.Scheme = "https"

	request.InstanceId = instanceId

	response, err = d.client.DescribeInstanceAttribute(request)
	return
}

func splitTag(tag string) *[]ecs.DescribeInstancesTag {
	tags := strings.Split(tag, ",")

	ret := make([]ecs.DescribeInstancesTag, 0)
	for _, t := range tags {
		t2 := strings.Split(t, ":")
		if len(t2) != 2{
			continue
		}
		tm := ecs.DescribeInstancesTag{
			Key: 	t2[0],
			Value:  t2[1],
		}
		ret = append(ret, tm)
	}
	return &ret
}

func(d *EcsClient)DescribeInstances(pageNumber, pageSize int) (response *ecs.DescribeInstancesResponse, err error) {
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	request.Status = "Running"
	request.InstanceNetworkType = "vpc"

	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)


	for _, f := range d.sdConfig.Filters {
		switch f.Name {
		case "InstanceIds":
			request.InstanceIds = f.Value
		case "Status":
			request.Status = f.Value
		case "Tag":
			request.Tag = splitTag(f.Value)
		case "InstanceName":
			request.InstanceName = f.Value
		}
	}

	response, err = d.client.DescribeInstances(request)
	return
}


func totalPage(record, size int) int {
	tp := record / size
	if record%size == 0 {
		return tp
	}
	return tp + 1
}

func (d *EcsClient)GetServers() (srvs *[]ecs.Instance, err error) {
	var servers []ecs.Instance
	const pageSize = 100
	pageIndex := 1
	response, err := d.DescribeInstances(pageIndex, pageSize)

	if err != nil {
		return nil, err
	}

	for ;pageIndex <= totalPage(response.TotalCount, pageSize); pageIndex++ {
		response, err = d.DescribeInstances(pageIndex, pageSize)
		if err != nil {
			return nil, err
		}
		servers = append(servers, response.Instances.Instance...)
	}

	srvs = &servers
	return
}

func NewEcsClient(config *config.SDConfig) *EcsClient {
	ecsClient, _ := ecs.NewClientWithAccessKey(config.Region, config.AccessKey, config.SecretKey)
	return &EcsClient{
		sdConfig: config,
		client:   *ecsClient,
	}
}
