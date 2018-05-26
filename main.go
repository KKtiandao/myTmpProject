package main

import (
	"github.com/spf13/pflag"
	"fmt"
	"crypto/hmac"
	"bufio"
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
	"strings"
	"context"
	gh "github.com/google/go-github/github"

	"crypto/sha1"
)

var (
	argJenkinsHost     = pflag.String("jenkins-server", "http://10.110.22.55:8080/", "jenkins server address default 'registry.cluster96.com:5000'")
	argJenkinsUsername = pflag.String("jenkins-username", "admin", "jenkins server username ")
	argJenkinsPassword = pflag.String("jenkins-password", "123456a?", "jenkins server password ")
)

var url = "http://127.0.0.1:8080/job/test1/"

//var CreatedAt *time.Time=
var hookEvents = []string{"push", "pull_request"}
var hookIsActive = bool(true)
var hookName = string("web")
var hookId int64 = 1

type HmacInfo struct {
	selfSaved string
	sceret    string
	password  []byte
}

func main() {
	/*
		jenkinsConf := gojenkins.JenkinsConf{*argJenkinsHost,*argJenkinsUsername, *argJenkinsPassword}
		jenkins := gojenkins.CreateJenkins(nil, jenkinsConf.Base, jenkinsConf.Username, jenkinsConf.Password)
		_, jenkinsInitErr := jenkins.Init()
		if (nil != jenkinsInitErr) {
			fmt.Print("main :jenkins Init Err")
			return
		}
		constructEntity := new(models.Construct)
		pipelineEntity := new(models.Pipeline)
		construct.CreateConstruct(constructEntity,jenkins,12)
		//生成token
			jobToken := jenkins.GetSecurerandom()
		//生成jobname
		jobName := "testForGitHub"

		//调用生成job xml
		jobConf,err := template.ParseJenkinsFile(constructEntity,pipelineEntity,jobToken)
		if  err != nil {
			fmt.Print("main :ParseJenkinsFile")
			return
		}

		//创建job
		_,joberr :=jenkins.CreateJob(jobConf,jobName)
		if  joberr != nil {
			fmt.Print("main :CreateJob")
			return
		}

	*/

	//创建GitHubClient，并注册ReposHook。
	r := bufio.NewReader(os.Stdin)
	fmt.Print("GitHub Username: ")
	username, _ := r.ReadString('\n')

	fmt.Print("GitHub Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	tp := gh.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client := gh.NewClient(tp.Client())
	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "KKtiandao")
	fmt.Printf("Client Info ,client :\n%#v\n,user :\n%#v\n, err: \n%#v\n", client, user, err)

	sceret := string("123456a?")
	var hookConfig = map[string]interface{}{"url": url, "content_type": "jason","secret":sceret}
	hmacProcess(sceret)
	input := &gh.Hook{Name: &(hookName), Events: hookEvents, Active: &hookIsActive, Config: hookConfig, ID: &hookId}

	fmt.Printf("Hook before created \n%#v\n", input)

	hook, response, err := client.Repositories.CreateHook(context.Background(), "KKtiandao", "Tesla", input)

	if err != nil {
		fmt.Printf("create hook err :\n%#v\n", err)
	}
	fmt.Printf("Hook after created :hook Info:\n%#v\n, response:\n%#v\n, err: \n%#v\n", hook, response, err)

}


func hmacProcess(sceret string) (HmacInfo) {
	hmacInfo := HmacInfo{"123456", sceret,nil}
	tmpKey := fmt.Sprintf("%s%s", hmacInfo.selfSaved, hmacInfo.sceret)
	fmt.Printf("HMAC key: \n%s\n", tmpKey)
	b := hmac.New(sha1.New, []byte(tmpKey))
	fmt.Printf("HMAC generated: \n%s\n", b)
	hmacInfo.password = b.Sum(nil)
	fmt.Printf("HMAC generated: \n%x\n",hmacInfo.password)
	return hmacInfo
}

/*
func  formatMilliSecond(int int64) string {
	var ss int64 = 1000
	var mi int64 = ss * 60
	var hh int64 = mi * 60
	var dd int64 = hh * 24

	var day = int / dd
	var hour = (int - day * dd) / hh
	var minute = (int - day * dd - hour * hh) / mi
	var second = (int - day * dd - hour * hh - minute * mi) / ss
	//var milliSecond = int - day * dd - hour * hh - minute * mi - second * ss

	var returnString = ""
	if(day > 0) {
		returnString = returnString+ strconv.FormatInt(day,10) +"天"
	}
	if(hour > 0) {
		returnString = returnString+ strconv.FormatInt(hour,10) +"小时"

	}
	if(minute > 0) {
		returnString = returnString+ strconv.FormatInt(minute,10) +"分"

	}
	if(second > 0) {
		returnString = returnString+ strconv.FormatInt(second,10) +"秒"

	}
	return returnString
}


func makeCmd(construct *models.Construct){
	if(construct.BaseImage == "spring_boot"){
		opts :=strings.Split(construct.CmdOpts, " ")

		construct.Cmd =  "CMD [\"java\", \"-jar\", \"/usr/lib/" + construct.ArtifactName +"\""
		for _,opt :=range opts{
			if(opt != ""){
				opt = strings.Replace(opt, "$","\\$",-1)
				construct.Cmd = construct.Cmd + ",\"" + opt +"\""
			}


		}
		construct.Cmd = construct.Cmd +"]"
	}
}
*/
