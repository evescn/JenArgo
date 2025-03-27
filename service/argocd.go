package service

import (
	"JenArgo/middleware"
	"JenArgo/model/bo"
	"JenArgo/model/po"
	"JenArgo/model/vo"
	"JenArgo/settings"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	neturl "net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

var ArgoCD argoCD

var mt sync.Mutex

type argoCD struct {
	Token          string    `json:"token"`
	LastUpdateTime time.Time `json:"-"`
}

// Session 获取 Token 信息
func (argo *argoCD) Session() error {
	url := fmt.Sprintf("%s/api/v1/session", settings.Conf.ArgoCD.ArgoCDUrl)
	userInfo := &po.ArgoCDUserInfo{
		Name:     settings.Conf.ArgoCD.Name,
		Password: settings.Conf.ArgoCD.Password,
	}

	body, err := middleware.Request.HttpRequest("POST", "argo", url, "", userInfo)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, argo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return err
	}

	argo.LastUpdateTime = time.Now().Local()

	return nil
}

// Applications 获取 app 信息
func (argo *argoCD) Applications(name, namespace string, page, size int) (*bo.AppsInfo, error) {
	url := fmt.Sprintf("%s/api/v1/applications?name=%s", settings.Conf.ArgoCD.ArgoCDUrl, name)

	body, err := middleware.Request.HttpRequest("GET", "argo", url, argo.Token, nil)
	if err != nil {
		return nil, err
	}

	appsInfo := &bo.AppsInfo{}

	if err := json.Unmarshal(body, appsInfo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return nil, err
	}

	// body 数据基于 namespace, page, size 分页
	if namespace != "" {
		// 过滤 namespace
		filteredItems := &bo.AppsInfo{}
		for _, item := range appsInfo.Items {
			if item.Spec.Destination.Namespace == namespace {
				filteredItems.Items = append(filteredItems.Items, item)
			}
		}

		appsInfo.Items = filteredItems.Items
	}

	// Implement pagination
	start := (page - 1) * size
	end := start + size

	if start < len(appsInfo.Items) {
		if end > len(appsInfo.Items) {
			end = len(appsInfo.Items)
		}
		appsInfo.Items = appsInfo.Items[start:end]
	} else {
		appsInfo.Items = nil
	}

	return appsInfo, nil
}

// Image 获取单个 app 信息
func (argo *argoCD) Image(name string) (*vo.ArgoCDImageResponse, error) {
	url := fmt.Sprintf("%s/api/v1/applications/%s", settings.Conf.ArgoCD.ArgoCDUrl, name)

	body, err := middleware.Request.HttpRequest("GET", "argo", url, argo.Token, nil)
	if err != nil {
		return nil, err
	}

	appInfo := &bo.AppInfo{}

	if err := json.Unmarshal(body, appInfo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return nil, err
	}

	// go 调用
	var wg sync.WaitGroup
	wg.Add(len(appInfo.Status.History))
	errs := make([]error, 0)
	data := &vo.ArgoCDImageResponse{}

	for _, history := range appInfo.Status.History {
		source := &vo.Source{
			AppName: appInfo.Metadata.Name,
			Helm: struct {
				ReleaseName string   `json:"releaseName"`
				ValueFiles  []string `json:"valueFiles"`
			}{
				ReleaseName: history.Source.Helm.ReleaseName,
				ValueFiles:  history.Source.Helm.ValueFiles,
			},
			Path:           history.Source.Path,
			RepoURL:        history.Source.RepoURL,
			TargetRevision: history.Revision,
		}
		imageRequest := &vo.ArgoCDImageRequest{
			AppName:    appInfo.Metadata.Name,
			AppProject: "default",
			AppID:      history.Id,
			Source:     source,
		}
		go func() {
			imageInfo, err := argo.appImage(imageRequest)
			if err != nil {
				errs = append(errs, err)
			}
			addMap(data, imageInfo)
			wg.Done()
		}()
	}
	wg.Wait()

	data = sortData(data)

	return data, nil
}

