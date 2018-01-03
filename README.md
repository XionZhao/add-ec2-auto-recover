# 增加Ec2实例的auto recover和sns通知

> 所在机器需要配置`aws_access_key_id`和`aws_secret_access_key`

## 依赖

```
go get -u github.com/aws/aws-sdk-go
```

## Build

```
$ go build .
```

## 执行命令
```
$ ./ec2-auto-recover
Usage of ./ec2-auto-recover:
  -accountid string
    	Aws account id
  -instanceid string
    	[Optional parameters] Specify a single instance
  -oversea
    	[Optional parameters] Whether for aws overseas area
  -region string
    	aws region name
  -snstopic string
    	sns topic name
```
- accountid="": 指定aws account id
- instanceid="": 可选参数, 当只改单个实例时使用
- oversea:可选参数, 是否为海外区域, 默认为否
- region: 指定region
- snstopic: sns topic名称 

### 针对国内Aws
```
./ec2-auto-recover -accountid=99999999 -region=cn-north-1 -snstopic=snstopic
```
### 针对海外Aws
```
./ec2-auto-recover -accountid=99999999 -region=ap-southeast-1 -snstopic=snstopic -oversea
```
