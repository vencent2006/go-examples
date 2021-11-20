/**
 * @Author: vincent
 * @Description:
 * @File:  middleware
 * @Version: 1.0.0
 * @Date: 2021/11/7 18:11
 */

package command

import (
	"bytes"
	"fmt"
	"go-examples/course/handwriting-web-inf/code_27/framework/cobra"
	"go-examples/course/handwriting-web-inf/code_27/framework/contract"
	"go-examples/course/handwriting-web-inf/code_27/framework/util"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/jianfengye/collection"

	"github.com/go-git/go-git/v5"

	"github.com/AlecAivazis/survey/v2"
)

// 初始化中间件相关命令
func initMiddlewareCommand() *cobra.Command {
	middlewareCommand.AddCommand(middlewareListCommand)
	middlewareCommand.AddCommand(middlewareMigrateCommand)
	middlewareCommand.AddCommand(middlewareCreateCommand)
	return middlewareCommand
}

// middlewareCommand 中间件二级命令
var middlewareCommand = &cobra.Command{
	Use:   "middleware",
	Short: "中间件相关命令",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

// middlewareAllCommand 显示所有安装的中间件
var middlewareListCommand = &cobra.Command{
	Use:   "list",
	Short: "显示所有中间件",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		middlewarePath := path.Join(appService.BaseFolder(), "app", "http", "middleware")

		// 读取文件夹
		files, err := ioutil.ReadDir(middlewarePath)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.IsDir() {
				fmt.Println(f.Name())
			}
		}
		return nil
	},
}

// 从gin-contrib中迁移中间件
var middlewareMigrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "迁移gin-contrib中间件，迁移地址: https://github.com/gin-contrib/[middleware].git",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("迁移一个Gin中间件")
		var repo string
		{
			prompt := &survey.Input{
				Message: "请输入中间件名称:",
			}
			err := survey.AskOne(prompt, &repo)
			if err != nil {
				return err
			}
		}

		// step2 : 下载git到一个目录中
		appService := container.MustMake(contract.AppKey).(contract.App)

		middlewarePath := appService.MiddlewareFolder()
		url := "https://github.com/gin-contrib/" + repo + ".git"
		fmt.Println("下载中间件 gin-contrib:")
		fmt.Println(url)
		_, err := git.PlainClone(path.Join(middlewarePath, repo), false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}

		// step3:删除不必要的文件 go.mod, go.sum. .git
		repoFolder := path.Join(middlewarePath, repo)
		fmt.Println("remove " + path.Join(repoFolder, "go.mod"))
		os.Remove(path.Join(repoFolder, "go.mod"))
		fmt.Println("remove " + path.Join(repoFolder, "go.sum"))
		os.Remove(path.Join(repoFolder, "go.sum"))
		fmt.Println("remove " + path.Join(repoFolder, ".git"))
		os.Remove(path.Join(repoFolder, ".git"))

		// step4: 替换关键词
		filepath.Walk(repoFolder, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				// 是目录
				return nil
			}

			if filepath.Ext(path) != ".go" {
				// 不是.go文件
				return nil
			}

			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			isContain := bytes.Contains(c, []byte("github.com/gin-gonic/gin"))
			if isContain {
				fmt.Println("更新文件:" + path)
				// todo: 替换文件的内容，主要是gin的引用
				c = bytes.ReplaceAll(c, []byte("github.com/gin-gonic/gin"), []byte("github.com/gohade/hade/framework/gin"))
				// todo: 注意这里的权限是0644
				err = ioutil.WriteFile(path, c, 0644)
				if err != nil {
					return err
				}
			}
			return nil
		})
		return nil
	},
}

// providerCreateCommand 创建一个新的服务，包括服务提供者，服务接口协议，服务实例
var middlewareCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"},
	Short:   "创建一个中间件",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println("创建一个中间件")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入中间件名称",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入中间件所在目录名称(默认:同中间件名称):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}

		if folder == "" {
			folder = name
		}

		app := container.MustMake(contract.AppKey).(contract.App)

		pFolder := app.MiddlewareFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		subCollection := collection.NewStrCollection(subFolders)
		if subCollection.Contains(folder) {
			fmt.Println("目录已经存在")
			return nil
		}

		// 开始创建文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}
		funcs := template.FuncMap{"title": strings.Title}
		{
			// 创建
			file := filepath.Join(pFolder, folder, "middleware.go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			t := template.Must(template.New("middleware").Funcs(funcs).Parse(middlewareTmp))
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}
		fmt.Println("创建中间件成功, 文件夹地址:", filepath.Join(pFolder, folder))
		return nil
	},
}

var middlewareTmp string = `package {{.}}

import "github.com/gohade/hade/framework/gin"

// {{.|title}}Middleware 代表中间件函数
func {{.|title}}Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
	}
}

`