func addMap(data *vo.ArgoCDImageResponse, imageInfo *bo.ImageInfo) {
	mt.Lock()
	defer mt.Unlock()
	// Check if image already exists
	for _, existing := range data.Images {
		if existing.Image == imageInfo.Image {
			return // Skip duplicate image
		}
	}
	data.Images = append(data.Images, &bo.ImageInfo{
		ID:    imageInfo.ID,
		Image: imageInfo.Image,
	})
}

// extractTag 提取 image 中的 tag 部分，提取后去掉 "t1-" 前缀，得到 "20240323-1123"
func extractTag(image string) string {
	parts := strings.Split(image, ":")
	if len(parts) < 2 {
		return ""
	}
	tag := parts[1]
	tag = tag[3:]
	return tag
}

// sortData 对数据进行排序，并返回排序后的数据
func sortData(data *vo.ArgoCDImageResponse) *vo.ArgoCDImageResponse {
	sort.Slice(data.Images, func(i, j int) bool {
		tag1 := extractTag(data.Images[i].Image)
		tag2 := extractTag(data.Images[j].Image)
		return tag1 > tag2
	})
	return data
}

// AppImage 获取 app 镜像信息
func (argo *argoCD) appImage(appInfo *vo.ArgoCDImageRequest) (*bo.ImageInfo, error) {
	url := fmt.Sprintf("%s/api/v1/repositories/%s/appdetails", settings.Conf.ArgoCD.ArgoCDUrl,
		neturl.QueryEscape(settings.Conf.ArgoCD.ArgoCDRope))

	body, err := middleware.Request.HttpRequest("POST", "argo", url, argo.Token, appInfo)
	if err != nil {
		return nil, err
	}

	image := &bo.Image{}

	if err := json.Unmarshal(body, image); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return nil, err
	}
	imageInfo := &bo.ImageInfo{
		ID: appInfo.AppID,
	}
	for _, value := range image.Helm.Parameters {
		if value.Name == "imageUrl" {
			imageInfo.Image = value.Value
		}
	}

	return imageInfo, nil
}

// Rollback 服务回滚
func (argo *argoCD) Rollback(name string, rollbackID int) error {
	url := fmt.Sprintf("%s/api/v1/applications/%s/rollback", settings.Conf.ArgoCD.ArgoCDUrl, name)
	data := &po.ArgoCDRolloutInfo{
		ID: rollbackID,
	}
	body, err := middleware.Request.HttpRequest("POST", "argo", url, argo.Token, data)
	if err != nil {
		return err
	}
	appInfo := &bo.AppInfo{}
	if err := json.Unmarshal(body, appInfo); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return err
	}

	return nil
}

// Log 服务日志
func (argo *argoCD) Log(appName, namespace string) (string, error) {
	padName, err := argo.resourceTree(appName)
	if err != nil {
		return "", err
	}
	name := strings.SplitN(appName, "-", 2)
	url := fmt.Sprintf("%s/api/v1/applications/%s/logs?container=%s-container&namespace=%s&follow=false&podName=%s&tailLines=%s",
		settings.Conf.ArgoCD.ArgoCDUrl, appName, name[1], namespace, padName, settings.Conf.ArgoCD.ArgoCDLogTailLines)

	body, err := middleware.Request.HttpRequest("GET", "argo", url, argo.Token, nil)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(strings.NewReader(string(body)))

	var logInfo string
	for {
		log := &bo.Log{}
		if err := decoder.Decode(log); err != nil {
			// 如果遇到 EOF 说明读取完毕
			if err.Error() == "EOF" {
				break
			}
			return "", err
		}
		logInfo += fmt.Sprintf("%+v\r\n", log.Result.Content)
	}

	return logInfo, nil
}

// resourceTree
func (argo *argoCD) resourceTree(name string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/applications/%s/resource-tree", settings.Conf.ArgoCD.ArgoCDUrl, name)
	body, err := middleware.Request.HttpRequest("GET", "argo", url, argo.Token, nil)
	if err != nil {
		return "", err
	}

	rTree := &bo.ResourceTree{}
	if err := json.Unmarshal(body, rTree); err != nil {
		zap.L().Error("JSON 解析失败: " + err.Error())
		return "", err
	}

	var podName = ""

	for _, value := range rTree.Nodes {
		if value.Kind == "Pod" {
			podName = value.Name
		}
	}

	return podName, nil
}
